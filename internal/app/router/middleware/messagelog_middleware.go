package middleware

import (
	"TuruBot-Go/internal/app/types"
	"github.com/sirupsen/logrus"
)

func MessageMiddleware() types.Middleware {
	return func(next types.CommandHandler) types.CommandHandler {
		return func(ctx *types.BotContext) error {
			logrus.WithFields(logrus.Fields{
				"Message": ctx.GetMessageString(),
				"From":    ctx.Event.Info.PushName,
				"Chat":    ctx.Event.Info.Chat,
			}).Info("New Message")
			return next(ctx)
		}
	}
}
