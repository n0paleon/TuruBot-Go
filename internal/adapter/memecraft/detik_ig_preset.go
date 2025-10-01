package memecraft

import (
	"TuruBot-Go/internal/port"
	"context"
)

func (c *MemeCraft) GenerateDetikIg(ctx context.Context, overlayImage string, args map[string]string) (*port.MemeCraftResponse, error) {
	return c.generate(ctx, c.BaseUrl+"/presets/kompas-ig-preset/memes", overlayImage, args)
}
