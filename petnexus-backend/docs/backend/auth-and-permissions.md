# Authentication and Permissions

## Password and token model

- Registration hashes passwords with bcrypt.
- Plain-text passwords and password hashes are never returned.
- Register and login return a JWT `accessToken`.
- JWT contains authenticated user ID, role, and expiry.
- Default expiry is controlled by `JWT_EXPIRES_IN`; deployed configuration can
  override it.
- There is no refresh-token or revocation/logout endpoint yet.

Protected requests require:

```http
Authorization: Bearer <accessToken>
```

Missing, malformed, invalid, or expired tokens return 401 `UNAUTHORIZED`.

## Roles

| Role | Current meaning |
| --- | --- |
| `owner` | Owner Mobile App account; owner profile and pet access |
| `clinic` | Canonical Clinic Web Dashboard account; clinic profile access |
| `clinic_staff` | Legacy-compatible clinic-side role; clinic profile access retained |
| `admin` | Reserved enum value; public registration forbidden |

Public registration currently accepts `owner`, `clinic`, and `clinic_staff`.
Owner behavior is unchanged by clinic role support.

## Public endpoints

- `GET /health`
- `GET /health/db`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/breeds`
- `GET /api/breeds?species=dog`
- `GET /api/breeds?species=cat`

## Authenticated endpoint

- `GET /api/me` — any valid authenticated role

## Owner-only endpoints

- `POST /api/owner/profile`
- `GET /api/owner/profile`
- `PATCH /api/owner/profile`
- `POST /api/pets`
- `GET /api/pets`
- `GET /api/pets/:id`
- `PATCH /api/pets/:id`
- `POST /api/owner/appointments`
- `GET /api/owner/appointments`
- `GET /api/owner/appointments/:id`
- `PATCH /api/owner/appointments/:id/cancel`

Clinic roles receive 403 `FORBIDDEN_ROLE` on these routes.

## Clinic-only endpoints

- `POST /api/clinic/profile`
- `GET /api/clinic/profile`
- `PATCH /api/clinic/profile`
- `GET /api/clinic/pet-lookup?pet_id=PNX-PET-XXXXXX`
- `GET /api/clinic/pet-lookup?owner_phone=<exact-phone>`
- `POST /api/clinic/appointments`
- `GET /api/clinic/appointments`
- `GET /api/clinic/appointments/:id`
- `PATCH /api/clinic/appointments/:id/status`
- `PATCH /api/clinic/appointments/:id/cancel`
- `GET /api/clinic/patients`
- `GET /api/clinic/patients/:petId`
- `POST /api/clinic/patients/:petId/medical-records`
- `GET /api/clinic/medical-records`
- `GET /api/clinic/medical-records/:recordId`
- `PATCH /api/clinic/medical-records/:recordId`

Allowed roles are canonical `clinic` and legacy-compatible `clinic_staff`.
Owner receives 403 `FORBIDDEN_ROLE`.

Clinic pet lookup accepts exactly one query parameter. It exposes only limited
pet identity, breed data, owner display name, and a masked phone number. It does
not grant clinic access to the pet or expose medical/private owner data.

Clinic patients are derived from non-cancelled appointments for the current
clinic profile. Cross-clinic or unrelated pet access returns 404.

Medical records are scoped by the JWT-derived clinic profile. The backend
derives `clinic_profile_id` and `created_by_user_id`; clients cannot choose
them. Cross-clinic records and unrelated pets return 404.

## Identity and ownership rules

- JWT middleware validates token and stores user ID/role in Gin context.
- Role middleware runs after authentication.
- Handlers read authenticated user ID from context.
- Request DTOs do not expose `user_id`.
- Pet DTOs do not expose `owner_profile_id`.
- Services resolve profiles and ownership through repositories.
- Another owner's pet returns 404 rather than exposing its existence.
- Owner and clinic appointment repositories always scope detail/update queries
  by the JWT-derived profile ID; cross-account appointments return 404.
- Clinic patient list/detail queries are scoped by the JWT-derived clinic
  profile ID; unrelated pets and other clinics' patients return 404.
- Medical record create/list/detail/update queries are scoped by the
  JWT-derived clinic profile ID; other clinics' records return 404.

## Common authorization-related status codes

| Status | Meaning |
| --- | --- |
| 400 | Invalid body, field, UUID, species, date, weight, or empty PATCH |
| 401 | Token missing, malformed, invalid, or expired |
| 403 | Authenticated role is not allowed |
| 404 | Profile/resource missing or pet ownership is blocked |
| 409 | Unique account/profile already exists |
| 422 | Auth registration/login input validation failed |

## Security rules for future work

- Never trust client-supplied ownership identifiers.
- Never log JWT, password, password hash, `DATABASE_URL`, or `JWT_SECRET`.
- New clinic/pet actions need explicit role and ownership decisions.
- QR/access/visit endpoints must not be added until authorization lifetime and
  audit requirements are designed.
