# PetNexus Project Context

## 1. Project Overview

**Project Name:** PetNexus

**Product Type:**
Digital Pet Passport / Digital Pet Identity Platform

**Target Users:**

1. Pet owners
2. Small-to-medium veterinary clinics
3. Clinic staff / veterinarians

**Core Idea:**
One pet should have one portable, verified, owner-controlled health identity.

PetNexus is not just a pet diary. It solves the problem that pet health data is scattered across paper vaccine books, LINE chats, clinic systems, PDF files, and owner memory.

The goal is to make pet health information portable, verifiable, and controlled by the pet owner.

---

## 2. Core Problem

Pet owners and clinics face several problems:

1. Pet health records are fragmented.
2. Owners often forget or lose vaccine / visit history.
3. Clinics waste time asking the same history repeatedly.
4. Vaccine records are difficult to verify when used with clinics, pet hotels, grooming services, or other pet services.
5. When pets change clinics or are referred to another hospital, the medical history does not travel with the pet.
6. Owners do not have clear control over who can access their pet’s data.

---

## 3. Core Product Value

PetNexus provides:

### 1. Portable Record

Pet health history follows the pet, not the clinic.

### 2. Clinic Verified Data

Important records such as clinic visits and vaccinations can be verified by clinics.

### 3. Owner-Controlled Authorization

Owners decide which clinic can access which pet data.

### 4. Fast QR Sharing

Owners can share pet information quickly using a QR code.

### 5. Clear Medical Timeline

Owners and authorized clinics can view a pet’s health history in chronological order.

---

## 4. MVP Goal

The MVP should focus only on the core flow:

Owner creates pet
→ Owner shares QR
→ Clinic scans QR
→ Clinic requests access
→ Owner approves access
→ Clinic creates verified visit
→ Owner sees updated timeline

This flow is the heart of PetNexus.

Do not overbuild the first version.

---

## 5. MVP Scope

### Build in MVP

* User authentication
* Owner profile setup
* Digital Pet creation
* Pet Passport
* QR sharing
* QR token validation
* Clinic access request
* Owner approval / rejection / revocation
* Clinic verified visit record
* Pet health timeline
* Basic notification system
* Basic audit log

### Do Not Build in MVP

* AI diagnosis
* OCR vaccine book scanning
* Complex medical analytics
* Payment system
* Chat system
* Multi-branch clinic system
* Pet hotel / grooming integration
* Advanced calendar system
* Full vaccine schedule engine
* PDF export
* Avatar / virtual pet animation
* Machine learning features
* Mobile offline mode

The MVP must stay realistic and focused.

---

## 6. Target Users

## 6.1 Pet Owner

Pet owners use the mobile app.

Main goals:

* Register / login
* Set up owner profile
* Create digital pet
* View pet passport
* Share QR code with clinic
* Approve or reject clinic access
* View health timeline
* Receive notification when clinic adds new record

---

## 6.2 Veterinary Clinic

Clinics use the web dashboard.

Main goals:

* Register / login as clinic staff
* Scan owner’s QR code
* View pet preview only
* Request access to pet record
* View full pet record after owner approval
* Create visit record
* Add clinic verified data to pet timeline

---

## 7. Product Platforms

PetNexus has 3 main technical parts:

1. **Mobile App**

   * For pet owners
   * Built with Flutter

2. **Clinic Web Dashboard**

   * For clinic staff
   * Built with Next.js

3. **Backend API**

   * Handles auth, database, QR, authorization, visits, and timeline
   * Built with Go

---

## 8. Recommended Tech Stack

## 8.1 Mobile App

**Stack:**

* Flutter
* Dart
* Riverpod
* GoRouter
* Dio
* flutter_secure_storage
* image_picker
* qr_flutter

**Responsibility:**

* Owner login / register
* Owner profile setup
* Pet onboarding
* Pet passport
* QR sharing
* Authorization approval
* Notification viewing
* Timeline viewing

---

## 8.2 Clinic Web Dashboard

**Stack:**

* Next.js
* TypeScript
* Tailwind CSS
* TanStack Query
* React Hook Form
* Zod
* QR scanner library

**Responsibility:**

* Clinic login / register
* Clinic dashboard
* QR scan page
* Access request
* Patient list
* Pet record view
* Create visit form
* Clinic timeline view

---

## 8.3 Backend

**Stack:**

* Go
* Gin
* PostgreSQL
* GORM
* JWT
* bcrypt
* godotenv
* golang-migrate
* Docker later if needed

**Responsibility:**

* Authentication
* Role-based authorization
* Owner profile API
* Pet API
* Breed API
* QR session API
* Clinic access request API
* Authorization approval API
* Visit creation API
* Timeline API
* Notification API
* Audit log

---

## 9. Current Backend Decision

The backend will not use Supabase for the current plan.

Current backend direction:

Go + Gin + PostgreSQL + GORM + JWT

Reason:

* The project owner wants to learn backend seriously.
* Go is a good learning path for building clean APIs.
* PostgreSQL fits the relational data model of PetNexus.
* The system has clear entities and relationships.
* The backend can grow later without being locked into a BaaS platform.

Important:

Keep the first backend version simple.
Do not use microservices.
Do not overuse complex clean architecture.
Do not build advanced permission engines yet.

---

## 10. High-Level Architecture

```txt
Flutter Owner App
        |
        | HTTP REST API
        v
Go Backend API
        |
        v
PostgreSQL Database


Next.js Clinic Web
        |
        | HTTP REST API
        v
Go Backend API
        |
        v
PostgreSQL Database
```

The backend is the central source of truth.

Frontend apps must not directly access the database.

---

## 11. Backend Architecture

Use beginner-friendly layered architecture.

```txt
handler      = receives HTTP request and returns response
service      = business logic
repository   = database access
model        = database entity
middleware   = auth and role checks
routes       = route registration
config       = environment configuration
database     = PostgreSQL connection
```

Recommended backend structure:

```txt
petnexus-backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   ├── handlers/
│   └── routes/
│
├── migrations/
├── docs/
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

---

## 12. Core Data Models

## 12.1 User

Represents login account.

Fields:

* id
* email
* phone
* password_hash
* role
* created_at
* updated_at

Roles:

* owner
* clinic_staff
* admin

---

## 12.2 OwnerProfile

Represents pet owner profile.

Fields:

* id
* user_id
* full_name
* nickname
* phone
* email
* address
* emergency_contact_name
* emergency_contact_phone
* created_at
* updated_at

Relationship:

* One user has one owner profile.
* One owner profile has many pets.

---

## 12.3 Pet

Represents one digital pet.

Fields:

* id
* owner_id
* pet_nexus_id
* name
* species
* breed_id
* gender
* birth_date
* weight_kg
* color_note
* allergy_note
* chronic_disease_note
* photo_url
* created_at
* updated_at

Species:

* dog
* cat

Relationship:

* One owner has many pets.
* One pet has many QR sessions.
* One pet has many authorizations.
* One pet has many visits.

---

## 12.4 Breed

Represents pet breed seed data.

Fields:

* id
* species
* name_th
* name_en
* image_url

Initial dog breeds:

* Golden Retriever
* Pomeranian
* French Bulldog
* Poodle

Initial cat breeds:

* British Shorthair
* Persian
* Scottish Fold
* Siamese

Also include:

* Not sure / Select later

---

## 12.5 Clinic

Represents veterinary clinic.

Fields:

* id
* name
* address
* phone
* email
* verified_status
* created_at
* updated_at

Verified status:

* pending
* verified
* rejected

For MVP, clinic verification can be simple and manual.

---

## 12.6 ClinicStaff

Represents staff account inside a clinic.

Fields:

* id
* clinic_id
* user_id
* role
* license_no
* created_at
* updated_at

Roles:

* vet
* assistant
* clinic_admin

Relationship:

* One clinic has many clinic staff.

---

## 12.7 QRSession

Represents temporary QR sharing session.

Fields:

* id
* pet_id
* owner_id
* token
* purpose
* expires_at
* used_at
* created_at

Purpose:

* clinic_checkin
* passport_share
* emergency

Important rules:

* QR must contain only a temporary secure token.
* QR must not contain full pet data.
* Token must expire.
* Backend must validate token before showing any pet data.

---

## 12.8 Authorization

Controls clinic access to pet data.

Fields:

* id
* pet_id
* clinic_id
* owner_id
* status
* permissions
* expires_at
* created_at
* updated_at

Status:

* pending
* approved
* rejected
* revoked

Permissions:

* view_profile
* view_history
* create_visit

Important rules:

* Clinic cannot view full pet data without approved authorization.
* Owner can approve, reject, or revoke access.
* Authorization can expire later.

---

## 12.9 Visit

Represents one clinic visit or medical record.

Fields:

* id
* pet_id
* clinic_id
* vet_id
* visit_date
* chief_complaint
* diagnosis
* treatment
* medication
* follow_up_date
* note
* verification_status
* created_at
* updated_at

Verification status:

* clinic_verified

Important rules:

* Visit created by clinic is marked as clinic_verified.
* Owner can view clinic-created visit.
* Owner should not edit clinic-created visit.
* Clinic should not edit another clinic’s visit.

---

## 12.10 Notification

Represents notification to user.

Fields:

* id
* user_id
* title
* message
* type
* is_read
* created_at

Types:

* access_request
* access_approved
* access_rejected
* visit_created
* appointment_reminder
* vaccine_due

---

## 12.11 AuditLog

Tracks important sensitive actions.

Fields:

* id
* actor_user_id
* action
* target_type
* target_id
* metadata
* created_at

Example actions:

* clinic_scanned_qr
* clinic_requested_access
* owner_approved_access
* owner_rejected_access
* owner_revoked_access
* clinic_viewed_pet
* clinic_created_visit

---

## 13. Core API Plan

## 13.1 Health

```txt
GET /health
```

Returns:

```json
{
  "status": "ok",
  "service": "petnexus-backend"
}
```

---

## 13.2 Auth

```txt
POST /api/auth/register
POST /api/auth/login
GET  /api/me
```

---

## 13.3 Owner Profile

```txt
POST /api/owner/profile
GET  /api/owner/profile
PATCH /api/owner/profile
```

---

## 13.4 Breeds

```txt
GET /api/breeds
GET /api/breeds?species=dog
GET /api/breeds?species=cat
```

---

## 13.5 Pets

```txt
POST  /api/pets
GET   /api/pets
GET   /api/pets/:petId
PATCH /api/pets/:petId
GET   /api/pets/:petId/passport
GET   /api/pets/:petId/timeline
```

---

## 13.6 QR Session

```txt
POST /api/pets/:petId/qr-session
POST /api/clinic/scan-qr
```

Important:

Clinic scan QR should return only pet preview first.

---

## 13.7 Authorization

```txt
POST /api/clinic/access-request
GET  /api/owner/access-requests
POST /api/authorizations/:id/approve
POST /api/authorizations/:id/reject
POST /api/authorizations/:id/revoke
```

---

## 13.8 Clinic

```txt
POST /api/clinic/register
POST /api/clinic/login
GET  /api/clinic/me
GET  /api/clinic/dashboard
GET  /api/clinic/patients
GET  /api/clinic/pets/:petId
GET  /api/clinic/pets/:petId/timeline
```

---

## 13.9 Visit

```txt
POST /api/clinic/pets/:petId/visits
GET  /api/pets/:petId/timeline
```

---

## 13.10 Notification

```txt
GET   /api/notifications
PATCH /api/notifications/:id/read
```

---

## 14. Permission Rules

## 14.1 Owner Can

* Create pet
* Edit pet basic info
* View own pet timeline
* Create QR session
* View pending clinic access requests
* Approve clinic access
* Reject clinic access
* Revoke clinic access
* View notifications

---

## 14.2 Clinic Can

* Login as clinic staff
* Scan QR
* See pet preview after valid QR token
* Request access to pet record
* View full pet record only after owner approval
* Create visit only if authorization includes create_visit
* View visits it is allowed to access

---

## 14.3 Clinic Cannot

* Delete pet
* Edit owner profile
* View full pet record without authorization
* Edit another clinic’s visit
* Create visit without approved permission
* Access expired or revoked pet records

---

## 15. UI Direction

## 15.1 Owner Mobile App UI

Style:

* Cute but clean
* Friendly
* Pastel cream background
* Teal primary color
* Soft pink / yellow accents
* Rounded cards
* Big tap targets
* Not crowded

Important:

* Onboarding should feel easy.
* Avoid making it feel like a hospital form.
* Separate owner-added data and clinic-verified data visually.
* Pet Passport should be simple and clear.

---

## 15.2 Clinic Web UI

Style:

* Clean
* Professional
* Readable
* Not too cute
* Left sidebar navigation
* White cards
* Teal / dark navy color
* Plenty of spacing

Sidebar pages:

* Dashboard
* QR Pet Data
* Patients
* Calendar
* Medical Records
* Reports
* Settings

---

## 16. Owner Mobile Main Pages

```txt
Welcome
Login
Register
Owner Profile Setup
Pet Type Selection
Breed Selection
Create Pet
Home
Pet Passport
Share QR
Health Timeline
Authorization Requests
Notifications
Settings
```

---

## 17. Clinic Web Main Pages

```txt
Clinic Login
Clinic Register
Dashboard
QR Scan / Check-in
Access Requests
Patient List
Pet Record View
Create Visit
Calendar
Settings
```

---

## 18. MVP Roadmap

## Sprint 1: Backend Foundation

* Set up Go project structure
* Set up Gin server
* Add GET /health
* Add config loading
* Add database connection placeholder
* Add error response format
* Add route registration structure

Goal:

Backend can run.

---

## Sprint 2: Auth Foundation

* Register
* Login
* Password hashing with bcrypt
* JWT generation
* JWT auth middleware
* Role-based middleware
* Owner / clinic_staff roles

Goal:

Users can authenticate and backend knows their role.

---

## Sprint 3: Pet Owner Core

* Owner profile
* Breed seed data
* Create pet
* Get pet list
* Get pet passport
* Upload pet photo later if needed

Goal:

Owner can create digital pet and view passport.

---

## Sprint 4: QR + Access Request

* Create QR session
* Generate temporary token
* Validate QR token
* Clinic sees pet preview only
* Clinic requests access
* Owner sees pending request

Goal:

Clinic can request access without seeing full record immediately.

---

## Sprint 5: Authorization

* Owner approves access
* Owner rejects access
* Owner revokes access
* Backend checks approved authorization before clinic views full record
* Add audit log for sensitive actions

Goal:

Owner controls clinic access.

---

## Sprint 6: Clinic Visit + Timeline

* Clinic views approved pet record
* Clinic creates visit
* Visit is marked clinic_verified
* Owner sees updated timeline
* Notification created when visit is added

Goal:

The full core MVP flow works.

---

## Sprint 7: Polish for Demo

* Loading states
* Empty states
* Better error handling
* Demo seed data
* Basic activity log
* Basic calendar placeholder
* Responsive clinic web
* Demo script

Goal:

Project is presentable.

---

## 19. Development Rules

1. Keep MVP small.
2. Build backend foundation first.
3. Do not add AI in v1.
4. Do not add OCR in v1.
5. Do not add advanced analytics in v1.
6. Do not allow clinic to view full pet data from QR alone.
7. QR must contain token only.
8. Owner approval is required before full record access.
9. Clinic-created visits must be marked clinic_verified.
10. Sensitive actions should be logged in audit_logs.
11. Backend must own the business logic.
12. Frontend should not decide permission by itself.
13. Do not over-engineer architecture.
14. Prefer readable code over clever code.
15. Build the core flow before polishing UI.

---

## 20. Current Development Stage

Current stage:

Planning and project structure.

The backend scaffold should be created first.

Do not implement full features yet.

First backend goal:

```txt
go run ./cmd/api
```

and

```txt
GET /health
```

should return:

```json
{
  "status": "ok",
  "service": "petnexus-backend"
}
```

After the backend skeleton is ready, continue with authentication.

---

## 21. Team / Learning Context

This project is also a learning project.

The developer is new to backend development and wants to learn by building a real product.

Backend code should be:

* Beginner-friendly
* Well-structured
* Easy to read
* Easy to extend
* Not overly abstract
* Not enterprise-heavy

Avoid:

* Microservices
* Kubernetes
* GraphQL
* Redis
* Message queues
* Complex permission engines
* Advanced clean architecture
* Overly generic repository patterns

Use simple layered architecture first.

---

## 22. Final Product Vision

PetNexus should become a trusted digital identity layer for pets.

The long-term vision is:

* Every pet has a portable digital health identity.
* Owners control access.
* Clinics can verify medical records.
* Pet health history becomes easier to transfer between clinics.
* QR sharing makes check-in and verification faster.
* Pet owners no longer rely only on paper books, chat messages, or memory.

But for MVP, focus only on proving one clear thing:

A pet owner can create a pet passport, share access with a clinic, approve the clinic, and receive a clinic-verified visit record in the pet timeline.
