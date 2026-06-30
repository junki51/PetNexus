# PetNexus Database Plan

## 1. Database Overview

PetNexus uses **PostgreSQL** as the main database.

The database must support the core MVP flow:

```txt
Owner creates pet
→ Owner shares QR
→ Clinic scans QR
→ Clinic requests access
→ Owner approves access
→ Clinic creates verified visit
→ Owner sees updated timeline
```

The database design should focus on:

* Pet ownership
* Clinic access control
* Temporary QR sharing
* Clinic verified medical records
* Timeline history
* Notifications
* Audit logs

Do not overbuild the database in v1.

---

## 2. Database Goals

The database should make these rules easy to enforce:

1. One owner can have many pets.
2. One pet belongs to one owner.
3. One clinic can have many staff.
4. A clinic cannot view full pet data without approved authorization.
5. QR code must store token only, not full pet data.
6. Clinic-created visits must be marked as `clinic_verified`.
7. Owner can view clinic-created visits but should not edit them.
8. Sensitive actions should be recorded in audit logs.

---

## 3. Recommended Database Stack

```txt
Database: PostgreSQL
ORM: GORM
Migration Tool: golang-migrate
ID Type: UUID
Time Fields: created_at, updated_at
Soft Delete: optional later, not required in MVP
```

Use UUID for primary keys because the app has mobile, web, and backend clients. UUID is safer than exposing incremental IDs.

---

## 4. Naming Convention

Use snake_case for database columns.

Examples:

```txt
pet_nexus_id
owner_id
clinic_id
created_at
updated_at
verification_status
```

Use plural table names:

```txt
users
owner_profiles
pets
breeds
clinics
clinic_staff
qr_sessions
authorizations
visits
notifications
audit_logs
```

---

## 5. Migration Order

Create migrations in this order:

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

Reason:

* Enums must exist before tables use them.
* Users should exist before owner profiles and clinic staff.
* Pets need owner profiles and breeds.
* QR sessions, authorizations, and visits need pets.
* Notifications and audit logs can be added after core entities.

---

## 6. Enum Types

## 6.1 user_role

```sql
CREATE TYPE user_role AS ENUM (
  'owner',
  'clinic_staff',
  'admin'
);
```

---

## 6.2 pet_species

```sql
CREATE TYPE pet_species AS ENUM (
  'dog',
  'cat'
);
```

---

## 6.3 pet_gender

```sql
CREATE TYPE pet_gender AS ENUM (
  'male',
  'female',
  'unknown'
);
```

---

## 6.4 clinic_verified_status

```sql
CREATE TYPE clinic_verified_status AS ENUM (
  'pending',
  'verified',
  'rejected'
);
```

---

## 6.5 clinic_staff_role

```sql
CREATE TYPE clinic_staff_role AS ENUM (
  'vet',
  'assistant',
  'clinic_admin'
);
```

---

## 6.6 qr_purpose

```sql
CREATE TYPE qr_purpose AS ENUM (
  'clinic_checkin',
  'passport_share',
  'emergency'
);
```

---

## 6.7 authorization_status

```sql
CREATE TYPE authorization_status AS ENUM (
  'pending',
  'approved',
  'rejected',
  'revoked'
);
```

---

## 6.8 visit_verification_status

```sql
CREATE TYPE visit_verification_status AS ENUM (
  'clinic_verified'
);
```

---

## 6.9 notification_type

```sql
CREATE TYPE notification_type AS ENUM (
  'access_request',
  'access_approved',
  'access_rejected',
  'visit_created',
  'appointment_reminder',
  'vaccine_due'
);
```

---

## 7. Tables

# 7.1 users

Represents login accounts.

This table is used for both pet owners and clinic staff.

## Fields

```txt
id
email
phone
password_hash
role
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  email VARCHAR(255) NOT NULL UNIQUE,
  phone VARCHAR(30),
  password_hash TEXT NOT NULL,

  role user_role NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Notes

* `email` must be unique.
* `password_hash` stores bcrypt hash, never plain password.
* `role` controls whether user is owner, clinic_staff, or admin.

---

# 7.2 owner_profiles

Represents pet owner profile.

One user can have one owner profile.

## Fields

```txt
id
user_id
full_name
nickname
phone
email
address
emergency_contact_name
emergency_contact_phone
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE owner_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

  full_name VARCHAR(255) NOT NULL,
  nickname VARCHAR(100),
  phone VARCHAR(30),
  email VARCHAR(255),
  address TEXT,

  emergency_contact_name VARCHAR(255),
  emergency_contact_phone VARCHAR(30),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Relationships

```txt
users 1--1 owner_profiles
owner_profiles 1--many pets
```

## Rules

* Only users with role `owner` should have owner profile.
* One user should not create more than one owner profile.

---

# 7.3 breeds

Represents dog and cat breeds.

## Fields

```txt
id
species
name_th
name_en
image_url
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE breeds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  species pet_species NOT NULL,
  name_th VARCHAR(255),
  name_en VARCHAR(255) NOT NULL,
  image_url TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE(species, name_en)
);
```

## Initial Seed Data

Dog:

```txt
Golden Retriever
Pomeranian
French Bulldog
Poodle
Not sure / Select later
```

Cat:

```txt
British Shorthair
Persian
Scottish Fold
Siamese
Not sure / Select later
```

## Notes

Breed should be seed data in MVP.

Do not build breed management UI yet.

---

# 7.4 pets

Represents one digital pet.

## Fields

```txt
id
owner_id
pet_nexus_id
name
species
breed_id
gender
birth_date
approx_age_text
weight_kg
color_note
allergy_note
chronic_disease_note
photo_url
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE pets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  owner_id UUID NOT NULL REFERENCES owner_profiles(id) ON DELETE CASCADE,
  pet_nexus_id VARCHAR(50) NOT NULL UNIQUE,

  name VARCHAR(255) NOT NULL,
  species pet_species NOT NULL,
  breed_id UUID REFERENCES breeds(id) ON DELETE SET NULL,

  gender pet_gender NOT NULL DEFAULT 'unknown',
  birth_date DATE,
  approx_age_text VARCHAR(100),

  weight_kg NUMERIC(5,2),

  color_note TEXT,
  allergy_note TEXT,
  chronic_disease_note TEXT,

  photo_url TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Relationships

```txt
owner_profiles 1--many pets
breeds 1--many pets
pets 1--many qr_sessions
pets 1--many authorizations
pets 1--many visits
```

## Rules

* Owner can only access their own pets.
* `pet_nexus_id` must be unique.
* `birth_date` can be null if owner only knows approximate age.
* `weight_kg` should not be negative.
* Owner-added pet data should be visually separated from clinic-verified data in frontend.

## Suggested PetNexus ID Format

```txt
PNX-XXXXXX
```

Example:

```txt
PNX-8F3K2A
```

Generate it in backend service, not frontend.

---

# 7.5 clinics

Represents veterinary clinics.

## Fields

```txt
id
name
address
phone
email
verified_status
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE clinics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  name VARCHAR(255) NOT NULL,
  address TEXT,
  phone VARCHAR(30),
  email VARCHAR(255),

  verified_status clinic_verified_status NOT NULL DEFAULT 'pending',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Notes

For MVP, clinic verification can stay simple.

Do not build full clinic verification admin system yet.

---

# 7.6 clinic_staff

Represents staff accounts inside a clinic.

## Fields

```txt
id
clinic_id
user_id
role
license_no
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE clinic_staff (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

  role clinic_staff_role NOT NULL DEFAULT 'assistant',
  license_no VARCHAR(100),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Relationships

```txt
users 1--1 clinic_staff
clinics 1--many clinic_staff
```

## Rules

* Only users with role `clinic_staff` should have clinic_staff row.
* One user should belong to one clinic in MVP.
* Multi-clinic staff can be postponed.

---

# 7.7 qr_sessions

Represents temporary QR sharing sessions.

## Fields

```txt
id
pet_id
owner_id
token
purpose
expires_at
used_at
created_at
```

## SQL Plan

```sql
CREATE TABLE qr_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
  owner_id UUID NOT NULL REFERENCES owner_profiles(id) ON DELETE CASCADE,

  token TEXT NOT NULL UNIQUE,
  purpose qr_purpose NOT NULL DEFAULT 'clinic_checkin',

  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Rules

* QR contains token only.
* Token must be secure random string.
* Token must expire.
* QR scan must validate token through backend.
* Clinic should only see pet preview before owner approval.
* QR must never contain full pet data.

## Recommended Expiration

For MVP:

```txt
clinic_checkin: 15 minutes
passport_share: 30 minutes
emergency: 60 minutes
```

Can adjust later.

---

# 7.8 authorizations

Controls clinic access to pet data.

## Fields

```txt
id
pet_id
clinic_id
owner_id
status
permissions
reason
expires_at
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE authorizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
  clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
  owner_id UUID NOT NULL REFERENCES owner_profiles(id) ON DELETE CASCADE,

  status authorization_status NOT NULL DEFAULT 'pending',

  permissions TEXT[] NOT NULL DEFAULT ARRAY['view_profile', 'view_history'],
  reason TEXT,

  expires_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Recommended Index / Constraint

```sql
CREATE INDEX idx_authorizations_pet_id ON authorizations(pet_id);
CREATE INDEX idx_authorizations_clinic_id ON authorizations(clinic_id);
CREATE INDEX idx_authorizations_status ON authorizations(status);
```

Optional partial unique index:

```sql
CREATE UNIQUE INDEX unique_pending_authorization
ON authorizations(pet_id, clinic_id)
WHERE status = 'pending';
```

## Permissions

Use text array for MVP:

```txt
view_profile
view_history
create_visit
```

## Rules

* Clinic cannot view full pet data without approved authorization.
* Owner can approve, reject, or revoke.
* Authorization can expire.
* Duplicate pending request should be blocked.
* Clinic can create visit only if approved authorization includes `create_visit`.

## Why TEXT[] Instead of Permission Table?

For MVP, `TEXT[]` is simpler.

Later, if permission gets complex, create:

```txt
permissions
authorization_permissions
```

But do not do that in v1.

---

# 7.9 visits

Represents clinic-created medical visit records.

## Fields

```txt
id
pet_id
clinic_id
vet_id
visit_date
chief_complaint
diagnosis
treatment
medication
follow_up_date
note
verification_status
created_at
updated_at
```

## SQL Plan

```sql
CREATE TABLE visits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
  clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
  vet_id UUID REFERENCES clinic_staff(id) ON DELETE SET NULL,

  visit_date DATE NOT NULL,

  chief_complaint TEXT,
  diagnosis TEXT,
  treatment TEXT,
  medication TEXT,
  follow_up_date DATE,
  note TEXT,

  verification_status visit_verification_status NOT NULL DEFAULT 'clinic_verified',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Rules

* Only authorized clinic can create visit.
* Clinic-created visit must be `clinic_verified`.
* Owner can view clinic-created visit.
* Owner should not edit clinic-created visit.
* Clinic should not edit another clinic’s visit.
* Timeline should read from visits table.

---

# 7.10 notifications

Represents notifications sent to users.

## Fields

```txt
id
user_id
title
message
type
is_read
created_at
```

## SQL Plan

```sql
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  title VARCHAR(255) NOT NULL,
  message TEXT NOT NULL,
  type notification_type NOT NULL,

  is_read BOOLEAN NOT NULL DEFAULT false,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Notification Examples

```txt
access_request
visit_created
```

## Rules

* When clinic requests access, create notification for owner.
* When clinic creates visit, create notification for owner.
* User can mark notification as read.

---

# 7.11 audit_logs

Tracks sensitive actions.

## Fields

```txt
id
actor_user_id
action
target_type
target_id
metadata
created_at
```

## SQL Plan

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

  action VARCHAR(100) NOT NULL,
  target_type VARCHAR(100) NOT NULL,
  target_id UUID,

  metadata JSONB,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Example Actions

```txt
clinic_scanned_qr
clinic_requested_access
owner_approved_access
owner_rejected_access
owner_revoked_access
clinic_viewed_pet
clinic_created_visit
```

## Example Metadata

```json
{
  "petId": "pet_id",
  "clinicId": "clinic_id",
  "authorizationId": "authorization_id"
}
```

## Rules

* Use audit logs for sensitive access-related actions.
* Do not build admin audit log viewer in MVP unless needed for demo.

---

## 8. Relationship Diagram

```txt
users
  ├── owner_profiles
  │     └── pets
  │           ├── qr_sessions
  │           ├── authorizations
  │           └── visits
  │
  └── clinic_staff
        └── clinics
              ├── authorizations
              └── visits
```

More detailed:

```txt
User 1--1 OwnerProfile
OwnerProfile 1--many Pet
Breed 1--many Pet

User 1--1 ClinicStaff
Clinic 1--many ClinicStaff

Pet 1--many QRSession
Pet 1--many Authorization
Clinic 1--many Authorization

Pet 1--many Visit
Clinic 1--many Visit
ClinicStaff 1--many Visit

User 1--many Notification
User 1--many AuditLog as actor
```

---

## 9. Core Query Patterns

## 9.1 Owner Gets Own Pets

```sql
SELECT *
FROM pets
WHERE owner_id = $1
ORDER BY created_at DESC;
```

---

## 9.2 Check If Owner Owns Pet

```sql
SELECT 1
FROM pets
WHERE id = $1
AND owner_id = $2;
```

Use this before owner updates or views pet.

---

## 9.3 Clinic Checks Approved Access

```sql
SELECT 1
FROM authorizations
WHERE pet_id = $1
AND clinic_id = $2
AND status = 'approved'
AND (expires_at IS NULL OR expires_at > now())
AND $3 = ANY(permissions);
```

`$3` example:

```txt
view_profile
view_history
create_visit
```

---

## 9.4 Clinic Gets Approved Patients

```sql
SELECT p.*
FROM pets p
JOIN authorizations a ON a.pet_id = p.id
WHERE a.clinic_id = $1
AND a.status = 'approved'
AND (a.expires_at IS NULL OR a.expires_at > now())
ORDER BY a.updated_at DESC;
```

---

## 9.5 Owner Gets Pending Access Requests

```sql
SELECT a.*
FROM authorizations a
JOIN pets p ON p.id = a.pet_id
WHERE p.owner_id = $1
AND a.status = 'pending'
ORDER BY a.created_at DESC;
```

---

## 9.6 Pet Timeline

```sql
SELECT *
FROM visits
WHERE pet_id = $1
ORDER BY visit_date DESC, created_at DESC;
```

For MVP, timeline can be visit-only first.

Later, timeline can combine:

```txt
visits
vaccines
appointments
documents
owner_notes
```

---

## 10. Index Plan

Add indexes for fields used often in WHERE clauses.

```sql
CREATE INDEX idx_pets_owner_id ON pets(owner_id);
CREATE INDEX idx_pets_pet_nexus_id ON pets(pet_nexus_id);

CREATE INDEX idx_qr_sessions_token ON qr_sessions(token);
CREATE INDEX idx_qr_sessions_pet_id ON qr_sessions(pet_id);
CREATE INDEX idx_qr_sessions_expires_at ON qr_sessions(expires_at);

CREATE INDEX idx_authorizations_pet_id ON authorizations(pet_id);
CREATE INDEX idx_authorizations_clinic_id ON authorizations(clinic_id);
CREATE INDEX idx_authorizations_owner_id ON authorizations(owner_id);
CREATE INDEX idx_authorizations_status ON authorizations(status);

CREATE INDEX idx_visits_pet_id ON visits(pet_id);
CREATE INDEX idx_visits_clinic_id ON visits(clinic_id);
CREATE INDEX idx_visits_visit_date ON visits(visit_date);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);

CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_target ON audit_logs(target_type, target_id);
```

---

## 11. Validation Rules

Database constraints are not enough. Also validate in Go service layer.

## User

```txt
email required
password required
role must be owner or clinic_staff in MVP
email must be unique
```

## OwnerProfile

```txt
full_name required
emergency contact optional
one profile per user
```

## Pet

```txt
name required
species required
species must be dog or cat
breed must match species if breed_id exists
weight_kg should be greater than or equal to 0
birth_date cannot be in future
```

## QRSession

```txt
pet must belong to owner
token must be unique
expires_at must be future
expired token cannot be used
```

## Authorization

```txt
clinic cannot request duplicate pending access
owner must own pet before approve/reject/revoke
status transition must be valid
permissions must only include allowed values
```

Valid status transitions:

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
```

## Visit

```txt
clinic must have approved authorization
authorization must include create_visit
visit_date required
verification_status must be clinic_verified
```

---

## 12. Authorization Logic

Backend should have service/helper functions:

```txt
CanOwnerAccessPet(ownerUserId, petId)
CanClinicAccessPet(clinicId, petId, permission)
GetOwnerProfileByUserId(userId)
GetClinicStaffByUserId(userId)
```

## CanOwnerAccessPet

Logic:

```txt
1. Get owner profile from current user.
2. Check pets.id = petId and pets.owner_id = ownerProfile.id.
3. If yes, allow.
4. If no, reject with 403.
```

## CanClinicAccessPet

Logic:

```txt
1. Get clinic staff from current user.
2. Get clinic_id.
3. Find authorization where:
   - pet_id matches
   - clinic_id matches
   - status = approved
   - expires_at is null or in the future
   - permissions include required permission
4. If found, allow.
5. If not found, reject with 403.
```

---

## 13. QR Token Design

QR should not contain full pet data.

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

## Token Generation

Generate secure random token in Go backend.

Recommended:

```txt
crypto/rand
base64url string
minimum 32 bytes before encoding
```

## Token Storage

Store only token or token hash?

For MVP:

```txt
Store plain token in database.
```

For stronger security later:

```txt
Store token hash instead of plain token.
```

## QR Expiration

Recommended MVP expiration:

```txt
15 minutes for clinic check-in
```

---

## 14. Timeline Design

In MVP, timeline is based only on `visits`.

Timeline item:

```txt
id
type
date
clinic_name
chief_complaint
diagnosis
treatment
medication
follow_up_date
verification_status
created_at
```

Later timeline can include:

```txt
vaccination records
medicine records
appointments
documents
owner notes
weight logs
```

Do not build separate timeline table yet.

For MVP, generate timeline from visits query.

---

## 15. Notification Design

Notifications are simple database rows.

No realtime push in MVP.

## Create Notification When Clinic Requests Access

Recipient:

```txt
pet owner user_id
```

Type:

```txt
access_request
```

Message example:

```txt
Happy Pet Clinic requested access to Milo's record.
```

## Create Notification When Clinic Creates Visit

Recipient:

```txt
pet owner user_id
```

Type:

```txt
visit_created
```

Message example:

```txt
Happy Pet Clinic added a verified visit record for Milo.
```

---

## 16. Audit Log Design

Audit log should be created from service layer.

Example:

When clinic scans QR:

```txt
action = clinic_scanned_qr
target_type = pet
target_id = pet_id
actor_user_id = clinic_staff_user_id
metadata = { clinicId, qrSessionId }
```

When owner approves access:

```txt
action = owner_approved_access
target_type = authorization
target_id = authorization_id
actor_user_id = owner_user_id
metadata = { petId, clinicId }
```

When clinic creates visit:

```txt
action = clinic_created_visit
target_type = visit
target_id = visit_id
actor_user_id = clinic_staff_user_id
metadata = { petId, clinicId }
```

---

## 17. MVP Database Cut List

Do not create these tables in MVP unless needed:

```txt
payments
subscriptions
chat_messages
clinic_branches
documents
vaccination_records
medicine_records
weight_logs
pet_hotels
groomers
admin_reviews
refresh_tokens
oauth_accounts
timeline_events
```

Some of these are useful later, but they will slow down MVP.

---

## 18. Future Tables

After MVP, consider adding:

## 18.1 vaccination_records

For structured vaccine tracking.

```txt
id
pet_id
clinic_id
vaccine_name
vaccination_date
next_due_date
lot_number
verification_status
created_at
```

## 18.2 documents

For vaccine book images, PDFs, lab results.

```txt
id
pet_id
uploaded_by_user_id
file_url
file_type
title
created_at
```

## 18.3 appointments

For calendar and reminders.

```txt
id
pet_id
clinic_id
appointment_date
title
type
status
note
created_at
```

## 18.4 weight_logs

For pet weight tracking over time.

```txt
id
pet_id
weight_kg
recorded_at
source
created_at
```

Do not build these before the core flow works.

---

## 19. GORM Model Notes

GORM models should be simple.

Example model style:

```go
type Pet struct {
    ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    OwnerID            uuid.UUID `gorm:"type:uuid;not null"`
    PetNexusID         string    `gorm:"uniqueIndex;not null"`
    Name               string    `gorm:"not null"`
    Species            string    `gorm:"not null"`
    BreedID            *uuid.UUID
    Gender             string
    BirthDate          *time.Time
    ApproxAgeText      string
    WeightKg           *float64
    ColorNote          string
    AllergyNote        string
    ChronicDiseaseNote string
    PhotoURL           string
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

Use UUID package:

```txt
github.com/google/uuid
```

Use pointer types for nullable fields:

```txt
*uuid.UUID
*time.Time
*float64
```

---

## 20. Seed Data

## Dog Breeds

```sql
INSERT INTO breeds (species, name_en, name_th)
VALUES
('dog', 'Golden Retriever', 'โกลเด้น รีทรีฟเวอร์'),
('dog', 'Pomeranian', 'ปอมเมอเรเนียน'),
('dog', 'French Bulldog', 'เฟรนช์ บูลด็อก'),
('dog', 'Poodle', 'พุดเดิ้ล'),
('dog', 'Not sure / Select later', 'ยังไม่แน่ใจ / เลือกทีหลัง');
```

## Cat Breeds

```sql
INSERT INTO breeds (species, name_en, name_th)
VALUES
('cat', 'British Shorthair', 'บริติช ชอร์ตแฮร์'),
('cat', 'Persian', 'เปอร์เซีย'),
('cat', 'Scottish Fold', 'สก็อตติช โฟลด์'),
('cat', 'Siamese', 'วิเชียรมาศ'),
('cat', 'Not sure / Select later', 'ยังไม่แน่ใจ / เลือกทีหลัง');
```

---

## 21. Local Development Database

Recommended local `.env`:

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

Recommended local database name:

```txt
petnexus
```

---

## 22. Docker Compose Later

Do not require Docker on day one if it slows development.

But later, use Docker Compose for PostgreSQL:

```yaml
services:
  postgres:
    image: postgres:16
    container_name: petnexus-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: petnexus
    ports:
      - "5432:5432"
    volumes:
      - petnexus_postgres_data:/var/lib/postgresql/data

volumes:
  petnexus_postgres_data:
```

---

## 23. Database Milestones

## Milestone 1: Foundation

```txt
Create database
Connect Go backend to PostgreSQL
Run health route
```

## Milestone 2: Core Tables

```txt
Create users
Create owner_profiles
Create breeds
Create pets
Seed breeds
```

## Milestone 3: Clinic Tables

```txt
Create clinics
Create clinic_staff
```

## Milestone 4: Access Control Tables

```txt
Create qr_sessions
Create authorizations
```

## Milestone 5: Medical Records

```txt
Create visits
Create timeline query
```

## Milestone 6: Support Tables

```txt
Create notifications
Create audit_logs
```

---

## 24. Database MVP Definition

Database MVP is complete when these are possible:

```txt
A user can exist as owner.
An owner can have a profile.
An owner can create a pet.
A clinic staff user can belong to a clinic.
A pet can have a QR session.
A clinic can request authorization for a pet.
An owner can approve authorization.
A clinic can create a clinic_verified visit.
An owner can see the visit in the pet timeline.
Notifications can be created.
Audit logs can be created.
```

If these are possible, the database supports the PetNexus MVP.
