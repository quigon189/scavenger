package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"data-service/internal/middlewares"
	"data-service/internal/services"
)

type FileHandlers struct {
	fileService *services.FileService
	authMiddleware *middlewares.AuthMiddleware
}

func NewFileHandlers(fileService *services.FileService, authMiddleware *middlewares.AuthMiddleware) *FileHandlers {
	return &FileHandlers{
		fileService: fileService,
		authMiddleware: authMiddleware,
	}
}

func (h *FileHandlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	filename, contentType, size, err := validateFileUpload(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	fileType := r.FormValue("file_type")

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("failed to get file: %v", err))
		return
	}
	defer file.Close()

	uploadedFile, err := h.fileService.UploadFile(r.Context(), filename, fileType, contentType, size, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, uploadedFile)
}

// GetFileInfo возвращает информацию о файле с подписанной ссылкой
func (h *FileHandlers) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	fileID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	file, err := h.fileService.GetFileWithSignedURL(r.Context(), fileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, file)
}

// DownloadFile скачивает файл напрямую
func (h *FileHandlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	file, reader, err := h.fileService.GetFile(r.Context(), fileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
	}
}

// GenerateSignedURL генерирует подписанную ссылку для файла
func (h *FileHandlers) GenerateSignedURL(w http.ResponseWriter, r *http.Request) {
	fileID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var expires time.Duration
	expiresStr := r.URL.Query().Get("expires")
	if expiresStr != "" {
		expiresSec, err := strconv.ParseInt(expiresStr, 10, 64)
		if err == nil {
			expires = time.Duration(expiresSec) * time.Second
		}
	}

	signedURL, err := h.fileService.GenerateSignedURL(r.Context(), fileID, expires)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, map[string]string{
		"signed_url": signedURL,
	})
}

// DeleteFile удаляет файл
func (h *FileHandlers) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = h.fileService.DeleteFile(r.Context(), fileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, map[string]string{
		"message": "File deleted successfully",
	})
}

// RegisterFileRoutes регистрирует маршруты для работы с файлами
func (h *FileHandlers) RegisterFileRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/files/upload", h.authMiddleware.RequireAuth(h.UploadFile))
	router.HandleFunc("GET /api/files/{id}", h.authMiddleware.RequireAuth(h.GetFileInfo))
	router.HandleFunc("GET /api/files/{id}/download", h.authMiddleware.RequireAuth(h.DownloadFile))
	router.HandleFunc("GET /api/files/{id}/url", h.authMiddleware.RequireAuth(h.GenerateSignedURL))
	router.HandleFunc("DELETE /api/files/{id}", h.authMiddleware.RequireAuth(h.DeleteFile))
}
