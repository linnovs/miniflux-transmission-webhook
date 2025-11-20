package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	httpLogger := log.NewWithOptions(os.Stderr, log.Options{Prefix: "http"})
	httpStdlog := httpLogger.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel})

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", getEnv("PORT", "8080")),
		Handler:      mux,
		ErrorLog:     httpStdlog,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("Starting server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
