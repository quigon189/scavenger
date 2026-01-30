package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type TeacherHandlers struct {
	teacherService *services.TeacherService
	authMiddleware *middlewares.AuthMiddleware
}

func NewTeacherHandlers(teacherService *services.TeacherService, authMiddleware *middlewares.AuthMiddleware) *TeacherHandlers {
	return &TeacherHandlers{
		teacherService: teacherService,
		authMiddleware: authMiddleware,
	}
}

func (h *TeacherHandlers) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := parseBody(r, &teacher); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err := h.teacherService.CreateTeacher(r.Context(), &teacher)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, teacher)
}

func (h *TeacherHandlers) GetTeacher(w http.ResponseWriter, r *http.Request) {
	teacherID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	teacher, err := h.teacherService.GetTeacher(r.Context(), teacherID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, teacher)
}

func (h *TeacherHandlers) GetTeacherWithDisciplines(w http.ResponseWriter, r *http.Request) {
	teacherID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	teacher, err := h.teacherService.GetTeacherWithDisciplines(r.Context(), teacherID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, teacher)
}

func (h *TeacherHandlers) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.teacherService.GetAllTeachers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, teachers)
}

func (h *TeacherHandlers) RegisterTeacherRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/teachers", h.authMiddleware.RequireAuth(h.CreateTeacher))
	router.HandleFunc("GET /api/teachers", h.authMiddleware.RequireAuth(h.GetAllTeachers))
	router.HandleFunc("GET /api/teachers/{id}", h.authMiddleware.RequireAuth(h.GetTeacher))
	router.HandleFunc("GET /api/teachers/{id}/with-disciplines", h.authMiddleware.RequireAuth(h.GetTeacherWithDisciplines))
}
