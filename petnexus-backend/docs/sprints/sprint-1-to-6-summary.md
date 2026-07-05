# Sprint 1–6 Backend Summary

## Sprint 1: Backend Foundation

**Added:** Go/Gin project structure, environment config loading, response
helper, and server startup.

**Main endpoint:** `GET /health`.

**Tables:** None.

**Access:** Public health check.

**Intentionally excluded:** Database, auth, profiles, pets, clinic workflows.

## Sprint 2: PostgreSQL Foundation

**Added:** Docker Compose PostgreSQL, GORM connection, local `DB_*` settings,
Render `DATABASE_URL` support, and database ping.

**Main endpoint:** `GET /health/db`.

**Tables:** No domain tables introduced by the foundation itself.

**Access:** Public database health check.

**Intentionally excluded:** Auth and business tables/features.

## Sprint 3: Auth Foundation

**Added:** User account schema, bcrypt hashing, JWT access tokens, public
register/login, current-user lookup, auth middleware, and role middleware.

**Main endpoints:**

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`

**Main table/type:** `users`, `user_role` enum.

**Access:** Register/login public; `/api/me` requires JWT. Public registration
supports owner and clinic-side roles while admin is blocked.

**Intentionally excluded:** Refresh tokens, password recovery, email
verification, owner/clinic domain profiles.

## Sprint 4: Owner Profile

**Added:** One editable owner identity/contact profile per owner account with
create, fetch, and partial update.

**Main endpoints:**

- `POST /api/owner/profile`
- `GET /api/owner/profile`
- `PATCH /api/owner/profile`

**Main table:** `owner_profiles`.

**Access:** JWT + role `owner`; profile is resolved from JWT user ID. Client
does not send `user_id`.

**Intentionally excluded:** Pet creation, QR, clinics, visits, and UI.

## Sprint 5: Breed + Pet Creation

**Added:** Dog/cat breed catalog with idempotent seed data and owner-controlled
basic pet create/list/detail/partial update.

**Main endpoints:**

- `GET /api/breeds`
- `GET /api/breeds?species=dog`
- `GET /api/breeds?species=cat`
- `POST /api/pets`
- `GET /api/pets`
- `GET /api/pets/:id`
- `PATCH /api/pets/:id`

**Main tables:** `breeds`, `pets`.

**Access:** Breed list public. Pet routes require JWT + role `owner` + existing
owner profile. Pet ownership comes from JWT user → owner profile; client does
not send `user_id` or `owner_profile_id`.

**Intentionally excluded:** Pet Passport, QR, clinic authorization, medical
records, visits, timeline, and real image upload.

## Sprint 6: Clinic Profile Foundation

**Added:** Clinic-side registration compatibility and one clinic settings
profile per clinic account with create, fetch, and partial update.

**Main endpoints:**

- `POST /api/clinic/profile`
- `GET /api/clinic/profile`
- `PATCH /api/clinic/profile`

**Main table:** `clinic_profiles`.

**Access:** JWT + canonical role `clinic`; legacy `clinic_staff` remains
compatible. Owner receives 403. Profile is resolved from JWT user ID and API
responses do not expose it.

**Intentionally excluded:** Staff management, QR, clinic access request,
patients, visits, medical records, timeline, calendar, and reports.

## Current backend boundary

The backend can identify owners, their profiles and pets, and clinic accounts
with clinic profiles. It cannot yet create a trusted relationship between a
clinic and an owner's pet. That authorization boundary is the recommended next
design topic.
