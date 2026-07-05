# API Reference

Base URLs:

```text
Local:  http://localhost:8080
Render: https://petnexus-api.onrender.com
```

Protected endpoints require:

```http
Authorization: Bearer <accessToken>
```

All responses use the envelope documented in [Architecture](./architecture.md).

## Shared response objects

### User

```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "phone": "0812345678",
  "role": "owner",
  "createdAt": "2026-07-05T08:00:00Z"
}
```

Auth responses place `{ "user": User, "accessToken": "..." }` in `data`.
`GET /api/me` places `{ "user": User }` in `data`. Auth uses `createdAt`
while other current DTOs use snake_case timestamps.

### Owner profile

```json
{
  "id": "owner-profile-uuid",
  "first_name": "Sunny",
  "last_name": "Example",
  "display_name": "Sunny Example",
  "gender": "male",
  "date_of_birth": "1995-05-20",
  "phone_number": "0812345678",
  "avatar_url": null,
  "address_line1": "123 Pet Street",
  "address_line2": null,
  "province": "Bangkok",
  "district": "Bang Rak",
  "subdistrict": "Si Phraya",
  "postal_code": "10500",
  "created_at": "2026-07-05T08:10:00Z",
  "updated_at": "2026-07-05T08:10:00Z"
}
```

### Breed

```json
{
  "id": "breed-uuid",
  "species": "dog",
  "name": "Golden Retriever",
  "name_th": null
}
```

### Pet

```json
{
  "id": "pet-uuid",
  "species": "dog",
  "name": "Milo",
  "gender": "male",
  "date_of_birth": "2022-05-10",
  "age_years": 4,
  "breed": {
    "id": "breed-uuid",
    "species": "dog",
    "name": "Golden Retriever",
    "name_th": null
  },
  "weight_kg": 12.5,
  "microchip_id": "MC-123456789",
  "avatar_url": null,
  "color": "Brown",
  "distinctive_marks": "White spot on chest",
  "is_neutered": true,
  "created_at": "2026-07-05T08:20:00Z",
  "updated_at": "2026-07-05T08:20:00Z"
}
```

### Clinic profile

```json
{
  "id": "clinic-profile-uuid",
  "clinic_name": "Happy Paws Clinic",
  "phone_number": "02-123-4567",
  "email": "contact@happypaws.example",
  "address": "123 Pet Street, Bangkok",
  "created_at": "2026-07-05T09:00:00Z",
  "updated_at": "2026-07-05T09:00:00Z"
}
```

## Health

### `GET /health`

- **Purpose:** Confirm Gin/backend process is online.
- **Auth/role:** Public; none.
- **Request body:** None.
- **Success:** 200; `data` is
  `{ "status": "ok", "service": "petnexus-backend" }`.
- **Common errors:** Network/service unavailable.
- **Backend notes:** Does not test PostgreSQL.

### `GET /health/db`

- **Purpose:** Confirm the running backend can ping PostgreSQL.
- **Auth/role:** Public; none.
- **Request body:** None.
- **Success:** 200; `data` is
  `{ "database": "postgresql", "status": "connected" }`.
- **Common errors:** 503 `DATABASE_UNAVAILABLE`.
- **Backend notes:** Uses the configured GORM connection and a timeout.

## Auth

### `POST /api/auth/register`

- **Purpose:** Create an account and return an access token.
- **Auth/role:** Public. Accepted public roles: `owner`, canonical `clinic`, and
  legacy-compatible `clinic_staff`; `admin` is rejected.
- **Request:**

```json
{
  "email": "owner@example.com",
  "phone": "0812345678",
  "password": "password123",
  "role": "owner"
}
```

- **Success:** 201; `data` is `{ "user": User, "accessToken": "<jwt>" }`.
- **Common errors:** 400 `INVALID_REQUEST`, 422 `VALIDATION_ERROR`, 403
  `FORBIDDEN_ROLE`, 409 `EMAIL_ALREADY_EXISTS`, 500 `INTERNAL_SERVER_ERROR`.
- **Backend notes:** Email is normalized; password length is 8–72 bytes and is
  stored only as a bcrypt hash.

### `POST /api/auth/login`

- **Purpose:** Validate credentials and return a new access token.
- **Auth/role:** Public; any existing role.
- **Request:**

```json
{
  "email": "owner@example.com",
  "password": "password123"
}
```

- **Success:** 200; `data` is `{ "user": User, "accessToken": "<jwt>" }`.
- **Common errors:** 400 `INVALID_REQUEST`, 422 `VALIDATION_ERROR`, 401
  `INVALID_CREDENTIALS`, 500.
- **Backend notes:** Unknown email and wrong password intentionally share one
  response.

### `GET /api/me`

- **Purpose:** Return the current authenticated account.
- **Auth/role:** JWT required; any authenticated role.
- **Request body:** None.
- **Success:** 200; `data` is `{ "user": User }`.
- **Common errors:** 401 `UNAUTHORIZED`, 404 `USER_NOT_FOUND`, 500.
- **Backend notes:** Useful for session/token verification.

## Owner Profile

All Owner Profile routes require JWT and role `owner`.

### `POST /api/owner/profile`

- **Purpose:** Create the current owner's single profile.
- **Request:**

```json
{
  "first_name": "Sunny",
  "last_name": "Example",
  "gender": "male",
  "date_of_birth": "1995-05-20",
  "phone_number": "0812345678",
  "avatar_url": "https://example.com/avatar.png",
  "address_line1": "123 Pet Street",
  "address_line2": "",
  "province": "Bangkok",
  "district": "Bang Rak",
  "subdistrict": "Si Phraya",
  "postal_code": "10500"
}
```

- **Success:** 201; `data` is an Owner profile object.
- **Common errors:** 400 `INVALID_REQUEST`/`VALIDATION_ERROR`, 401, 403, 409
  `OWNER_PROFILE_ALREADY_EXISTS`, 500.
- **Backend notes:** `user_id` is not accepted; backend uses JWT.

### `GET /api/owner/profile`

- **Purpose:** Fetch the current owner's profile.
- **Request body:** None.
- **Success:** 200; `data` is an Owner profile object.
- **Common errors:** 401, 403, 404 `OWNER_PROFILE_NOT_FOUND`, 500.
- **Backend notes:** A 404 is the normal state before owner onboarding.

### `PATCH /api/owner/profile`

- **Purpose:** Partially update the current owner's profile.
- **Request:**

```json
{
  "first_name": "Sunny Updated",
  "phone_number": "0899999999"
}
```

- **Success:** 200; `data` is the complete updated Owner profile.
- **Common errors:** 400 for malformed/empty/invalid body, 401, 403, 404, 500.
- **Backend notes:** Only supplied fields change; required names/phone cannot be
  changed to empty.

## Breeds

### `GET /api/breeds`

- **Purpose:** List all seeded dog and cat breeds.
- **Auth/role:** Public; none.
- **Request body:** None.
- **Success:** 200; `data` is an array of Breed objects.
- **Common errors:** 500.
- **Backend notes:** Current seed contains 8 dog and 8 cat breeds.

### `GET /api/breeds?species=dog`

- **Purpose:** List dog breeds only.
- **Auth/role:** Public; none.
- **Success:** 200; `data` is an array of dog Breed objects.
- **Common errors:** 400 `VALIDATION_ERROR` for unsupported species, 500.
- **Backend notes:** Query filter is optional and normalized.

### `GET /api/breeds?species=cat`

- **Purpose:** List cat breeds only.
- **Auth/role:** Public; none.
- **Success:** 200; `data` is an array of cat Breed objects.
- **Common errors:** 400 `VALIDATION_ERROR` for unsupported species, 500.
- **Backend notes:** Only `dog` and `cat` are supported.

## Pets

All Pet routes require JWT, role `owner`, and an existing owner profile.

### `POST /api/pets`

- **Purpose:** Create a pet for the current owner's profile.
- **Request:**

```json
{
  "species": "dog",
  "name": "Milo",
  "breed_id": "breed-uuid",
  "gender": "male",
  "date_of_birth": "2022-05-10",
  "weight_kg": 12.5,
  "microchip_id": "MC-123456789",
  "avatar_url": "https://example.com/milo.png",
  "color": "Brown",
  "distinctive_marks": "White spot on chest",
  "is_neutered": true
}
```

- **Success:** 201; `data` is a Pet object.
- **Common errors:** 400 `INVALID_REQUEST`, `VALIDATION_ERROR`, or
  `BREED_SPECIES_MISMATCH`; 401; 403; 404 `OWNER_PROFILE_REQUIRED` or
  `BREED_NOT_FOUND`; 500.
- **Backend notes:** `user_id` and `owner_profile_id` are not accepted. Breed is
  optional but must match species when supplied.

### `GET /api/pets`

- **Purpose:** List pets belonging to the current owner.
- **Request body:** None.
- **Success:** 200; `data` is an array of Pet objects; empty list is `[]`.
- **Common errors:** 401, 403, 404 `OWNER_PROFILE_REQUIRED`, 500.
- **Backend notes:** Repository lookup is scoped by resolved owner profile.

### `GET /api/pets/:id`

- **Purpose:** Fetch one current-owner pet.
- **Request body:** None; `:id` must be UUID.
- **Success:** 200; `data` is a Pet object.
- **Common errors:** 400 `INVALID_PET_ID`, 401, 403, 404
  `OWNER_PROFILE_REQUIRED`/`PET_NOT_FOUND`, 500.
- **Backend notes:** Another owner's pet also returns 404.

### `PATCH /api/pets/:id`

- **Purpose:** Partially update one current-owner pet.
- **Request:**

```json
{
  "name": "Milo Updated",
  "weight_kg": 13.2
}
```

- **Success:** 200; `data` is the complete updated Pet object.
- **Common errors:** 400 invalid UUID/body/validation/breed mismatch, 401, 403,
  404 owner profile/breed/pet not found, 500.
- **Backend notes:** Empty PATCH is rejected. An empty `breed_id` clears breed;
  changing species requires clearing or replacing an incompatible breed.

## Clinic Profile

All Clinic Profile routes require JWT and a clinic-side role: canonical
`clinic`, with `clinic_staff` retained for compatibility.

### `POST /api/clinic/profile`

- **Purpose:** Create the current clinic user's single profile.
- **Request:**

```json
{
  "clinic_name": "Happy Paws Clinic",
  "phone_number": "02-123-4567",
  "email": "contact@happypaws.example",
  "address": "123 Pet Street, Bangkok"
}
```

- **Success:** 201; `data` is a Clinic profile object.
- **Common errors:** 400 invalid body/validation, 401, 403, 409
  `CLINIC_PROFILE_ALREADY_EXISTS`, 500.
- **Backend notes:** `clinic_name` is required; `user_id` is not accepted.

### `GET /api/clinic/profile`

- **Purpose:** Fetch the current clinic user's profile.
- **Request body:** None.
- **Success:** 200; `data` is a Clinic profile object.
- **Common errors:** 401, 403, 404 `CLINIC_PROFILE_NOT_FOUND`, 500.
- **Backend notes:** Owner role receives 403.

### `PATCH /api/clinic/profile`

- **Purpose:** Partially update the current clinic user's profile.
- **Request:**

```json
{
  "clinic_name": "Happy Paws Bangkok",
  "phone_number": "02-999-9999"
}
```

- **Success:** 200; `data` is the complete updated Clinic profile.
- **Common errors:** 400 malformed/empty/invalid body, 401, 403, 404, 500.
- **Backend notes:** Optional phone/email/address can be cleared with empty
  string; unspecified fields remain unchanged.
