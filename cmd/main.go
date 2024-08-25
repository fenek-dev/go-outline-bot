package main

import (
	"context"
	"fmt"
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
	"github.com/fenek-dev/go-outline-bot/internal/worker"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
)

func main() {
	ctx := context.Background()
	cfg := configs.MustLoad()

	logger := slog.Default()

	// Graceful shutdown
	stopSignal := make(chan os.Signal)

	storage := pg.New(
		ctx,
		cfg.DB,
		pg.WithMaxConnections(100),
		pg.WithMinConnections(10),
	)
	log.Println("Db connected")

	paymentClient := payment_service.NewClient(
		"",
		payment_service.WithLogger(logger),
	)

	service := services.New(
		storage,
		paymentClient,
		cfg,
		services.WithLogger(logger),
	)

	httpServer := server.New(
		cfg.Port,
		service,
		stopSignal,
		server.WithLogger(logger),
	)

	httpServer.Handle(http.MethodGet, "/health", httpServer.HealthHandler)
	httpServer.Handle(http.MethodPost, "/payment/webhook", httpServer.PaymentWebhookHandler)
	go httpServer.Run()

	worker := worker.New(service, stopSignal)
	go worker.Run()

	tgHandlers := handlers.New(service)

	bot, err := telegram.InitBot(&cfg.Tg, tgHandlers)
	log.Println("Telegram bot api inited")

	if err != nil {
		panic(fmt.Errorf("can not connect to telegram bot: %w", err))
	}

	bot.Start()

}
