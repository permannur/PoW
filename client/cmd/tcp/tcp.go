package main

import (
	"log"
	"test/faraway/client/config"
	"test/faraway/client/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Client config error: %s", err)
	}
	app.Run(cfg)
}
