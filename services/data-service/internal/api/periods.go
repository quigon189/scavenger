package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type PeriodHandlers struct {
	periodService *services.PeriodService
	authMiddleware *middlewares.AuthMiddleware
}

func NewPeriodHandlers(periodService *services.PeriodService, authMiddleware *middlewares.AuthMiddleware) *PeriodHandlers {
	return &PeriodHandlers{
		periodService: periodService,
		authMiddleware: authMiddleware,
	}
}

func (h *PeriodHandlers) CreatePeriod(w http.ResponseWriter, r *http.Request) {
	var period models.Period
	if err := parseBody(r, &period); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err := h.periodService.CreatePeriod(r.Context(), &period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, period)
}

func (h *PeriodHandlers) GetPeriod(w http.ResponseWriter, r *http.Request) {
	periodID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	period, err := h.periodService.GetPeriod(r.Context(), periodID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, period)
}

func (h *PeriodHandlers) GetActivePeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.periodService.GetActivePeriod(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, period)
}

func (h *PeriodHandlers) UpdatePeriod(w http.ResponseWriter, r *http.Request) {
	periodID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	var period models.Period
	if err := parseBody(r, &period); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	period.ID = periodID
	
	err = h.periodService.UpdatePeriod(r.Context(), &period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, period)
}

func (h *PeriodHandlers) DeletePeriod(w http.ResponseWriter, r *http.Request) {
	periodID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err = h.periodService.DeletePeriod(r.Context(), periodID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, map[string]string{
		"message": "Period deleted successfully",
	})
}

func (h *PeriodHandlers) GetAllPeriods(w http.ResponseWriter, r *http.Request) {
	periods, err := h.periodService.GetAllPeriods(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, periods)
}

func (h *PeriodHandlers) RegisterPeriodRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/periods", h.authMiddleware.RequireAuth(h.CreatePeriod))
	router.HandleFunc("GET /api/periods", h.authMiddleware.RequireAuth(h.GetAllPeriods))
	router.HandleFunc("GET /api/periods/active", h.authMiddleware.RequireAuth(h.GetActivePeriod))
	router.HandleFunc("GET /api/periods/{id}", h.authMiddleware.RequireAuth(h.GetPeriod))
	router.HandleFunc("PUT /api/periods/{id}", h.authMiddleware.RequireAuth(h.UpdatePeriod))
	router.HandleFunc("DELETE /api/periods/{id}", h.authMiddleware.RequireAuth(h.DeletePeriod))
}
