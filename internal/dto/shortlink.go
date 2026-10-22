package dto

type CreateShortLinkRequest struct {
	ShortCode      string `json:"short_code" binding:"required" example:"abc123"`
	DestinationURL string `json:"destination_url" binding:"required" example:"https://example.com/article"`
}

type ShortLinkObject struct {
	ID             string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ShortCode      string `json:"short_code" binding:"required" example:"abc123"`
	DestinationURL string `json:"destination_url" binding:"required" example:"https://example.com/article"`
	CreatedAt      string `json:"created_at" binding:"required" example:"2026-10-22T12:00:00Z"`
}

type UpdateDestinationURLRequest struct {
	NewURL string `json:"destination_url" binding:"required" example:"https://example.com/updated-article"`
}
