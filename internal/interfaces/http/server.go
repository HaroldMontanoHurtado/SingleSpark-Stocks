package httpapi

import (
    //"context"
    "log"
    "net/http"
    "time"
)

func Serve(addr string, handler http.Handler) error {
    srv := &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Printf("starting server on %s", addr)
    return srv.ListenAndServe()
}
