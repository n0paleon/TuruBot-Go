package main

import (
	"TuruBot-Go/internal/app/commands"
	"TuruBot-Go/internal/app/router"
	"TuruBot-Go/internal/app/router/middleware"
	"TuruBot-Go/internal/app/whatsapp"
	"TuruBot-Go/internal/config"
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
	cmd := commands.Init(r)

	// route middleware
	r.Use(middleware.MessageMiddleware())
	r.Use(middleware.AllowedPrefixMiddleware())

	r.Handle(cmd.PingHandler).
		SetCmd("ping").
		SetDescription("Bot availability status check")
	r.Handle(cmd.StatusHandler).
		SetCmd("status").
		SetAliases("stat", "stats").
		SetDescription("Bot system status check")
	r.Handle(cmd.GenerateStickerByImage).
		SetCmd("sticker").
		SetAliases("st", "stiker").
		SetDescription("Generate sticker with image")
	r.Handle(cmd.ShowMenu).
		SetCmd("menu").
		SetDescription("Show bot menu")

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

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	routeSet := botRoutes()
	wa, err := whatsapp.NewClient(pool, routeSet, 20, cfg)
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
