package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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

func computeHMAC(body []byte, secret string) string {
	mac := hmac.New(crypto.SHA256.New, []byte(secret))
	mac.Write(body)

	return fmt.Sprintf("%x", mac.Sum(nil))
}

func minifluxValidateSignatureMiddleware(secret string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Miniflux-Signature")

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		if err := r.Body.Close(); err != nil { // Close the original body
			http.Error(w, "Failed to close request body", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		hash := computeHMAC(bodyBytes, secret)
		if hash != signature {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
