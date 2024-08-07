package main

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/notifier"
	"github.com/fenek-dev/go-outline-bot/internal/worker"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/fenek-dev/go-outline-bot/configs"
	"github.com/fenek-dev/go-outline-bot/internal/server"
	"github.com/fenek-dev/go-outline-bot/internal/services"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/internal/telegram"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/handlers"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
)

func main() {
	ctx := context.Background()
	cfg := configs.MustLoad()

	// Graceful shutdown
	stopSignal := make(chan os.Signal)

	storage := pg.New(
		ctx,
		cfg.DbUrl,
		pg.WithMaxConnections(100),
		pg.WithMinConnections(10),
	)
	log.Println("Db connected")

	payment := payment_service.NewClient(
		"",
		payment_service.WithLogger(slog.Default()),
	)

	notify := notifier.New(nil)

	service := services.New(storage, payment, notify, cfg)

	httpServer := server.New(
		cfg.Port,
		service,
		stopSignal,
		server.WithLogger(slog.Default()),
	)

	httpServer.Handle(http.MethodGet, "/health", httpServer.HealthHandler)
	httpServer.Handle(http.MethodPost, "/payment/webhook", httpServer.PaymentWebhookHandler)
	go httpServer.Run()

	tgHandlers := handlers.New(service)

	bot, err := telegram.InitBot(&cfg.Tg, tgHandlers)
	log.Println("Telegram bot api inited")

	notify.SetBot(bot)

	if err != nil {
		panic(fmt.Errorf("can not connect to telegram bot: %w", err))
	}

	worker := worker.New(service, stopSignal)
	go worker.Run()

	bot.Start()

}
