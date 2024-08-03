package main

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/configs"
	"github.com/fenek-dev/go-outline-bot/internal/services"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/internal/telegram"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/handlers"
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

	service := services.New(storage)

	tgHandlers := handlers.New(service)

	bot, err := telegram.InitBot(&cfg.Tg, tgHandlers)

	if err != nil {
		panic(fmt.Sprintf("Can not connect to telegram bot: %e", err))
	}

	bot.Start()

}
