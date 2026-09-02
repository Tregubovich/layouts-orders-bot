package main

import (
	"layouts-orders-bot/internal/app"
	"layouts-orders-bot/internal/config"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Fatal("service is stopped: ", err)
	}
}
