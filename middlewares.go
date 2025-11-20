package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func rndID() string {
	const byteSize = 8

	rndBytes := make([]byte, byteSize)
	rand.Read(rndBytes) // nolint:errcheck,gosec

	return base64.RawURLEncoding.EncodeToString(rndBytes)
}

func loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := rndID()
		w = newLoggingResponseWriter(w)
		logger := log.With(
			"requestID", id,
			"remoteAddr", r.RemoteAddr,
			"host", r.Host,
			"method", r.Method,
			"uri", r.RequestURI,
			"userAgent", r.UserAgent(),
		)
		r = r.WithContext(context.WithValue(r.Context(), log.ContextKey, logger))

		handler.ServeHTTP(w, r)
		logger.Info(
			"request completed",
			"status", w.(*loggingResponseWriter).statusCode,
			"ressponeTime", time.Since(start).String(),
		)
	})
}
