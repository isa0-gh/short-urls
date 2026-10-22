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
