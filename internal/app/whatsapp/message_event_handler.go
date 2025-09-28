package whatsapp

import (
	"TuruBot-Go/internal/app/types"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"time"
)

func (wa *WAClient) MessageEventHandler(msg *events.Message) {
	if msg.IsEdit || msg.Info.IsFromMe {
		return
	}

	if time.Now().Sub(msg.Info.Timestamp) > 10*time.Second {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	botCtx := &types.BotContext{
		Pool:    wa.WorkerPool,
		Context: ctx,
		WAC:     wa,
		Event:   msg,
	}

	msgString := botCtx.GetMessageString()

	parts := strings.Fields(msgString)
	if len(parts) == 0 {
		return
	}

	firstWord := parts[0]
	if len(firstWord) < 2 {
		return
	}

	cmdName := strings.ToLower(firstWord[1:])
	start := time.Now()

	err := wa.Router.Exec(cmdName, botCtx)
	if err != nil {
		if err.Error() == "COMMAND_NOT_FOUND" {
			return
		}
		logrus.Errorf("failed to execute command: %v", err)
	}

	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{
		"route_key": cmdName,
		"duration":  fmt.Sprintf("%dms", duration.Milliseconds()),
		"error":     err != nil,
	}).Info("Execution Result")
}
