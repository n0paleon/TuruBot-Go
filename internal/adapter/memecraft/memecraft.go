package memecraft

import (
	"TuruBot-Go/internal/port"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type MemeCraft struct {
	BaseUrl    string
	HttpClient *resty.Client
}

func NewMemeCraft() *MemeCraft {
	return &MemeCraft{
		BaseUrl:    "https://wj1o54jcd0.execute-api.ap-southeast-5.amazonaws.com/prod/MemeCraft",
		HttpClient: resty.New(),
	}
}

func (c *MemeCraft) generate(ctx context.Context, presetUrl, overlayImage string, args map[string]string) (*port.MemeCraftResponse, error) {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetBody(map[string]any{
			"overlay":     overlayImage,
			"resize_mode": "fill",
			"text":        args,
		}).
		Post(presetUrl)

	if err != nil {
		return nil, fmt.Errorf("error generating detik ig preset: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error generating detik ig preset, status: %s", resp.Status())
	}

	var response port.MemeCraftResponse
	if err := sonic.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("error parsing detik ig preset: %w", err)
	}

	return &response, nil
}
