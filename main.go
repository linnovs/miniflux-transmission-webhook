package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	httpLogger := log.NewWithOptions(os.Stderr, log.Options{Prefix: "http.Server"})
	httpStdlog := httpLogger.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel})

	cfg := loadConfig()
	mux := http.NewServeMux()
	middlewares := []func(http.Handler) http.Handler{
		loggingMiddleware,
		newMinifluxValidateSignatureMiddleware(cfg.minifluxSecret),
	}

	mux.Handle("POST /", chainMiddlewares(nil, middlewares...))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      mux,
		ErrorLog:     httpStdlog,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("starting server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
