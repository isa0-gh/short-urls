package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/isa0-gh/short-urls/internal/dto"
	"github.com/isa0-gh/short-urls/internal/middleware"
	"github.com/isa0-gh/short-urls/internal/response"
	"github.com/isa0-gh/short-urls/internal/service"
)

type ShortLinkHandler struct {
	linkService *service.ShortLinkService
	// telemetryService *service.TelemetryService
}

func NewShortLinkHandler(linkService *service.ShortLinkService) *ShortLinkHandler {
	return &ShortLinkHandler{
		linkService: linkService,
	}
}

// RedirectToDestination redirects a short code to its destination URL.
// @Summary Follow a short link
// @Tags shortlinks
// @Produce plain
// @Param shortCode path string true "Short code"
// @Success 302
// @Failure 404 {string} string "Short link not found"
// @Router /s/{shortCode} [get]
func (h *ShortLinkHandler) RedirectToDestination(w http.ResponseWriter, r *http.Request) {
	destination, err := h.linkService.GetDestinationUrl(r.Context(), r.PathValue("shortCode"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 - Not Found")
		return
	}
	// TODO: implement telemetry or remove telemetry
	http.Redirect(w, r, destination, http.StatusFound)
}

// CreateShortlink creates a short link for the authenticated user.
// @Summary Create a short link
// @Tags shortlinks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateShortLinkRequest true "Short link details"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/create [post]
func (h *ShortLinkHandler) CreateShortlink(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	h.linkService.CreateShortLink(r.Context(), userID, req.ShortCode, req.DestinationURL)
}

// ListShortLinks returns the authenticated user's short links.
// @Summary List short links
// @Tags shortlinks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ShortLinkObject
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/links [get]
func (h *ShortLinkHandler) ListShortLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	data, err := h.linkService.ListShortLinks(r.Context(), userID, 50)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch links")
		return
	}
	var links []dto.ShortLinkObject
	for _, link := range data {
		links = append(links, dto.ShortLinkObject{
			ID:             link.ID.String(),
			ShortCode:      link.ShortCode,
			DestinationURL: link.DestinationUrl,
			CreatedAt:      link.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	response.JSON(w, http.StatusOK, links)
}

// DeleteShortLink deletes a short link owned by the authenticated user.
// @Summary Delete a short link
// @Tags shortlinks
// @Produce json
// @Security BearerAuth
// @Param linkID path string true "Link ID (UUID)" format(uuid)
// @Success 200
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/link/{linkID} [delete]
func (h *ShortLinkHandler) DeleteShortLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	err := h.linkService.DeleteShortLink(r.Context(), userID, r.PathValue("linkID"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete link")
		return
	}
	response.JSON(w, http.StatusOK, nil)
}

// UpdateShortLink changes a short link's destination URL.
// @Summary Update a short link
// @Tags shortlinks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param linkID path string true "Link ID (UUID)" format(uuid)
// @Param request body dto.UpdateDestinationURLRequest true "Updated destination URL"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/link/{linkID} [post]
func (h *ShortLinkHandler) UpdateShortLink(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDestinationURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	err := h.linkService.UpdateDestination(r.Context(), userID, r.PathValue("linkID"), req.NewURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update link")
		return
	}
	response.JSON(w, http.StatusOK, nil)
}
