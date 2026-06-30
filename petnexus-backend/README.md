# PetNexus Backend

PetNexus is a digital pet passport and owner-controlled pet identity platform. This directory contains its Go REST API.

## Stack

Sprint 2 uses Go, Gin, godotenv, PostgreSQL, GORM, and Docker Compose. JWT, bcrypt, and golang-migrate are reserved for later sprints.

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

## Start PostgreSQL

```bash
docker compose up -d
```

Check that the container is running and healthy:

```bash
docker compose ps
```

The container exposes PostgreSQL on `localhost:5432` and stores its data in a named Docker volume.

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

## Current status

Sprint 2 adds a PostgreSQL container, a verified GORM connection at startup, and `GET /health/db`. The original `GET /health` endpoint remains available.

No tables or migrations are created yet. Registration, login, JWT, password hashing, pet CRUD, QR sessions, authorization, clinic visits, timelines, notifications, and audit-log logic are deliberately not implemented.

## Recommended next step

Follow `docs/database-plan.md` to design and add versioned PostgreSQL migrations with `golang-migrate`. Keep schema work separate from authentication and feature APIs.
