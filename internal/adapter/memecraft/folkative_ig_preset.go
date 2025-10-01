package memecraft

import (
	"TuruBot-Go/internal/port"
	"context"
)

func (c *MemeCraft) GenerateFolkativeIg(ctx context.Context, overlayImage string, args map[string]string) (*port.MemeCraftResponse, error) {
	return c.generate(ctx, c.BaseUrl+"/presets/folkative-ig-preset/memes", overlayImage, args)
}
