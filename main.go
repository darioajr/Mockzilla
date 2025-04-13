package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	responseMessage      = flag.String("response-message", `{"status":"Received"}`, "Custom response message (optional)")
	responseContentType  = flag.String("content-type", "application/json", "Custom Content-Type for the response (optional)")
	maxBodySize          = int64(64 * 1024) // 64KB
	messageCount         int64
	responseMessageBytes = []byte(*responseMessage)
	responseCode         *int
)

func handler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Error reading request body")
		return
	}
	defer r.Body.Close()

	if r.ContentLength == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		log.Info().Msg("Received empty request body")
		return
	}

	currentCount := atomic.AddInt64(&messageCount, 1)

	log.Info().Msgf("Received message #%d: %v", currentCount, string(bodyBytes))

	w.Header().Set("Content-Type", *responseContentType)
	w.WriteHeader(*responseCode)

	if _, err := w.Write(responseMessageBytes); err != nil {
		log.Info().Err(err).Msg("Error writing response")
	}
}

func main() {
	port := flag.String("port", "8080", "Server port")
	certFile := flag.String("cert", "", "Path to the SSL certificate (optional)")
	keyFile := flag.String("key", "", "Path to the SSL key (optional)")
	responseCode = flag.Int("response-code", 200, "HTTP response code to return (optional)")
	flag.Parse()

	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(zerolog.New(writer).With().Timestamp().Logger())

	mux := http.NewServeMux()
	mux.HandleFunc("/message", handler)
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	server := &http.Server{
		Addr:              ":" + *port,
		Handler:           mux,
		ReadTimeout:       120 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if *certFile != "" && *keyFile != "" {
			log.Info().Msgf("Mockzilla v0.4.0 Server started on https://localhost:%s (HTTP/2 enabled)", *port)
			if err := server.ListenAndServeTLS(*certFile, *keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Server error")
			}
		} else {
			log.Info().Msgf("Mockzilla v0.4.0 Server started on http://localhost:%s (HTTP/1.1 only)", *port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Server error")
			}
		}
	}()

	<-done
	log.Info().Msgf("Server stopped - shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
	log.Info().Msgf("Server exited properly")

}
