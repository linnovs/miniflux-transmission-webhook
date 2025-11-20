package main

import "github.com/charmbracelet/log"

type config struct {
	port           string
	debug          bool
}

func loadConfig() *config {
	cfg := &config{
		port:           getEnv("PORT", "8080"),
		debug:          getEnv("DEBUG", "false") == "true",
	}

	if cfg.debug {
		log.SetLevel(log.DebugLevel)
	}

	if cfg.minifluxSecret == "" {
		log.Fatal("MINIFLUX_SECRET environment variable is required")
	}

	return cfg
}
