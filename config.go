package main

type config struct {
	port           string
}

func loadConfig() *config {
	return &config{
		port:           getEnv("PORT", "8080"),
	}
}
