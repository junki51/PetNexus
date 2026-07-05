# Sprint 6: Clinic Profile Foundation

Updated: 5 July 2026

## Goal

Add the minimum backend foundation for the Clinic Web Dashboard Settings page:
one clinic profile bound to the currently authenticated clinic-side user.

## Role assumption

`clinic` is the canonical Clinic Web Dashboard role. PostgreSQL's existing
`user_role` enum and guarded startup migration already contain both `clinic`
and the legacy `clinic_staff` value, so no enum migration was needed. Public
registration and Clinic Profile routes accept `clinic`; `clinic_staff` remains
accepted for backward compatibility.

## Completed

- Added safe startup/manual migration for `clinic_profiles`.
- Added a unique `user_id` index and guarded foreign key to `users(id)`.
- Added ClinicProfile model and create/update/response DTOs.
- Added repository, service, handler, and route wiring.
- Added clinic-staff-only POST, GET, and PATCH profile endpoints.
- Bound every operation to the JWT user ID; request DTOs do not contain
  `user_id` and responses do not expose it.
- Added required clinic-name validation, optional phone/email/address
  normalization, email format validation, and partial update behavior.
- Added duplicate protection in service logic and the database unique index.
- Added canonical `clinic` public registration and Clinic Profile route access;
  legacy `clinic_staff` remains compatible and owner remains forbidden.
- Added migration, service, route-role, and identity-spoofing regression tests.

## Database

Migration: `migrations/005_create_clinic_profiles.sql`

```text
clinic_profiles
- id UUID primary key
- user_id UUID not null, unique, references users(id)
- clinic_name varchar(200) not null
- phone_number varchar(30) nullable
- email varchar(255) nullable
- address text nullable
- created_at timestamptz not null
- updated_at timestamptz not null
```

Startup remains idempotent and does not use GORM AutoMigrate or destructive
constraint changes.

## Endpoints

```text
POST  /api/clinic/profile
GET   /api/clinic/profile
PATCH /api/clinic/profile
```

All endpoints require JWT authentication and a clinic-side role (`clinic`, or
legacy-compatible `clinic_staff`). Owner remains forbidden.

| Situation | Status/code |
| --- | --- |
| Create profile | 201 |
| Fetch or update profile | 200 |
| Invalid JSON/validation/empty PATCH | 400 |
| Missing or invalid token | 401 `UNAUTHORIZED` |
| Owner or other disallowed role | 403 `FORBIDDEN_ROLE` |
| Profile missing | 404 `CLINIC_PROFILE_NOT_FOUND` |
| Duplicate profile | 409 `CLINIC_PROFILE_ALREADY_EXISTS` |
| Unexpected database/server error | 500 `INTERNAL_SERVER_ERROR` |

## Main files

Created or implemented:

- `migrations/005_create_clinic_profiles.sql`
- `internal/models/clinic.go`
- `internal/dto/clinic_dto.go`
- `internal/repositories/clinic_repository.go`
- `internal/services/clinic_service.go`
- `internal/handlers/clinic_handler.go`
- `internal/services/clinic_service_test.go`
- `internal/routes/clinic_routes_test.go`

Integrated or updated:

- `internal/database/migrate.go`
- `internal/database/migrate_test.go`
- `internal/routes/routes.go`
- `cmd/api/main.go`
- `README.md`
- `migrations/README.md`
- `docs/progress/README.md`

## Verification

- `gofmt` completed.
- `go test ./...` passed after implementation and Sprint 6 tests were added.
- Existing Sprint 1–5 regression tests remain green.
- Fresh role-specific PostgreSQL/API smoke test passed: registration, login,
  and `/api/me` returned role `clinic`; clinic create/get/PATCH succeeded; owner
  access returned 403.
- Fresh local PostgreSQL/API smoke test passed:
  - health and database health
  - unauthenticated clinic profile request returned 401
  - clinic GET before create returned 404
  - clinic staff create/get/PATCH succeeded
  - duplicate create returned 409
  - owner role returned 403
  - spoofed `user_id` was ignored and not exposed
  - existing owner profile and pet creation still succeeded
  - restart against the same database completed migration idempotently
- The temporary smoke database and Clinic test artifacts were removed after
  verification; existing development data was not modified.

## Intentionally excluded

- QR sharing or scanning
- Clinic access request and owner approval/rejection
- Visits, medical records, and timeline
- Calendar or appointments
- Analytics and reports
- Staff-member management
- Clinic frontend/UI

## Next step

Run the README curl flow locally, deploy/redeploy to Render, confirm startup
migration logs are clean, and repeat clinic-profile create/get/PATCH/role tests
against the Render API before starting another feature sprint.
