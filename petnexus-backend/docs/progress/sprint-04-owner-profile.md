# Sprint 4: Owner Profile

Updated: 2 July 2026

## Goal

Add the owner identity/profile backend on top of Sprint 3 authentication,
without starting pet, clinic, QR, visit, timeline, or UI features.

## Completed

- Added the `owner_profiles` PostgreSQL table with UUID keys and timestamps.
- Added an idempotent unique index on `user_id` and a guarded foreign key to
  `users(id)`.
- Kept startup migrations as explicit SQL; GORM `AutoMigrate` is not used.
- Added OwnerProfile model and create/update/response DTOs.
- Added repository, service, and handler layers.
- Added owner-only `POST`, `GET`, and `PATCH /api/owner/profile` routes.
- Profile ownership always comes from the authenticated JWT user ID.
- Added validation, whitespace normalization, gender allow-list, date checks,
  URL validation, partial updates, and computed `display_name`.
- Added automated service, route authorization, identity-spoofing, and
  migration tests.
- Updated the project and migration READMEs with PowerShell manual tests.

## API behavior

```text
POST  /api/owner/profile  -> 201 or 409 when a profile already exists
GET   /api/owner/profile  -> 200 or 404 when no profile exists
PATCH /api/owner/profile  -> 200; an empty patch returns 400
```

All three endpoints require JWT authentication and role `owner`. Missing or
invalid authentication returns 401; other roles return 403.

## Verification

- `gofmt` completed.
- `go test ./...` passed, including Owner Profile service and route tests.
- Docker/PostgreSQL API smoke testing was not run on 2 July 2026 because the
  local Docker daemon was not running.

## Intentionally not implemented

- Pet, Breed, Pet Passport, QR Sharing
- Clinic Access or authorization flows
- Clinic Visit or Timeline
- Flutter owner UI or clinic web UI
- Any changes to the existing Auth API behavior

## Recommended next step

Run the README PowerShell smoke test against local Docker PostgreSQL, redeploy
the backend to Render, and verify the same Owner Profile flow there. After that,
plan the Pet/Breed database foundation as a separate sprint.
