# Sprint 5: Breed + Pet Creation Backend

Implemented: 4 July 2026  
Documentation reviewed: 5 July 2026

## Goal

Add the first owner-controlled pet profile flow and a reusable dog/cat breed
catalog. This sprint intentionally stops before Pet Passport, QR, clinic
authorization, visit, and timeline behavior.

## Completed

- Added `breeds` and `pets` models and safe startup/manual migrations.
- Added guarded species/gender checks, foreign keys, and required indexes.
- Seeded 8 dog and 8 cat breeds idempotently.
- Added public breed listing with optional dog/cat filtering.
- Added owner-only pet create, list, detail, and partial-update endpoints.
- Resolved pet ownership only through JWT user ID and `owner_profiles`.
- Enforced breed existence and breed/pet species consistency.
- Returned 404 for cross-owner pet access to avoid leaking pet existence.
- Added validation for dates, weight, gender, URLs, string lengths, and empty
  PATCH bodies.
- Added service, route-security, migration, and identity-spoofing tests.

## Database

Migration: `migrations/004_create_breeds_and_pets.sql`

- `breeds` stores `species`, English name, optional Thai name, and timestamps.
- `pets` belongs to `owner_profiles`; `breed_id` is optional.
- Supported species are `dog` and `cat`.
- Supported pet genders are `male`, `female`, and `unknown`.
- Foreign keys connect pets to owner profiles and breeds.
- Indexes cover breed lookup and pet owner/breed/species lookup.
- Breed seed data uses `ON CONFLICT DO NOTHING`, so startup is repeatable.
- Startup uses guarded SQL and does not introduce GORM `AutoMigrate`.

## Endpoints

```text
GET   /api/breeds
POST  /api/pets
GET   /api/pets
GET   /api/pets/:id
PATCH /api/pets/:id
```

`GET /api/breeds` is public and accepts `?species=dog` or `?species=cat`.
Every `/api/pets` endpoint requires a valid owner JWT and an existing owner
profile.

## API rules

- Ownership comes only from the JWT user ID and its `owner_profiles` row.
- Request DTOs do not contain `user_id` or `owner_profile_id`.
- Cross-owner detail/update requests return 404 instead of revealing that the
  pet exists.
- A supplied breed must exist and have the same species as the pet.
- PATCH updates only supplied fields and rejects an empty JSON object.
- Sending an empty `breed_id` in PATCH clears the current breed.
- Responses omit owner/user IDs and compute `age_years` from `date_of_birth`.

Expected status codes:

| Situation | Status |
| --- | --- |
| Pet created | 201 |
| Breed/pet fetched or pet updated | 200 |
| Invalid JSON, UUID, species, breed match, date, weight, or empty PATCH | 400 |
| Missing or invalid JWT | 401 |
| Authenticated non-owner | 403 |
| Missing owner profile, breed, pet, or cross-owner access | 404 |
| Unexpected database/server failure | 500 |

## Validation

- Pet name and species are required and strings are whitespace-trimmed.
- Pet/breed species is restricted to dog or cat.
- Date of birth uses `YYYY-MM-DD` and cannot be in the future.
- Weight must be greater than 0 and at most 200 kg.
- Avatar accepts valid HTTP/HTTPS URLs.
- String lengths are bounded to match the database schema.
- Nullable optional fields can be cleared without overwriting unspecified
  fields.

## Main files

Created or implemented:

- `migrations/004_create_breeds_and_pets.sql`
- `internal/models/breed.go`
- `internal/models/pet.go`
- `internal/dto/pet_dto.go`
- `internal/repositories/breed_repository.go`
- `internal/repositories/pet_repository.go`
- `internal/services/breed_service.go`
- `internal/services/pet_service.go`
- `internal/handlers/breed_handler.go`
- `internal/handlers/pet_handler.go`
- `internal/routes/pet_routes_test.go`
- `internal/services/pet_service_test.go`

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
- `go test ./...` passed again on 5 July 2026.
- Sprint 5 service, route-security, identity-spoofing, and migration tests are
  included in the repository.
- Fresh PostgreSQL smoke test passed:
  - health and database health
  - 16 total breeds, 8 dog, 8 cat
  - existing auth and owner profile flow
  - pet create/list/detail/patch
  - 401 without token and 403 for clinic staff
  - 400 invalid species and breed mismatch
  - 404 cross-owner pet access
  - restart migration stayed idempotent with 16 breeds

Render has not been re-verified for Sprint 5 yet; deployment verification is
the remaining operational step.

## Assumptions

- Breed listing is public because it contains no sensitive data.
- `breed_id` is optional and an empty value in PATCH clears the breed.
- `microchip_id` is not unique until the product defines its ownership and
  collision policy.
- Seeded `name_th` values remain null until translations are product-approved.

## Intentionally excluded

Pet Passport, QR sharing, clinic access/authorization, clinic visits, timeline,
notifications, file upload, Flutter UI, and clinic web UI.

## Before commit

1. Review `git diff` and confirm only intended Sprint 5/documentation files are
   included.
2. Run `gofmt` on changed Go files.
3. Run `go test ./...`.
4. Commit the migration, implementation, tests, and documentation together so
   startup migration and routes cannot be deployed separately.
5. Do not commit `.env`, database URLs, passwords, JWT secrets, or test tokens.

## Render deployment checklist

1. Push the Sprint 5 commit and deploy/redeploy the Render web service.
2. Confirm startup logs contain `database migration completed successfully`
   and no migration or constraint errors.
3. Verify `GET /health` and `GET /health/db`.
4. Verify 16 breeds total, including 8 dogs and 8 cats.
5. Repeat register, login, owner-profile, pet create/list/detail/PATCH, 401,
   403, breed mismatch, and cross-owner 404 checks from the root README.
6. Restart/redeploy once more and confirm the seed remains at 16 breeds.

## Next development step

Complete the Render checklist above, then scope Pet Passport or QR sharing as a
separate sprint. Do not mix clinic access, visits, or timeline behavior into
this completed basic pet profile foundation.
