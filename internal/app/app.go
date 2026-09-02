package app

import (
	"context"
	"layouts-orders-bot/internal/calculator"
	"layouts-orders-bot/internal/clients/client"
	"layouts-orders-bot/internal/config"
	event_consumer "layouts-orders-bot/internal/consumer"
	"layouts-orders-bot/internal/handlers/answers"
	"layouts-orders-bot/internal/handlers/commands"
	processor "layouts-orders-bot/internal/processor/telegram"
	sessionsrepo "layouts-orders-bot/internal/repo/inmemory/sessions"
	ordersrepo "layouts-orders-bot/internal/repo/sqlite/orders"
	"log"
)

func Run(cfg *config.Config) error {
	ordersStorage, err := ordersrepo.New(cfg.SQLite.Path)
	if err != nil {
		return err
	}
	if err := ordersStorage.Init(context.Background()); err != nil {
		log.Fatal("can't init sqlite: ", err)
	}

	sessionsStorage := sessionsrepo.New()

	defer ordersStorage.Close()

	costCalculator := calculator.New()

	commandsHandler := commands.NewHandler(ordersStorage, costCalculator)
	answersHandler := answers.NewHandler(sessionsStorage, costCalculator)

	tgProcessor := processor.New(client.New(cfg.TG.Host, cfg.TG.TgBotToken), commandsHandler, answersHandler)

	log.Print("service started")

	consumer := event_consumer.New(tgProcessor, cfg.BatchSize)
	if err := consumer.Start(); err != nil {
		return err
	}
	return nil
}
