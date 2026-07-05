# PetNexus Backend

PetNexus is a digital pet passport and owner-controlled pet identity platform. This directory contains its Go REST API.

## Stack

Sprint 6 uses Go, Gin, godotenv, PostgreSQL, GORM, Docker Compose, bcrypt, and JWT access tokens. Safe SQL migrations run automatically at startup; the versioned SQL files can also be applied manually with `psql`.

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
table, `owner_profiles`, `breeds`, and `pets`. Unique indexes enforce one
account per email, one owner profile per user, and unique breed names per
species. Startup stops immediately
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
Get-Content .\migrations\004_create_breeds_and_pets.sql | docker exec -i petnexus-postgres psql -v ON_ERROR_STOP=1 -U postgres -d petnexus
Get-Content .\migrations\005_create_clinic_profiles.sql | docker exec -i petnexus-postgres psql -v ON_ERROR_STOP=1 -U postgres -d petnexus
```

If Docker shows a different container name, replace `petnexus-postgres` in the
commands. These migrations create the currently implemented auth, owner
profile, breed catalog, basic pet profile, and clinic profile schema only.

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

## Sprint 5: Breed + Pet Creation Backend

Sprint 5 adds the breed reference catalog and owner-controlled basic pet
profiles. It supports this relationship without adding passport or clinic
behavior:

```text
users 1:1 owner_profiles
owner_profiles 1:N pets
breeds 1:N pets
```

New tables are `breeds` and `pets`. The startup migration safely seeds 8 dog
and 8 cat breeds with `ON CONFLICT DO NOTHING`. `breed_id` is optional, but a
provided breed must exist and match the pet species.

Endpoints:

```text
GET   /api/breeds
GET   /api/breeds?species=dog
GET   /api/breeds?species=cat
POST  /api/pets
GET   /api/pets
GET   /api/pets/:id
PATCH /api/pets/:id
```

Breed listing is public. Every pet endpoint requires JWT authentication, role
`owner`, and an existing owner profile. Ownership is always resolved from the
JWT; request DTOs do not accept `user_id` or `owner_profile_id`. Accessing
another owner's pet returns 404.

Example create request:

```json
{
  "species": "dog",
  "name": "Milo",
  "breed_id": "optional-breed-uuid",
  "gender": "male",
  "date_of_birth": "2022-05-10",
  "weight_kg": 12.5,
  "microchip_id": "MC-123456789",
  "avatar_url": "https://example.com/milo.png",
  "color": "Brown",
  "distinctive_marks": "White spot on chest",
  "is_neutered": true
}
```

The response contains basic pet data, computed `age_years`, and a safe breed
object. It does not contain owner or user IDs.

### Sprint 5 PowerShell test

Start Docker and `go run ./cmd/api`, then run:

```powershell
$baseUrl = "http://localhost:8080"
Invoke-RestMethod "$baseUrl/health"
Invoke-RestMethod "$baseUrl/health/db"

$suffix = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$email = "sprint5.owner.$suffix@example.com"
$password = "password123"
$register = @{ email=$email; phone="0812345678"; password=$password; role="owner" } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/register" -ContentType "application/json" -Body $register
$login = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/login" -ContentType "application/json" `
  -Body (@{ email=$email; password=$password } | ConvertTo-Json)
$headers = @{ Authorization = "Bearer $($login.data.accessToken)" }

$profile = @{ first_name="Sunny"; last_name="Example"; phone_number="0812345678" } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/owner/profile" -Headers $headers -ContentType "application/json" -Body $profile

$breeds = Invoke-RestMethod "$baseUrl/api/breeds"
$dogs = Invoke-RestMethod "$baseUrl/api/breeds?species=dog"
$cats = Invoke-RestMethod "$baseUrl/api/breeds?species=cat"
$breedId = ($dogs.data | Where-Object name -eq "Poodle" | Select-Object -First 1).id

$petBody = @{
  species="dog"; name="Milo"; breed_id=$breedId; gender="male"
  date_of_birth="2022-05-10"; weight_kg=12.5; microchip_id="MC-123456789"
  avatar_url="https://example.com/milo.png"; color="Brown"
  distinctive_marks="White spot on chest"; is_neutered=$true
} | ConvertTo-Json
$pet = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/pets" -Headers $headers -ContentType "application/json" -Body $petBody

Invoke-RestMethod -Method Get -Uri "$baseUrl/api/pets" -Headers $headers
Invoke-RestMethod -Method Get -Uri "$baseUrl/api/pets/$($pet.data.id)" -Headers $headers
$patch = @{ name="Milo Updated"; weight_kg=13.2 } | ConvertTo-Json
Invoke-RestMethod -Method Patch -Uri "$baseUrl/api/pets/$($pet.data.id)" -Headers $headers -ContentType "application/json" -Body $patch
```

Negative checks:

```powershell
# No token: expect 401
try { Invoke-RestMethod "$baseUrl/api/pets" } catch { $_.Exception.Response.StatusCode.value__ }

# Invalid species: expect 400
$invalid = @{ species="bird"; name="Invalid" } | ConvertTo-Json
try { Invoke-RestMethod -Method Post -Uri "$baseUrl/api/pets" -Headers $headers -ContentType "application/json" -Body $invalid } catch { $_.Exception.Response.StatusCode.value__ }

# Dog with a cat breed: expect 400
$catBreedId = ($cats.data | Select-Object -First 1).id
$mismatch = @{ species="dog"; name="Mismatch"; breed_id=$catBreedId } | ConvertTo-Json
try { Invoke-RestMethod -Method Post -Uri "$baseUrl/api/pets" -Headers $headers -ContentType "application/json" -Body $mismatch } catch { $_.Exception.Response.StatusCode.value__ }

# Empty PATCH: expect 400
try { Invoke-RestMethod -Method Patch -Uri "$baseUrl/api/pets/$($pet.data.id)" -Headers $headers -ContentType "application/json" -Body '{}' } catch { $_.Exception.Response.StatusCode.value__ }

# Invalid pet UUID: expect 400
try { Invoke-RestMethod -Method Get -Uri "$baseUrl/api/pets/not-a-uuid" -Headers $headers } catch { $_.Exception.Response.StatusCode.value__ }
```

Clinic staff must receive 403:

```powershell
$clinicEmail = "sprint5.clinic.$suffix@example.com"
$clinicRegister = @{
  email=$clinicEmail; phone="0823456789"; password=$password; role="clinic_staff"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/register" -ContentType "application/json" -Body $clinicRegister
$clinicLogin = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/login" -ContentType "application/json" `
  -Body (@{ email=$clinicEmail; password=$password } | ConvertTo-Json)
$clinicHeaders = @{ Authorization = "Bearer $($clinicLogin.data.accessToken)" }
try {
  Invoke-RestMethod -Method Get -Uri "$baseUrl/api/pets" -Headers $clinicHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 403
}
```

A second owner must create an owner profile before creating pets, and must not
be able to read the first owner's pet:

```powershell
$otherEmail = "sprint5.other-owner.$suffix@example.com"
$otherRegister = @{
  email=$otherEmail; phone="0834567890"; password=$password; role="owner"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/register" -ContentType "application/json" -Body $otherRegister
$otherLogin = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/auth/login" -ContentType "application/json" `
  -Body (@{ email=$otherEmail; password=$password } | ConvertTo-Json)
$otherHeaders = @{ Authorization = "Bearer $($otherLogin.data.accessToken)" }

# Owner without owner profile: expect 404
$minimalPet = @{ species="dog"; name="No Profile Yet" } | ConvertTo-Json
try {
  Invoke-RestMethod -Method Post -Uri "$baseUrl/api/pets" -Headers $otherHeaders -ContentType "application/json" -Body $minimalPet
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 404
}

$otherProfile = @{
  first_name="Second"; last_name="Owner"; phone_number="0834567890"
} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "$baseUrl/api/owner/profile" -Headers $otherHeaders -ContentType "application/json" -Body $otherProfile

# Another owner's pet is hidden: expect 404
try {
  Invoke-RestMethod -Method Get -Uri "$baseUrl/api/pets/$($pet.data.id)" -Headers $otherHeaders
} catch {
  $_.Exception.Response.StatusCode.value__ # Expected: 404
}
```

## Sprint 6: Clinic Profile Foundation

Sprint 6 adds one settings/identity profile for each authenticated clinic-side
account. The existing role `clinic_staff` is used; no new enum value is needed.
This foundation does not include staff management, owners/patients, QR,
authorization, visits, medical records, calendar, or reports.

New table:

```text
users 1:1 clinic_profiles
```

Clinic profile fields:

```text
id, user_id, clinic_name, phone_number, email, address, created_at, updated_at
```

`clinic_name` is required. Phone, email, and address are optional. `user_id`
always comes from JWT and is never accepted or returned by the profile API.

Endpoints:

```text
POST  /api/clinic/profile
GET   /api/clinic/profile
PATCH /api/clinic/profile
```

All three endpoints require `Authorization: Bearer <token>` and role
`clinic_staff`. Missing/invalid authentication returns 401; an owner token
returns 403. GET/PATCH before profile creation returns 404. Creating a second
profile returns 409.

Example create request:

```json
{
  "clinic_name": "Happy Paws Veterinary Clinic",
  "phone_number": "02-123-4567",
  "email": "contact@happypaws.example",
  "address": "123 Pet Street, Bangkok"
}
```

Example success response:

```json
{
  "success": true,
  "message": "Clinic profile created successfully",
  "data": {
    "id": "65eedfa2-4514-4ec8-8e14-c12b3354b762",
    "clinic_name": "Happy Paws Veterinary Clinic",
    "phone_number": "02-123-4567",
    "email": "contact@happypaws.example",
    "address": "123 Pet Street, Bangkok",
    "created_at": "2026-07-05T09:00:00Z",
    "updated_at": "2026-07-05T09:00:00Z"
  }
}
```

Example partial update:

```json
{
  "clinic_name": "Happy Paws Clinic Bangkok",
  "phone_number": "02-999-9999"
}
```

### Sprint 6 curl test

Start PostgreSQL and `go run ./cmd/api`, then register/login a clinic staff
account. Use a unique email if repeating the test.

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"clinic@example.com","phone":"021234567","password":"password123","role":"clinic_staff"}'

curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"clinic@example.com","password":"password123"}'
```

Copy `data.accessToken` from login and replace `<clinic_access_token>`:

```bash
# Before create: expect 404
curl http://localhost:8080/api/clinic/profile \
  -H "Authorization: Bearer <clinic_access_token>"

# Create: expect 201
curl -X POST http://localhost:8080/api/clinic/profile \
  -H "Authorization: Bearer <clinic_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"clinic_name":"Happy Paws Veterinary Clinic","phone_number":"02-123-4567","email":"contact@happypaws.example","address":"123 Pet Street, Bangkok"}'

# Get: expect 200
curl http://localhost:8080/api/clinic/profile \
  -H "Authorization: Bearer <clinic_access_token>"

# Patch: expect 200
curl -X PATCH http://localhost:8080/api/clinic/profile \
  -H "Authorization: Bearer <clinic_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"clinic_name":"Happy Paws Clinic Bangkok","phone_number":"02-999-9999"}'

# Repeat create: expect 409
curl -X POST http://localhost:8080/api/clinic/profile \
  -H "Authorization: Bearer <clinic_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"clinic_name":"Duplicate Clinic"}'

# No token: expect 401
curl http://localhost:8080/api/clinic/profile
```

To verify owner-role protection, login as an `owner` and call the same GET with
the owner token; expect HTTP 403.

## Current status

Sprint 6 adds the `clinic_profiles` schema and clinic-staff-only create/get/PATCH
profile APIs. Sprint 1–5 health, auth, owner profile, breed, and pet endpoints
remain unchanged.

Pet Passport, QR sharing, clinic access requests, authorization decisions, visits, timelines, notifications, real file uploads, Flutter UI, and clinic web UI are deliberately not implemented in this sprint.

รายละเอียดสิ่งที่ทำแล้วและข้อมูลส่งต่องานอยู่ที่ [`docs/progress/README.md`](docs/progress/README.md)

## Recommended next step

Deploy/redeploy to Render and repeat the Sprint 6 clinic profile smoke flow.
Plan QR, clinic access, visits, or staff management only as separate,
explicitly scoped sprints.
