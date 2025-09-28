package whatsapp

import (
	"context"
	"go.mau.fi/whatsmeow"
)

func (wa *WAClient) Upload(ctx context.Context, plaintext []byte, appInfo whatsmeow.MediaType) (whatsmeow.UploadResponse, error) {
	return wa.Client.Upload(ctx, plaintext, appInfo)
}

func (wa *WAClient) Download(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
	return wa.Client.Download(ctx, msg)
}
