# PetNexus Backend Roadmap

## 1. Backend Goal

The backend is the core system of PetNexus.

It is responsible for:

* Authentication
* Role-based access control
* Pet data management
* QR session generation and validation
* Clinic access request
* Owner approval / rejection / revocation
* Clinic verified visit creation
* Pet health timeline
* Notifications
* Audit logs

The backend must be the source of truth.

Frontend apps must not decide sensitive permission rules by themselves.

---

## 2. Current Backend Stack

Backend stack for MVP:

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

Optional later:

```txt
Docker
Cloud object storage
Redis
Background jobs
Email / push notification service
```

Do not add optional tools until the core MVP flow works.

---

## 3. Core MVP Flow

The backend must support this flow first:

```txt
Owner registers / logs in
→ Owner creates profile
→ Owner creates pet
→ Owner generates QR session
→ Clinic scans QR token
→ Backend validates QR token
→ Clinic sees pet preview only
→ Clinic requests access
→ Owner approves access
→ Backend creates authorization
→ Clinic views full pet record
→ Clinic creates clinic_verified visit
→ Owner sees updated timeline
```

This is the main success path.

If this flow works clearly, the MVP is strong.

---

## 4. Backend Development Principles

### 4.1 Keep It Simple

This is a beginner-friendly backend project.

Avoid:

```txt
Microservices
GraphQL
Kubernetes
Event-driven architecture
Message queues
Complex permission engine
Over-abstracted clean architecture
```

Use simple layered architecture first.

---

### 4.2 Use Layered Architecture

```txt
handler      = receives HTTP request and returns response
service      = business logic and permission checks
repository   = database access
model        = database entity
middleware   = auth and role checks
routes       = route registration
config       = environment config
database     = PostgreSQL connection
```

Rule:

```txt
Handler should not directly talk to database.
Handler → Service → Repository → Database
```

---

### 4.3 Backend Owns Permission Logic

Frontend can hide buttons, but backend must still check permission.

Important examples:

```txt
Clinic cannot view full pet record without approved authorization.
Clinic cannot create visit without create_visit permission.
Owner cannot edit clinic_verified visit.
QR token must not expose full pet data.
```

---

### 4.4 Build by Vertical Slices

Do not build all models first and then all APIs later.

Build feature by feature:

```txt
Auth works
→ Pet works
→ QR works
→ Authorization works
→ Visit works
→ Timeline works
```

Each slice should be testable.

---

## 5. Recommended Backend Folder Structure

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
│   ├── 001_create_users.sql
│   ├── 002_create_owner_profiles.sql
│   ├── 003_create_breeds.sql
│   ├── 004_create_pets.sql
│   ├── 005_create_clinics.sql
│   ├── 006_create_clinic_staff.sql
│   ├── 007_create_qr_sessions.sql
│   ├── 008_create_authorizations.sql
│   ├── 009_create_visits.sql
│   ├── 010_create_notifications.sql
│   └── 011_create_audit_logs.sql
│
├── docs/
│   ├── context.md
│   ├── backend-roadmap.md
│   ├── api-plan.md
│   └── database-plan.md
│
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

---

## 6. Sprint Plan

## Sprint 1: Backend Foundation

### Goal

Make the backend project runnable.

### Build

```txt
Go module setup
Gin server setup
Config loading
Database connection placeholder
Route registration
Health check endpoint
Basic response format
Basic error format
```

### Endpoint

```txt
GET /health
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

### Definition of Done

```txt
go run ./cmd/api works
GET /health returns status ok
.env.example exists
README explains how to run backend
Folder structure is ready
```

---

## Sprint 2: Database Foundation

### Goal

Connect backend to PostgreSQL and prepare core tables.

### Build

```txt
PostgreSQL connection
GORM setup
Migration system
Base models
Breed seed data
Database README
```

### Tables

```txt
users
owner_profiles
breeds
pets
clinics
clinic_staff
qr_sessions
authorizations
visits
notifications
audit_logs
```

### Definition of Done

```txt
Backend connects to PostgreSQL
Migrations can run
Breeds are seeded
Database structure matches MVP entities
```

---

## Sprint 3: Auth Foundation

### Goal

Users can register and login.

### Build

```txt
Register API
Login API
Password hashing with bcrypt
JWT generation
JWT middleware
GET /api/me
Role-based middleware
```

### Roles

```txt
owner
clinic_staff
admin
```

### Endpoints

```txt
POST /api/auth/register
POST /api/auth/login
GET  /api/me
```

### Definition of Done

```txt
Owner can register
Clinic staff can register
User can login
Login returns JWT token
Protected route requires token
Backend can read current user from token
```

---

## Sprint 4: Owner Profile

### Goal

Owner can create and view profile.

### Build

```txt
Create owner profile
Get owner profile
Update owner profile
Validate owner role
```

### Endpoints

```txt
POST  /api/owner/profile
GET   /api/owner/profile
PATCH /api/owner/profile
```

### Definition of Done

```txt
Only owner role can create owner profile
One user can have only one owner profile
Owner profile is linked to user id
```

---

## Sprint 5: Breed + Pet Core

### Goal

Owner can create digital pet.

### Build

```txt
Get breeds
Create pet
Get pet list
Get pet detail
Update pet basic info
Generate PetNexus ID
```

### Endpoints

```txt
GET   /api/breeds
GET   /api/breeds?species=dog
GET   /api/breeds?species=cat

POST  /api/pets
GET   /api/pets
GET   /api/pets/:petId
PATCH /api/pets/:petId
GET   /api/pets/:petId/passport
```

### Definition of Done

```txt
Owner can create dog or cat
Owner can select breed
Pet gets unique PetNexus ID
Owner can only see own pets
Pet passport returns important info
```

---

## Sprint 6: QR Session

### Goal

Owner can create QR token and clinic can scan it.

### Build

```txt
Create QR session
Generate secure random token
Set expiration time
Validate QR token
Return pet preview only
Mark usedAt if needed
Add audit log for clinic scan
```

### Endpoints

```txt
POST /api/pets/:petId/qr-session
POST /api/clinic/scan-qr
```

### Important Rules

```txt
QR must contain token only.
QR must not contain full pet data.
Expired token must be rejected.
Invalid token must be rejected.
Clinic sees only pet preview before approval.
```

### Pet Preview Should Include

```txt
petId
petNexusId
petName
species
breed
photoUrl
ownerNickname
```

Do not include:

```txt
allergyNote
chronicDiseaseNote
full timeline
diagnosis
treatment
medication
owner full address
```

### Definition of Done

```txt
Owner can create QR session
QR token can be scanned by clinic
Expired token is rejected
Clinic receives preview only
Audit log records clinic_scanned_qr
```

---

## Sprint 7: Access Request

### Goal

Clinic can request access to pet record.

### Build

```txt
Create access request
Prevent duplicate pending request
Notify owner
Add audit log
```

### Endpoints

```txt
POST /api/clinic/access-request
GET  /api/owner/access-requests
```

### Access Request Status

```txt
pending
approved
rejected
revoked
```

### Definition of Done

```txt
Clinic can request access after scanning QR
Owner can see pending request
Duplicate pending request is blocked
Notification is created for owner
Audit log records clinic_requested_access
```

---

## Sprint 8: Owner Authorization Decision

### Goal

Owner can approve, reject, or revoke clinic access.

### Build

```txt
Approve authorization
Reject authorization
Revoke authorization
Permission list
Authorization expiration placeholder
Notify clinic later if needed
Audit log
```

### Endpoints

```txt
POST /api/authorizations/:id/approve
POST /api/authorizations/:id/reject
POST /api/authorizations/:id/revoke
```

### Permissions

```txt
view_profile
view_history
create_visit
```

### Definition of Done

```txt
Only pet owner can approve request
Only pet owner can reject request
Only pet owner can revoke approved access
Approved clinic can view full pet record
Rejected clinic cannot view full pet record
Revoked clinic cannot view full pet record
Audit logs are created
```

---

## Sprint 9: Clinic Patient Access

### Goal

Clinic can view approved pet records.

### Build

```txt
Get clinic patients
Get pet record after authorization
Get pet timeline after authorization
Authorization check middleware/helper
```

### Endpoints

```txt
GET /api/clinic/patients
GET /api/clinic/pets/:petId
GET /api/clinic/pets/:petId/timeline
```

### Definition of Done

```txt
Clinic sees only approved patients
Clinic cannot access pets without authorization
Clinic can view allergy/chronic note only after approval
Clinic can view timeline only after approval
```

---

## Sprint 10: Clinic Visit

### Goal

Clinic can create verified visit record.

### Build

```txt
Create visit
Check clinic authorization
Check create_visit permission
Mark visit as clinic_verified
Notify owner
Add audit log
Update timeline
```

### Endpoint

```txt
POST /api/clinic/pets/:petId/visits
```

### Visit Fields

```txt
visitDate
chiefComplaint
diagnosis
treatment
medication
followUpDate
note
```

### Important Rules

```txt
Clinic-created visit is clinic_verified.
Owner can view but cannot edit clinic-created visit.
Clinic cannot edit another clinic’s visit.
Clinic cannot create visit without approved authorization.
```

### Definition of Done

```txt
Approved clinic can create visit
Unapproved clinic cannot create visit
Visit appears in pet timeline
Visit has clinic_verified status
Owner receives notification
Audit log records clinic_created_visit
```

---

## Sprint 11: Timeline

### Goal

Owner and authorized clinic can view pet timeline.

### Build

```txt
Timeline API
Sort visits by date
Show clinic verified badge
Separate owner-added and clinic-created data later
```

### Endpoints

```txt
GET /api/pets/:petId/timeline
GET /api/clinic/pets/:petId/timeline
```

### Definition of Done

```txt
Owner can see own pet timeline
Authorized clinic can see pet timeline
Unauthorized clinic cannot see timeline
Timeline shows verified visit records
```

---

## Sprint 12: Notifications

### Goal

Owner can see important system notifications.

### Build

```txt
Create notification when clinic requests access
Create notification when clinic creates visit
Get notifications
Mark notification as read
```

### Endpoints

```txt
GET   /api/notifications
PATCH /api/notifications/:id/read
```

### Definition of Done

```txt
Owner sees access request notification
Owner sees visit created notification
User can mark notification as read
```

---

## Sprint 13: Audit Log

### Goal

Sensitive actions are tracked.

### Build

```txt
Create audit log helper
Log important actions
Admin read endpoint later
```

### Actions

```txt
clinic_scanned_qr
clinic_requested_access
owner_approved_access
owner_rejected_access
owner_revoked_access
clinic_viewed_pet
clinic_created_visit
```

### Definition of Done

```txt
Sensitive actions are stored
Each log has actorUserId
Each log has targetType and targetId
Metadata can store extra JSON
```

---

## Sprint 14: Demo Polish

### Goal

Make backend stable enough for frontend and demo.

### Build

```txt
Consistent error responses
Input validation
Pagination for list APIs
Basic logging
Seed demo owner
Seed demo clinic
Seed demo pet
Seed demo visit
README API examples
Postman collection or Thunder Client collection
```

### Definition of Done

```txt
Demo data is available
Frontend can consume APIs
Errors are readable
Core flow can be tested from beginning to end
```

---

## 7. Backend Priority Summary

Build in this order:

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

Do not skip QR and authorization rules.

That is the heart of PetNexus.

---

## 8. MVP Cut List

Do not build these in backend v1:

```txt
AI diagnosis
OCR
Chat
Payment
Multi-branch clinic
Complex vaccine engine
Advanced appointment system
Export PDF
Realtime notification
Push notification
Admin approval dashboard
Full clinic verification system
File document management
Grooming / pet hotel integration
```

Add these only after the core flow is done.

---

## 9. Backend Risk Checklist

### Risk 1: Backend Becomes Too Big

Solution:

```txt
Build one flow at a time.
Do not implement future features early.
Keep service methods focused.
```

---

### Risk 2: Permission Bugs

Solution:

```txt
Every clinic pet access must check authorization.
Never trust frontend.
Write helper function:
CanClinicAccessPet(clinicId, petId, permission)
```

---

### Risk 3: QR Token Leaks Data

Solution:

```txt
QR contains token only.
Token expires.
Scan endpoint returns preview only.
Full data requires approved authorization.
```

---

### Risk 4: Confusing Data Ownership

Solution:

```txt
Owner owns pet.
Clinic owns clinic-created visit.
Owner can view clinic visit.
Owner cannot edit clinic_verified visit.
Clinic cannot edit other clinic’s visit.
```

---

### Risk 5: New Backend Developer Gets Lost

Solution:

```txt
Keep folders simple.
Comment important files.
Write README.
Create API examples.
Avoid clever abstractions.
```

---

## 10. First Backend Milestone

The first milestone is intentionally small.

Command:

```bash
go run ./cmd/api
```

Test:

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

After this works, continue to database connection and auth.

---

## 11. Final Backend MVP Definition

Backend MVP is complete when this scenario works:

```txt
Owner registers
Owner creates profile
Owner creates pet
Owner creates QR session
Clinic scans QR
Clinic sees pet preview
Clinic requests access
Owner approves
Clinic views full pet record
Clinic creates visit
Visit is marked clinic_verified
Owner sees visit in timeline
Owner receives notification
Audit logs are created
```

If this works, PetNexus has a strong MVP backend.
