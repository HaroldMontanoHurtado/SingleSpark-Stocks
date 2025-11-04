package httpapi

import (
	"log"
	"net/http"
	"os"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stock"
)

// SetupHTTPServer inicializa la capa HTTP
func SetupHTTPServer() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL no configurado")
	}

	// Inicializar el repositorio PostgreSQL
	repo, err := db.NewPostgresRepository(connString)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Inicializar el cliente externo (para APIs o scraping)
    apiKey := os.Getenv("EXTERN_API_KEY")
    baseURL := os.Getenv("EXTERN_BASE_URL")

    if apiKey == "" || baseURL == "" {
        log.Fatal("Faltan variables de entorno EXTERN_API_KEY o EXTERN_BASE_URL")
    }

    client := extern.NewClient(apiKey, baseURL)

	// Crear el caso de uso principal de ingesta
	ingestor := stockusecase.NewIngestor(client, repo)

	// Registrar rutas HTTP
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
            return
        }
        ctx := r.Context()
        _, err := ingestor.Ingest(ctx)
        if err != nil {
            http.Error(w, "Error en la ingesta: "+err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write([]byte("Ingesta completada correctamente"))
    })


	log.Println("Servidor HTTP iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
