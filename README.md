# Movie API Server - Concepto de API General

Una API REST moderna construida con Go, Fiber y SQLite para gestionar una base de datos de películas con géneros, directores, actores y reviews.

## 🎬 Características

- **Framework**: Fiber (Go)
- **Base de Datos**: SQLite con GORM
- **Validación**: go-playground/validator
- **Entidades**: Películas, Géneros, Directores, Actores, Usuarios, Reviews
- **Relaciones**: Many-to-many entre películas y actores
- **Datos de Prueba**: Incluye películas populares con datos completos

## 📋 Endpoints Disponibles

### Health Check
- `GET /api/v1/health` - Estado de la API

### Películas
- `GET /api/v1/movies` - Listar películas (con filtros y paginación)
- `GET /api/v1/movies/:id` - Obtener película por ID
- `POST /api/v1/movies` - Crear nueva película
- `PUT /api/v1/movies/:id` - Actualizar película
- `DELETE /api/v1/movies/:id` - Eliminar película
- `GET /api/v1/movies/search?title=inception` - Buscar películas por título
- `GET /api/v1/movies/top-rated?limit=10` - Películas mejor valoradas

### Géneros
- `GET /api/v1/genres/:id/movies` - Películas por género

### Directores
- `GET /api/v1/directors/:id/movies` - Películas por director

### Actores
- `GET /api/v1/actors/:id/movies` - Películas por actor

### Parámetros de Consulta
- `page` - Número de página (default: 1)
- `limit` - Elementos por página (default: 10, max: 100)
- `genre_id` - Filtrar por género
- `director_id` - Filtrar por director
- `min_rating` - Filtrar por rating mínimo
- `title` - Buscar por título (para endpoint search)

## 🛠️ Instalación y Uso

### Prerrequisitos
- Go 1.21 o superior

### Instalación
```bash
# Clonar el repositorio
git clone <repository-url>
cd api-server

# Instalar dependencias
go mod tidy

# Ejecutar la aplicación
go run main.go
```

La API estará disponible en `http://localhost:4444`

### Datos de Prueba
La aplicación incluye datos de ejemplo:
- **6 géneros**: Acción, Comedia, Drama, Terror, Ciencia Ficción, Romance
- **3 directores**: Christopher Nolan, Quentin Tarantino, Greta Gerwig
- **4 actores**: Leonardo DiCaprio, Margot Robbie, Tom Hardy, Emma Stone
- **4 películas**: Inception, Barbie, Pulp Fiction, Poor Things
- **3 usuarios** con reviews

## 📝 Ejemplos de Uso

### Listar todas las películas
```bash
curl http://localhost:4444/api/v1/movies
```

### Buscar películas por título
```bash
curl "http://localhost:4444/api/v1/movies/search?title=inception"
```

### Filtrar películas por género
```bash
curl "http://localhost:4444/api/v1/movies?genre_id=1&min_rating=8.0"
```

### Obtener películas mejor valoradas
```bash
curl "http://localhost:4444/api/v1/movies/top-rated?limit=5"
```

### Crear una nueva película
```bash
curl -X POST http://localhost:4444/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "description": "Un programador descubre que la realidad es una simulación",
    "release_year": 1999,
    "duration": 136,
    "rating": 8.7,
    "poster_url": "https://example.com/matrix.jpg",
    "trailer_url": "https://example.com/matrix-trailer.mp4",
    "genre_id": 5,
    "director_id": 1,
    "actor_ids": [1, 3]
  }'
```

### Obtener película por ID
```bash
curl http://localhost:4444/api/v1/movies/1
```

### Películas por director
```bash
curl http://localhost:4444/api/v1/directors/1/movies
```

## 🗄️ Estructura de la Base de Datos

### Películas
- `id` - ID único
- `title` - Título de la película
- `description` - Descripción
- `release_year` - Año de lanzamiento
- `duration` - Duración en minutos
- `rating` - Calificación (0-10)
- `poster_url` - URL del poster
- `trailer_url` - URL del trailer
- `genre_id` - ID del género
- `director_id` - ID del director
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

### Géneros
- `id` - ID único
- `name` - Nombre del género
- `description` - Descripción
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

### Directores
- `id` - ID único
- `name` - Nombre del director
- `biography` - Biografía
- `birth_date` - Fecha de nacimiento
- `nationality` - Nacionalidad
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

### Actores
- `id` - ID único
- `name` - Nombre del actor
- `biography` - Biografía
- `birth_date` - Fecha de nacimiento
- `nationality` - Nacionalidad
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

### Reviews
- `id` - ID único
- `movie_id` - ID de la película
- `user_id` - ID del usuario
- `rating` - Calificación (1-10)
- `comment` - Comentario
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

### Usuarios
- `id` - ID único
- `username` - Nombre de usuario
- `email` - Email
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

## 🔧 Configuración

La aplicación usa SQLite por defecto. El archivo de base de datos se crea automáticamente como `api_server.db` en el directorio raíz.

## 🧪 Próximas Mejoras

- [ ] Autenticación JWT
- [ ] Sistema de watchlist
- [ ] Recomendaciones basadas en preferencias
- [ ] Documentación Swagger
- [ ] Tests unitarios
- [ ] Cache con Redis
- [ ] Búsqueda avanzada (por actor, director)
- [ ] Subida de archivos (posters, trailers)
- [ ] Logs estructurados
- [ ] Métricas y monitoreo
- [ ] API para géneros, directores y actores (CRUD completo)

## 📄 Licencia

Este proyecto es para fines educativos y de demostración.


