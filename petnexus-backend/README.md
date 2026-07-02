# PetNexus Backend

PetNexus is a digital pet passport and owner-controlled pet identity platform. This directory contains its Go REST API.

## Stack

Sprint 4 uses Go, Gin, godotenv, PostgreSQL, GORM, Docker Compose, bcrypt, and JWT access tokens. Safe SQL migrations run automatically at startup; the versioned SQL files can also be applied manually with `psql`.

## Architecture

The code follows one dependency direction:

```text
handler -> service -> repository -> database
```

- `cmd/api` starts the application.
- `internal/config` loads environment configuration.
- `internal/routes` registers endpoints.
- `internal/handlers` receives HTTP requests and sends responses.
- `internal/services` contains business and permission rules.
- `internal/repositories` contains database access.
- `internal/models` contains database entities.
- `internal/dto` contains API request and response shapes.
- `internal/middleware` contains authentication and role checks.
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

## Run startup migrations

The backend automatically runs a safe, idempotent SQL startup migration before
registering routes. It ensures `pgcrypto`, the `user_role` enum, the `users`
table, and the Sprint 4 `owner_profiles` table. Unique indexes enforce one
account per email and one owner profile per user. Startup stops immediately
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
Get-Content .\migrations\003_create_owner_profiles.sql | docker exec -i petnexus-postgres psql -v ON_ERROR_STOP=1 -U postgres -d petnexus
```

If Docker shows a different container name, replace `petnexus-postgres` in the commands. These migrations create only the implemented auth and owner-profile schema.

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

## Sprint 4: Owner Profile

An owner profile stores the pet owner's identity and contact details. It is
separate from `users` so authentication data (email, password hash, and role)
does not become coupled to editable profile data. The `owner_profiles` table
has a unique foreign key to `users.id`, giving each account at most one owner
profile.

All owner-profile endpoints require a valid JWT and the `owner` role:

```text
POST  /api/owner/profile
GET   /api/owner/profile
PATCH /api/owner/profile
```

The backend always gets `user_id` from the JWT. Clients must not send or choose
it. `display_name` is computed from `first_name` and `last_name` and is not
stored in PostgreSQL.

Example create request:

```json
{
  "first_name": "Sunny",
  "last_name": "Example",
  "gender": "male",
  "date_of_birth": "2008-01-01",
  "phone_number": "0812345678",
  "avatar_url": "https://example.com/avatar.png",
  "address_line1": "123 Pet Street",
  "address_line2": "",
  "province": "Bangkok",
  "district": "Bang Rak",
  "subdistrict": "Si Phraya",
  "postal_code": "10500"
}
```

Example response data:

```json
{
  "id": "4ec32b2d-16c8-4ac9-99e1-988b08a8cb42",
  "first_name": "Sunny",
  "last_name": "Example",
  "display_name": "Sunny Example",
  "gender": "male",
  "date_of_birth": "2008-01-01",
  "phone_number": "0812345678",
  "avatar_url": "https://example.com/avatar.png",
  "address_line1": "123 Pet Street",
  "address_line2": null,
  "province": "Bangkok",
  "district": "Bang Rak",
  "subdistrict": "Si Phraya",
  "postal_code": "10500",
  "created_at": "2026-07-02T08:00:00Z",
  "updated_at": "2026-07-02T08:00:00Z"
}
```

### PowerShell manual test

Start PostgreSQL and the API first, then run these commands in a separate
PowerShell terminal.

Health checks:

```powershell
$baseUrl = "http://localhost:8080"
Invoke-RestMethod -Method Get -Uri "$baseUrl/health"
Invoke-RestMethod -Method Get -Uri "$baseUrl/health/db"
```

Register an owner, log in, and save the access token:

```powershell
$suffix = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$ownerEmail = "owner.$suffix@example.com"
$password = "password123"

$registerBody = @{
  email = $ownerEmail
  phone = "0812345678"
  password = $password
  role = "owner"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/register" -ContentType "application/json" -Body $registerBody

$loginBody = @{ email = $ownerEmail; password = $password } | ConvertTo-Json
$login = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/login" -ContentType "application/json" -Body $loginBody
$token = $login.data.accessToken
$ownerHeaders = @{ Authorization = "Bearer $token" }
```

Create, fetch, and patch the profile:

```powershell
$profileBody = @{
  first_name = "Sunny"
  last_name = "Example"
  gender = "male"
  date_of_birth = "2008-01-01"
  phone_number = "0812345678"
  avatar_url = "https://example.com/avatar.png"
  address_line1 = "123 Pet Street"
  address_line2 = ""
  province = "Bangkok"
  district = "Bang Rak"
  subdistrict = "Si Phraya"
  postal_code = "10500"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/owner/profile" -Headers $ownerHeaders -ContentType "application/json" -Body $profileBody

Invoke-RestMethod -Method Get -Uri "$baseUrl/api/owner/profile" -Headers $ownerHeaders

$patchBody = @{ first_name = "Sunny Updated"; phone_number = "0899999999" } | ConvertTo-Json
Invoke-RestMethod -Method Patch -Uri "$baseUrl/api/owner/profile" -Headers $ownerHeaders -ContentType "application/json" -Body $patchBody
```

Creating the same profile again must return HTTP 409:

```powershell
try {
  Invoke-RestMethod -Method Post -Uri "$baseUrl/api/owner/profile" -Headers $ownerHeaders -ContentType "application/json" -Body $profileBody
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 409
}
```

A request without a token must return HTTP 401:

```powershell
try {
  Invoke-RestMethod -Method Get -Uri "$baseUrl/api/owner/profile"
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 401
}
```

A `clinic_staff` token must return HTTP 403:

```powershell
$clinicEmail = "clinic.$suffix@example.com"
$clinicRegisterBody = @{
  email = $clinicEmail
  phone = "0823456789"
  password = $password
  role = "clinic_staff"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/register" -ContentType "application/json" -Body $clinicRegisterBody

$clinicLoginBody = @{ email = $clinicEmail; password = $password } | ConvertTo-Json
$clinicLogin = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/login" -ContentType "application/json" -Body $clinicLoginBody
$clinicHeaders = @{ Authorization = "Bearer $($clinicLogin.data.accessToken)" }
try {
  Invoke-RestMethod -Method Get -Uri "$baseUrl/api/owner/profile" -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 403
}
```

## Current status

Sprint 4 adds the `owner_profiles` schema and owner-only create, get, and partial-update APIs on top of the completed authentication foundation. Existing auth and health endpoints remain unchanged.

Pet CRUD, breeds, pet passports, QR sharing, clinic access requests, authorization decisions, visits, timelines, Flutter UI, and clinic web UI are deliberately not implemented in this sprint.

รายละเอียดสิ่งที่ทำแล้วและข้อมูลส่งต่องานอยู่ที่ [`docs/progress/README.md`](docs/progress/README.md)

## Recommended next step

Continue with the next agreed database/domain sprint, preferably Pet and Breed foundations, without mixing in QR, clinic authorization, visit, or timeline behavior prematurely.
