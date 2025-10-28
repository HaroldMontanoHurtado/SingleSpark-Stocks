import (
    "database/sql"
    "net/http"
    "os"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib" // driver pgx compatible
    "github.com/your/module/internal/infrastructure/ingestor"
    httpint "github.com/your/module/internal/interfaces/http"
)

func main() {
    // lee variables de entorno
    dbURL := os.Getenv("DATABASE_URL") // ejemplo: postgresql://root@cockroach:26257/defaultdb?sslmode=disable
    apiURL := os.Getenv("EXTERNAL_API_URL") // https://api.karenai.click/swechallenge/list
    apiToken := os.Getenv("EXTERNAL_API_TOKEN")

    db, err := sql.Open("pgx", dbURL)
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    defer db.Close()

    ing := ingestor.New(apiURL, apiToken)
    ingestHandler := httpint.NewIngestHandler(ing, db)

    mux := http.NewServeMux()
    mux.Handle("/internal/fetch", ingestHandler)
    // ejemplo: endpoint público para listar stocks (implementar en tu repo)
    // mux.HandleFunc("/api/stocks", handlerListStocks)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Printf("listening on :%s", port)
    http.ListenAndServe(":"+port, mux)
}
