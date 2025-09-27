package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
)

func GenerateThumbnail(imgData []byte, size int, bgTransparent bool) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	scale := float64(size) / float64(max(w, h))
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))

	if !bgTransparent {
		draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
	}

	offsetX := (size - newW) / 2
	offsetY := (size - newH) / 2

	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH), resized, image.Point{}, draw.Over)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, canvas); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
