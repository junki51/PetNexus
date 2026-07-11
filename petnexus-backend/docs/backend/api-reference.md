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
  "public_pet_id": "PNX-PET-A1B2C3",
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
  optional but must match species when supplied. `public_pet_id` is generated
  by the backend and cannot be supplied by the client.

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

## Clinic Pet Lookup

Clinic Pet Lookup requires JWT and a clinic-side role: canonical `clinic`, with
`clinic_staff` retained for compatibility. It does not create clinic access or
authorization.

### `GET /api/clinic/pet-lookup?pet_id=PNX-PET-A1B2C3`

- **Purpose:** Find one pet by its backend-generated public pet ID.
- **Success:** 200; `data` is one limited pet lookup object.
- **Common errors:** 400 `VALIDATION_ERROR`, 401, 403, 404 `PET_NOT_FOUND`, 500.
- **Backend notes:** Public ID input is normalized to uppercase.

### `GET /api/clinic/pet-lookup?owner_phone=0812345678`

- **Purpose:** Find pets whose owner profile phone exactly matches the query.
- **Success:** 200; `data` is `{ "items": [...] }`; no matches return
  `{ "items": [] }`.
- **Common errors:** 400 `VALIDATION_ERROR`, 401, 403, 500.
- **Backend notes:** Partial phone matching is not supported. Supply exactly one
  of `pet_id` and `owner_phone`; supplying neither or both returns 400.

Limited lookup item:

    {
      "id": "pet-uuid",
      "public_pet_id": "PNX-PET-A1B2C3",
      "name": "Milo",
      "species": "dog",
      "breed": {
        "id": "breed-uuid",
        "species": "dog",
        "name": "Golden Retriever",
        "name_th": null
      },
      "gender": "male",
      "date_of_birth": "2022-05-10",
      "avatar_url": null,
      "owner": {
        "display_name": "Sunny Example",
        "masked_phone": "081****678"
      }
    }

The response deliberately excludes owner address, email/user ID, pet
microchip, weight, distinctive marks, and all medical data.

## Owner Appointments

All owner appointment routes require JWT role `owner` and an existing owner
profile. The backend derives owner ownership from JWT.

### `POST /api/owner/appointments`

Creates a `requested` appointment for one pet owned by the current owner.

```json
{
  "clinic_profile_id": "clinic-profile-uuid",
  "pet_id": "pet-uuid",
  "title": "Annual checkup",
  "appointment_type": "checkup",
  "scheduled_at": "2026-07-10T10:00:00+07:00",
  "duration_minutes": 30,
  "note": "Bring vaccine card"
}
```

- **Success:** 201.
- **Errors:** 400 invalid UUID/type/time/duration/body; 401; 403; 404 owner
  profile, clinic, or owned pet not found; 500.

### `GET /api/owner/appointments`

Returns the current owner's appointments sorted by scheduled time ascending.
Optional filters: `date_from=YYYY-MM-DD`, `date_to=YYYY-MM-DD`, and
`status`.

### `GET /api/owner/appointments/:id`

Returns one current-owner appointment. Another owner's appointment returns 404.

### `PATCH /api/owner/appointments/:id/cancel`

Cancels one current-owner appointment and sets `cancelled_at`. Repeating the
request is idempotent.

## Clinic Appointments

Clinic appointment routes require JWT role `clinic`; legacy `clinic_staff`
remains compatible. The backend derives `clinic_profile_id` from JWT.

### `POST /api/clinic/appointments`

Creates a `scheduled` appointment. Supply exactly one pet lookup field:

```json
{
  "public_pet_id": "PNX-PET-A1B2C3",
  "appointment_type": "vaccination",
  "scheduled_at": "2026-07-11T09:30:00+07:00",
  "duration_minutes": 45,
  "note": "Booster dose"
}
```

`pet_id` may be used instead. Both or neither pet field returns 400.

### `GET /api/clinic/appointments`

Returns all appointments belonging to the current clinic, sorted by
`scheduled_at` ascending. Optional filters:

- `date=YYYY-MM-DD`
- `date_from=YYYY-MM-DD&date_to=YYYY-MM-DD`
- `status=requested|scheduled|checked_in|completed|cancelled`
- `appointment_type=<allowed-type>`

Calendar dates use UTC. `date_to` includes its complete day. Exact `date`
cannot be combined with range fields. With no date filter, all clinic
appointments are returned.

### `GET /api/clinic/appointments/:id`

Returns one appointment under the current clinic profile. Another clinic's
appointment returns 404.

### `PATCH /api/clinic/appointments/:id/status`

```json
{
  "status": "checked_in"
}
```

All five appointment statuses are accepted. This sprint intentionally has no
complex transition engine. Setting a non-cancelled status clears
`cancelled_at`.

### `PATCH /api/clinic/appointments/:id/cancel`

Cancels one current-clinic appointment and is idempotent.

### Appointment response

```json
{
  "id": "appointment-uuid",
  "title": "Annual checkup",
  "appointment_type": "checkup",
  "scheduled_at": "2026-07-10T03:00:00Z",
  "duration_minutes": 30,
  "status": "requested",
  "note": "Bring vaccine card",
  "created_by_role": "owner",
  "cancelled_at": null,
  "created_at": "2026-07-08T08:00:00Z",
  "updated_at": "2026-07-08T08:00:00Z",
  "pet": {
    "id": "pet-uuid",
    "public_pet_id": "PNX-PET-A1B2C3",
    "name": "Milo",
    "species": "dog",
    "avatar_url": null,
    "breed": null
  },
  "owner": {
    "display_name": "Sunny Example",
    "masked_phone": "081****678"
  },
  "clinic": {
    "id": "clinic-profile-uuid",
    "clinic_name": "Happy Paws Clinic",
    "phone_number": "02-123-4567",
    "email": "clinic@example.com"
  }
}
```

Responses omit internal owner/clinic profile ownership IDs, creator user ID,
password data, JWT claims, private owner profile data, and medical data.

## Clinic Patients

Clinic patient routes require JWT and a clinic-side role: canonical `clinic`,
with legacy `clinic_staff` retained for compatibility.

For Sprint 9, a clinic patient is a unique pet that has at least one
non-cancelled appointment with the authenticated clinic profile. No separate
`patients` table exists.

### `GET /api/clinic/patients`

- **Purpose:** List patients for the current clinic's Patients page.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Success:** 200; `data` is an array of clinic patient list items.
- **Common errors:** 400 invalid query, 401, 403, 404
  `CLINIC_PROFILE_REQUIRED`, 500.
- **Backend notes:** The backend resolves `clinic_profile_id` from JWT and
  scopes the list by appointments owned by that clinic.

Supported query parameters:

| Query | Description |
| --- | --- |
| `q` | Optional search by pet name or `public_pet_id`; max 100 characters |
| `species` | Optional `dog` or `cat` |
| `status` | Optional latest computed non-cancelled appointment status |
| `limit` | Optional positive integer; default 20, max 100 |
| `offset` | Optional non-negative integer; default 0 |
| `sort` | Optional sort key; default `latest_appointment_desc` |

Supported `sort` values:

```text
latest_appointment_desc
latest_appointment_asc
name_asc
name_desc
next_appointment_asc
```

Example list item:

```json
{
  "pet": {
    "id": "pet-uuid",
    "public_pet_id": "PNX-PET-A1B2C3",
    "name": "Milo",
    "species": "dog",
    "avatar_url": null,
    "breed": {
      "id": "breed-uuid",
      "species": "dog",
      "name": "Golden Retriever",
      "name_th": null
    }
  },
  "owner": {
    "display_name": "Sunny Example",
    "masked_phone": "081****678"
  },
  "appointment_summary": {
    "total_appointments": 2,
    "last_appointment_at": "2026-07-12T03:00:00Z",
    "next_appointment_at": "2026-07-15T03:00:00Z",
    "latest_status": "scheduled"
  },
  "first_seen_at": "2026-07-01T03:00:00Z"
}
```

### `GET /api/clinic/patients/:petId`

- **Purpose:** Fetch one patient detail summary for the current clinic.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Path:** `:petId` must be a UUID.
- **Success:** 200; `data` is a clinic patient detail object.
- **Common errors:** 400 `INVALID_PET_ID`, 401, 403, 404
  `CLINIC_PROFILE_REQUIRED`/`CLINIC_PATIENT_NOT_FOUND`, 500.
- **Backend notes:** A pet from another clinic or a pet with no non-cancelled
  appointment relationship to this clinic returns 404.

Example detail data:

```json
{
  "pet": {
    "id": "pet-uuid",
    "public_pet_id": "PNX-PET-A1B2C3",
    "name": "Milo",
    "species": "dog",
    "gender": "male",
    "date_of_birth": "2022-05-10",
    "weight_kg": 12.5,
    "microchip_id": "MC-123456789",
    "avatar_url": null,
    "color": "Brown",
    "distinctive_marks": "White spot on chest",
    "is_neutered": true,
    "breed": null
  },
  "owner": {
    "display_name": "Sunny Example",
    "masked_phone": "081****678"
  },
  "clinic_relationship": {
    "first_appointment_at": "2026-07-01T03:00:00Z",
    "last_appointment_at": "2026-07-12T03:00:00Z",
    "next_appointment_at": "2026-07-15T03:00:00Z",
    "total_appointments": 2
  },
  "recent_appointments": [
    {
      "id": "appointment-uuid",
      "scheduled_at": "2026-07-12T03:00:00Z",
      "appointment_type": "checkup",
      "status": "scheduled",
      "title": "Annual checkup"
    }
  ]
}
```

Clinic Patients responses deliberately omit `user_id`, `owner_profile_id`,
`clinic_profile_id`, password data, JWT claims, full owner address, medical
records, visits, reports, notifications, payment, and staff schedule data.

## Medical Records

Medical Record routes require JWT and a clinic-side role: canonical `clinic`,
with legacy `clinic_staff` retained for compatibility.

Sprint 10 implements one clinic-owned table: `medical_records`. It does not
implement visits, owner timeline, files, lab results, vaccination records,
prescription tables, audit logs, or version history.

### `POST /api/clinic/patients/:petId/medical-records`

- **Purpose:** Create a medical record for one current-clinic patient.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Path:** `:petId` must be a UUID.
- **Success:** 201; `data` is a Medical Record detail object.
- **Common errors:** 400 invalid body/UUID/date/vitals, 401, 403, 404
  `CLINIC_PROFILE_REQUIRED`/`CLINIC_PATIENT_NOT_FOUND`/`APPOINTMENT_NOT_FOUND`,
  409 `APPOINTMENT_MEDICAL_RECORD_EXISTS`, 500.
- **Backend notes:** `clinic_profile_id`, `created_by_user_id`, and `pet_id`
  are derived by the backend. They are not accepted in the request body.

Request:

```json
{
  "appointmentId": "optional-appointment-uuid",
  "visitAt": "2026-07-11T08:00:00Z",
  "chiefComplaint": "Coughing and reduced appetite",
  "clinicalFindings": "Mild fever",
  "diagnosis": "Upper respiratory infection",
  "treatmentPlan": "Rest and medication",
  "medications": "Medication notes as free text",
  "followUpInstructions": "Follow up in 7 days",
  "nextFollowUpAt": "2026-07-18T08:00:00Z",
  "weightKg": 12.5,
  "temperatureC": 38.2,
  "notes": "Owner informed"
}
```

`appointmentId` is optional. When supplied, it must belong to the authenticated
clinic, match `:petId`, not be cancelled, and not already have a medical
record.

### `GET /api/clinic/medical-records`

- **Purpose:** List medical records for the current clinic.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Success:** 200; `data` is `{ "items": [...], "pagination": {...} }`.
- **Common errors:** 400 invalid query, 401, 403, 404
  `CLINIC_PROFILE_REQUIRED`, 500.

Supported query parameters:

| Query | Description |
| --- | --- |
| `pet_id` | Optional pet UUID |
| `from` | Optional visit date lower bound, `YYYY-MM-DD` |
| `to` | Optional visit date upper bound, `YYYY-MM-DD`; includes the full day |
| `page` | Optional positive integer; default 1 |
| `limit` | Optional positive integer; default 20, max 100 |

Records are sorted by `visitAt` descending.

Example list data:

```json
{
  "items": [
    {
      "id": "medical-record-uuid",
      "visitAt": "2026-07-11T08:00:00Z",
      "chiefComplaint": "Coughing and reduced appetite",
      "diagnosis": "Upper respiratory infection",
      "createdAt": "2026-07-11T08:10:00Z",
      "updatedAt": "2026-07-11T08:10:00Z",
      "pet": {
        "id": "pet-uuid",
        "publicPetId": "PNX-PET-A1B2C3",
        "name": "Milo",
        "species": "dog",
        "breed": null
      },
      "owner": {
        "id": "owner-profile-uuid",
        "fullName": "Sunny Example",
        "phoneNumber": "0812345678"
      },
      "appointment": {
        "id": "appointment-uuid",
        "scheduledAt": "2026-07-11T08:00:00Z",
        "status": "checked_in"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "totalPages": 1
  }
}
```

### `GET /api/clinic/medical-records/:recordId`

- **Purpose:** Fetch one current-clinic medical record.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Path:** `:recordId` must be a UUID.
- **Success:** 200; `data` is a Medical Record detail object.
- **Common errors:** 400 `INVALID_MEDICAL_RECORD_ID`, 401, 403, 404
  `MEDICAL_RECORD_NOT_FOUND`, 500.
- **Backend notes:** Another clinic's record returns 404.

### `PATCH /api/clinic/medical-records/:recordId`

- **Purpose:** Update clinical fields on one current-clinic medical record.
- **Auth/role:** JWT required; `clinic` or `clinic_staff`.
- **Success:** 200; `data` is the updated Medical Record detail object.
- **Common errors:** 400 invalid body/date/vitals/empty patch, 401, 403, 404,
  500.
- **Backend notes:** The API does not allow changing `id`, `clinic_profile_id`,
  `pet_id`, `appointment_id`, `created_by_user_id`, or `created_at`.

Allowed patch fields:

```json
{
  "visitAt": "2026-07-11T08:30:00Z",
  "chiefComplaint": "Updated complaint",
  "clinicalFindings": "Updated findings",
  "diagnosis": "Updated diagnosis",
  "treatmentPlan": "Updated plan",
  "medications": "Updated medication notes",
  "followUpInstructions": "Updated follow-up",
  "nextFollowUpAt": "2026-07-18T08:30:00Z",
  "weightKg": 13.1,
  "temperatureC": 38.0,
  "notes": "Updated notes"
}
```

Medical Record detail includes all list fields plus clinical detail fields and
minimal `createdBy` info:

```json
{
  "id": "medical-record-uuid",
  "visitAt": "2026-07-11T08:00:00Z",
  "chiefComplaint": "Coughing and reduced appetite",
  "clinicalFindings": "Mild fever",
  "diagnosis": "Upper respiratory infection",
  "treatmentPlan": "Rest and medication",
  "medications": "Medication notes as free text",
  "followUpInstructions": "Follow up in 7 days",
  "nextFollowUpAt": "2026-07-18T08:00:00Z",
  "weightKg": 12.5,
  "temperatureC": 38.2,
  "notes": "Owner informed",
  "createdAt": "2026-07-11T08:10:00Z",
  "updatedAt": "2026-07-11T08:10:00Z",
  "pet": {
    "id": "pet-uuid",
    "publicPetId": "PNX-PET-A1B2C3",
    "name": "Milo",
    "species": "dog",
    "breed": null
  },
  "owner": {
    "id": "owner-profile-uuid",
    "fullName": "Sunny Example",
    "phoneNumber": "0812345678"
  },
  "appointment": null,
  "createdBy": {
    "id": "user-uuid",
    "email": "clinic@example.com",
    "role": "clinic"
  }
}
```
