package commands

import (
	"TuruBot-Go/internal/app/types"
	"TuruBot-Go/internal/app/utils"

	"github.com/sirupsen/logrus"
)

func (cmd *Command) GenerateStickerByImage(ctx *types.BotContext) error {
	imageData := ctx.GetImageMessage()

	if imageData == nil {
		return ctx.Reply("kirim/reply gambar dengan caption: /stiker")
	}

	imageBytes, err := ctx.Download(ctx.Context, imageData)
	if err != nil {
		logrus.Errorf("failed to download image: %v", err)
		return ctx.Reply("error nih anjing")
	}

	stickerCh := make(chan []byte)
	stickerErrCh := make(chan error)
	thumbCh := make(chan []byte)
	thumbErrCh := make(chan error)

	_ = ctx.Pool.Submit(func() {
		stickerImg, err := utils.ImageToStickerViaFFMPEG(imageBytes)
		if err != nil {
			stickerErrCh <- err
			return
		}
		stickerCh <- stickerImg
	})

	_ = ctx.Pool.Submit(func() {
		thumbnail, err := utils.GenerateThumbnail(imageBytes, 96, true)
		if err != nil {
			thumbErrCh <- err
			return
		}
		thumbCh <- thumbnail
	})

	var stickerImg, thumbnail []byte

	select {
	case s := <-stickerCh:
		stickerImg = s
	case e := <-stickerErrCh:
		logrus.Errorf("failed to convert image to sticker: %v", e)
		return ctx.Reply("error nih anjing, gambar lu bikin error kocak, anjing")
	}

	select {
	case t := <-thumbCh:
		thumbnail = t
	case e := <-thumbErrCh:
		logrus.Errorf("failed to generate thumbnail: %v", e)
	}

	if err := ctx.ReplyWithSticker(&types.ImageSticker{
		Image:        stickerImg,
		PNGThumbnail: thumbnail,
	}); err != nil {
		logrus.Errorf("failed to reply with sticker: %v", err)
		return ctx.Reply("error nih anjing")
	}

	return nil
}
