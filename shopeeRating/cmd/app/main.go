package main

import (
	"log"

	"github.com/TranQuocToan1996/shopeerating/config"
	"github.com/TranQuocToan1996/shopeerating/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
