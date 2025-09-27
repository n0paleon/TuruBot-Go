package commands

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"time"
)

func PingHandler(ctx *types.BotContext) error {
	start := time.Now()
	err := ctx.Reply("Pong!")
	if err != nil {
		return err
	}
	latency := time.Since(start).Milliseconds()
	return ctx.Reply(fmt.Sprintf("⏱ Bot latency: %dms", latency))
}
