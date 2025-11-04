package httpapi

import (
	"log"
	"net/http"
	"time"
)

// Serve lanza el servidor HTTP con configuración de timeout.
func Serve(addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", addr)
	return srv.ListenAndServe()
}
