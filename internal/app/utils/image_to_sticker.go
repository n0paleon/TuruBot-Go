package utils

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/oklog/ulid/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func ImageToSticker(imgData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode source image: %w", err)
	}

	const size = 512

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	scale := float64(size) / float64(max(w, h))
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))

	offsetX := (size - newW) / 2
	offsetY := (size - newH) / 2

	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH), resized, image.Point{}, draw.Over)

	buf := new(bytes.Buffer)
	if err := webp.Encode(buf, canvas, &webp.Options{
		Lossless: true,
		Quality:  100,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode webp lossless: %w", err)
	}

	return buf.Bytes(), nil
}

func ImageToStickerViaFFMPEG(imgData []byte) ([]byte, error) {
	uid := ulid.Make().String()
	inFile := filepath.Join(os.TempDir(), "wa_sticker_in_"+uid+".png")
	outFile := filepath.Join(os.TempDir(), "wa_sticker_out_"+uid+".webp")

	if err := os.WriteFile(inFile, imgData, 0600); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}
	defer func() { _ = os.Remove(inFile) }()
	defer func() { _ = os.Remove(outFile) }()

	err := ffmpeg.
		Input(inFile).
		Filter("scale", ffmpeg.Args{"512:512:force_original_aspect_ratio=decrease"}).
		Filter("pad", ffmpeg.Args{"512:512:(ow-iw)/2:(oh-ih)/2:color=0x00000000"}). // transparent padding
		Output(outFile,
			ffmpeg.KwArgs{
				"pix_fmt":           "yuva420p",
				"lossless":          "1",
				"compression_level": "6",
				"q:v":               "80", // quality
				"f":                 "webp",
				"preset":            "picture",
			},
		).
		OverWriteOutput().
		Silent(true).
		Run()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	result, err := os.ReadFile(outFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return result, nil
}
