package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DBHost string
    DBPort string
    DBUser string
    DBPass string
    DBName string
    DBSSL  string

    ExternalAPIURL string
    ExternalAPIKey string

    HTTPPort string
}

func Load() (*Config, error) {
    _ = godotenv.Load() // carga .env si existe; no falla si no

    c := &Config{
        DBHost:         getenv("DB_HOST", "localhost"),
        DBPort:         getenv("DB_PORT", "26257"),
        DBUser:         getenv("DB_USER", "root"),
        DBPass:         getenv("DB_PASS", ""),
        DBName:         getenv("DB_NAME", "singlespark"),
        DBSSL:          getenv("DB_SSLMODE", "disable"),
        ExternalAPIURL: getenv("EXTERNAL_API_URL", "https://api.karenai.click/swechallenge/list"),
        ExternalAPIKey: getenv("EXTERNAL_API_KEY", ""),
        HTTPPort:       getenv("HTTP_PORT", "8080"),
    }
    return c, nil
}

func getenv(key, fallback string) string {
    v := os.Getenv(key)
    if v == "" {
        return fallback
    }
    return v
}

func (c *Config) PostgresConnString() string {
    // pgx connection string
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBSSL)
}
