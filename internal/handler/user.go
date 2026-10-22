package handler

import (
	"net/http"

	"github.com/isa0-gh/short-urls/internal/dto"
	"github.com/isa0-gh/short-urls/internal/middleware"
	"github.com/isa0-gh/short-urls/internal/response"
	"github.com/isa0-gh/short-urls/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserDetails returns the authenticated user's profile.
// @Summary Get user profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserDetailsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/profile [get]
func (h *UserHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch user details")
		return
	}
	response.JSON(w, http.StatusOK, dto.UserDetailsResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})
}
