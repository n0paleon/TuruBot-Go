package commands

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"time"
)

func (cmd *Command) PingHandler(ctx *types.BotContext) error {
	start := time.Now()
	ctx.SetQueueMode(types.QueueBlocking)
	err := ctx.Reply("Pong!")
	if err != nil {
		return err
	}

	latency := time.Since(start).Milliseconds()

	return ctx.Reply(fmt.Sprintf("⏱ Bot latency: %dms", latency))
}
