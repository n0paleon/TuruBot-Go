package main

import (
	"TuruBot-Go/internal/adapter/memecraft"
	"TuruBot-Go/internal/adapter/storage"
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

	// external adapter
	storageAdapter := storage.NewCatboxMoeStorage()
	memeCraftAdapter := memecraft.NewMemeCraft()

	// command initiator
	cmd := commands.Init(r, storageAdapter, memeCraftAdapter)

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
		SetAliases("st", "stiker", "s").
		SetDescription("Generate sticker with image")
	r.Handle(cmd.ShowMenu).
		SetCmd("menu").
		SetDescription("Show bot menu")
	r.Handle(cmd.GenerateTimpaDetik).
		SetCmd("timpadetik").
		SetHelp("/timpadetik --main-content=Prabowo Akan Berikan MBG Untuk Anak Rajin").
		SetDescription("Generate Timpa Detik by memecraft.cettalabs.com")
	r.Handle(cmd.GenerateTimpaCnn).
		SetCmd("timpacnn").
		SetHelp("/timpacnn --headline=Jokowi Putuskan Akan Kembali ke Solo").
		SetDescription("Generate Timpa CNN by memecraft.cettalabs.com")
	r.Handle(cmd.GenerateTimpaFolkative).
		SetCmd("timpafolkative").
		SetHelp("/timpafolkative --headline=Doraemon Ternyata Asli Solo, Kata Pakar --media-name=TukangTimpa").
		SetDescription("Generate Timpa Folkative by memecraft.cettalabs.com")

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
	wa, err := whatsapp.NewClient(pool, routeSet, 30, cfg)
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
