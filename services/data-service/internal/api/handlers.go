package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response представляет стандартный ответ API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// writeJSON записывает JSON ответ
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeError записывает ошибку в формате JSON
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

// writeSuccess записывает успешный ответ в формате JSON
func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// parseIntParam парсит целочисленный параметр из URL
func parseIntParam(r *http.Request, param string) (int, error) {
	value := r.PathValue(param)
	if value == "" {
		return 0, fmt.Errorf("parameter %s is required", param)
	}
	
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", param, value)
	}
	
	return id, nil
}

// getSessionID извлекает session_id из заголовков или куков
func getSessionID(r *http.Request) (string, error) {
	// Пробуем получить из заголовка
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID != "" {
		return sessionID, nil
	}
	
	// Пробуем получить из куки
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	
	return "", errors.New("session ID not found")
}

// parseBody парсит тело запроса в структуру
func parseBody(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	
	defer r.Body.Close()
	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	return nil
}

// validateFileUpload проверяет загружаемый файл
func validateFileUpload(r *http.Request) (string, string, int64, error) {
	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return "", "", 0, errors.New("content type must be multipart/form-data")
	}
	
	// Парсим multipart форму
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to parse multipart form: %v", err)
	}
	
	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", "", 0, fmt.Errorf("file not found in request: %v", err)
	}
	defer file.Close()
	
	// Проверяем размер файла
	if header.Size > 50<<20 { // 50 MB
		return "", "", 0, errors.New("file size exceeds 50 MB limit")
	}
	
	// Проверяем тип файла
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"application/zip": true,
		"text/plain":      true,
		"text/markdown":   true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
	}
	
	if !allowedTypes[header.Header.Get("Content-Type")] {
		return "", "", 0, errors.New("file type not allowed")
	}
	
	return header.Filename, header.Header.Get("Content-Type"), header.Size, nil
}

// Middleware для логирования
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Создаем ResponseWriter для отслеживания статуса
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		
		log.Printf("[%s] %s %s - %d - %v", r.Method, r.URL.Path, r.RemoteAddr, rw.statusCode, duration)
	})
}

// responseWriter для отслеживания статуса ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware для CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника (в продакшене нужно указать конкретные домены)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 часа
		
		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
