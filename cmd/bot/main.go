package main

import (
	"TuruBot-Go/internal/app/commands"
	"TuruBot-Go/internal/app/router"
	"TuruBot-Go/internal/app/router/middleware"
	"TuruBot-Go/internal/infra/whatsapp"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
}

func botRoutes() *router.Router {
	r := router.New()

	// route middleware
	r.Use(middleware.MessageMiddleware())
	r.Use(middleware.AllowedPrefixMiddleware())

	r.Handle(commands.PingHandler).
		Cmd("ping").
		Description("Bot availability status check")
	r.Handle(commands.StatusHandler).
		Cmd("status").
		Aliases("stat", "stats").
		Description("Bot system status check")
	r.Handle(commands.GenerateStickerByImage).
		Cmd("sticker").
		Aliases("st", "stiker").
		Description("Generate sticker with image")

	// print all routes
	r.PrintRoutes()

	return r
}

func main() {
	pool, err := ants.NewPool(10000,
		ants.WithPreAlloc(true),
		ants.WithNonblocking(true),
	)
	if err != nil {
		logrus.Fatalf("failed to initialize ants.Pool: %v", err)
	}

	routeSet := botRoutes()
	wa, err := whatsapp.NewClient(pool, routeSet)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := wa.Connect(); err != nil {
		logrus.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logrus.Println("✅ Bot started. Press CTRL+C to stop...")

	<-sig
	logrus.Println("🛑 Shutting down...")
	wa.Disconnect()
}
