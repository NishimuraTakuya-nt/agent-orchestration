package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/go-chi/chi/v5"
	
	"agent-orchestration/entities"
	"agent-orchestration/use_cases"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUseCase *use_cases.UserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *use_cases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	user, err := h.userUseCase.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		switch err {
		case entities.ErrUserAlreadyExists:
			h.writeError(w, http.StatusConflict, err.Error())
		case entities.ErrUserNameRequired, entities.ErrUserEmailRequired:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}
	
	h.writeJSON(w, http.StatusCreated, user)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	user, err := h.userUseCase.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			h.writeError(w, http.StatusNotFound, err.Error())
		case entities.ErrInvalidID:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "failed to get user")
		}
		return
	}
	
	h.writeJSON(w, http.StatusOK, user)
}

// UpdateUser handles PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	user, err := h.userUseCase.UpdateUser(r.Context(), id, req.Name, req.Email)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			h.writeError(w, http.StatusNotFound, err.Error())
		case entities.ErrInvalidID, entities.ErrUserNameRequired, entities.ErrUserEmailRequired:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "failed to update user")
		}
		return
	}
	
	h.writeJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	err = h.userUseCase.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			h.writeError(w, http.StatusNotFound, err.Error())
		case entities.ErrInvalidID:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "failed to delete user")
		}
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.ListUsers(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	
	h.writeJSON(w, http.StatusOK, users)
}

// writeJSON writes JSON response
func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes error response
func (h *UserHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}