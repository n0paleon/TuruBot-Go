package port

import "context"

type MemeCraftResponse struct {
	ImageUrl    string `json:"image_url"`
	ContentType string `json:"content_type"`
	Size        string `json:"size"`
}

type MemeCraft interface {
	GenerateDetikIg(ctx context.Context, overlayImage string, args map[string]string) (*MemeCraftResponse, error)
	GenerateCnnBreakingNews(ctx context.Context, overlayImage string, args map[string]string) (*MemeCraftResponse, error)
	GenerateFolkativeIg(ctx context.Context, overlayImage string, args map[string]string) (*MemeCraftResponse, error)
}
