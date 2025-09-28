package types

import (
	"context"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

type Messenger interface {
	EnqueueMessage(ctx context.Context, chatJID types.JID, msg *waE2E.Message) error
	EnqueueMessageNonBlocking(ctx context.Context, chatJID types.JID, msg *waE2E.Message) error
	Download(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error)
	Upload(ctx context.Context, plaintext []byte, appInfo whatsmeow.MediaType) (whatsmeow.UploadResponse, error)
}
