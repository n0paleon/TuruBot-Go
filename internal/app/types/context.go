package types

import (
	"context"
	"github.com/bytedance/sonic"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type BotContext struct {
	Context context.Context
	Client  *whatsmeow.Client
	Event   *events.Message
	Pool    WorkerPool
}

func (ctx *BotContext) GetMessageString() string {
	msg := ctx.Event.Message
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

func (ctx *BotContext) Reply(msg string) error {
	_, err := ctx.Client.SendMessage(ctx.Context, ctx.Event.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(msg),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(ctx.Event.Info.ID),
				Participant:   proto.String(ctx.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: ctx.Event.Message,
			},
		},
	})

	return err
}

func (ctx *BotContext) GetEventMessageJson() (string, error) {
	bytes, err := sonic.MarshalIndent(ctx.Event, "", "	")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (ctx *BotContext) GetImageMessage() *waE2E.ImageMessage {
	msg := ctx.Event.Message

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

func (ctx *BotContext) ReplyWithSticker(sticker *ImageSticker) error {
	upload, err := ctx.Client.Upload(ctx.Context, sticker.Image, whatsmeow.MediaImage)
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
				StanzaID:      proto.String(ctx.Event.Info.ID),
				Participant:   proto.String(ctx.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: ctx.Event.Message,
			},
		},
	}

	_, err = ctx.Client.SendMessage(ctx.Context, ctx.Event.Info.Chat, message)
	return err
}
