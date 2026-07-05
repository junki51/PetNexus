# Backend Module Map

## Health module

- **Purpose:** Report process and PostgreSQL availability.
- **Main tables:** None directly.
- **Endpoints:** `GET /health`, `GET /health/db`.
- **Access:** Public.
- **Code:** health handler plus PostgreSQL ping helper.

## Auth module

- **Purpose:** Account registration/login, JWT creation/validation, current user.
- **Main table:** `users`; PostgreSQL enum `user_role`.
- **Endpoints:** `POST /api/auth/register`, `POST /api/auth/login`,
  `GET /api/me`.
- **Access:** Register/login public; `/api/me` authenticated.
- **Rules:** Public roles are `owner`, canonical `clinic`, and legacy-compatible
  `clinic_staff`; public `admin` registration is forbidden. Passwords are
  bcrypt hashes and never returned.

## Owner Profile module

- **Purpose:** Store editable identity/contact information separately from auth.
- **Main table:** `owner_profiles`.
- **Endpoints:** POST/GET/PATCH `/api/owner/profile`.
- **Access:** Owner only.
- **Rules:** One user has one owner profile; user ID comes from JWT; duplicate
  create returns conflict.

## Breed module

- **Purpose:** Supply reference breeds for pet creation.
- **Main table:** `breeds`.
- **Endpoints:** `GET /api/breeds`, optionally `?species=dog` or `?species=cat`.
- **Access:** Public.
- **Rules:** Dog/cat only; startup seeds 8 dog and 8 cat breeds idempotently;
  English breed name is available and Thai name is nullable.

## Pet module

- **Purpose:** Manage owner-controlled basic pet identity.
- **Main tables:** `pets`, with optional relation to `breeds` and required
  relation to `owner_profiles`.
- **Endpoints:** POST/GET `/api/pets`, GET/PATCH `/api/pets/:id`.
- **Access:** Owner only.
- **Rules:** Owner profile is resolved from JWT; owners see only their pets;
  supplied breed must exist and match species; PATCH changes supplied fields
  only.

## Clinic Profile module

- **Purpose:** Store clinic identity/settings for Clinic Web Dashboard foundation.
- **Main table:** `clinic_profiles`.
- **Endpoints:** POST/GET/PATCH `/api/clinic/profile`.
- **Access:** Canonical `clinic` role and legacy-compatible `clinic_staff`.
- **Rules:** One user has one clinic profile; user ID comes from JWT; owner is
  forbidden; duplicate create returns conflict.

## Shared infrastructure

- JWT authentication middleware injects user ID and role into Gin context.
- Role middleware checks route allow-lists.
- Response helper enforces one JSON envelope.
- Typed application errors keep internal causes out of responses.
- Startup SQL creates only currently implemented Sprint 1–6 schema.
