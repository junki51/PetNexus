# PetNexus Backend Setup Checklist

## 1. Purpose

This checklist is used to verify that the PetNexus backend foundation is set up correctly before building real features.

Use this after Codex creates or modifies the backend scaffold.

---

## 2. Expected Project Root

Expected folder:

```txt
petnexus-backend/
```

Expected command:

```bash
go run ./cmd/api
```

Expected first endpoint:

```txt
GET /health
```

---

## 3. Required Folder Structure

Check that these folders exist:

```txt
cmd/
cmd/api/

internal/
internal/config/
internal/database/
internal/middleware/
internal/models/
internal/repositories/
internal/services/
internal/handlers/
internal/routes/
internal/dto/
internal/utils/

migrations/
docs/
```

If any folder is missing, add it before building features.

---

## 4. Required Root Files

Check that these files exist:

```txt
go.mod
README.md
.env.example
.gitignore
```

Optional later:

```txt
docker-compose.yml
Makefile
```

Do not add optional files too early if they slow progress.

---

## 5. Required Docs

The docs folder should include:

```txt
docs/context.md
docs/backend-roadmap.md
docs/api-plan.md
docs/database-plan.md
docs/backend-codex-rules.md
docs/backend-setup-checklist.md
```

Optional later:

```txt
docs/test-flow.md
docs/deployment-plan.md
docs/demo-script.md
```

---

## 6. Go Module Checklist

Run:

```bash
go mod tidy
```

Check:

```txt
go.mod exists
go.sum exists after dependencies are added
module name is reasonable
dependencies are not excessive
```

Recommended module name example:

```txt
github.com/your-name/petnexus-backend
```

Do not worry if the final GitHub username is not ready yet. It can be changed later.

---

## 7. Required Dependencies for Foundation

For Sprint 1, expected dependencies:

```txt
github.com/gin-gonic/gin
github.com/joho/godotenv
```

For database sprint later:

```txt
gorm.io/gorm
gorm.io/driver/postgres
github.com/google/uuid
```

For auth sprint later:

```txt
github.com/golang-jwt/jwt/v5
golang.org/x/crypto/bcrypt
```

Do not add all dependencies at once if they are not used yet.

---

## 8. Environment Checklist

`.env.example` should include:

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

`.env` should not be committed.

`.gitignore` should include:

```txt
.env
*.exe
tmp/
dist/
bin/
.DS_Store
```

---

## 9. Config Checklist

`internal/config/config.go` should be responsible for:

```txt
loading env variables
providing default values
returning config struct
```

Config should include:

```txt
AppEnv
Port
DBHost
DBPort
DBUser
DBPassword
DBName
DBSSLMode
JWTSecret
JWTExpiresIn
```

For Sprint 1, database and JWT config can exist as placeholders.

Do not connect to database from config package.

---

## 10. Database Checklist

`internal/database/postgres.go` should be responsible for:

```txt
building PostgreSQL DSN
connecting to PostgreSQL
returning *gorm.DB
```

For Sprint 1, database connection may be placeholder or not called yet.

For Sprint 2, it must connect for real.

Check later:

```txt
database connection uses env variables
database connection error is handled
database instance is passed to repositories
```

Do not hardcode database credentials.

---

## 11. Route Checklist

`internal/routes/routes.go` should register routes.

Minimum Sprint 1 route:

```txt
GET /health
```

Future route groups:

```txt
/api/auth
/api/owner
/api/pets
/api/breeds
/api/clinic
/api/authorizations
/api/notifications
```

Do not put business logic in route registration.

---

## 12. Handler Checklist

`internal/handlers/health_handler.go` should contain health handler.

Expected route:

```txt
GET /health
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

Handler should use shared response helper if available.

---

## 13. Response Helper Checklist

`internal/utils/response.go` should provide consistent response functions.

Recommended helpers:

```txt
Success(c, statusCode, message, data)
Error(c, statusCode, message, code, details)
```

Response format must be consistent.

Do not return random response shapes in different handlers.

---

## 14. Middleware Checklist

Sprint 1 placeholders:

```txt
auth_middleware.go
role_middleware.go
```

They can contain TODO comments.

Do not fake auth behavior.

Do not implement insecure temporary auth.

When auth sprint starts, middleware should:

```txt
read Authorization header
validate JWT
extract user id
put current user into context
reject invalid token
```

---

## 15. Model Checklist

Model files can be placeholders in Sprint 1.

Expected model files:

```txt
user.go
owner_profile.go
pet.go
breed.go
clinic.go
clinic_staff.go
qr_session.go
authorization.go
visit.go
notification.go
audit_log.go
```

Each file should include short comment explaining what the model represents.

Do not add incomplete broken GORM tags unless implementing database sprint.

---

## 16. Repository Checklist

Repository files can be placeholders in Sprint 1.

Expected repository files:

```txt
user_repository.go
owner_repository.go
pet_repository.go
breed_repository.go
clinic_repository.go
qr_repository.go
authorization_repository.go
visit_repository.go
notification_repository.go
audit_log_repository.go
```

Repositories should later own database operations.

Do not put HTTP request logic in repository.

---

## 17. Service Checklist

Service files can be placeholders in Sprint 1.

Expected service files:

```txt
auth_service.go
owner_service.go
pet_service.go
clinic_service.go
qr_service.go
authorization_service.go
visit_service.go
timeline_service.go
notification_service.go
audit_log_service.go
```

Services should later own business rules.

Important future service rules:

```txt
PetService checks owner ownership.
QRService creates temporary secure tokens.
AuthorizationService handles approve/reject/revoke.
VisitService checks create_visit permission.
TimelineService returns only allowed timeline.
```

---

## 18. DTO Checklist

DTO files can be placeholders in Sprint 1.

Expected DTO files:

```txt
auth_dto.go
pet_dto.go
clinic_dto.go
qr_dto.go
authorization_dto.go
visit_dto.go
```

DTOs should keep API request/response separate from DB models.

Do not expose password_hash in response DTOs.

---

## 19. Migration Checklist

Migration folder should exist.

For Sprint 1:

```txt
migrations/README.md
```

For Sprint 2:

```txt
001_create_enums.sql
002_create_users.sql
003_create_owner_profiles.sql
004_create_breeds.sql
005_create_pets.sql
006_create_clinics.sql
007_create_clinic_staff.sql
008_create_qr_sessions.sql
009_create_authorizations.sql
010_create_visits.sql
011_create_notifications.sql
012_create_audit_logs.sql
013_seed_breeds.sql
```

Do not write random migrations without matching `database-plan.md`.

---

## 20. Health Test Checklist

Start server:

```bash
go run ./cmd/api
```

Test:

```bash
curl http://localhost:8080/health
```

Expected:

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

If this does not work, do not continue to auth yet.

---

## 21. Sprint 1 Definition of Done

Sprint 1 is complete when:

```txt
Project runs with go run ./cmd/api
GET /health works
Folder structure exists
Config loading exists
Routes are registered from routes package
Response helper exists
README explains how to run
.env.example exists
No full feature is implemented yet
```

Do not continue to Sprint 2 until Sprint 1 is stable.

---

## 22. Sprint 2 Definition of Done

Sprint 2 is complete when:

```txt
PostgreSQL database exists
Backend can connect to database
GORM is configured
Migration tool is chosen
Core migrations are prepared
Breed seed data exists
Database connection errors are handled
README explains database setup
```

Do not implement auth before database foundation is ready unless instructed.

---

## 23. Sprint 3 Definition of Done

Sprint 3 is complete when:

```txt
Owner can register
Clinic staff can register
User can login
Password is hashed with bcrypt
JWT token is returned
GET /api/me works
Protected route rejects missing token
Role middleware can restrict owner / clinic_staff
```

Do not build pet features until auth works.

---

## 24. Sprint 4 Definition of Done

Sprint 4 is complete when:

```txt
Owner can create profile
Owner can fetch profile
Owner can update profile
Only owner role can use owner profile endpoints
One user cannot create duplicate owner profile
```

---

## 25. Sprint 5 Definition of Done

Sprint 5 is complete when:

```txt
Breeds are seeded
Owner can fetch breeds
Owner can create pet
Pet gets unique PetNexus ID
Owner can list own pets
Owner can fetch own pet detail
Owner cannot fetch another owner’s pet
```

---

## 26. Sprint 6 Definition of Done

Sprint 6 is complete when:

```txt
Owner can create QR session for own pet
QR session has secure token
QR session expires
Clinic can scan QR token
Expired token is rejected
Invalid token is rejected
Clinic sees pet preview only
Audit log records clinic scan
```

---

## 27. Sprint 7 Definition of Done

Sprint 7 is complete when:

```txt
Clinic can request access after scanning QR
Duplicate pending access request is blocked
Owner can see pending access requests
Notification is created for owner
Audit log records request
```

---

## 28. Sprint 8 Definition of Done

Sprint 8 is complete when:

```txt
Owner can approve access
Owner can reject access
Owner can revoke approved access
Only pet owner can decide authorization
Approved clinic can access full pet record
Rejected clinic cannot access full pet record
Revoked clinic cannot access full pet record
```

---

## 29. Sprint 9 Definition of Done

Sprint 9 is complete when:

```txt
Clinic can list approved patients
Clinic can view full approved pet record
Clinic cannot view unapproved pet
Clinic cannot view expired authorization
Clinic cannot view revoked authorization
```

---

## 30. Sprint 10 Definition of Done

Sprint 10 is complete when:

```txt
Clinic can create visit for approved pet
Authorization must include create_visit
Visit is marked clinic_verified
Owner receives notification
Audit log records clinic_created_visit
Visit appears in timeline
```

---

## 31. Full MVP Backend Test Flow

The full backend MVP is ready when this flow works:

```txt
1. Owner registers
2. Owner logs in
3. Owner creates profile
4. Owner creates pet
5. Owner creates QR session
6. Clinic registers
7. Clinic logs in
8. Clinic scans QR
9. Clinic sees pet preview only
10. Clinic requests access
11. Owner sees access request
12. Owner approves access
13. Clinic views full pet record
14. Clinic creates clinic_verified visit
15. Owner views timeline
16. Owner sees visit notification
17. Audit logs are created
```

This is the main demo path.

---

## 32. Stop Conditions

Stop and fix before continuing if:

```txt
/health does not work
go run fails
response format is inconsistent
handlers directly query database
permission logic is in frontend only
QR returns full pet data
clinic can view pet without approval
clinic can create visit without permission
owner can edit clinic_verified visit
.env is committed
README is outdated
```

Do not build more features on top of broken foundation.

---

## 33. Review Checklist Before Every Commit

Before committing, check:

```txt
Does the app still run?
Does /health still work?
Did I keep the correct layer separation?
Did I update README if needed?
Did I avoid adding out-of-scope features?
Did I avoid hardcoded secrets?
Did I keep response format consistent?
Did I protect owner/clinic permission rules?
```

If any answer is no, fix before commit.

---

## 34. Final Reminder

PetNexus backend should be simple, secure, and understandable.

Good MVP backend is not the one with the most features.

Good MVP backend is the one where the core flow works clearly and safely.
