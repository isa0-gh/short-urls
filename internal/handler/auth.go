package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isa0-gh/short-urls/internal/dto"
	"github.com/isa0-gh/short-urls/internal/response"
	"github.com/isa0-gh/short-urls/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewAuthHandler(userService *service.UserService, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	user, err := h.userService.CreateUser(r.Context(), req.Username, req.Email, string(hashBytes))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	response.JSON(w, http.StatusCreated, dto.UserDetailsResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})

}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, err := h.authService.CheckPassword(r.Context(), req.Password, req.Email)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create token")
		return
	}
	response.JSON(w, http.StatusCreated, dto.LoginResponse{
		Token: token,
	})
}
