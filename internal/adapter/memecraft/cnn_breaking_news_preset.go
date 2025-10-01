package memecraft

import (
	"TuruBot-Go/internal/port"
	"context"
)

func (c *MemeCraft) GenerateCnnBreakingNews(ctx context.Context, overlayImage string, args map[string]string) (*port.MemeCraftResponse, error) {
	return c.generate(ctx, c.BaseUrl+"/presets/cnn-breaking-news-preset/memes", overlayImage, args)
}
