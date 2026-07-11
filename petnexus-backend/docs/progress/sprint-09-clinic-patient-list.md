# Sprint 9: Clinic Patient List Backend

Date: 2026-07-11

## Goal

Make the Clinic Web Patients page use real backend data.

For Sprint 9, a clinic patient means:

```text
A unique pet that has at least one non-cancelled appointment with the current clinic profile.
```

No separate `patients` table was added. Patient data is derived from:

- `appointments`
- `pets`
- `owner_profiles`
- `breeds`

## Added endpoints

```text
GET /api/clinic/patients
GET /api/clinic/patients/:petId
```

Both endpoints require:

- valid JWT
- clinic-side role (`clinic`, with `clinic_staff` still supported)
- an existing clinic profile for the current authenticated user

## Patient list behavior

`GET /api/clinic/patients` returns unique pets for the authenticated clinic.
Only appointments with `status <> cancelled` are used for the patient
relationship and appointment summary.

Supported query parameters:

```text
q        optional search by pet name or public_pet_id
species  optional dog|cat
status   optional latest computed appointment status
limit    optional, default 20, max 100
offset   optional, default 0
sort     optional, default latest_appointment_desc
```

Supported sort values:

```text
latest_appointment_desc
latest_appointment_asc
name_asc
name_desc
next_appointment_asc
```

## Patient detail behavior

`GET /api/clinic/patients/:petId` returns:

- pet identity/detail summary
- owner display name and masked phone
- clinic relationship summary
- latest recent non-cancelled appointments for this clinic and pet

It does not return:

- `user_id`
- `owner_profile_id`
- `clinic_profile_id`
- password data
- JWT claims
- full owner address
- medical records

## Ownership and privacy rules

- The service uses `currentUserID` from JWT only.
- The service resolves the current clinic profile by user ID.
- The repository scopes all patient data by `clinic_profile_id`.
- Cross-clinic or unrelated pet access returns 404.
- Owner role receives 403 from route middleware.
- Missing/invalid token receives 401 from auth middleware.
- Owner phone is masked using the existing clinic lookup masking behavior.

## Files created

- `internal/dto/clinic_patient_dto.go`
- `internal/repositories/clinic_patient_repository.go`
- `internal/services/clinic_patient_service.go`
- `internal/handlers/clinic_patient_handler.go`
- `internal/services/clinic_patient_service_test.go`
- `internal/routes/clinic_patient_routes_test.go`
- `docs/progress/sprint-09-clinic-patient-list.md`
- `docs/sprints/sprint-1-to-9-summary.md`

## Files changed

- `cmd/api/main.go`
- `internal/routes/routes.go`
- `README.md`
- `docs/README.md`
- `docs/progress/README.md`
- `docs/backend/api-reference.md`
- `docs/backend/auth-and-permissions.md`
- `docs/backend/backend-overview.md`
- `docs/backend/database-schema.md`
- `docs/backend/module-map.md`
- `docs/backend/roadmap.md`
- `docs/backend/testing-guide.md`
- `migrations/README.md`

## Migration

No migration was added for Sprint 9.

Reason: clinic patients are derived from existing appointment relationships.
The existing `appointments` schema already has the required clinic, pet, owner,
status, and scheduled time fields.

## Verification

Automated checks run locally:

```powershell
gofmt -w <changed Go files>
go test ./...
```

Result: passed.

## Not included

- Medical records
- Visit records
- Dashboard aggregation
- Reports
- Notifications
- Staff schedule
- Payment
- Google Calendar sync
- Frontend
- Separate patients table

## Assumptions

- Cancelled-only appointment relationships should not create a patient row.
- Recent appointments on patient detail also use non-cancelled appointments.
- `status` filter matches the latest computed non-cancelled appointment status.
- Calendar date/time values remain UTC at the backend layer until clinic
  timezone support exists.

## Risks before deployment

- If production data has many appointments, later pagination may need a total
  count endpoint or response wrapper.
- Patient search currently covers pet name and `public_pet_id` only.
- There is no medical-record access model yet, so the Patients page must not
  assume medical records exist.
