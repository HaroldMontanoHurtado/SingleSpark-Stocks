# ---------- Stage: builder ----------
FROM golang:1.21-alpine AS builder
# herramientas
RUN apk add --no-cache git ca-certificates

WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Copiamos módulos y descargamos dependencias (usa la raiz del repo)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el repo (necesario porque la main/entry está en cmd/)
COPY . .

# Build del binario apuntando a tu main. Ajusta el path si tu entry difiere.
RUN go build -o /app/server ./cmd/server

# ---------- Stage: dev (para hot-reload) ----------
FROM golang:1.21-alpine AS dev
RUN apk add --no-cache bash git curl
# instalar air (hot reload) en GOPATH/bin
RUN go install github.com/cosmtrek/air@latest

WORKDIR /src
ENV CGO_ENABLED=0

# Copiamos solo lo necesario para que cuando se use target=dev haga mount de volúmenes y funcione
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ENTRYPOINT por defecto es air (override en compose puede también usar air)
CMD ["air", "-c", ".air.toml"]

# ---------- Stage: runtime (producción) ----------
FROM gcr.io/distroless/static:nonroot AS runtime
WORKDIR /app
# Copiamos el binario compilado en builder
COPY --from=builder /app/server /usr/local/bin/server
# Copiamos archivos de configuración si existen (opcional)
COPY config ./config

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/server"]
