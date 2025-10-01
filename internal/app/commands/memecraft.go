package commands

import (
	"TuruBot-Go/internal/adapter/logstream"
	"TuruBot-Go/internal/app/types"
	"TuruBot-Go/internal/port"
	"TuruBot-Go/pkg/stringparser"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (cmd *Command) generateMeme(
	ctx *types.BotContext,
	generator func(ctx context.Context, imageURL string, args map[string]string) (*port.MemeCraftResponse, error),
) error {
	logger, err := logstream.NewWartaLogStream(ctx.Context)
	if err != nil {
		logrus.Errorf("WartaLogStream failed to create: %v", err)
		return err
	}
	defer logger.Close()

	_ = ctx.Pool.Submit(func() {
		_ = ctx.Reply(fmt.Sprintf("Pemrosesan dimulai, log bisa dilihat melalui link berikut:\n\n%s", logger.GetStreamUrl()))
	})
	logger.PushLog("Parsing input args")
	inputArgs := stringparser.ParseArgs(ctx.GetMessageString())
	imageMessage := ctx.GetImageMessage()
	if imageMessage == nil {
		logger.PushLog("Input image is not found")
		return ctx.Reply("gambarnya mana kntl")
	}

	logger.PushLog("downloading the input image...")
	imageByte, err := ctx.Download(ctx.Context, imageMessage)
	if err != nil {
		logger.PushLog("failed to download the input image")
		logrus.Printf("error downloading image: %v", err)
		return ctx.Reply("Error while downloading image!")
	}

	logger.PushLog("uploading the input image...")
	uploadResult, err := cmd.storageAdapter.UploadBytes(ctx.Context, imageByte)
	if err != nil {
		logger.PushLog("failed to upload the input image")
		logrus.Printf("error uploading image: %v", err)
		return ctx.Reply("Error while uploading image!")
	}

	logger.PushLog("creating meme with memecraft.cettalabs.com")
	memeCraftResponse, err := generator(ctx.Context, uploadResult.DirectURL, inputArgs)
	if err != nil {
		logger.PushLog("failed to generate meme")
		logrus.Printf("error generating meme: %v", err)
		return ctx.Reply("Error while generating meme!")
	}

	logger.PushLog("writing meme to bytes")
	finalImageByte, err := cmd.storageAdapter.DownloadToBytes(ctx.Context, memeCraftResponse.ImageUrl)
	if err != nil {
		logger.PushLog("failed to download final meme")
		logrus.Printf("error downloading final image: %v", err)
		return ctx.Reply("Error while downloading final image!")
	}

	logger.PushLog("sending meme result to user")
	logger.PushLog("finished")
	logger.SetNote("request completed successfully")
	return ctx.ReplyWithImage(
		finalImageByte.ByteContent,
		"image/jpeg",
		"_Powered by Cettalabs_\n\nhttps://memecraft.cettalabs.com/",
	)
}

// Khusus untuk Detik
func (cmd *Command) GenerateTimpaDetik(ctx *types.BotContext) error {
	return cmd.generateMeme(ctx, cmd.memeCraftAdapter.GenerateDetikIg)
}

// Khusus untuk CNN
func (cmd *Command) GenerateTimpaCnn(ctx *types.BotContext) error {
	return cmd.generateMeme(ctx, cmd.memeCraftAdapter.GenerateCnnBreakingNews)
}

// Khusus untuk Folkative
func (cmd *Command) GenerateTimpaFolkative(ctx *types.BotContext) error {
	return cmd.generateMeme(ctx, cmd.memeCraftAdapter.GenerateFolkativeIg)
}
