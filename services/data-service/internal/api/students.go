package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type StudentHandlers struct {
	studentService *services.StudentService
	authMiddleware *middlewares.AuthMiddleware
}

func NewStudentHandlers(studentService *services.StudentService, authMiddleware *middlewares.AuthMiddleware) *StudentHandlers {
	return &StudentHandlers{
		studentService: studentService,
		authMiddleware: authMiddleware,
	}
}

func (h *StudentHandlers) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := parseBody(r, &student); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err := h.studentService.CreateStudent(r.Context(), &student)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, student)
}

func (h *StudentHandlers) GetStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	student, err := h.studentService.GetStudent(r.Context(), studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, student)
}

func (h *StudentHandlers) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var student models.Student
	if err := parseBody(r, &student); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	student.ID = studentID

	err = h.studentService.UpdateStudent(r.Context(), &student)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, student)
}

func (h *StudentHandlers) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = h.studentService.DeleteStudent(r.Context(), studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, map[string]string{
		"message": "Student deleted successfully",
	})
}

func (h *StudentHandlers) GetStudentsByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseIntParam(r, "group_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	students, err := h.studentService.GetStudentsByGroup(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, students)
}

func (h *StudentHandlers) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.studentService.GetAllStudents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, students)
}

// RegisterStudentRoutes регистрирует маршруты для работы со студентами
func (h *StudentHandlers) RegisterStudentRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/students", h.authMiddleware.RequireAuth(h.CreateStudent))
	router.HandleFunc("GET /api/students", h.authMiddleware.RequireAuth(h.GetAllStudents))
	router.HandleFunc("GET /api/students/{id}", h.authMiddleware.RequireAuth(h.GetStudent))
	router.HandleFunc("PUT /api/students/{id}", h.authMiddleware.RequireAuth(h.UpdateStudent))
	router.HandleFunc("DELETE /api/students/{id}", h.authMiddleware.RequireAuth(h.DeleteStudent))
	router.HandleFunc("GET /api/students/group/{group_id}", h.authMiddleware.RequireAuth(h.GetStudentsByGroup))
}
