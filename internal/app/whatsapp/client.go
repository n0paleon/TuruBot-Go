package whatsapp

import (
	"TuruBot-Go/internal/app/router"
	"TuruBot-Go/internal/app/types"
	"context"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"golang.org/x/time/rate"
	"os"
)

type WAClient struct {
	Client     *whatsmeow.Client
	WorkerPool types.WorkerPool
	Router     *router.Router

	sendQueue chan sendTask
	limiter   *rate.Limiter
}

func NewClient(workerpool types.WorkerPool, r *router.Router, maxPerSecond int) (*WAClient, error) {
	ctx := context.Background()
	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:session_store.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	wa := &WAClient{
		Client:     client,
		WorkerPool: workerpool,
		Router:     r,
		sendQueue:  make(chan sendTask, 1000),
		limiter:    rate.NewLimiter(rate.Limit(maxPerSecond), 5), // burst = maxPerSecond
	}
	client.AddEventHandler(wa.EventHandler)

	_ = wa.WorkerPool.Submit(func() {
		wa.worker()
	})

	return wa, nil
}

func (wa *WAClient) Connect() error {

	if wa.Client.Store.ID == nil {
		qrChan, _ := wa.Client.GetQRChannel(context.Background())
		_ = wa.WorkerPool.Submit(func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.H, os.Stdout)
					return
				} else {
					fmt.Println("QR event:", evt.Event)
				}
			}
			return
		})
	}
	return wa.Client.Connect()
}

func (wa *WAClient) Disconnect() {
	wa.Client.Disconnect()
}
