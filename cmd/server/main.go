package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"context"
	"time"
)

// handler de prueba
func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("SingleSpark Stocks backend OK\n"))
}

func main() {
	// puerto por defecto
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// arrancar servidor en goroutine para poder hacer shutdown limpio
	go func() {
		log.Printf("Server listening on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// esperar señal para shutdown (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
