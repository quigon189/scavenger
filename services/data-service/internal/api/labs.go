package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type LabHandlers struct {
	labService *services.LabService
	authMiddleware *middlewares.AuthMiddleware
}

func NewLabHandlers(labService *services.LabService, authMiddleware *middlewares.AuthMiddleware) *LabHandlers {
	return &LabHandlers{
		labService: labService,
		authMiddleware: authMiddleware,
	}
}

func (h *LabHandlers) CreateLab(w http.ResponseWriter, r *http.Request) {
	var lab models.Lab
	if err := parseBody(r, &lab); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err := h.labService.CreateLab(r.Context(), &lab)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, lab)
}

func (h *LabHandlers) GetLab(w http.ResponseWriter, r *http.Request) {
	labID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	lab, err := h.labService.GetLab(r.Context(), labID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, lab)
}

func (h *LabHandlers) GetLabsByDiscipline(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := parseIntParam(r, "discipline_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	labs, err := h.labService.GetLabsByDiscipline(r.Context(), disciplineID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, labs)
}

func (h *LabHandlers) UpdateLab(w http.ResponseWriter, r *http.Request) {
	labID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	var lab models.Lab
	if err := parseBody(r, &lab); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	lab.ID = labID
	
	err = h.labService.UpdateLab(r.Context(), &lab)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, lab)
}

func (h *LabHandlers) DeleteLab(w http.ResponseWriter, r *http.Request) {
	labID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err = h.labService.DeleteLab(r.Context(), labID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, map[string]string{
		"message": "Lab deleted successfully",
	})
}

func (h *LabHandlers) AddFileToLab(w http.ResponseWriter, r *http.Request) {
	labID, err := parseIntParam(r, "lab_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	var request struct {
		FileID int `json:"file_id"`
	}
	
	if err := parseBody(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err = h.labService.AddFileToLab(r.Context(), labID, request.FileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, map[string]string{
		"message": "File added to lab successfully",
	})
}

func (h *LabHandlers) RegisterLabRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/labs", h.authMiddleware.RequireAuth(h.CreateLab))
	router.HandleFunc("GET /api/labs/{id}", h.authMiddleware.RequireAuth(h.GetLab))
	router.HandleFunc("GET /api/labs/discipline/{discipline_id}", h.authMiddleware.RequireAuth(h.GetLabsByDiscipline))
	router.HandleFunc("PUT /api/labs/{id}", h.authMiddleware.RequireAuth(h.UpdateLab))
	router.HandleFunc("DELETE /api/labs/{id}", h.authMiddleware.RequireAuth(h.DeleteLab))
	router.HandleFunc("POST /api/labs/{lab_id}/files", h.authMiddleware.RequireAuth(h.AddFileToLab))
}
