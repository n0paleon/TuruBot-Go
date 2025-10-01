package whatsapp

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"time"
)

type sendTask struct {
	ctx     context.Context
	chatJID types.JID
	message *waE2E.Message
	done    chan error
}

func (wa *WAClient) worker() {
	for task := range wa.sendQueue {
		err := wa.limiter.Wait(task.ctx)
		if err != nil {
			task.done <- err
			continue
		}

		_, err = wa.Client.SendMessage(task.ctx, task.chatJID, task.message)
		task.done <- err
	}
}

// EnqueueMessage add new message task to sendQueue channel (blocking, 1 message processed by 1 goroutine)
func (wa *WAClient) EnqueueMessage(ctx context.Context, chatJID types.JID, msg *waE2E.Message) error {
	done := make(chan error, 1)
	wa.sendQueue <- sendTask{
		ctx:     ctx,
		chatJID: chatJID,
		message: msg,
		done:    done,
	}
	return <-done
}

// EnqueueMessageNonBlocking (non-blocking version oof EnqueueMessage)
func (wa *WAClient) EnqueueMessageNonBlocking(_ context.Context, chatJID types.JID, msg *waE2E.Message) error {
	return wa.WorkerPool.Submit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := wa.limiter.Wait(ctx); err != nil {
			logrus.WithFields(logrus.Fields{
				"SendTo":   chatJID,
				"ErrorMsg": err.Error(),
			}).Warn("Rate limiter blocked send")
			return
		}

		_, err := wa.Client.SendMessage(ctx, chatJID, msg)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"SendTo":   chatJID,
				"ErrorMsg": err.Error(),
			}).Error("EnqueueMessageNonBlocking error")
		}
	})
}
