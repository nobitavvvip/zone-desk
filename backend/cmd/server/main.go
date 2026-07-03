package main

import (
	"log"
	"os"

	"zonedesk/internal/app"
	"zonedesk/internal/config"
)

func main() {
	configPath := "./config/config.yaml"
	if v := os.Getenv("ZONEDESK_CONFIG"); v != "" {
		configPath = v
	}
	for i, a := range os.Args {
		if a == "--config" && i+1 < len(os.Args) {
			configPath = os.Args[i+1]
			break
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	r := app.New(cfg)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("ZoneDesk started at %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
