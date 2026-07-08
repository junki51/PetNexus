# Sprint 8: Appointment Calendar Foundation

Updated: 8 July 2026

## Goal

Provide real backend appointment data for owner scheduling and the Clinic Web
Calendar without adding medical, staff, payment, report, or notification scope.

## Completed

- Added guarded/idempotent `appointments` migration and model.
- Added owner create/list/detail/cancel endpoints.
- Added clinic create/list/detail/status/cancel endpoints.
- Added clinic creation by UUID pet ID or public pet ID.
- Added exact owner/clinic ownership scoping with cross-scope 404 responses.
- Added UTC date, inclusive date-range, status, and appointment-type filters.
- Added privacy-limited appointment responses with pet, owner, and clinic
  summaries.
- Added service, route, migration, validation, role, privacy, and regression
  tests.

## Product decisions

- Owner-created appointments begin as `requested`.
- Clinic-created appointments begin as `scheduled`.
- Clinic create accepts exactly one of `pet_id` and `public_pet_id`; both or
  neither returns 400.
- No date filter returns all scoped appointments sorted by `scheduled_at`
  ascending.
- `date`, `date_from`, and `date_to` are interpreted as UTC calendar
  dates. `date_to` includes the complete final day.
- `date` cannot be combined with `date_from` or `date_to`.
- Cancellation is idempotent. Clinic status updates deliberately use a simple
  allowed-value model; no workflow engine is included.

## Migration

File: `migrations/007_create_appointments.sql`

The migration creates the table only if missing, creates indexes with
`IF NOT EXISTS`, and adds named constraints only when absent. It does not use
GORM AutoMigrate and does not drop or rewrite Sprint 1–7 schema.

## Verification

- `gofmt` completed for changed Go files.
- `go test ./...` passed.
- Fresh PostgreSQL/API smoke test passed:
  - owner and clinic appointment creation
  - owner and clinic list/detail
  - clinic day/status/type filtering
  - clinic status update
  - owner and clinic cancellation
  - missing token 401
  - wrong roles 403
  - cross-owner/cross-clinic access 404
  - invalid type/status/time/duration 400
  - sensitive ownership fields absent
- Temporary PostgreSQL database was removed.

## Main files

Created:

- `migrations/007_create_appointments.sql`
- `internal/models/appointment.go`
- `internal/dto/appointment_dto.go`
- `internal/repositories/appointment_repository.go`
- owner/clinic appointment services and handlers
- appointment service and route tests

Updated:

- startup migration runner and migration tests
- clinic/pet repositories
- route registration and process wiring
- README and backend documentation

## Risks before deployment

- Render database role must allow table, index, constraint, and foreign-key
  creation.
- Calendar dates currently use UTC; a future clinic-timezone setting may be
  needed.
- There is no overlap detection, capacity management, pagination, or workflow
  transition engine yet.
- Concurrent startup instances can still contend while adding constraints;
  deploy with one migration-running instance and observe Render startup logs.

## Intentionally excluded

- Medical records
- Dashboard aggregation
- Reports and notifications
- Staff schedules
- Payment
- Google Calendar sync
- Full QR/access workflow
- Frontend changes
