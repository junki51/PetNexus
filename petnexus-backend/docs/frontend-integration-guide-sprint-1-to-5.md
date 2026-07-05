# PetNexus Frontend Integration Guide — Sprint 1 to Sprint 5

Updated: 5 July 2026

This document describes the backend behavior that is implemented today. It is
intended for the Owner Mobile App team and as a boundary reference for a future
Clinic Web Dashboard. Request and response examples match the current Go code.

## 1. Overview

The Sprint 1–5 backend currently supports:

- Gin server startup and environment-based configuration
- standard JSON success/error envelopes
- application and PostgreSQL health checks
- owner or clinic-staff account registration and login
- bcrypt password hashing and JWT access tokens
- authentication and role middleware
- current-user lookup
- one owner profile per owner account
- public dog/cat breed catalog
- owner-only pet creation, list, detail, and partial update

Current owner app flow:

```text
Open app
→ Register or Login
→ Store accessToken
→ Send Authorization: Bearer <accessToken>
→ Fetch or create owner profile
→ Select dog/cat
→ Load breeds for selected species
→ Create pet
→ View/list/update pet
```

The backend does not yet support Pet Passport, QR sharing, clinic access,
medical visits, timeline, notifications, social login, or real file upload.
Those UI areas must remain mocked, disabled, or clearly marked as unavailable.

### Standard response envelope

Every successful endpoint uses this envelope:

```json
{
  "success": true,
  "message": "Human-readable message",
  "data": {}
}
```

Every error uses this envelope:

```json
{
  "success": false,
  "message": "Human-readable message",
  "error": {
    "code": "MACHINE_READABLE_CODE",
    "details": "Safe details for the client"
  }
}
```

Frontend code should branch on the HTTP status and `success`. Use
`error.code` for stable application decisions and `error.details` for a
developer-friendly or localized message mapping. Do not parse `message` to
make application decisions.

Important naming detail: Auth user responses currently use `createdAt`, while
Owner Profile and Pet responses use snake_case fields such as `created_at`.
Frontend models must match these real field names.

## 2. Base URLs

Local development:

```text
http://localhost:8080
```

Render deployment:

```text
<replace-with-render-backend-url>
```

Replace the placeholder with the real deployed Render web-service URL. Keep
the base URL in app environment/configuration; do not scatter or hardcode it
inside screens, repositories, or widgets.

Example route construction:

```text
{baseUrl}/api/auth/login
{baseUrl}/api/owner/profile
{baseUrl}/api/pets
```

## 3. Authentication Rule

Both register and login currently return an `accessToken`. The frontend must
store it securely and send it to every protected endpoint:

```http
Authorization: Bearer <accessToken>
```

Rules:

- Missing, malformed, invalid, or expired token returns HTTP 401 with code
  `UNAUTHORIZED`.
- A valid non-owner token on Owner Profile or Pet endpoints returns HTTP 403
  with code `FORBIDDEN_ROLE`.
- The default configured lifetime is 24 hours, but deployed configuration may
  override it. The client must rely on the 401 response rather than assuming a
  fixed lifetime.
- There is no refresh-token endpoint yet. On 401, clear the unusable token and
  return the user to login.
- Do not log tokens or include them in analytics/crash metadata.

Recommended storage:

- Flutter/mobile: secure storage backed by Keychain/Keystore.
- Web development: use the team's approved secure token strategy. Avoid
  exposing tokens through logs or URLs.

## 4. Role Rules

Current roles relevant to frontend integration:

| Role | Auth APIs | Owner Profile APIs | Pet APIs | Breed APIs |
| --- | --- | --- | --- | --- |
| `owner` | Yes | Yes | Yes | Yes (public) |
| `clinic` | Yes | No (403) | No (403) | Yes (public) |
| `clinic_staff` | Yes | No (403) | No (403) | Yes (public) |
| `admin` | Public registration blocked | No | No | Yes (public) |

Owner and pet ownership rules:

- Frontend must never send `user_id` or `owner_profile_id`.
- Backend reads the user ID from JWT claims.
- Backend finds the matching `owner_profiles` row.
- Every new pet is linked to that profile.
- An owner can list, fetch, or update only their own pets.
- Access to another owner's pet returns 404 to avoid revealing its existence.
- Pet creation/list/detail/update requires an existing owner profile; otherwise
  the backend returns 404 `OWNER_PROFILE_REQUIRED`.

Clinic note: `clinic` is the canonical Clinic Web Dashboard role;
`clinic_staff` remains accepted for compatibility. Sprint 6 adds Clinic Profile
settings, but patient lists, access requests, medical records, and other clinic
workflow APIs do not exist yet.

## 5. UI Label Mapping

Thai labels must map to backend values before requests are sent:

| UI label | Backend `species` |
| --- | --- |
| สุนัข | `dog` |
| แมว | `cat` |

Only these exact normalized values are supported:

```text
dog
cat
```

The backend trims whitespace and lowercases input, but frontend state should
still use canonical values. Invalid species returns HTTP 400
`VALIDATION_ERROR`.

Pet gender mapping:

| Suggested UI label | Backend value |
| --- | --- |
| เพศผู้ | `male` |
| เพศเมีย | `female` |
| ไม่ทราบ | `unknown` |

Owner-profile gender supports `male`, `female`, `prefer_not_to_say`, and
`other`; this differs from pet gender.

## 6. Sprint 1–2 — Health and Database Integration

Sprint 1 provides the Gin server, configuration loading, response helper, and
application health endpoint. Sprint 2 adds GORM/PostgreSQL, local Docker
Compose, database health, and Render `DATABASE_URL` support.

Health endpoints are useful for a developer/debug screen, CI smoke check, or
support diagnostics. The normal mobile app flow does not need to poll them.

### `GET /health`

1. **Purpose:** Verify that the backend process and Gin router are online.
2. **Auth:** Not required.
3. **Headers:** None required; `Accept: application/json` is optional.
4. **Request body:** None.
5. **Success:** HTTP 200.

```json
{
  "success": true,
  "message": "PetNexus backend is running",
  "data": {
    "service": "petnexus-backend",
    "status": "ok"
  }
}
```

6. **Common errors:** Network/DNS failure or service unavailable. The handler
   itself has no application-level validation errors.
7. **Frontend notes:** Use to distinguish “cannot reach backend” from an
   authentication or form error.

### `GET /health/db`

1. **Purpose:** Verify that the backend can reach PostgreSQL.
2. **Auth:** Not required.
3. **Headers:** None required.
4. **Request body:** None.
5. **Success:** HTTP 200.

```json
{
  "success": true,
  "message": "Database connection is healthy",
  "data": {
    "database": "postgresql",
    "status": "connected"
  }
}
```

6. **Common error:** HTTP 503.

```json
{
  "success": false,
  "message": "Database connection is unhealthy",
  "error": {
    "code": "DATABASE_UNAVAILABLE",
    "details": "Unable to reach PostgreSQL"
  }
}
```

7. **Frontend notes:** Useful for debugging only. Do not expose database
   internals or retry aggressively from every mobile screen.

Local PostgreSQL starts with:

```powershell
docker compose up -d
```

The local database is `petnexus` on port `5432`. Render uses `DATABASE_URL`;
frontend applications never receive or use database credentials.

## 7. Sprint 3 — Auth Integration

Passwords are hashed with bcrypt and never returned. JWT is the only current
session mechanism. Public registration accepts `owner`, `clinic`, and
legacy-compatible `clinic_staff`; the Owner Mobile App must always send
`owner`.

### `POST /api/auth/register`

1. **Purpose:** Create a new account and issue an access token.
2. **Auth:** Not required.
3. **Headers:** `Content-Type: application/json`.
4. **Request body:**

```json
{
  "email": "owner@example.com",
  "phone": "0812345678",
  "password": "password123",
  "role": "owner"
}
```

`email`, `password`, and `role` are required. `phone` is accepted and limited
to 30 characters but is not currently required by Auth registration.

5. **Success:** HTTP 201. Register already returns an access token.

```json
{
  "success": true,
  "message": "Registered successfully",
  "data": {
    "user": {
      "id": "87ddab6b-5250-4df3-9ca4-fc557280084c",
      "email": "owner@example.com",
      "phone": "0812345678",
      "role": "owner",
      "createdAt": "2026-07-05T08:00:00Z"
    },
    "accessToken": "<jwt-access-token>"
  }
}
```

6. **Common errors:**

| Status | Code | Meaning |
| --- | --- | --- |
| 400 | `INVALID_REQUEST` | Malformed/non-JSON body |
| 422 | `VALIDATION_ERROR` | Invalid email, missing field, password under 8 characters or over 72 bytes |
| 403 | `FORBIDDEN_ROLE` | Public registration attempted with unsupported role such as `admin` |
| 409 | `EMAIL_ALREADY_EXISTS` | Account already exists |
| 500 | `INTERNAL_SERVER_ERROR` | Unexpected server/database error |

7. **Frontend notes:**

- Confirm-password validation is frontend-only; do not send a
  `confirm_password` field.
- Terms/privacy checkbox is frontend-only because there is no
  `accepted_terms` field yet.
- Social registration is not supported.
- Register returns a usable token. Calling login immediately afterward is
  optional; if product flow requires a separate login step, it remains valid.
- Email is normalized to lowercase by the backend.

### `POST /api/auth/login`

1. **Purpose:** Authenticate credentials and issue a new access token.
2. **Auth:** Not required.
3. **Headers:** `Content-Type: application/json`.
4. **Request body:**

```json
{
  "email": "owner@example.com",
  "password": "password123"
}
```

5. **Success:** HTTP 200.

```json
{
  "success": true,
  "message": "Logged in successfully",
  "data": {
    "user": {
      "id": "87ddab6b-5250-4df3-9ca4-fc557280084c",
      "email": "owner@example.com",
      "phone": "0812345678",
      "role": "owner",
      "createdAt": "2026-07-05T08:00:00Z"
    },
    "accessToken": "<jwt-access-token>"
  }
}
```

6. **Common errors:** HTTP 400 `INVALID_REQUEST`, HTTP 422
   `VALIDATION_ERROR`, HTTP 401 `INVALID_CREDENTIALS`, and HTTP 500
   `INTERNAL_SERVER_ERROR`.
7. **Frontend notes:** Read the token from `data.accessToken`, store it
   securely, and never place it in query parameters. The backend deliberately
   uses the same invalid-credentials response for unknown email and wrong
   password.

### `GET /api/me`

1. **Purpose:** Fetch the current authenticated account and verify a stored
   token during app startup.
2. **Auth:** Required; any authenticated role is allowed.
3. **Headers:** `Authorization: Bearer <accessToken>`.
4. **Request body:** None.
5. **Success:** HTTP 200. Note the `data.user` nesting.

```json
{
  "success": true,
  "message": "Current user fetched successfully",
  "data": {
    "user": {
      "id": "87ddab6b-5250-4df3-9ca4-fc557280084c",
      "email": "owner@example.com",
      "phone": "0812345678",
      "role": "owner",
      "createdAt": "2026-07-05T08:00:00Z"
    }
  }
}
```

6. **Common errors:** HTTP 401 `UNAUTHORIZED`; HTTP 404 `USER_NOT_FOUND` if
   the token refers to a deleted account; HTTP 500 for unexpected failures.
7. **Frontend notes:** Call after restoring a token. On 401, clear session and
   navigate to login. On success, route by the returned role.

## 8. Sprint 4 — Owner Profile Integration

`users` contains account/auth data. `owner_profiles` contains editable owner
identity/contact data. One user can have only one owner profile. All three
profile endpoints require JWT and role `owner`; no request accepts `user_id`.

Supported fields:

| UI concept | Backend field | Support |
| --- | --- | --- |
| First name | `first_name` | Supported; required on create |
| Last name | `last_name` | Supported; required on create |
| Gender | `gender` | Supported |
| Age | — | Do not send; calculate from `date_of_birth` |
| Date of birth | `date_of_birth` | Supported as `YYYY-MM-DD` |
| Phone | `phone_number` | Supported; required on create |
| Profile image URL | `avatar_url` | URL reference supported |
| Real image upload | — | Not supported; mock/skip for now |
| Address | `address_line1`, `address_line2` | Supported |
| Province | `province` | Supported |
| District | `district` | Supported |
| Subdistrict | `subdistrict` | Supported |
| Postal code | `postal_code` | Supported |

Allowed owner genders: `male`, `female`, `prefer_not_to_say`, `other`.

### Owner Profile response shape

Create, fetch, and update return this object directly in `data`:

```json
{
  "id": "118c08c9-c80f-4354-bdaa-a17f95fc69ae",
  "first_name": "Sunny",
  "last_name": "Example",
  "display_name": "Sunny Example",
  "gender": "male",
  "date_of_birth": "1995-05-20",
  "phone_number": "0812345678",
  "avatar_url": "https://example.com/avatar.png",
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

`display_name` is computed by the backend and is not accepted in requests.
Nullable fields may be `null`.

### `POST /api/owner/profile`

1. **Purpose:** Create the current owner's one profile.
2. **Auth:** JWT required; role `owner` only.
3. **Headers:** `Content-Type: application/json` and
   `Authorization: Bearer <accessToken>`.
4. **Request body:**

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

5. **Success:** HTTP 201.

```json
{
  "success": true,
  "message": "Owner profile created successfully",
  "data": {
    "id": "118c08c9-c80f-4354-bdaa-a17f95fc69ae",
    "first_name": "Sunny",
    "last_name": "Example",
    "display_name": "Sunny Example",
    "gender": "male",
    "date_of_birth": "1995-05-20",
    "phone_number": "0812345678",
    "avatar_url": "https://example.com/avatar.png",
    "address_line1": "123 Pet Street",
    "address_line2": null,
    "province": "Bangkok",
    "district": "Bang Rak",
    "subdistrict": "Si Phraya",
    "postal_code": "10500",
    "created_at": "2026-07-05T08:10:00Z",
    "updated_at": "2026-07-05T08:10:00Z"
  }
}
```

6. **Common errors:** HTTP 400 `INVALID_REQUEST` or `VALIDATION_ERROR`, HTTP
   401 `UNAUTHORIZED`, HTTP 403 `FORBIDDEN_ROLE`, HTTP 409
   `OWNER_PROFILE_ALREADY_EXISTS`, and HTTP 500.
7. **Frontend notes:** Never send `user_id`. If GET already returned a profile,
   use PATCH instead of POST. Empty optional strings are normalized to null.

### `GET /api/owner/profile`

1. **Purpose:** Fetch the profile belonging to the current JWT user.
2. **Auth:** JWT required; role `owner` only.
3. **Headers:** `Authorization: Bearer <accessToken>`.
4. **Request body:** None.
5. **Success:** HTTP 200 with the Owner Profile response shape above and
   message `Owner profile fetched successfully`.
6. **Common errors:** HTTP 401 `UNAUTHORIZED`, HTTP 403 `FORBIDDEN_ROLE`, HTTP
   404 `OWNER_PROFILE_NOT_FOUND`, and HTTP 500.
7. **Frontend notes:** Use this after `/api/me`. A 404 is an expected onboarding
   state: show the owner-profile setup screen rather than a generic fatal error.

### `PATCH /api/owner/profile`

1. **Purpose:** Partially update the current owner's profile.
2. **Auth:** JWT required; role `owner` only.
3. **Headers:** `Content-Type: application/json` and
   `Authorization: Bearer <accessToken>`.
4. **Request body:** All fields are optional, but at least one must be present.

```json
{
  "first_name": "Sunny Updated",
  "phone_number": "0899999999"
}
```

5. **Success:** HTTP 200 with the complete updated Owner Profile object and
   message `Owner profile updated successfully`.
6. **Common errors:** HTTP 400 `INVALID_REQUEST` or `VALIDATION_ERROR` (including
   `{}`), HTTP 401, HTTP 403, HTTP 404 `OWNER_PROFILE_NOT_FOUND`, and HTTP 500.
7. **Frontend notes:** Send only dirty fields. Unspecified fields remain
   unchanged. Required fields cannot be changed to empty. Optional string/date
   fields can be cleared with an empty string.

## 9. Sprint 5 — Breed Integration

The `breeds` table is reference data for the Select Pet and Fill Pet Profile
screens. The backend currently seeds 8 dog breeds and 8 cat breeds. Breed
listing is public because it contains no owner data.

### `GET /api/breeds`

1. **Purpose:** Fetch all available dog and cat breeds.
2. **Auth:** Not required.
3. **Headers:** None required.
4. **Request body:** None.
5. **Success:** HTTP 200. `data` is an array, not `{ "breeds": [...] }`.

```json
{
  "success": true,
  "message": "Breeds fetched successfully",
  "data": [
    {
      "id": "3f14345f-6553-40f1-98ec-d056b6855192",
      "species": "dog",
      "name": "Golden Retriever",
      "name_th": null
    },
    {
      "id": "657fbf27-0d15-4f60-90b2-a552bdb32ddd",
      "species": "cat",
      "name": "Siamese",
      "name_th": null
    }
  ]
}
```

6. **Common errors:** HTTP 500 `INTERNAL_SERVER_ERROR` if breed lookup fails.
7. **Frontend notes:** Cache for the current session if useful. `name_th` is
   currently null in seed data, so show `name` as the fallback label.

### `GET /api/breeds?species=dog`

1. **Purpose:** Fetch dog breeds after the user selects สุนัข.
2. **Auth:** Not required.
3. **Headers:** None required.
4. **Request body:** None; query parameter is `species=dog`.
5. **Success:** HTTP 200 with the same envelope and only dog items.

```json
{
  "success": true,
  "message": "Breeds fetched successfully",
  "data": [
    {
      "id": "3f14345f-6553-40f1-98ec-d056b6855192",
      "species": "dog",
      "name": "Golden Retriever",
      "name_th": null
    }
  ]
}
```

6. **Common errors:** HTTP 400 `VALIDATION_ERROR` only if an unsupported value
   is used; HTTP 500 for unexpected database failure.
7. **Frontend notes:** Clear any previously selected cat `breed_id` when the
   species changes to dog.

### `GET /api/breeds?species=cat`

1. **Purpose:** Fetch cat breeds after the user selects แมว.
2. **Auth:** Not required.
3. **Headers:** None required.
4. **Request body:** None; query parameter is `species=cat`.
5. **Success:** HTTP 200 with the same envelope and only cat items.

```json
{
  "success": true,
  "message": "Breeds fetched successfully",
  "data": [
    {
      "id": "657fbf27-0d15-4f60-90b2-a552bdb32ddd",
      "species": "cat",
      "name": "Siamese",
      "name_th": null
    }
  ]
}
```

6. **Common errors:** HTTP 400 `VALIDATION_ERROR` for values other than dog or
   cat; HTTP 500 for unexpected database failure.
7. **Frontend notes:** Clear any previously selected dog `breed_id` when the
   species changes to cat.

Invalid query example:

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "species query must be dog or cat"
  }
}
```

## 10. Sprint 5 — Pet Integration

All Pet endpoints require JWT, role `owner`, and an existing owner profile.
Frontend must never send `user_id` or `owner_profile_id`; neither field exists
in the request DTOs. Backend resolves ownership from JWT.

Supported create fields:

| Field | Required | Rules |
| --- | --- | --- |
| `species` | Yes | `dog` or `cat` |
| `name` | Yes | Non-empty, maximum 100 characters |
| `breed_id` | No | UUID; breed must exist and match species |
| `gender` | No | `male`, `female`, `unknown` |
| `date_of_birth` | No | `YYYY-MM-DD`, not in future |
| `weight_kg` | No | Greater than 0 and at most 200 |
| `microchip_id` | No | Maximum 100 characters; not currently unique |
| `avatar_url` | No | Valid HTTP/HTTPS URL reference |
| `color` | No | Maximum 100 characters |
| `distinctive_marks` | No | Maximum 1000 characters |
| `is_neutered` | No | Boolean |

Real pet image upload is not implemented. `avatar_url` is only a URL/reference
field and can point to a placeholder or a separately hosted asset.

### Pet response shape

Create/detail/update return this object directly in `data`; list returns an
array of these objects:

```json
{
  "id": "b00d17fd-5467-42fc-b9a6-a464460fd37c",
  "species": "dog",
  "name": "Milo",
  "gender": "male",
  "date_of_birth": "2022-05-10",
  "age_years": 4,
  "breed": {
    "id": "3f14345f-6553-40f1-98ec-d056b6855192",
    "species": "dog",
    "name": "Golden Retriever",
    "name_th": null
  },
  "weight_kg": 12.5,
  "microchip_id": "MC-123456789",
  "avatar_url": "https://example.com/milo.png",
  "color": "Brown",
  "distinctive_marks": "White spot on chest",
  "is_neutered": true,
  "created_at": "2026-07-05T08:20:00Z",
  "updated_at": "2026-07-05T08:20:00Z"
}
```

`age_years` is calculated from `date_of_birth`; do not store it as editable
frontend source data. `breed` and other optional values can be null.

### `POST /api/pets`

1. **Purpose:** Create a basic pet profile for the current owner.
2. **Auth:** JWT required; role `owner` only; owner profile required.
3. **Headers:** `Content-Type: application/json` and
   `Authorization: Bearer <accessToken>`.
4. **Request body:**

```json
{
  "species": "dog",
  "name": "Milo",
  "breed_id": "3f14345f-6553-40f1-98ec-d056b6855192",
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

5. **Success:** HTTP 201 with message `Pet created successfully` and the Pet
   response object in `data`.
6. **Common errors:**

| Status | Code | Meaning |
| --- | --- | --- |
| 400 | `INVALID_REQUEST` | Malformed JSON |
| 400 | `VALIDATION_ERROR` | Invalid/missing field, date, URL, UUID, weight, or gender |
| 400 | `BREED_SPECIES_MISMATCH` | Breed and pet species differ |
| 401 | `UNAUTHORIZED` | Token missing/invalid/expired |
| 403 | `FORBIDDEN_ROLE` | Authenticated user is not owner |
| 404 | `OWNER_PROFILE_REQUIRED` | Owner onboarding profile is missing |
| 404 | `BREED_NOT_FOUND` | Supplied breed does not exist |
| 500 | `INTERNAL_SERVER_ERROR` | Unexpected backend/database error |

7. **Frontend notes:** Load breeds after species selection. `breed_id` is
   optional for MVP. After 201, navigate to success/detail and use `data.id` as
   the pet route ID.

### `GET /api/pets`

1. **Purpose:** List pets belonging to the current owner profile.
2. **Auth:** JWT required; role `owner` only; owner profile required.
3. **Headers:** `Authorization: Bearer <accessToken>`.
4. **Request body:** None.
5. **Success:** HTTP 200. Empty owners receive `data: []`.

```json
{
  "success": true,
  "message": "Pets fetched successfully",
  "data": [
    {
      "id": "b00d17fd-5467-42fc-b9a6-a464460fd37c",
      "species": "dog",
      "name": "Milo",
      "gender": "male",
      "date_of_birth": "2022-05-10",
      "age_years": 4,
      "breed": {
        "id": "3f14345f-6553-40f1-98ec-d056b6855192",
        "species": "dog",
        "name": "Golden Retriever",
        "name_th": null
      },
      "weight_kg": 12.5,
      "microchip_id": "MC-123456789",
      "avatar_url": "https://example.com/milo.png",
      "color": "Brown",
      "distinctive_marks": "White spot on chest",
      "is_neutered": true,
      "created_at": "2026-07-05T08:20:00Z",
      "updated_at": "2026-07-05T08:20:00Z"
    }
  ]
}
```

6. **Common errors:** HTTP 401, HTTP 403, HTTP 404
   `OWNER_PROFILE_REQUIRED`, and HTTP 500.
7. **Frontend notes:** Treat an empty array as a normal “no pets yet” state and
   show the create-pet CTA.

### `GET /api/pets/:id`

1. **Purpose:** Fetch one pet belonging to the current owner.
2. **Auth:** JWT required; role `owner` only; owner profile required.
3. **Headers:** `Authorization: Bearer <accessToken>`.
4. **Request body:** None. Replace `:id` with the pet UUID.
5. **Success:** HTTP 200 with message `Pet fetched successfully` and the Pet
   response object in `data`.
6. **Common errors:** HTTP 400 `INVALID_PET_ID`, HTTP 401, HTTP 403, HTTP 404
   `OWNER_PROFILE_REQUIRED`, HTTP 404 `PET_NOT_FOUND`, and HTTP 500. Another
   owner's valid pet ID also returns `PET_NOT_FOUND`.
7. **Frontend notes:** Do not infer ownership from cached data; always handle
   404. Use the returned complete object to refresh the detail screen.

### `PATCH /api/pets/:id`

1. **Purpose:** Partially update one pet belonging to the current owner.
2. **Auth:** JWT required; role `owner` only; owner profile required.
3. **Headers:** `Content-Type: application/json` and
   `Authorization: Bearer <accessToken>`.
4. **Request body:** Send only fields that changed; body cannot be empty.

```json
{
  "name": "Milo Updated",
  "weight_kg": 13.2
}
```

5. **Success:** HTTP 200 with message `Pet updated successfully` and the full
   updated Pet response in `data`.
6. **Common errors:** HTTP 400 `INVALID_PET_ID`, `INVALID_REQUEST`,
   `VALIDATION_ERROR`, or `BREED_SPECIES_MISMATCH`; HTTP 401; HTTP 403; HTTP
   404 `OWNER_PROFILE_REQUIRED`, `BREED_NOT_FOUND`, or `PET_NOT_FOUND`; HTTP
   500.
7. **Frontend notes:**

- Unspecified fields are not overwritten.
- If species changes, clear or replace an incompatible breed in the same
  request.
- Sending `"breed_id": ""` clears breed.
- Empty strings clear nullable string/date fields.
- `false` is a valid `is_neutered` update and must not be omitted by frontend
  serialization.
- The current DTO does not provide an explicit way to clear `weight_kg` or
  `is_neutered` back to database null. Do not invent a client convention.

## 11. Recommended Frontend Flow

### A. First app launch

1. Show the open/welcome screen.
2. User chooses Login or Register.
3. For registration, call `POST /api/auth/register` with role `owner`.
4. Save the returned token, or call `POST /api/auth/login` if the product
   intentionally requires a separate login step.
5. Store `data.accessToken` securely.
6. Call `GET /api/me` to verify restored/current token.
7. If `/api/me` returns 401, clear session and show Login.

### B. Owner profile

1. Call `GET /api/owner/profile`.
2. On 200, continue to the pet/home flow.
3. On 404 `OWNER_PROFILE_NOT_FOUND`, show profile setup.
4. Submit setup with `POST /api/owner/profile`.
5. Use `PATCH /api/owner/profile` for later edits.
6. If POST returns 409, another request already created the profile; refetch it.

### C. Pet setup

1. User selects สุนัข or แมว.
2. Map to `dog` or `cat`.
3. Call `GET /api/breeds?species=dog` or `?species=cat`.
4. Show returned breeds; display `name_th ?? name`.
5. Keep `breed_id` optional and clear it when species changes.
6. Submit with `POST /api/pets`.
7. On 201, navigate to the done/detail screen using `data.id`.
8. Fetch `GET /api/pets/:id` when a fresh server representation is needed.

### D. Pet list and detail

1. Call `GET /api/pets` for the owner's pet cards.
2. Show an empty state when `data` is `[]`.
3. Call `GET /api/pets/:id` for pet detail or Pet ID screen.
4. Call `PATCH /api/pets/:id` to update basic pet data.
5. Pet Passport/QR/medical/timeline sections must use mock data or disabled
   states because those APIs do not exist yet.

## 12. Error Handling Guide

| HTTP status | Current meaning | Recommended frontend behavior |
| --- | --- | --- |
| 200 OK | Fetch/list/login/update success | Parse `data`, render success state |
| 201 Created | Register/profile/pet created | Save returned object/token and continue |
| 400 Bad Request | Malformed JSON, invalid species/gender/UUID/date/weight, empty PATCH, breed mismatch | Keep form state; map `error.code/details` to field or form message |
| 401 Unauthorized | Missing, malformed, invalid, or expired token; invalid login credentials | For protected APIs clear session and open Login; for login show credential error |
| 403 Forbidden | Wrong role or disallowed registration role | Show permission message; do not retry automatically |
| 404 Not Found | User/profile/breed/pet missing; cross-owner pet hidden | Handle onboarding 404 separately; otherwise show not-found state |
| 409 Conflict | Email/profile already exists | Offer login or refetch existing profile |
| 422 Unprocessable Entity | Auth registration/login validation failure | Show validation details without clearing user input |
| 500 Internal Server Error | Unexpected backend/database error | Show retryable generic error; retain form data |
| 503 Service Unavailable | Database health check failed | Show backend/database unavailable in diagnostics |

Example validation error:

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "weight_kg must be greater than 0 and at most 200"
  }
}
```

Suggested client strategy:

- Decode the error envelope even for non-2xx responses.
- Centralize 401 handling in the API client/interceptor.
- Do not convert every 404 to one generic screen: owner-profile 404 means
  onboarding, while pet 404 means missing/not-owned.
- Preserve submitted form values after 400/422/500.
- Show friendly localized messages, but retain `error.code` in non-sensitive
  debug context.
- Use bounded retry/backoff for network or 500 errors. Never retry validation,
  401, 403, 404, or 409 in a loop.

## 13. Currently Mocked / Not Implemented Yet

The frontend must not call invented endpoints for these features:

- Google login
- Apple login
- Facebook login
- Forgot password
- Refresh token and logout/revocation
- Email/phone verification
- Real owner profile image upload
- Real pet image upload
- QR sharing
- Pet Passport backend
- Clinic access request
- Owner approval/rejection for clinic access
- Clinic profile
- Clinic patient list
- Medical records and verified visits
- Timeline
- Calendar
- Notifications
- Reports and analytics
- Export PDF
- Data backup
- Pet deletion
- Owner-profile deletion

For mock-only screens, use a clearly separated mock repository/data source so
it cannot accidentally be mistaken for a live backend integration. Disable
destructive or trust-sensitive actions that have no backend enforcement.

## 14. Frontend Integration Checklist

- [ ] Backend local server is running.
- [ ] PostgreSQL/Docker database is running.
- [ ] `GET /health` works.
- [ ] `GET /health/db` works.
- [ ] Register owner works.
- [ ] Login owner works.
- [ ] `accessToken` is stored securely.
- [ ] Authorization header is sent correctly.
- [ ] `GET /api/me` works.
- [ ] Invalid/expired token clears session and redirects to Login.
- [ ] Create owner profile works.
- [ ] Get owner profile works.
- [ ] Owner-profile 404 opens onboarding.
- [ ] Patch owner profile works.
- [ ] Select dog/cat maps correctly to `dog`/`cat`.
- [ ] `GET /api/breeds` works.
- [ ] `GET /api/breeds?species=dog` works.
- [ ] `GET /api/breeds?species=cat` works.
- [ ] Breed dropdown falls back from `name_th` to `name`.
- [ ] Create pet works.
- [ ] List my pets works, including empty state.
- [ ] Get pet detail works.
- [ ] Patch pet updates only changed fields.
- [ ] Breed selection is cleared/reloaded when species changes.
- [ ] 400/422 errors are shown clearly without losing form input.
- [ ] 401 redirects user to Login.
- [ ] 403 shows a permission error.
- [ ] 404 shows the correct onboarding/not-found state.
- [ ] Mock-only features do not call nonexistent backend endpoints.
- [ ] Tokens and sensitive data are absent from logs and analytics.

## 15. PowerShell Test Examples

Set the local base URL:

```powershell
$baseUrl = "http://localhost:8080"
```

### Health

```powershell
Invoke-RestMethod -Method GET "$baseUrl/health"
Invoke-RestMethod -Method GET "$baseUrl/health/db"
```

### Register an owner

Use a unique email when repeating the test:

```powershell
$email = "owner.$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())@example.com"
$password = "password123"
$registerBody = @{
  email = $email
  phone = "0812345678"
  password = $password
  role = "owner"
} | ConvertTo-Json

$register = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/register" `
  -ContentType "application/json" `
  -Body $registerBody
```

### Login and save token

```powershell
$loginBody = @{
  email = $email
  password = $password
} | ConvertTo-Json

$login = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/auth/login" `
  -ContentType "application/json" `
  -Body $loginBody

$token = $login.data.accessToken
$headers = @{ Authorization = "Bearer $token" }
```

### Protected request

```powershell
Invoke-RestMethod -Method GET `
  -Uri "$baseUrl/api/me" `
  -Headers $headers
```

### Create owner profile

```powershell
$profileBody = @{
  first_name = "Sunny"
  last_name = "Example"
  gender = "male"
  date_of_birth = "1995-05-20"
  phone_number = "0812345678"
  province = "Bangkok"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/owner/profile" `
  -ContentType "application/json" `
  -Headers $headers `
  -Body $profileBody
```

### Breeds

```powershell
$breeds = Invoke-RestMethod -Method GET "$baseUrl/api/breeds"
$dogs = Invoke-RestMethod -Method GET "$baseUrl/api/breeds?species=dog"
$cats = Invoke-RestMethod -Method GET "$baseUrl/api/breeds?species=cat"
$breedId = ($dogs.data | Select-Object -First 1).id
```

### Create pet

```powershell
$petBody = @{
  species = "dog"
  name = "Milo"
  breed_id = $breedId
  gender = "male"
  date_of_birth = "2022-05-10"
  weight_kg = 12.5
  microchip_id = "MC-123456789"
  avatar_url = "https://example.com/milo.png"
  color = "Brown"
  distinctive_marks = "White spot on chest"
  is_neutered = $true
} | ConvertTo-Json

$pet = Invoke-RestMethod -Method POST `
  -Uri "$baseUrl/api/pets" `
  -ContentType "application/json" `
  -Headers $headers `
  -Body $petBody
```

### List, detail, and update pet

```powershell
Invoke-RestMethod -Method GET -Uri "$baseUrl/api/pets" -Headers $headers
Invoke-RestMethod -Method GET -Uri "$baseUrl/api/pets/$($pet.data.id)" -Headers $headers

$petPatch = @{ name = "Milo Updated"; weight_kg = 13.2 } | ConvertTo-Json
Invoke-RestMethod -Method PATCH `
  -Uri "$baseUrl/api/pets/$($pet.data.id)" `
  -ContentType "application/json" `
  -Headers $headers `
  -Body $petPatch
```

For Render testing, change only `$baseUrl`; never paste the production token
or database URL into committed files.

## 16. Notes for Mobile App Team

- Do not hardcode `localhost` for real-device builds.
- Android emulator commonly reaches the host machine at `10.0.2.2`, so the
  local URL may be `http://10.0.2.2:8080` depending on network setup.
- iOS Simulator may resolve host localhost, but verify the team's environment.
- A physical device must use the development computer's LAN IP and both devices
  must be reachable on the same network, or use the Render URL.
- Local HTTP may require Android/iOS development network-security settings.
- Store base URL in environment/flavor/configuration.
- Centralize JSON envelope parsing and Authorization header injection.
- Always implement loading, empty, success, offline, timeout, and error states.
- Cancel or ignore stale breed requests if the user quickly switches species.
- Never assume mock-only features have live endpoints.
- Never store a plain-text password after the login/register request finishes.
- Never expose JWT, database URL, or backend secrets in screenshots/logs.

## 17. Final Summary

Backend Sprint 1–5 currently supports:

- server and database health
- owner registration/login
- JWT authentication and role enforcement
- current-user lookup
- owner profile create/fetch/update
- public dog/cat breed list and filtering
- owner pet create/list/detail/update

Frontend should integrate only these real backend capabilities today. Other UI
features must stay mocked, disabled, or marked “coming soon” until future
backend sprints provide the required data model, permission rules, and APIs.

Canonical integration path:

```text
Auth → Owner Profile → Species → Breeds → Create Pet → List/Detail/Patch Pet
```
