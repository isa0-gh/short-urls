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

func (h *ShortLinkHandler) DeleteShortLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}
	err := h.linkService.DeleteShortLink(r.Context(), userID, r.PathValue("shortCode"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete link")
		return
	}
	response.JSON(w, http.StatusOK, nil)
}

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
	err := h.linkService.UpdateDestination(r.Context(), userID, r.PathValue("shortCode"), req.NewURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete link")
		return
	}
	response.JSON(w, http.StatusOK, nil)
}
