package main

import (
	"log"
	"test/faraway/server/config"
	"test/faraway/server/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Server config error: %s", err)
	}
	app.Run(cfg)
}
