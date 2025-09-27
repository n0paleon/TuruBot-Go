package commands

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"time"
)

func PingHandler(ctx *types.BotContext) error {
	latency := time.Since(ctx.Event.Info.Timestamp).Milliseconds()
	return ctx.Reply(fmt.Sprintf("Pong!\n\n_%dms_", latency))
}
