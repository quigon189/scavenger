package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/services"
)

type ReportHandlers struct {
	reportService  *services.ReportService
	authMiddleware *middlewares.AuthMiddleware
}

func NewReportHandlers(reportService *services.ReportService, authMiddleware *middlewares.AuthMiddleware) *ReportHandlers {
	return &ReportHandlers{
		reportService:  reportService,
		authMiddleware: authMiddleware,
	}
}

func (h *ReportHandlers) CreateReport(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LabID   int    `json:"lab_id"`
		Comment string `json:"comment"`
	}

	if err := parseBody(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	report, err := h.reportService.CreateReport(r.Context(), request.LabID, request.Comment)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, report)
}

func (h *ReportHandlers) GetReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	report, err := h.reportService.GetReport(r.Context(), reportID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, report)
}

func (h *ReportHandlers) GetReportsByLab(w http.ResponseWriter, r *http.Request) {
	labID, err := parseIntParam(r, "lab_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reports, err := h.reportService.GetReportsByLab(r.Context(), labID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, reports)
}

func (h *ReportHandlers) GetReportsByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIntParam(r, "student_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reports, err := h.reportService.GetReportsByStudent(r.Context(), studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, reports)
}

func (h *ReportHandlers) GradeReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var request struct {
		Grade       int    `json:"grade"`
		TeacherNote string `json:"teacher_note"`
	}

	if err := parseBody(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = h.reportService.GradeReport(r.Context(), reportID, request.Grade, request.TeacherNote)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, map[string]string{
		"message": "Report graded successfully",
	})
}

func (h *ReportHandlers) AddFileToReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := parseIntParam(r, "report_id")
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

	err = h.reportService.AddFileToReport(r.Context(), reportID, request.FileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, map[string]string{
		"message": "File added to report successfully",
	})
}

func (h *ReportHandlers) RegisterReportRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/reports", h.authMiddleware.RequireAuth(h.CreateReport))
	router.HandleFunc("GET /api/reports/{id}", h.authMiddleware.RequireAuth(h.GetReport))
	router.HandleFunc("GET /api/reports/lab/{lab_id}", h.authMiddleware.RequireAuth(h.GetReportsByLab))
	router.HandleFunc("GET /api/reports/student/{student_id}", h.authMiddleware.RequireAuth(h.GetReportsByStudent))
	router.HandleFunc("POST /api/reports/{id}/grade", h.authMiddleware.RequireAuth(h.GradeReport))
	router.HandleFunc("POST /api/reports/{report_id}/files", h.authMiddleware.RequireAuth(h.AddFileToReport))
}
