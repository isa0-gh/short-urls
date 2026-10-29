package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isa0-gh/short-urls/internal/dto"
	"github.com/isa0-gh/short-urls/internal/response"
	"github.com/isa0-gh/short-urls/internal/service"
	"github.com/isa0-gh/short-urls/internal/validation"
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

// CreateUser registers a user account.
// @Summary Register a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User details"
// @Success 201 {object} dto.UserDetailsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/register [post]
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validation.User(req.Username, req.Email, req.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	user, err := h.userService.CreateUser(r.Context(), req.Username, req.Email, string(hashBytes))
	if err != nil {
		response.Error(w, response.DatabaseErrorStatus(err), "Failed to create user")
		return
	}
	response.JSON(w, http.StatusCreated, dto.UserDetailsResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})

}

// LoginUser authenticates a user and creates a session cookie.
// @Summary Log in
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginUserRequest true "Login credentials"
// @Success 201
// @Header 201 {string} Set-Cookie "HttpOnly session cookie"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/login [post]
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validation.Login(req.Email, req.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := h.authService.CheckPassword(r.Context(), req.Password, req.Email)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := h.authService.GenerateToken(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create token")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	response.JSON(w, http.StatusCreated, nil)
}

// LogoutUser revokes the current session and clears its cookie.
// @Summary Log out
// @Tags auth
// @Produce json
// @Security SessionCookie
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /api/v1/logout [post]
func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(service.SessionCookieName); err == nil {
		if err := h.authService.RevokeToken(r.Context(), cookie.Value); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to log out")
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}
