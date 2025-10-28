# ---------- Stage: builder ----------
FROM node:18-alpine AS builder
WORKDIR /app

# Instalar dependencias
COPY ui/package*.json ./
# si tu package.json está directamente en ./ui
RUN npm ci

# Copiar el resto del front
COPY ui/ ./

# Construir producción (salida en /app/dist)
RUN npm run build

# ---------- Stage: dev (para hot reload) ----------
FROM node:18-alpine AS dev
WORKDIR /app
RUN npm install -g pnpm || true
# No copiamos ficheros: en dev usaremos volumes (override) para montar el código
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# ---------- Stage: runtime (nginx) ----------
FROM nginx:stable-alpine AS runtime
COPY --from=builder /app/dist /usr/share/nginx/html
# Si usas history mode en Vue Router copia configuración nginx
# COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
