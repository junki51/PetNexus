# PetNexus Backend

PetNexus is a digital pet passport and owner-controlled pet identity platform. This directory contains its Go REST API.

## Stack

Sprint 3 uses Go, Gin, godotenv, PostgreSQL, GORM, Docker Compose, bcrypt, and JWT access tokens. Versioned migrations are currently applied with `psql`; golang-migrate can be introduced in a later database-tooling sprint.

## Architecture

The code follows one dependency direction:

```text
handler -> service -> repository -> database
```

- `cmd/api` starts the application.
- `internal/config` loads environment configuration.
- `internal/routes` registers endpoints.
- `internal/handlers` receives HTTP requests and sends responses.
- `internal/services` will contain business and permission rules.
- `internal/repositories` will contain database access.
- `internal/models` will contain database entities.
- `internal/dto` will contain API request and response shapes.
- `internal/middleware` will contain authentication and role checks.
- `internal/utils` contains shared helpers.
- `migrations` will contain versioned PostgreSQL schema changes.

## Prerequisites

Install:

- Go 1.22 or newer
- Docker Desktop with Docker Compose

Then open a terminal in this directory.

## Install dependencies

```bash
go mod tidy
```

Optionally copy `.env.example` to `.env`. On PowerShell:

```powershell
Copy-Item .env.example .env
```

The example settings match the local PostgreSQL container in `docker-compose.yml`.
Use a strong, private `JWT_SECRET` outside local development. Never commit a real secret to `.env.example`.

### Database configuration

- Render/production: set `DATABASE_URL` to the managed PostgreSQL connection string.
- Local Docker: leave `DATABASE_URL` empty and continue using `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and `DB_SSLMODE=disable`.

When both are present, `DATABASE_URL` takes precedence. The `GET /health/db` endpoint checks whichever connection is active.

## Start PostgreSQL

```bash
docker compose up -d
```

Check that the container is running and healthy:

```bash
docker compose ps
```

The container exposes PostgreSQL on `localhost:5432` and stores its data in a named Docker volume.

## Run Sprint 3 migrations

The backend automatically runs a safe, idempotent SQL startup migration before
registering routes. It ensures `pgcrypto`, the `user_role` enum, the `users`
table, and the `idx_users_email_unique` unique index. Startup stops immediately
with a clear error if schema migration fails. This supports fresh Render
PostgreSQL databases and avoids GORM `AutoMigrate` constraint rewrites on
existing databases.

The commands below remain available when you want to apply or inspect the SQL
manually during local development.

First confirm the container name:

```powershell
docker ps
```

If the container is named `petnexus-postgres`, run:

```powershell
Get-Content .\migrations\001_create_enums.sql | docker exec -i petnexus-postgres psql -v ON_ERROR_STOP=1 -U postgres -d petnexus
Get-Content .\migrations\002_create_users.sql | docker exec -i petnexus-postgres psql -v ON_ERROR_STOP=1 -U postgres -d petnexus
```

If Docker shows a different container name, replace `petnexus-postgres` in both commands. These migrations create only `user_role` and `users`.

## Run the backend

```bash
go run ./cmd/api
```

The backend connects to PostgreSQL before starting the HTTP server. It exits with a clear error if the database is unavailable.

The default API address is `http://localhost:8080`. Set `PORT` to use another port.

## Test the health endpoints

```bash
curl http://localhost:8080/health
curl http://localhost:8080/health/db
```

Expected response from `GET /health`:

```json
{
  "success": true,
  "message": "PetNexus backend is running",
  "data": {
    "status": "ok",
    "service": "petnexus-backend"
  }
}
```

Expected response from `GET /health/db`:

```json
{
  "success": true,
  "message": "Database connection is healthy",
  "data": {
    "database": "postgresql",
    "status": "connected"
  }
}
```

Stop the local database when finished:

```bash
docker compose down
```

## Test authentication

PowerShell has a `curl` alias on some versions, so the examples call `curl.exe` explicitly.

### Register an owner

```powershell
curl.exe -X POST http://localhost:8080/api/auth/register `
  -H "Content-Type: application/json" `
  -d "{\"email\":\"owner@example.com\",\"phone\":\"0812345678\",\"password\":\"password123\",\"role\":\"owner\"}"
```

Expected: HTTP 201, `success: true`, a safe user object, and `accessToken`. The response must not contain `passwordHash`.

Run the same command again to verify it returns HTTP 409 with `EMAIL_ALREADY_EXISTS`.

### Login

```powershell
curl.exe -X POST http://localhost:8080/api/auth/login `
  -H "Content-Type: application/json" `
  -d "{\"email\":\"owner@example.com\",\"password\":\"password123\"}"
```

Expected: `success: true` and a new `accessToken`.

### Get the current user

Copy the login token into the command:

```powershell
curl.exe http://localhost:8080/api/me `
  -H "Authorization: Bearer <access_token>"
```

Expected: the current user without `passwordHash`.

Without a token, the same endpoint must return HTTP 401:

```powershell
curl.exe http://localhost:8080/api/me
```

### Verify public admin registration is blocked

```powershell
curl.exe -X POST http://localhost:8080/api/auth/register `
  -H "Content-Type: application/json" `
  -d "{\"email\":\"admin@example.com\",\"phone\":\"0800000000\",\"password\":\"password123\",\"role\":\"admin\"}"
```

Expected: HTTP 403 with `FORBIDDEN_ROLE`. Public registration supports only `owner` and `clinic_staff`.

## Current status

Sprint 3 adds the `users` schema, bcrypt password storage, JWT access tokens, register/login APIs, authentication middleware, role middleware, and protected `GET /api/me`. Both health endpoints remain public and unchanged.

Owner profiles, pet CRUD, breeds, QR sessions, clinic access requests, authorization decisions, visits, timelines, notifications, audit logs, refresh tokens, and password recovery are deliberately not implemented.

รายละเอียดสิ่งที่ทำแล้วและข้อมูลส่งต่องานอยู่ที่ [`docs/progress/README.md`](docs/progress/README.md)

## Recommended next step

Continue with Sprint 4: Owner Profile. Keep owner profile tables, repository, service, and routes separate from the completed authentication foundation.
