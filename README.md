# URL Shortener

Servicio simple de acortador de URLs escrito en Go. Proporciona endpoints para crear, obtener, actualizar, borrar y consultar estadísticas de un short code.

**Estado:** prueba / desarrollo

**Tecnologías:** Go, Echo (web framework), SQLite (opcional), GitHub Actions (CI)

**Estructura básica**
- `main.go` : punto de entrada, configura rutas y dependencias.
- `handlers.go` : handlers HTTP (POST/GET/PUT/DELETE/Stats).
- `service.go` : lógica de negocio (servicio que usa `Store`).
- `SQLStore.go`, `inMemoryStore.go` : implementaciones de `Store`.
- `types.go` : DTOs y la interfaz `Store`.
- `util.go` : helpers (generación de short code, transformaciones).
- `migrations/` : migraciones SQL para SQLite.
- `service_test.go`, `handlers_test.go` : tests unitarios e integración de handlers.
- `.github/workflows/ci.yml` : CI (gofmt, go vet, go test).

**Requisitos**
- Go 1.20+ instalado
- (Opcional) SQLite3 si usas `SQLStore`

Instalación y ejecución local

1. Clona el repositorio:

```bash
git clone <tu-repo-url>
cd urlShortener
```

2. Descargar dependencias:

```bash
go mod download
```

3. Ejecutar la aplicación (usa SQLite por defecto):

```bash
go run main.go
```

La API empieza en `:8000` por defecto (ver `main.go`).

Endpoints
- POST `/shorten` (form-url-encoded): `url` => crea short code. Respuesta 201 con JSON.
- GET `/shorten/:shortCode` => devuelve el recurso y aumenta contador de accesos.
- PUT `/shorten/:shortCode` (form-url-encoded): `url` => actualiza la URL original.
- DELETE `/shorten/:shortCode` => borra el short code.
- GET `/shorten/:shortCode/stats` => devuelve estadísticas (`accessCount`, etc.).

Tests

Ejecutar la suite de tests:

```bash
go test ./... -v
```

CI

Hay un workflow de GitHub Actions en `.github/workflows/ci.yml` que corre `gofmt -l .`, `go vet` y `go test` en Go 1.20/1.21/1.22 en push/PR a `main`.

Notas y próximos pasos sugeridos
- Añadir `golangci-lint` para linteo más estricto.
- Añadir tests de integración que usen una base SQLite temporal y apliquen `migrations/` antes de correr.
- Añadir badge de estado del CI en este `README.md` después de hacer push al remoto.

Si quieres que cree el badge del CI o que escriba la sección de contribución/PR flow, dímelo y lo agrego.
