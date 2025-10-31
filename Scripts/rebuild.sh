#!/bin/bash
echo "🧹 Limpiando entorno Docker..."
docker compose down -v --remove-orphans
docker image prune -a -f
echo "⚙️ Reconstruyendo contenedores..."
docker compose build --no-cache
docker compose up
