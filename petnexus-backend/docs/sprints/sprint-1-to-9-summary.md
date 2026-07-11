# PetNexus Backend Sprint 1-9 Summary

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

## Current clinic-side backend coverage

Implemented:

- Settings / Clinic Profile
- QR Pet Data / Pet Lookup foundation
- Calendar / Appointments foundation
- Patients list/detail foundation

Not implemented yet:

- Medical Records
- Visit Records
- Dashboard aggregation
- Reports
- Notifications
- Staff schedule
- Payment
- Google Calendar sync
- Frontend

## Sprint 9 endpoints

```text
GET /api/clinic/patients
GET /api/clinic/patients/:petId
```

Both require JWT authentication and clinic-side role.

## Sprint 9 patient definition

```text
clinic patient = unique pet with at least one non-cancelled appointment
                 for the authenticated clinic profile
```

No `patients` table was added.

## Validation status

Latest local automated checks:

```powershell
gofmt -w <changed Go files>
go test ./...
```

Result: passed.

## Recommended next step

Deploy/redeploy to Render and run the Sprint 9 patient smoke flow after the
existing owner, pet, clinic profile, and appointment setup flow.
