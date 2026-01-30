package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type DisciplineHandlers struct {
	disciplineService *services.DisciplineService
	authMiddleware *middlewares.AuthMiddleware
}

func NewDisciplineHandlers(disciplineService *services.DisciplineService, authMiddleware *middlewares.AuthMiddleware) *DisciplineHandlers {
	return &DisciplineHandlers{
		disciplineService: disciplineService,
		authMiddleware: authMiddleware,
	}
}

func (h *DisciplineHandlers) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	var discipline models.Discipline
	if err := parseBody(r, &discipline); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err := h.disciplineService.CreateDiscipline(r.Context(), &discipline)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, discipline)
}

func (h *DisciplineHandlers) GetDiscipline(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	discipline, err := h.disciplineService.GetDiscipline(r.Context(), disciplineID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, discipline)
}

func (h *DisciplineHandlers) GetDisciplinesByTeacher(w http.ResponseWriter, r *http.Request) {
	teacherID, err := parseIntParam(r, "teacher_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	disciplines, err := h.disciplineService.GetDisciplinesByTeacher(r.Context(), teacherID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, disciplines)
}

func (h *DisciplineHandlers) GetDisciplinesByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseIntParam(r, "group_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	disciplines, err := h.disciplineService.GetDisciplinesByGroup(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, disciplines)
}

func (h *DisciplineHandlers) UpdateDiscipline(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	var discipline models.Discipline
	if err := parseBody(r, &discipline); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	discipline.ID = disciplineID
	
	err = h.disciplineService.UpdateDiscipline(r.Context(), &discipline)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, discipline)
}

func (h *DisciplineHandlers) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = h.disciplineService.DeleteDiscipline(r.Context(), disciplineID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, map[string]string{
		"message": "Discipline deleted successfully",
	})
}

func (h *DisciplineHandlers) GetAllDisciplines(w http.ResponseWriter, r *http.Request) {
	disciplines, err := h.disciplineService.GetAllDisciplines(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, disciplines)
}

func (h *DisciplineHandlers) RegisterDisciplineRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/disciplines", h.authMiddleware.RequireAuth(h.CreateDiscipline))
	router.HandleFunc("GET /api/disciplines", h.authMiddleware.RequireAuth(h.GetAllDisciplines))
	router.HandleFunc("GET /api/disciplines/{id}", h.authMiddleware.RequireAuth(h.GetDiscipline))
	router.HandleFunc("GET /api/disciplines/teacher/{teacher_id}", h.authMiddleware.RequireAuth(h.GetDisciplinesByTeacher))
	router.HandleFunc("GET /api/disciplines/group/{group_id}", h.authMiddleware.RequireAuth(h.GetDisciplinesByGroup))
	router.HandleFunc("PUT /api/disciplines/{id}", h.authMiddleware.RequireAuth(h.UpdateDiscipline))
	router.HandleFunc("DELETE /api/disciplines/{id}", h.authMiddleware.RequireAuth(h.DeleteDiscipline))
}
