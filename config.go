package main

type config struct {
	port           string
	debug          bool
}

func loadConfig() *config {
	return &config{
		port:           getEnv("PORT", "8080"),
		debug:          getEnv("DEBUG", "false") == "true",
	}
}
