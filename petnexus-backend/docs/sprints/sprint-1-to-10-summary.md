# PetNexus Backend Sprint 1-10 Summary

Last updated: 2026-07-11

## Completed backend scope

### Sprint 1: Backend Foundation

- Go + Gin API scaffold
- Config loading
- Standard response helper
- `GET /health`

### Sprint 2: PostgreSQL Foundation

- Local PostgreSQL with Docker Compose
- GORM PostgreSQL connection
- `DATABASE_URL` support for Render
- `GET /health/db`

### Sprint 3: Auth Foundation

- `users` table
- `user_role` enum
- Register/login
- bcrypt password hashing
- JWT access token
- Auth middleware
- Role middleware
- `GET /api/me`

### Sprint 4: Owner Profile

- `owner_profiles` table
- Owner-only profile create/get/patch
- Profile ownership resolved from JWT

### Sprint 5: Breed + Pet Creation

- `breeds` table
- `pets` table
- Public breed list
- Owner-only pet create/list/get/patch
- Backend owns `owner_profile_id`

### Sprint 6: Clinic Profile Foundation

- `clinic` role supported
- `clinic_profiles` table
- Clinic profile create/get/patch
- Clinic profile ownership resolved from JWT

### Sprint 7: Public Pet ID + Clinic Lookup

- Backend-generated `public_pet_id`
- Existing pet backfill
- Privacy-limited clinic pet lookup by public pet ID or exact owner phone
- QR remains future optional shortcut only

### Sprint 8: Appointment Calendar Foundation

- `appointments` table
- Owner appointment create/list/get/cancel
- Clinic appointment create/list/get/status/cancel
- Clinic calendar filters
- Appointment responses include safe pet, owner, and clinic summaries

### Sprint 9: Clinic Patient List Backend

- Clinic patient list endpoint
- Clinic patient detail endpoint
- Patients derived from non-cancelled appointments
- Search/filter/pagination basics for Clinic Web Patients page
- Cross-clinic access hidden as 404

### Sprint 10: Medical Records Foundation Backend

- `medical_records` table
- Clinic creates medical records for existing patients
- Optional one-to-one appointment link
- Clinic medical record list/detail/update endpoints
- Pet/date pagination filters
- Cross-clinic medical record access hidden as 404

## Current clinic-side backend coverage

Implemented:

- Settings / Clinic Profile
- QR Pet Data / Pet Lookup foundation
- Calendar / Appointments foundation
- Patients list/detail foundation
- Medical Records foundation

Not implemented yet:

- Visit Records
- Dashboard aggregation
- Reports
- Notifications
- Staff schedule
- Staff account management
- Payment
- Google Calendar sync
- Full QR system
- Access request / owner approval
- Owner medical timeline
- File uploads
- Lab results
- Vaccination records
- Prescription tables
- Frontend

## Sprint 10 endpoints

```text
POST  /api/clinic/patients/:petId/medical-records
GET   /api/clinic/medical-records
GET   /api/clinic/medical-records/:recordId
PATCH /api/clinic/medical-records/:recordId
```

All require JWT authentication and clinic-side role.

## Sprint 10 schema

`medical_records` belongs to:

```text
clinic_profiles 1:N medical_records
pets 1:N medical_records
appointments 0..1:0..1 medical_records
users 1:N medical_records as creator
```

No owner profile ID is stored directly.

## Validation status

Latest local automated checks:

```powershell
gofmt -w <changed Go files>
go test ./...
```

Result: passed.

## Recommended next step

Deploy/redeploy to Render and run the Sprint 10 smoke flow after the existing
owner, pet, clinic profile, patient, and appointment setup flow.
