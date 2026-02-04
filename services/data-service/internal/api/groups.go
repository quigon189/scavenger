package api

import (
	"net/http"

	"data-service/internal/middlewares"
	"data-service/internal/models"
	"data-service/internal/services"
)

type GroupHandlers struct {
	groupService *services.GroupService
	authMiddleware *middlewares.AuthMiddleware
}

func NewGroupHandlers(groupService *services.GroupService, authMiddleware *middlewares.AuthMiddleware) *GroupHandlers {
	return &GroupHandlers{
		groupService: groupService,
		authMiddleware: authMiddleware,
	}
}

func (h *GroupHandlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := parseBody(r, &group); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err := h.groupService.CreateGroup(r.Context(), &group)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, group)
}

func (h *GroupHandlers) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	group, err := h.groupService.GetGroup(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, group)
}

func (h *GroupHandlers) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	var group models.Group
	if err := parseBody(r, &group); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	group.ID = groupID
	
	err = h.groupService.UpdateGroup(r.Context(), &group)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, group)
}

func (h *GroupHandlers) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	
	err = h.groupService.DeleteGroup(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, map[string]string{
		"message": "Group deleted successfully",
	})
}

func (h *GroupHandlers) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupService.GetAllGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	writeSuccess(w, groups)
}

func (h *GroupHandlers) RegisterGroupRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/groups", h.authMiddleware.RequireAuth(h.CreateGroup))
	router.HandleFunc("GET /api/groups", h.GetAllGroups)
	router.HandleFunc("GET /api/groups/{id}", h.authMiddleware.RequireAuth(h.GetGroup))
	router.HandleFunc("PUT /api/groups/{id}", h.authMiddleware.RequireAuth(h.UpdateGroup))
	router.HandleFunc("DELETE /api/groups/{id}", h.authMiddleware.RequireAuth(h.DeleteGroup))
}
