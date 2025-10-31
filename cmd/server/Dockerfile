# Ubicación del archivo: /cmd/server/Dockefile
# Dockerfile para el backend
# ======================================================
# Etapa de compilación (builder)
# ======================================================
FROM golang:1.21-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

# Copiamos dependencias primero para aprovechar la caché
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el proyecto (el go.mod está en raíz)
COPY . .

# Compilamos el backend (ubicado en cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/backend ./cmd/server

# ======================================================
# Etapa de ejecución
# ======================================================
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl
WORKDIR /app

# Copiamos el binario
COPY --from=builder /app/bin/backend /usr/local/bin/backend

# Copiamos las migraciones si existen
COPY migrations /app/migrations

EXPOSE 8080
ENV HTTP_PORT=8080

CMD ["/usr/local/bin/backend"]
