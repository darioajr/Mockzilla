package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	responseMessage      = flag.String("response-message", `{"status":"Received"}`, "Custom response message (optional)")
	responseContentType  = flag.String("content-type", "application/json", "Custom Content-Type for the response (optional)")
	maxBodySize          = int64(64 * 1024) // 64KB
	messageCount         int64
	responseMessageBytes = []byte(*responseMessage)
	responseCode         *int // Pointer to response code
)

//type Message struct {
//	Data map[string]interface{} `json:"data"`
//}

//var messagePool = sync.Pool{
//	New: func() interface{} {
//		return new(Message)
//	},
//}

func handler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Lê o corpo da requisição
	//bodyBytes, err := io.ReadAll(r.Body)
	//if err != nil {
	//	http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	//	log.Printf("Error reading request body: %v", err)
	//		return
	//	}
	//	defer r.Body.Close()

	// Verifica se o corpo da requisição está vazio
	if r.ContentLength == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		log.Printf("Received empty request body")
		return
	}

	// Decodifica o JSON do corpo da requisição
	//if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
	//		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	//		log.Printf("Error decoding JSON: %v", err)
	//		return
	//	}

	// Incrementa o contador de mensagens
	currentCount := atomic.AddInt64(&messageCount, 1)

	// Log da mensagem recebida e contador
	log.Printf("Received message #%d: %v", currentCount, "") //string(bodyBytes))

	// Define o Content-Type e o código de status configurados
	w.Header().Set("Content-Type", *responseContentType)
	w.WriteHeader(*responseCode)

	// Escreve a mensagem de resposta personalizada
	if _, err := w.Write(responseMessageBytes); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	port := flag.String("port", "8080", "Server port")
	certFile := flag.String("cert", "", "Path to the SSL certificate (optional)")
	keyFile := flag.String("key", "", "Path to the SSL key (optional)")
	responseCode = flag.Int("response-code", 200, "HTTP response code to return (optional)")
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
		if *certFile != "" && *keyFile != "" {
			log.Printf("Server started on https://localhost:%s (HTTP/2 enabled)", *port)
			if err := server.ListenAndServeTLS(*certFile, *keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		} else {
			log.Printf("Server started on http://localhost:%s (HTTP/1.1 only)", *port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}
	}()

	<-done
	log.Print("Server stopped - shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Print("Server exited properly")
}
