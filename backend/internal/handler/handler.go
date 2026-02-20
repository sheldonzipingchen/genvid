package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/genvid/backend/internal/middleware"
	"github.com/genvid/backend/internal/model"
	"github.com/genvid/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required", nil)
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters", nil)
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == service.ErrEmailExists {
			respondError(w, http.StatusConflict, "EMAIL_EXISTS", "Email already registered", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user", nil)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required", nil)
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to login", nil)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required", nil)
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired refresh token", nil)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	profile, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get profile", nil)
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(profile))
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	profile, err := h.authService.UpdateProfile(r.Context(), userID, updates)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update profile", nil)
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(profile))
}

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req model.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.ProductName == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Product name is required", nil)
		return
	}

	project, err := h.projectService.Create(r.Context(), userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create project", nil)
		return
	}

	respondJSON(w, http.StatusCreated, model.SuccessResponse(project))
}

func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Project ID is required", nil)
		return
	}

	project, err := h.projectService.GetByID(r.Context(), projectID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(project))
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	projects, total, err := h.projectService.ListByUser(r.Context(), userID, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list projects", nil)
		return
	}

	meta := &model.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}

	respondJSON(w, http.StatusOK, model.SuccessResponseWithMeta(projects, meta))
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Project ID is required", nil)
		return
	}

	if err := h.projectService.Delete(r.Context(), projectID, userID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) GenerateVideo(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Project ID is required", nil)
		return
	}

	var req model.GenerateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.Script == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Script is required", nil)
		return
	}

	project, err := h.projectService.GenerateVideo(r.Context(), projectID, userID, &req)
	if err != nil {
		if err == service.ErrInsufficientCredits {
			respondError(w, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "No credits remaining", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate video", nil)
		return
	}

	respondJSON(w, http.StatusAccepted, model.SuccessResponse(project))
}

func getUserIDFromContext(r *http.Request) string {
	return middleware.GetUserID(r)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code, message string, details map[string]string) {
	respondJSON(w, status, model.ErrorResponse(code, message, details))
}
