package types

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"golang.org/x/time/rate"
	"sync"
)

type MessageSendTask struct {
	Ctx      *BotContext
	Message  *waE2E.Message
	DoneChan chan error
}

type MessageQueue struct {
	Client     *whatsmeow.Client
	sendQueue  chan MessageSendTask
	limiter    *rate.Limiter
	once       sync.Once
	workerStop chan struct{}
}

func NewMessageQueue(client *whatsmeow.Client, maxMessagesPerSecond int) *MessageQueue {
	mq := &MessageQueue{
		Client:     client,
		sendQueue:  make(chan MessageSendTask, 1000),
		limiter:    rate.NewLimiter(rate.Limit(maxMessagesPerSecond), 1),
		workerStop: make(chan struct{}),
	}
	go mq.worker()
	return mq
}

func (mq *MessageQueue) worker() {
	for {
		select {
		case task := <-mq.sendQueue:
			err := mq.limiter.Wait(task.Ctx.Context)
			if err != nil {
				task.DoneChan <- err
				continue
			}

			_, err = mq.Client.SendMessage(task.Ctx.Context, task.Ctx.Event.Info.Chat, task.Message)
			task.DoneChan <- err

		case <-mq.workerStop:
			return
		}
	}
}

func (mq *MessageQueue) EnqueueMessage(ctx *BotContext, message *waE2E.Message) error {
	done := make(chan error, 1)
	mq.sendQueue <- MessageSendTask{Ctx: ctx, Message: message, DoneChan: done}
	return <-done
}

func (mq *MessageQueue) Stop() {
	mq.once.Do(func() {
		close(mq.workerStop)
	})
}
