# Sprint 10: Medical Records Foundation Backend

Date: 2026-07-11

## Goal

Allow Clinic Web to create, list, filter, view, and update medical records for
the authenticated clinic's patients.

This sprint is backend only.

## Added table

`medical_records`

Important fields:

- `clinic_profile_id`
- `pet_id`
- optional `appointment_id`
- `created_by_user_id`
- `visit_at`
- `chief_complaint`
- clinical free-text fields
- optional vitals
- timestamps

No `owner_profile_id` is stored because the owner can be derived through:

```text
medical_records -> pets -> owner_profiles
```

## Added endpoints

```text
POST  /api/clinic/patients/:petId/medical-records
GET   /api/clinic/medical-records
GET   /api/clinic/medical-records/:recordId
PATCH /api/clinic/medical-records/:recordId
```

All endpoints require JWT and clinic-side role.

## Permission rules

- `clinic_profile_id` is resolved from the authenticated JWT user.
- `created_by_user_id` is resolved from the authenticated JWT user.
- Create uses `petId` from the URL path.
- Clinic can create records only for pets that are already patients of the
  current clinic.
- A patient means a pet with at least one non-cancelled appointment with the
  current clinic.
- Cross-clinic records return 404.
- Owner role receives 403.

## Appointment rules

`appointmentId` is optional.

When present, the service verifies:

- appointment exists
- appointment belongs to the current clinic profile
- appointment belongs to the same pet in the URL
- appointment is not cancelled
- appointment does not already have a medical record

An appointment can have at most one medical record through the partial unique
index:

```text
idx_medical_records_appointment_id_unique
```

## Validation

- `petId` and `recordId` must be UUIDs.
- `appointmentId` must be UUID when supplied.
- `visitAt` is required and must use RFC3339.
- `chiefComplaint` is required and cannot be blank.
- `weightKg` and `temperatureC` must be greater than zero when supplied.
- `nextFollowUpAt` cannot be earlier than `visitAt`.
- `from` and `to` filters use `YYYY-MM-DD`.
- `from` cannot be after `to`.
- `page` and `limit` must be positive integers.
- `limit` is capped at 100.

## Files created

- `internal/models/medical_record.go`
- `internal/dto/medical_record_dto.go`
- `internal/repositories/medical_record_repository.go`
- `internal/services/medical_record_service.go`
- `internal/handlers/medical_record_handler.go`
- `internal/services/medical_record_service_test.go`
- `internal/routes/medical_record_routes_test.go`
- `migrations/008_create_medical_records.sql`
- `docs/progress/sprint-10-medical-records-foundation.md`
- `docs/sprints/sprint-1-to-10-summary.md`

## Files changed

- `cmd/api/main.go`
- `internal/database/migrate.go`
- `internal/database/migrate_test.go`
- `internal/routes/routes.go`
- `README.md`
- `migrations/README.md`
- `docs/README.md`
- `docs/progress/README.md`
- `docs/backend/api-reference.md`
- `docs/backend/auth-and-permissions.md`
- `docs/backend/backend-overview.md`
- `docs/backend/database-schema.md`
- `docs/backend/module-map.md`
- `docs/backend/roadmap.md`
- `docs/backend/testing-guide.md`

## Verification

Automated checks:

```powershell
gofmt -w <changed Go files>
go test ./...
```

Result: passed.

## Not included

- Dashboard Summary
- Reports
- Notifications
- Payment
- Staff scheduling
- Staff account management
- Google Calendar sync
- Full QR system
- Access request or owner approval
- Frontend
- Owner medical timeline
- Cross-clinic record sharing
- File uploads
- Images
- Lab results
- Vaccination records
- Allergy system
- Separate prescription tables
- Medication catalog
- Audit log system
- Medical record version history
- AI diagnosis
- Hard delete endpoint

## Assumptions

- Medical record request/response fields use camelCase for Clinic Web:
  `visitAt`, `chiefComplaint`, `appointmentId`.
- Medical record clinical fields are free text for this foundation.
- Walk-in/historical records are allowed when the pet is already a clinic
  patient, even without `appointmentId`.
- Owner medical timeline and owner visibility are future sprints.

## Risks before deployment

- Render and local PostgreSQL are separate databases; verify migration and API
  flows on both.
- The free-text `medications` field is intentionally not a prescription system.
- There is no record version history or audit log yet, so future regulatory
  requirements may need additional schema.
