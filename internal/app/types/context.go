package types

import (
	"TuruBot-Go/internal/app/utils"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type BotContext struct {
	Context   context.Context
	WAC       Messenger
	RawWAC    *whatsmeow.Client
	Event     *events.Message
	Pool      WorkerPool
	queueMode QueueMode
}

type QueueMode int

const (
	QueueNonBlocking QueueMode = iota // default
	QueueBlocking
)

// sendMessage this is internal function that control how message is sent (blocking or non-blocking)
func (c *BotContext) sendMessage(ctx context.Context, chatJID types.JID, msg *waE2E.Message) error {
	switch {
	case c.queueMode == QueueBlocking:
		return c.WAC.EnqueueMessage(ctx, chatJID, msg)
	default:
		return c.WAC.EnqueueMessageNonBlocking(ctx, chatJID, msg)
	}
}

func (c *BotContext) SetQueueMode(mode QueueMode) {
	c.queueMode = mode
}

func (c *BotContext) GetMessageString() string {
	msg := c.unwrapMessage(c.Event.Message)

	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()
	case msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.GetText() != "":
		return msg.ExtendedTextMessage.GetText()
	case msg.ImageMessage != nil && msg.ImageMessage.GetCaption() != "":
		return msg.ImageMessage.GetCaption()
	case msg.VideoMessage != nil && msg.VideoMessage.GetCaption() != "":
		return msg.VideoMessage.GetCaption()
	case msg.DocumentMessage != nil && msg.DocumentMessage.GetCaption() != "":
		return msg.DocumentMessage.GetCaption()
	default:
		return ""
	}
}

func (c *BotContext) GetMessageSender() types.JID {
	return c.Event.Info.Chat
}

func (c *BotContext) Reply(msg string) error {
	message := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(msg),
			ContextInfo: &waE2E.ContextInfo{
				Expiration:    proto.Uint32(c.GetExpiration()),
				StanzaID:      proto.String(c.Event.Info.ID),
				Participant:   proto.String(c.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: c.Event.Message,
			},
		},
	}

	return c.sendMessage(c.Context, c.GetMessageSender(), message)
}

func (c *BotContext) GetEventMessageJson() (string, error) {
	bytes, err := sonic.MarshalIndent(c.Event, "", "	")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *BotContext) unwrapMessage(msg *waE2E.Message) *waE2E.Message {
	if msg == nil {
		return nil
	}

	if msg.GetDeviceSentMessage() != nil {
		return c.unwrapMessage(msg.GetDeviceSentMessage().Message)
	}

	if msg.GetEphemeralMessage() != nil {
		return c.unwrapMessage(msg.GetEphemeralMessage().Message)
	}

	if msg.GetViewOnceMessage() != nil {
		return c.unwrapMessage(msg.GetViewOnceMessage().Message)
	}

	return msg
}

func (c *BotContext) GetImageMessage() *waE2E.ImageMessage {
	msg := c.unwrapMessage(c.Event.Message)

	switch {
	case msg.GetImageMessage() != nil:
		return msg.GetImageMessage()
	case msg.ExtendedTextMessage != nil &&
		msg.ExtendedTextMessage.ContextInfo != nil &&
		msg.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.ExtendedTextMessage.ContextInfo.QuotedMessage.GetImageMessage() != nil:
		return msg.ExtendedTextMessage.ContextInfo.QuotedMessage.GetImageMessage()
	case msg.ViewOnceMessage != nil &&
		msg.ViewOnceMessage.Message != nil &&
		msg.ViewOnceMessage.Message.GetImageMessage() != nil:
		return msg.ViewOnceMessage.Message.GetImageMessage()
	case msg.ViewOnceMessageV2 != nil &&
		msg.ViewOnceMessageV2.Message != nil &&
		msg.ViewOnceMessageV2.Message.GetImageMessage() != nil:
		return msg.ViewOnceMessageV2.Message.GetImageMessage()
	default:
		return nil
	}
}

func (c *BotContext) Download(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
	return c.WAC.Download(ctx, msg)
}

func (c *BotContext) Upload(ctx context.Context, plaintext []byte, appInfo whatsmeow.MediaType) (whatsmeow.UploadResponse, error) {
	return c.WAC.Upload(ctx, plaintext, appInfo)
}

func (c *BotContext) GetExpiration() uint32 {
	if c.Event.Info.IsGroup {
		return 30 * 86400
	} else {
		return 30 * 86400
	}
}

func (c *BotContext) ReplyWithSticker(sticker *ImageSticker) error {
	upload, err := c.Upload(c.Context, sticker.Image, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	message := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(upload.URL),
			FileSHA256:    upload.FileSHA256,
			FileEncSHA256: upload.FileEncSHA256,
			MediaKey:      upload.MediaKey,
			Mimetype:      proto.String("image/webp"),
			DirectPath:    &upload.DirectPath,
			FileLength:    &upload.FileLength,
			PngThumbnail:  sticker.PNGThumbnail,
			IsAnimated:    proto.Bool(false),
			ContextInfo: &waE2E.ContextInfo{
				Expiration:    proto.Uint32(c.GetExpiration()),
				StanzaID:      proto.String(c.Event.Info.ID),
				Participant:   proto.String(c.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: c.Event.Message,
			},
		},
	}

	return c.sendMessage(c.Context, c.GetMessageSender(), message)
}

func (c *BotContext) ReplyWithImage(image []byte, mimetype, caption string) error {
	thumbnail, err := utils.GenerateThumbnail(image, 120, false)
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %v", err)
	}

	upload, err := c.Upload(c.Context, image, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}

	message := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &upload.URL,
			Mimetype:      proto.String(mimetype),
			Caption:       proto.String(caption),
			FileSHA256:    upload.FileSHA256,
			FileEncSHA256: upload.FileEncSHA256,
			FileLength:    &upload.FileLength,
			MediaKey:      upload.MediaKey,
			DirectPath:    &upload.DirectPath,
			JPEGThumbnail: thumbnail,
			ContextInfo: &waE2E.ContextInfo{
				Expiration:    proto.Uint32(c.GetExpiration()),
				StanzaID:      proto.String(c.Event.Info.ID),
				Participant:   proto.String(c.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: c.Event.Message,
			},
		},
	}

	return c.sendMessage(c.Context, c.GetMessageSender(), message)
}
