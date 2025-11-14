# /cmd/server/Dockerfile
# Dockerfile para el backend (multi-stage)
# --------------------------------------------------------

# Etapa de compilación (builder)
FROM golang:1.21-alpine AS builder

WORKDIR /src
# herramientas necesarias
RUN apk add --no-cache git build-base

# Copiamos go.mod/go.sum para cached deps
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del proyecto
COPY . .

# Compilamos el binario del backend (paquete en cmd/server)
# Salida directa al path final esperado por docker-compose.override (/usr/local/bin/backend)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /usr/local/bin/backend ./cmd/server

# Etapa runtime
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl

# Crear directorio de trabajo
WORKDIR /app

# Copiamos binario ya compilado
COPY --from=builder /usr/local/bin/backend /usr/local/bin/backend

# Copiamos migraciones (solo si existen)
COPY migrations /app/migrations

# Puerto por defecto en contenedor — puede ser sobrescrito por docker-compose/.env
ENV HTTP_PORT=8081
EXPOSE 8081

# Usuario no-root opcional (si quieres seguridad extra)
# RUN addgroup -S app && adduser -S -G app app
# USER app

ENTRYPOINT ["/usr/local/bin/backend"]
