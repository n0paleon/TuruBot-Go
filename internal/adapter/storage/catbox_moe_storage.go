package storage

import (
	"TuruBot-Go/internal/port"
	"bytes"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"strings"
)

type CatboxMoeProvider struct {
	ApiUrl string
	Client *resty.Client
}

func (s *CatboxMoeProvider) Upload(ctx context.Context, data io.Reader) (*port.UploadResult, error) {
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return s.UploadBytes(ctx, dataBytes)
}

func (s *CatboxMoeProvider) UploadBytes(ctx context.Context, data []byte) (*port.UploadResult, error) {
	contentType := http.DetectContentType(data)
	filename := GenerateFileName(contentType)

	resp, err := s.Client.R().
		SetContext(ctx).
		SetFileReader("fileToUpload", filename, bytes.NewReader(data)).
		SetFormData(map[string]string{
			"reqtype": "fileupload",
		}).
		Post(s.ApiUrl)

	if err != nil {
		return nil, fmt.Errorf("upload file error: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status code %d", resp.StatusCode())
	}

	directURL := strings.TrimSpace(string(resp.Body()))
	contentLength := int64(len(data))

	return &port.UploadResult{
		DirectURL:     directURL,
		ContentType:   contentType,
		Bytes:         contentLength,
		BytesReadable: ByteCountSI(contentLength),
	}, nil
}

func (s *CatboxMoeProvider) GetStorageName() string {
	return "catbox.moe"
}

func (s *CatboxMoeProvider) DownloadToBytes(ctx context.Context, url string) (*port.DownloadResult, error) {
	resp, err := s.Client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("download error: %v", err)
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("download failed with status code %d", resp.StatusCode())
	}

	data, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	contentLength := int64(len(data))

	return &port.DownloadResult{
		DirectURL:     url,
		ByteContent:   data,
		ContentType:   contentType,
		Bytes:         contentLength,
		BytesReadable: ByteCountSI(contentLength),
	}, nil
}

func NewCatboxMoeStorage() *CatboxMoeProvider {
	return &CatboxMoeProvider{
		ApiUrl: "https://catbox.moe/user/api.php",
		Client: resty.New(),
	}
}
