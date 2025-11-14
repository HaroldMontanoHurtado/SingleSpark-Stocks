package main

import (
	"context"
	"log"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/config"
	extern "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
	db "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
	httpapi "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/interfaces/http"
	stockusecase "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stockusecase"
)

func main() {
	// Cargar configuración desde .env (si existe) y entorno.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	// Conectar a la DB (context pasado)
	pg, err := db.NewPGRepo(context.Background(), cfg.PostgresConnString())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pg.Close()

	// Cliente externo y usecase de ingesta
	client := extern.NewClient(cfg.ExternalAPIURL, cfg.ExternalAPIKey)
	ing := stockusecase.NewIngestor(client, pg)

	// Handlers HTTP
	handlers := &httpapi.Handlers{
		Repo:     pg,
		Ingestor: ing,
	}

	// Router y servidor (usa el router que ya tienes en internal/interfaces/http)
	router := httpapi.NewRouter(handlers)

	// Usamos el Serve que ya existe en internal/interfaces/http/server.go para timeouts si quieres:
	if err := httpapi.Serve(":"+cfg.HTTPPort, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}
