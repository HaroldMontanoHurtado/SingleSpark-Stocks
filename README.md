## SingleSpark-Stocks

SingleSpark Stocks es un proyecto desarrollado en Golang (backend) y Vue 3 + TypeScript (frontend) que tiene como objetivo recuperar, almacenar, analizar y visualizar información bursátil en tiempo real.

El sistema se conecta a una API externa de datos financieros, procesa y guarda la información en CockroachDB, y ofrece al usuario una interfaz moderna, intuitiva y responsiva para explorar las acciones más relevantes del mercado.
Además, implementa un módulo de análisis y recomendación que sugiere las mejores acciones para invertir, basándose en criterios cuantitativos obtenidos desde la API y posibles fuentes complementarias.

### Objetivos del proyecto

1. Conectarse a la API bursátil externa para obtener datos actualizados sobre diferentes acciones.

2. Almacenar los datos en una base de datos distribuida (CockroachDB).

3. Proporcionar una API interna en Go que exponga los datos de forma estructurada.

4. Crear una interfaz de usuario con Vue 3 que permita buscar, filtrar y visualizar las acciones.

5. Implementar un algoritmo de recomendación que identifique las mejores oportunidades de inversión diarias.

6. Garantizar estabilidad y mantenibilidad, aplicando principios de Clean Architecture y patrones de diseño adecuados.

### Arquitectura del sistema

El proyecto sigue el enfoque de Arquitectura Limpia (Clean Architecture) o Hexagonal, con el fin de mantener una clara separación entre las capas de dominio, aplicación, infraestructura y presentación.

~~~
Frontend (Vue 3 + TypeScript + Pinia + Tailwind)
        ↑
[REST API interna en Go]
        ↑
Use Cases / Services (Lógica de negocio)
        ↑
Repositories (Acceso a datos)
        ↑
CockroachDB (Persistencia)
        ↑
API Externa (Datos bursátiles)
~~~

**Componentes principales:**

* Backend:
  + Lenguaje: Go (Golang)
  + Framework sugerido: Gin o Fiber
  + Estructura basada en Clean Architecture
  + Módulos:
      - Conector de API externa
      - Repositorios (CockroachDB)
      - Casos de uso y servicios
      - Servidor HTTP con endpoints REST
      - Lógica de recomendación
* Frontend:
  + Framework: Vue 3 + TypeScript
  + Estado global: Pinia
  + Estilos: Tailwind CSS
  + Funciones:
      - Visualización y filtrado de acciones
      - Sección de recomendaciones
      - Componentes reutilizables (tablas, tarjetas, buscadores)

### Patrones y principios aplicados

- **Clean Architecture:** separación de responsabilidades y capas desacopladas.
- **Repository Pattern:** abstracción del acceso a la base de datos.
- **Use Case Pattern:** encapsulación de la lógica de negocio.
- **Dependency Injection:** facilita pruebas y reduce el acoplamiento.
- **Adapter Pattern:** permite cambiar la fuente de datos sin modificar el núcleo.
- **Observer Pattern (Frontend):** reactivo con Pinia para actualización en tiempo real.

### Tecnologías utilizadas
| Área                  | Tecnología                 | Descripción                                        |
|-----------------------|----------------------------|----------------------------------------------------|
| Backend               | Go (Golang)                | API interna, lógica de negocio y conexión a datos. |
| Framework Backend     | Gin / Fiber                | Manejo de rutas y controladores HTTP.              |
| Base de datos         | CockroachDB                | Base SQL distribuida, compatible con PostgreSQL.   |
| ORM / Query Builder   | GORM / sqlx                | Abstracción del acceso a datos.                    |
| Frontend              | Vue 3 + TypeScript + Pinia | Interfaz reactiva y tipada.                        |
| Estilos               | Tailwind CSS               | Diseño limpio y responsive.                        |
| Infraestructura local | Docker Compose             | Orquestación de servicios en entorno local.        |
| Pruebas               | Go Test / Vitest           | Validación de la lógica del sistema.               |

### Flujo general de ejecución

1. El backend obtiene los datos bursátiles desde la API externa.
2. Los datos se transforman y se guardan en CockroachDB.
3. Se exponen endpoints REST para que el frontend los consuma.
4. El usuario accede a la interfaz web donde puede visualizar y filtrar las acciones.
5. El sistema ejecuta un algoritmo de recomendación que sugiere inversiones.

### Ejecución local

Requisitos previos

* Docker y Docker Compose instalados
* Go 1.22+
* Node.js 18+

### Iniciar Entorno

~~~
# Levantar servicios
docker compose up -d

# Ejecutar backend
cd cmd/api
go run main.go

# Ejecutar frontend
cd ui
npm run dev
~~~

### Estado del proyecto

Este proyecto se encuentra en desarrollo activo y está optimizado para funcionar en un único equipo local, sin dependencias externas innecesarias.
El código sigue principios de mantenibilidad y escalabilidad, facilitando la incorporación de nuevas fuentes de datos o módulos de análisis en el futuro.

### **Autor**

**Desarrollado por:** *Harold Montano*

**Rol:** *Developer*
