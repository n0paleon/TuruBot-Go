package middleware

import (
	"TuruBot-Go/internal/app/types"
	"strings"
)

var (
	allowedPrefix = []string{"!", "/", ".", "#", "?"}
)

func AllowedPrefixMiddleware() types.Middleware {
	return func(next types.CommandHandler) types.CommandHandler {
		return func(ctx *types.BotContext) error {
			msgString := ctx.GetMessageString()
			if msgString == "" {
				return nil
			}

			if ctx.Event.Info.PushName == "status@broadcast" {
				return nil
			}

			for _, p := range allowedPrefix {
				if strings.HasPrefix(msgString, p) {
					return next(ctx)
				}
			}

			// message not started with known prefix
			return nil
		}
	}
}
