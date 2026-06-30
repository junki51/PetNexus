# Codex Prompt: PetNexus Backend Sprint 1

You are a senior Go backend engineer and software architect.

I am building the backend for an MVP project called **PetNexus**.

PetNexus is a Digital Pet Passport / Digital Pet Identity Platform for pet owners and small-to-medium veterinary clinics.

Core MVP flow:

```txt
Owner creates pet
→ Owner shares QR
→ Clinic scans QR
→ Clinic requests access
→ Owner approves
→ Clinic creates verified visit
→ Owner sees updated timeline
```

For this task, do **Sprint 1 only**.

Do not implement full features yet.

---

## Tech Stack

Use:

```txt
Go
Gin
godotenv
```

Prepare placeholders for:

```txt
PostgreSQL
GORM
JWT
bcrypt
golang-migrate
```

Do not fully implement database, auth, QR, authorization, or visit yet.

---

## Task Goal

Create the backend foundation only.

The backend should be runnable with:

```bash
go run ./cmd/api
```

And this endpoint should work:

```txt
GET /health
```

Expected JSON response:

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

---

## Required Folder Structure

Create this structure:

```txt
petnexus-backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── role_middleware.go
│   │
│   ├── models/
│   │   ├── user.go
│   │   ├── owner_profile.go
│   │   ├── pet.go
│   │   ├── breed.go
│   │   ├── clinic.go
│   │   ├── clinic_staff.go
│   │   ├── qr_session.go
│   │   ├── authorization.go
│   │   ├── visit.go
│   │   ├── notification.go
│   │   └── audit_log.go
│   │
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── owner_repository.go
│   │   ├── pet_repository.go
│   │   ├── breed_repository.go
│   │   ├── clinic_repository.go
│   │   ├── qr_repository.go
│   │   ├── authorization_repository.go
│   │   ├── visit_repository.go
│   │   ├── notification_repository.go
│   │   └── audit_log_repository.go
│   │
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── owner_service.go
│   │   ├── pet_service.go
│   │   ├── clinic_service.go
│   │   ├── qr_service.go
│   │   ├── authorization_service.go
│   │   ├── visit_service.go
│   │   ├── timeline_service.go
│   │   ├── notification_service.go
│   │   └── audit_log_service.go
│   │
│   ├── handlers/
│   │   ├── health_handler.go
│   │   ├── auth_handler.go
│   │   ├── owner_handler.go
│   │   ├── pet_handler.go
│   │   ├── breed_handler.go
│   │   ├── clinic_handler.go
│   │   ├── qr_handler.go
│   │   ├── authorization_handler.go
│   │   ├── visit_handler.go
│   │   ├── timeline_handler.go
│   │   └── notification_handler.go
│   │
│   ├── dto/
│   │   ├── auth_dto.go
│   │   ├── pet_dto.go
│   │   ├── clinic_dto.go
│   │   ├── qr_dto.go
│   │   ├── authorization_dto.go
│   │   └── visit_dto.go
│   │
│   ├── utils/
│   │   ├── response.go
│   │   ├── password.go
│   │   ├── jwt.go
│   │   ├── token.go
│   │   └── validator.go
│   │
│   └── routes/
│       └── routes.go
│
├── migrations/
│   └── README.md
│
├── docs/
│   ├── context.md
│   ├── backend-roadmap.md
│   ├── api-plan.md
│   ├── database-plan.md
│   ├── backend-codex-rules.md
│   └── backend-setup-checklist.md
│
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

---

## Implementation Rules

Follow these rules strictly:

1. Implement Sprint 1 only.
2. Do not implement register/login yet.
3. Do not implement database connection for real yet unless it is only a safe placeholder.
4. Do not implement pet creation yet.
5. Do not implement QR logic yet.
6. Do not implement authorization logic yet.
7. Do not implement visit logic yet.
8. Create placeholder files with clear comments.
9. Keep the code beginner-friendly.
10. Use consistent response format.
11. Do not hardcode secrets.
12. Do not add unnecessary dependencies.
13. Do not add microservices, GraphQL, Redis, Docker, or Kubernetes.
14. Do not over-engineer.

---

## Layer Explanation Comments

Add useful comments in placeholder files.

Use this meaning:

```txt
handler      = receives HTTP request and returns response
service      = business logic and permission checks
repository   = database access
model        = database entity
middleware   = auth and role checks
routes       = route registration
config       = environment config
database     = PostgreSQL connection
dto          = request and response shape
utils        = shared helpers
```

---

## Required Working Code

## 1. main.go

Should:

```txt
load config
create Gin router
register routes
run server on configured port
```

Default port:

```txt
8080
```

---

## 2. config.go

Should:

```txt
load .env if available
read PORT
default PORT to 8080
prepare DB and JWT config fields for later
```

Do not fail if `.env` does not exist.

---

## 3. response.go

Should provide helpers:

```txt
Success(c, statusCode, message, data)
Error(c, statusCode, message, code, details)
```

Use consistent JSON format.

---

## 4. health_handler.go

Should provide:

```txt
HealthCheck
```

Response:

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

---

## 5. routes.go

Should register:

```txt
GET /health
```

Future route groups can be added as TODO comments only.

---

## 6. postgres.go

For Sprint 1, create a placeholder function and comments explaining that PostgreSQL connection will be implemented in Sprint 2.

Do not force the app to require PostgreSQL yet.

The app should run even without database.

---

## 7. auth_middleware.go

Create placeholder only.

Comment:

```txt
JWT authentication middleware will be implemented in Sprint 3.
```

Do not fake authentication.

---

## 8. role_middleware.go

Create placeholder only.

Comment:

```txt
Role-based middleware will protect owner and clinic_staff routes in Sprint 3.
```

Do not fake authorization.

---

## 9. Placeholder Feature Files

Each placeholder file should contain:

```txt
package name
short comment explaining future responsibility
```

Do not leave empty files if Go build fails.

Make sure all packages compile.

---

## 10. .env.example

Create:

```txt
APP_ENV=development
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=petnexus
DB_SSLMODE=disable

JWT_SECRET=change_me_in_production
JWT_EXPIRES_IN=24h
```

---

## 11. .gitignore

Include:

```txt
.env
*.exe
bin/
tmp/
dist/
.DS_Store
```

---

## 12. README.md

Create a beginner-friendly README.

Include:

```txt
Project name
Backend stack
How to install dependencies
How to run server
How to test /health
Current sprint status
What is intentionally not implemented yet
Next sprint
```

Example run command:

```bash
go run ./cmd/api
```

Example test command:

```bash
curl http://localhost:8080/health
```

---

## 13. migrations/README.md

Explain:

```txt
Migrations will be added in Sprint 2.
Database will use PostgreSQL.
Migration tool will be golang-migrate.
Do not create random tables before database-plan.md is followed.
```

---

## Sprint 1 Definition of Done

Sprint 1 is complete only if:

```txt
go run ./cmd/api works
GET /health returns expected JSON
Folder structure exists
Config loading works
Response helper exists
Routes package registers health route
README explains how to run
.env.example exists
No full feature is implemented yet
```

---

## After Making Changes

Summarize:

```txt
Files created
Files changed
How to run
How to test /health
What was intentionally not implemented
Recommended next step
```

Do not claim Sprint 2 is complete.

Do not say auth/database/pet/QR works yet.

Only Sprint 1 foundation should work.
