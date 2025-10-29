package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/config"
    extern "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
    db "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
    stockusecase "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stock"
    httpapi "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/interfaces/http"
    "github.com/go-chi/chi/v5"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    // optional flag to override port
    var port string
    flag.StringVar(&port, "port", cfg.HTTPPort, "http server port")
    flag.Parse()

    // connect to DB
    pg, err := db.NewPostgresRepository(cfg.PostgresConnString())
    if err != nil {
        log.Fatalf("connect db: %v", err)
    }
    // create external client
    client := extern.NewClient(cfg.ExternalAPIURL, cfg.ExternalAPIKey)
    ing := stockusecase.NewIngestor(client, pg)

    handlers := &httpapi.Handlers{
        Repo:     pg,
        Ingestor: ing,
    }

    // router
    r := chi.NewRouter()
    // mount our handlers
    r.Mount("/", httpapi.NewRouter(handlers))

    // graceful shutdown
    srvAddr := fmt.Sprintf(":%s", port)
    srv := &http.Server{
        Addr:    srvAddr,
        Handler: r,
    }

    go func() {
        log.Printf("server listening on %s", srvAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen error: %v", err)
        }
    }()

    // trap signals
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("server shutdown error: %v", err)
    }
    log.Println("server stopped")
}
