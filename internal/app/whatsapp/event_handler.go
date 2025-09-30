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
	case *events.IdentityChange:
		logrus.Printf("events.IdentityChange is detected")
	case *events.HistorySync:
		logrus.Printf("events.HistorySync is detected, the phone has sent a blob of historical messages")
	case *events.AppStateSyncComplete:
		logrus.WithFields(logrus.Fields{
			"WAPatchName": v.Name,
		}).Info("events.AppStateSyncComplete is detected")
	case *events.OfflineSyncCompleted:
		logrus.WithFields(logrus.Fields{
			"count": v.Count,
		}).Info("events.OfflineSyncComplete is detected")
	case *events.OfflineSyncPreview:
		logrus.WithFields(logrus.Fields{
			"Total":          v.Total,
			"AppDataChanges": v.AppDataChanges,
			"Messages":       v.Messages,
			"Notifications":  v.Notifications,
			"Receipts":       v.Receipts,
		}).Info("events.OfflineSyncPreview is detected, the server is going to send events that the client missed during downtime")
	case *events.UndecryptableMessage:
		logrus.Printf("events.UndecryptableMessage is detected, the client will automatically ask the sender to retry")
	}
}
