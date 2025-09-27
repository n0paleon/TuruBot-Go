package app

import (
	"TuruBot-Go/internal/app/router"
	"TuruBot-Go/internal/app/types"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"time"
)

var (
	allowedPrefix = []string{"!", "/", ".", "#", "?"}
)

func HandleMessage(client *whatsmeow.Client, evt *events.Message, r *router.Router, pool types.WorkerPool, mq *types.MessageQueue) {
	if evt.IsEdit || evt.Info.IsFromMe {
		return
	}

	if time.Now().Sub(evt.Info.Timestamp) > 10*time.Second {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	botCtx := &types.BotContext{
		Pool:    pool,
		Context: ctx,
		Client:  client,
		Event:   evt,
		Queue:   mq,
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

	err := r.Exec(cmdName, botCtx)
	if err != nil {
		logrus.Errorf("failed to execute command: %v", err)
	}

	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{
		"route_key": cmdName,
		"duration":  fmt.Sprintf("%dms", duration.Milliseconds()),
		"error":     err != nil,
	}).Info("Execution Result")
}

func HasAllowedPrefix(msg string, prefixes []string) (string, bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(msg, p) {
			return p, true
		}
	}
	return "", false
}
