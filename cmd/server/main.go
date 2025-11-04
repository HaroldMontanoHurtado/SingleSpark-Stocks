package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/config"
    extern "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
    db "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
    stockusecase "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stockusecase"
    httpapi "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/interfaces/http"
    "github.com/go-chi/chi/v5"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    var port string
    flag.StringVar(&port, "port", cfg.HTTPPort, "http server port")
    flag.Parse()

    // Conexión a la base de datos
    pg, err := db.NewPostgresRepository(cfg.PostgresConnString())
    if err != nil {
        log.Fatalf("connect db: %v", err)
    }

    // Cliente externo y usecase de ingesta
    client := extern.NewClient(cfg.ExternalAPIURL, cfg.ExternalAPIKey)
    ing := stockusecase.NewIngestor(client, pg)

    // Handlers HTTP
    handlers := &httpapi.Handlers{
        Repo:     pg,
        Ingestor: ing,
    }

    // Router principal
    r := chi.NewRouter()
    r.Mount("/", httpapi.NewRouter(handlers))

    // Dirección del servidor
    srvAddr := fmt.Sprintf(":%s", port)

    // Servidor HTTP y apagado limpio
    srv := &http.Server{
        Addr:    srvAddr,
        Handler: r,
    }

    go func() {
        log.Printf("Server running on %s", srvAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %v", err)
        }
    }()

    // Esperar señal de cierre
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited properly")
}
