package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	responseMessage = []byte(`{"status":"Received"}`)
	maxBodySize     = int64(64 * 1024) // 64KB
)

type Message struct {
	Data map[string]interface{} `json:"data"`
}

var messagePool = sync.Pool{
	New: func() interface{} {
		return new(Message)
	},
}

func handler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	msg := messagePool.Get().(*Message)
	defer messagePool.Put(msg)

	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(responseMessage); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	log.Printf("Server started on port %s", *port)

	<-done
	log.Print("Server stopped - shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Print("Server exited properly")
}
