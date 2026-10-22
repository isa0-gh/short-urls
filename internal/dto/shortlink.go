package dto

type CreateShortLinkRequest struct {
	ShortCode      string `json:"short_code"`
	DestinationURL string `json:"destination_url"`
}

type ShortLinkObject struct {
	ID             string `json:"id"`
	ShortCode      string `json:"short_code"`
	DestinationURL string `json:"destination_url"`
	CreatedAt      string `json:"created_at"`
}

type UpdateDestinationURLRequest struct {
	NewURL string `json:"destination_url"`
}
