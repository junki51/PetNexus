# PetNexus Backend

PetNexus is a digital pet passport and owner-controlled pet identity platform. This directory contains its Go REST API.

## Stack

Sprint 1 uses Go, Gin, and godotenv. PostgreSQL, GORM, JWT, bcrypt, and golang-migrate are planned for later sprints and are intentionally inactive.

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

Install Go 1.22 or newer, then open a terminal in this directory.

## Install dependencies

```bash
go mod tidy
```

Optionally copy `.env.example` to `.env`. On PowerShell:

```powershell
Copy-Item .env.example .env
```

## Run

```bash
go run ./cmd/api
```

The default address is `http://localhost:8080`. Set `PORT` to use another port.

## Test the health endpoint

```bash
curl http://localhost:8080/health
```

Expected response:

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

## Current status

Sprint 1 provides the server foundation, environment loading, consistent response helpers, route registration, and `/health` only.

Registration, login, JWT, password hashing, PostgreSQL/GORM, pet CRUD, QR sessions, authorization, clinic visits, timelines, notifications, and audit-log logic are deliberately not implemented. No fake security behavior is included.

## Recommended next step

Complete the database foundation in `docs/database-plan.md`: add PostgreSQL, GORM, and golang-migrate while keeping database access inside repositories.
