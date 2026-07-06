# Sprint 7: Pet Public ID + Clinic Pet Lookup

Updated: 6 July 2026

## Goal

Make Clinic Pet Lookup the primary pet discovery foundation. QR remains an
optional future shortcut and no QR token/access workflow is implemented.

## Completed

- Added backend-generated `pets.public_pet_id` values in
  `PNX-PET-XXXXXX` format.
- Added a unique index and non-null enforcement.
- Added safe startup and manual backfill for existing pets.
- Added retry handling for rare public-ID collisions.
- Included `public_pet_id` in owner Pet responses without accepting it in
  create/update requests.
- Added clinic lookup by public pet ID.
- Added clinic lookup by exact `owner_profiles.phone_number`.
- Added limited clinic DTO with masked owner phone and no sensitive IDs/data.
- Added JWT + clinic role route protection.
- Added generator, migration, service, route, privacy, and regression tests.

## Endpoint

```text
GET /api/clinic/pet-lookup?pet_id=PNX-PET-8F3K2A
GET /api/clinic/pet-lookup?owner_phone=0812345678
```

Query behavior:

- neither query: 400
- both queries: 400
- unknown `pet_id`: 404
- exact phone with no pets: 200 and `items: []`
- partial phone matching is not performed

## Privacy boundary

Clinic lookup returns only:

- pet UUID and public pet ID
- name, species, breed
- gender, date of birth, avatar URL
- owner display name and masked phone

It does not return user ID, owner-profile ID, raw owner phone, password/auth
data, weight, microchip, color, distinctive marks, neuter status, or any
medical/timeline data.

## Migration and ID generation

Migration: `migrations/006_add_public_pet_id.sql`

Startup order:

1. Add nullable column if missing.
2. Normalize empty values to null.
3. Create unique index (PostgreSQL permits multiple nulls during backfill).
4. Generate and assign IDs in application startup code, retrying unique
   collisions.
5. Enforce `NOT NULL` only after all pets are assigned.

New pet creation generates the ID before repository create and retries if the
database unique index reports a collision. The alphabet is uppercase
alphanumeric with confusing characters omitted.

## Main files

Created:

- `migrations/006_add_public_pet_id.sql`
- `internal/utils/public_pet_id.go`
- `internal/utils/public_pet_id_test.go`
- `internal/services/clinic_pet_lookup_service.go`
- `internal/services/clinic_pet_lookup_service_test.go`
- `internal/handlers/clinic_pet_lookup_handler.go`
- `internal/routes/clinic_pet_lookup_routes_test.go`

Updated:

- `internal/models/pet.go`
- `internal/dto/pet_dto.go`
- `internal/repositories/pet_repository.go`
- `internal/services/pet_service.go`
- `internal/services/pet_service_test.go`
- `internal/database/migrate.go`
- `internal/database/migrate_test.go`
- `internal/routes/routes.go`
- `cmd/api/main.go`
- backend README/docs and migration docs

## Verification

- `gofmt` completed for changed Go files.
- `go test ./...` passed after Sprint 7 tests were added.
- Fresh PostgreSQL smoke test simulated a Sprint 5 pet without public ID:
  - legacy pet received a valid ID
  - new pet received a valid ID
  - every pet ID was non-null and unique
  - lookup by ID and exact phone succeeded
  - partial phone returned no items
  - no-token returned 401
  - owner role returned 403
  - missing/both query returned 400
  - unknown pet returned 404
  - sensitive lookup fields were absent
  - Owner Profile regression succeeded
  - restart migration remained idempotent
- Temporary database and test artifacts were removed.

## Assumptions

- `owner_profiles.phone_number` is the exact phone source.
- Both lookup parameters together are rejected with 400.
- Unknown phone is represented as a successful empty list.
- Lookup normalizes public pet ID to uppercase but does not perform partial ID
  matching.
- Legacy `clinic_staff` remains compatible with canonical `clinic` role.

## Deployment risks

- Existing pets are updated during startup; deploy when migration can finish
  before traffic is served.
- Six-character IDs have finite space, but database uniqueness plus retry
  prevents accepted duplicates.
- Multiple startup instances are safe at assignment time because the unique
  index exists before backfill; a collision is retried.
- Verify Render database privileges allow `ALTER TABLE`, index creation, and
  updates to existing pets.
- Monitor migration time if the pet table becomes large in future.

## Intentionally excluded

- QR token/expiry/revocation system
- Clinic access request and owner approval
- Medical records and visits
- Timeline, calendar, reports, and notifications
- Frontend changes

## Next step

Redeploy to Render, verify startup migration logs, confirm existing pet IDs were
backfilled, and repeat the documented lookup/privacy/status tests. Design clinic
access authorization as a separate sprint.
