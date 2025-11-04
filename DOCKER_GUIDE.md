## 🐳 DOCKER_GUIDE.md
Guía completa para ejecutar SingleSpark Stocks con Docker
### 📘 Introducción

Esta guía explica cómo construir, ejecutar y administrar el entorno de desarrollo del proyecto SingleSpark Stocks utilizando Docker y Docker Compose.
El objetivo es mantener una configuración reproducible, aislada y fácil de desplegar para el backend (Go), el frontend (Vue 3 + TypeScript) y la base de datos (CockroachDB).

### 🧩 Estructura del Proyecto

La raíz del proyecto debe tener la siguiente estructura mínima:

~~~
SingleSpark-Stocks/
├── cmd/
│   └── backend/
│       └── Dockerfile
├── internal/
├── pkg/
├── config/
├── ui/
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── DOCKER_GUIDE.md
└── README.md
~~~

### ⚙️ Archivos principales
📄 `docker-compose.yml`

Orquesta los contenedores del backend, frontend y base de datos.
Cada servicio está conectado mediante la red singlespark-net.

📄 `cmd/backend/Dockerfile`

Compila y ejecuta el backend en Go con Clean Architecture.

📄 `ui/Dockerfile`

Levanta el entorno del frontend (Vue 3 + TypeScript) para desarrollo o despliegue.

### 🚀 Primeros pasos
1️⃣ Instalar Docker Desktop

Descarga e instala Docker Desktop para Windows desde:

🔗 https://www.docker.com/products/docker-desktop

Asegúrate de que los comandos docker y docker compose funcionen desde la terminal (PowerShell o Git Bash).

Prueba:

~~~
docker --version
docker compose version
~~~

### 🏗️ Construir y levantar el proyecto

Desde la raíz del proyecto:

~~~
docker compose up --build
~~~

📌 Este comando:

Construye las imágenes de backend y frontend.

Inicia los contenedores (backend, frontend, base de datos).

Conecta todo en la red singlespark-net.

Cuando termine, verás logs en consola indicando que los tres servicios están corriendo.

### 🌐 Accesos rápidos
| Servicio             | Puerto local          | Descripción                                  |
|----------------------|-----------------------|----------------------------------------------|
| Frontend (Vue 3)     | http://localhost:5173 | Interfaz web del sistema.                    |
| Backend (Go API)     | http://localhost:8080 | API principal de la aplicación.              |
| CockroachDB Admin UI | http://localhost:8081 | Panel de administración de la base de datos. |

### 🧱 Comandos esenciales

Reconstruye la imagen desde cero
~~~
docker compose build --no-cache
~~~

**▶️ Levantar el proyecto**
~~~
docker compose up
~~~

**🛠️ Reconstruir imágenes (por cambios en código o Dockerfile)**
~~~
docker compose up --build
~~~

**⏹️ Detener los contenedores**
~~~
docker compose down
~~~

**🧼 Detener y eliminar todos los contenedores, imágenes y volúmenes**
~~~
docker compose down --volumes --rmi all
~~~

**🔄 Reiniciar un servicio específico**
~~~
docker compose restart backend
~~~

**🧾 Ver logs en tiempo real**
~~~
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f db
~~~

### 🧠 Uso avanzado

**🧩 Entrar a un contenedor**

Por ejemplo, para acceder al backend:
~~~
docker exec -it singlespark-backend sh
~~~

Para CockroachDB (modo shell SQL):
~~~
docker exec -it singlespark-db ./cockroach sql --insecure
~~~

### 💾 Persistencia de datos

La base de datos CockroachDB almacena sus datos en el volumen:
~~~
volumes:
    cockroach-data:
~~~

👉 Esto asegura que la información no se pierda incluso si detienes los contenedores.

Para limpiar completamente los datos:
~~~
docker volume rm singlespark-stocks_cockroach-data
~~~

### 🔧 Variables de entorno (opcional)

Puedes crear un archivo .env en la raíz del proyecto para manejar las variables del backend:
~~~
DB_HOST=db
DB_PORT=26257
DB_USER=root
DB_NAME=singlespark
DB_SSLMODE=disable
~~~

Y modificar docker-compose.yml para incluir:
~~~
env_file:
    - .env
~~~

### 🧩 Flujo de desarrollo recomendado

1. Modifica código localmente (en internal/, ui/, etc.).
Los cambios se reflejan automáticamente dentro del contenedor gracias a los volúmenes montados.

2. Usa Docker Compose para ejecutar todo el sistema completo sin necesidad de instalar Go, Node o CockroachDB localmente.

3. Reinicia servicios individuales cuando modifiques dependencias o archivos de configuración.

### 🔍 Verificación final

Para confirmar que todo está funcionando correctamente:
~~~
docker ps
~~~

Deberías ver algo similar:
~~~
CONTAINER ID   NAME                 STATUS          PORTS
a12bc34de567   singlespark-frontend Up 2 minutes    0.0.0.0:5173->5173/tcp
b45de67fa890   singlespark-backend  Up 2 minutes    0.0.0.0:8080->8080/tcp
c89fa12bc345   singlespark-db       Up 2 minutes    26257/tcp, 0.0.0.0:8081->8080/tcp
~~~

### 🧩 Limpieza total

Para liberar espacio y eliminar cualquier residuo de compilación o base de datos:
~~~
docker compose down --volumes --rmi all
docker system prune -a
~~~

### ✅ Conclusión

Esta configuración permite:
- Ejecutar backend, frontend y base de datos en un entorno controlado.
- Facilitar el desarrollo local sin conflictos de dependencias.
- Preparar el proyecto para un despliegue en producción con solo pequeños ajustes.
