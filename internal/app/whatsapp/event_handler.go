package whatsapp

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/types/events"
)

func (wa *WAClient) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msg := v
		_ = wa.WorkerPool.Submit(func() {
			wa.MessageEventHandler(msg)
		})
	case *events.Connected:
		logrus.Println("Client connected successfully")
	case *events.Disconnected:
		logrus.Println("Client disconnected by the server side")
	case *events.ConnectFailure:
		logrus.Printf("Client connect receive ConnectFailure events, message: %s", evt.(*events.ConnectFailure).Message)
	case *events.LoggedOut:
		logrus.Printf("Client logged out, reason: %s", evt.(*events.LoggedOut).Reason.String())
	}
}
