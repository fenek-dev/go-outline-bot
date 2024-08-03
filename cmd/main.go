package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/fenek-dev/go-outline-bot/configs"
	"github.com/fenek-dev/go-outline-bot/internal/services"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/internal/telegram"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/handlers"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
)

func main() {
	ctx := context.Background()
	cfg := configs.MustLoad()

	storage := pg.New(
		ctx,
		cfg.DbUrl,
		pg.WithMaxConnections(100),
		pg.WithMinConnections(10),
	)
	log.Println("Db connected")

	paymentClient := payment_service.NewClient(
		"",
		payment_service.WithLogger(slog.Default()),
	)

	service := services.New(storage, paymentClient)
	tgHandlers := handlers.New(service)

	bot, err := telegram.InitBot(&cfg.Tg, tgHandlers)
	log.Println("Telegram bot api inited")

	if err != nil {
		panic(fmt.Sprintf("Can not connect to telegram bot: %e", err))
	}

	bot.Start()

}
