# PetNexus Backend Codex Rules

## 1. Purpose

This document defines strict coding rules for Codex when working on the **PetNexus Go Backend**.

PetNexus is a Digital Pet Passport / Digital Pet Identity Platform.

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

Codex must always protect this core flow.

Do not overbuild.

Do not add features outside the current sprint.

Do not change architecture without a clear reason.

---

## 2. Tech Stack

Use this backend stack:

```txt
Go
Gin
PostgreSQL
GORM
JWT
bcrypt
godotenv
golang-migrate
```

Do not replace the stack unless explicitly instructed.

Do not add extra frameworks without permission.

---

## 3. Architecture Rule

Use beginner-friendly layered architecture.

```txt
handler      = receives HTTP request and returns response
service      = business logic and permission checks
repository   = database access
model        = database entity
middleware   = authentication and role checks
routes       = route registration
config       = environment configuration
database     = PostgreSQL connection
dto          = request and response structs
utils        = shared helpers
```

Main rule:

```txt
Handler → Service → Repository → Database
```

Never skip layers without a strong reason.

---

## 4. Folder Responsibility

## 4.1 cmd/api

Contains the application entry point.

Allowed:

```txt
main.go
```

Responsibility:

* Load config
* Connect database
* Set up Gin router
* Register routes
* Start server

Do not put business logic here.

---

## 4.2 internal/config

Responsibility:

* Read environment variables
* Provide config struct
* Provide default values for development

Do not connect to database here.

---

## 4.3 internal/database

Responsibility:

* Connect to PostgreSQL
* Return database instance
* Configure GORM

Do not write business logic here.

---

## 4.4 internal/models

Responsibility:

* Define database entities
* Define GORM model structs
* Define model relationships

Do not define request body structs here unless they are truly database models.

---

## 4.5 internal/dto

Responsibility:

* Define request DTOs
* Define response DTOs
* Separate external API shape from database models

Example:

```txt
RegisterRequest
LoginRequest
CreatePetRequest
PetResponse
CreateVisitRequest
```

Do not return raw database models directly if sensitive fields may leak.

---

## 4.6 internal/repositories

Responsibility:

* Query database
* Create rows
* Update rows
* Delete rows if allowed
* Return models or domain data

Do not check business permission here unless it is a direct database condition needed by service.

Repositories should not know HTTP status codes.

---

## 4.7 internal/services

Responsibility:

* Business logic
* Permission checks
* Ownership checks
* Authorization checks
* QR token validation
* Visit creation rules
* Notification side effects
* Audit log side effects

Important:

Sensitive rules must live here.

Do not place important permission logic only in frontend or handlers.

---

## 4.8 internal/handlers

Responsibility:

* Bind JSON
* Validate required fields at basic level
* Call service
* Return consistent response
* Convert service errors to HTTP status codes

Handlers should be thin.

Bad:

```txt
Handler directly queries database.
Handler checks full authorization rules.
Handler generates QR token directly.
```

Good:

```txt
Handler parses request.
Handler calls service.
Service performs logic.
Handler returns response.
```

---

## 4.9 internal/middleware

Responsibility:

* Auth middleware
* Role middleware
* Request context setup
* Current user extraction

Do not put feature business logic here.

---

## 4.10 internal/routes

Responsibility:

* Register routes
* Group routes by feature
* Attach middleware

Routes should be easy to read.

---

## 4.11 internal/utils

Responsibility:

* Password hashing
* JWT helper
* Random token generator
* Response helper
* Error helper
* Validator helper

Do not place feature-specific logic here.

---

## 5. Response Format Rule

All API responses must use a consistent format.

## 5.1 Success Response

```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": {}
}
```

## 5.2 Error Response

```json
{
  "success": false,
  "message": "Something went wrong",
  "error": {
    "code": "ERROR_CODE",
    "details": "Readable error detail"
  }
}
```

## 5.3 List Response

```json
{
  "success": true,
  "message": "Data fetched successfully",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

Do not invent new response formats for each endpoint.

---

## 6. HTTP Status Code Rule

Use these status codes:

```txt
200 OK                    successful GET / PATCH
201 Created               successful create
400 Bad Request           invalid request body
401 Unauthorized          missing or invalid token
403 Forbidden             logged in but not allowed
404 Not Found             resource not found
409 Conflict              duplicate or conflicting action
422 Unprocessable Entity  validation error
500 Internal Server Error unexpected server error
```

Do not return `200 OK` for errors.

---

## 7. Auth Rule

Use JWT for protected routes.

Protected routes must require:

```txt
Authorization: Bearer <access_token>
```

Public routes:

```txt
GET  /health
POST /api/auth/register
POST /api/auth/login
POST /api/clinic/register
```

All other routes should be protected unless explicitly stated.

---

## 8. Role Rule

Supported MVP roles:

```txt
owner
clinic_staff
admin
```

Owner-only endpoints must reject clinic staff.

Clinic-only endpoints must reject owners.

Never trust frontend role checking.

Backend must check role.

---

## 9. Permission Rule

Permission is not the same as role.

A user can be `clinic_staff`, but still cannot view a pet unless the clinic has approved authorization.

Important backend helper:

```txt
CanOwnerAccessPet(ownerUserId, petId)
CanClinicAccessPet(clinicId, petId, permission)
```

Required permissions:

```txt
view_profile
view_history
create_visit
```

---

## 10. Pet Ownership Rule

Owner can access only their own pets.

Every owner pet endpoint must check:

```txt
current user
→ owner profile
→ pet.owner_id == owner_profile.id
```

If not owner:

```txt
403 NOT_PET_OWNER
```

---

## 11. Clinic Access Rule

Clinic can view full pet data only after owner approval.

Approved authorization must satisfy:

```txt
pet_id matches
clinic_id matches
status = approved
expires_at is null or in the future
permissions include required permission
```

If no valid authorization:

```txt
403 CLINIC_ACCESS_NOT_APPROVED
```

---

## 12. QR Security Rule

QR code must contain only a temporary secure token.

Good QR payload:

```json
{
  "type": "petnexus_qr_session",
  "token": "secure_random_token"
}
```

Bad QR payload:

```json
{
  "petName": "Milo",
  "allergyNote": "Chicken allergy",
  "ownerPhone": "0812345678"
}
```

Never put full pet data in QR.

Never put medical data in QR.

Never put owner private data in QR.

---

## 13. QR Scan Rule

When clinic scans QR, backend should return pet preview only.

Allowed pet preview:

```txt
petId
petNexusId
petName
species
breedName
photoUrl
```

Do not return before approval:

```txt
allergyNote
chronicDiseaseNote
full timeline
diagnosis
treatment
medication
owner address
owner emergency contact
```

Full pet record requires approved authorization.

---

## 14. Authorization Status Rule

Allowed statuses:

```txt
pending
approved
rejected
revoked
```

Allowed transitions:

```txt
pending → approved
pending → rejected
approved → revoked
```

Invalid transitions:

```txt
rejected → approved
revoked → approved
approved → pending
rejected → pending
revoked → pending
```

If invalid:

```txt
409 INVALID_AUTHORIZATION_STATUS_TRANSITION
```

---

## 15. Visit Rule

Clinic-created visits must be marked:

```txt
clinic_verified
```

Clinic can create visit only if:

```txt
clinic has approved authorization
authorization includes create_visit
authorization is not expired
authorization is not revoked
```

Owner can view clinic-created visits.

Owner should not edit clinic-created visits.

Clinic should not edit another clinic’s visit.

---

## 16. Audit Log Rule

Sensitive actions should create audit logs.

Important actions:

```txt
clinic_scanned_qr
clinic_requested_access
owner_approved_access
owner_rejected_access
owner_revoked_access
clinic_viewed_pet
clinic_created_visit
```

Audit log should include:

```txt
actor_user_id
action
target_type
target_id
metadata
created_at
```

Do not expose audit log API in MVP unless explicitly asked.

---

## 17. Notification Rule

Create notification when:

```txt
clinic requests access
clinic creates visit
```

For MVP, notification is only database row.

Do not add push notifications yet.

Do not add realtime notifications yet.

---

## 18. Database Rule

Use PostgreSQL.

Use UUID primary keys.

Use snake_case column names.

Use plural table names.

Use migrations.

Do not rely only on GORM AutoMigrate for production-like schema.

Use migration files for database structure.

---

## 19. GORM Rule

Use GORM for MVP.

Keep models simple.

Use pointer types for nullable fields.

Examples:

```txt
*uuid.UUID
*time.Time
*float64
```

Do not create overcomplicated generic repositories.

Do not hide all SQL behavior behind unnecessary abstractions.

---

## 20. Validation Rule

Validate in service layer and handler layer.

Examples:

```txt
email is required
password is required
role must be valid
pet name is required
species must be dog or cat
birth date cannot be in future
weight must not be negative
QR token must not be expired
permissions must be allowed values only
```

Do not trust raw client input.

---

## 21. Error Naming Rule

Use clear error codes.

Examples:

```txt
INVALID_REQUEST
VALIDATION_ERROR
UNAUTHORIZED
FORBIDDEN_ROLE
NOT_PET_OWNER
PET_NOT_FOUND
OWNER_PROFILE_NOT_FOUND
CLINIC_ACCESS_NOT_APPROVED
MISSING_PERMISSION
QR_SESSION_NOT_FOUND
QR_SESSION_EXPIRED
ACCESS_REQUEST_ALREADY_PENDING
ACCESS_ALREADY_APPROVED
AUTHORIZATION_NOT_FOUND
INVALID_AUTHORIZATION_STATUS_TRANSITION
INTERNAL_SERVER_ERROR
```

Do not return vague errors like:

```txt
ERROR
FAILED
BAD
```

---

## 22. Sprint Scope Rule

Codex must only implement the current sprint.

If working on Sprint 1, do not implement Sprint 2.

If working on Sprint 2, do not implement QR yet.

If working on auth, do not implement visit yet.

Build in order:

```txt
1. Health route
2. Database connection
3. Auth
4. Owner profile
5. Breed
6. Pet
7. QR session
8. Access request
9. Authorization decision
10. Clinic patient access
11. Clinic visit
12. Timeline
13. Notifications
14. Audit logs
15. Demo polish
```

---

## 23. No Overbuilding Rule

Do not build these in MVP unless explicitly asked:

```txt
AI diagnosis
OCR
Chat
Payment
Subscription
Multi-branch clinic
Complex vaccine engine
Advanced appointment system
PDF export
Realtime notification
Push notification
Admin approval dashboard
Full clinic verification system
File document management
Grooming integration
Pet hotel integration
Microservices
GraphQL
Redis
Kubernetes
Message queues
Event-driven architecture
```

---

## 24. Security Rule

Never log:

```txt
plain password
password hash
JWT secret
JWT token
private owner data unnecessarily
medical details in QR
```

Use bcrypt for password hashing.

Use environment variables for secrets.

Do not hardcode secrets.

---

## 25. Environment Rule

Use `.env.example`.

Required variables:

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

Do not commit real `.env`.

---

## 26. README Rule

Every major sprint should update README with:

```txt
what was implemented
how to run
how to test
new endpoints
required env variables
known limitations
next step
```

Do not leave README outdated.

---

## 27. Comment Rule

Use comments to explain beginner-important logic.

Good comment:

```txt
// Clinic can only see full pet data after owner has approved authorization.
```

Bad comment:

```txt
// create variable
// call function
```

Explain why, not obvious what.

---

## 28. Testing Rule

At minimum, manually test with Postman or Thunder Client.

Each sprint should include a small test checklist.

Do not claim a feature works without test steps.

---

## 29. Codex Output Rule

After making changes, Codex should summarize:

```txt
Files created
Files changed
How to run
How to test
What was intentionally not implemented
Recommended next step
```

This helps the developer learn and review.

---

## 30. Final Rule

Always protect the MVP.

The first complete backend demo should prove only this:

```txt
Owner registers
Owner creates pet
Owner creates QR
Clinic scans QR
Clinic requests access
Owner approves
Clinic creates clinic_verified visit
Owner sees timeline
```

Everything else is secondary.
