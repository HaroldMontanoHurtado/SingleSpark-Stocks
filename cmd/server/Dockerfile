# Etapa de compilación
FROM golang:1.21-alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0

# Copiar solo los archivos de dependencias desde la raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente del servidor
COPY cmd/server ./cmd/server

# Compilar el binario
RUN go build -ldflags="-s -w" -o /app/server ./cmd/server

# Etapa final
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
