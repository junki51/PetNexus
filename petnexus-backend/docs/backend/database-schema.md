# Database Schema

PostgreSQL schema is created by guarded startup SQL in
`internal/database/migrate.go`. Manual equivalents are numbered in
`migrations/001` through `008`. GORM AutoMigrate is intentionally not used for
existing tables. Sprint 9 adds no migration because clinic patients are derived
from appointment relationships.

## Relationship overview

```text
users 1 ── 0..1 owner_profiles
users 1 ── 0..1 clinic_profiles
owner_profiles 1 ── N pets
breeds 1 ── N pets (optional from pet side)
owner_profiles 1 ── N appointments
clinic_profiles 1 ── N appointments
pets 1 ── N appointments
```

## `users`

**Purpose:** Login identity, password hash, and role.

| Field | Type/rule |
| --- | --- |
| `id` | UUID primary key; `gen_random_uuid()` default |
| `email` | varchar(255), required |
| `phone` | varchar(30), nullable |
| `password_hash` | text, required |
| `role` | `user_role` enum, required |
| `created_at`, `updated_at` | timestamptz, required |

Important database rules:

- unique idempotent index `idx_users_email_unique`
- role lookup index `idx_users_role`
- enum values include `owner`, `clinic`, legacy `clinic_staff`, and `admin`
- enum values are added with guarded idempotent SQL

Passwords are bcrypt hashes. Models must never be serialized directly to API
responses.

## `owner_profiles`

**Purpose:** Editable pet-owner identity and contact data, separate from auth.

Important fields:

- UUID `id`
- UUID `user_id`
- required `first_name`, `last_name`, `phone_number`
- nullable `gender`, `date_of_birth`, `avatar_url`
- nullable `address_line1`, `address_line2`, `province`, `district`,
  `subdistrict`, `postal_code`
- timestamps

Constraints/indexes:

- unique index `idx_owner_profiles_user_id_unique`
- guarded foreign key `user_id → users(id)`
- unique `user_id` enforces at most one owner profile per account

## `clinic_profiles`

**Purpose:** Clinic identity/settings for the Clinic Web Dashboard foundation.

| Field | Type/rule |
| --- | --- |
| `id` | UUID primary key |
| `user_id` | UUID, required |
| `clinic_name` | varchar(200), required |
| `phone_number` | varchar(30), nullable |
| `email` | varchar(255), nullable |
| `address` | text, nullable |
| `created_at`, `updated_at` | timestamptz, required |

Constraints/indexes:

- unique index `idx_clinic_profiles_user_id_unique`
- guarded foreign key `user_id → users(id)`
- unique `user_id` enforces at most one clinic profile per account

## `breeds`

**Purpose:** Dog/cat breed reference catalog.

| Field | Type/rule |
| --- | --- |
| `id` | UUID primary key |
| `species` | varchar(20), required; dog/cat check |
| `name` | varchar(100), required |
| `name_th` | varchar(100), nullable |
| `created_at`, `updated_at` | timestamptz, required |

Constraints/indexes:

- unique composite index `idx_breeds_species_name_unique`
- species index `idx_breeds_species`
- guarded check allows only `dog` or `cat`
- startup seed inserts 8 dog and 8 cat breeds using
  `ON CONFLICT (species, name) DO NOTHING`

## `pets`

**Purpose:** Owner-controlled basic pet identity. Passport and medical data are
not stored here.

| Field | Type/rule |
| --- | --- |
| `id` | UUID primary key |
| `public_pet_id` | varchar(50), required; backend-generated public identifier |
| `owner_profile_id` | UUID, required |
| `breed_id` | UUID, nullable |
| `species` | varchar(20), required |
| `name` | varchar(100), required |
| `gender` | varchar(30), nullable |
| `date_of_birth` | date, nullable |
| `weight_kg` | numeric(6,2), nullable |
| `microchip_id` | varchar(100), nullable |
| `avatar_url` | text, nullable |
| `color` | varchar(100), nullable |
| `distinctive_marks` | text, nullable |
| `is_neutered` | boolean, nullable |
| `created_at`, `updated_at` | timestamptz, required |

Constraints/indexes:

- index `idx_pets_owner_profile_id`
- unique index `idx_pets_public_pet_id_unique`
- index `idx_pets_breed_id`
- index `idx_pets_species`
- guarded foreign key `owner_profile_id → owner_profiles(id)`
- guarded foreign key `breed_id → breeds(id)`
- species check allows only `dog` or `cat`
- gender check allows null, `male`, `female`, or `unknown`

`microchip_id` is not unique yet because collision and ownership policy has not
been finalized.

`public_pet_id` uses the format `PNX-PET-XXXXXX`, where the suffix is six
uppercase alphanumeric characters. The backend generates it and retries on a
database uniqueness collision. Startup migration 006 safely backfills existing
pets before enforcing `NOT NULL`.

## Ownership resolution

The client never supplies database ownership keys:

```text
JWT user ID
→ owner_profiles.user_id
→ owner_profiles.id
→ pets.owner_profile_id
```

Clinic profile resolution is:

```text
JWT user ID → clinic_profiles.user_id
```

Repository lookups and service rules enforce these relationships before data is
returned or updated.

## `appointments`

**Purpose:** Scheduling foundation for owners and the Clinic Web Calendar.

Important fields:

- UUID `id`
- required `owner_profile_id`, `clinic_profile_id`, and `pet_id`
- optional `title` and `note`
- required `appointment_type`, `scheduled_at`, `duration_minutes`, and
  `status`
- nullable `created_by_user_id`, required canonical `created_by_role`
- nullable `cancelled_at`
- timestamps

Rules:

- type: `checkup`, `vaccination`, `consultation`, `follow_up`,
  `grooming`, `emergency`, or `other`
- status: `requested`, `scheduled`, `checked_in`, `completed`, or
  `cancelled`
- creator role: `owner` or `clinic`
- duration: 5–480 minutes
- guarded foreign keys reference owner profiles, clinic profiles, pets, and
  optional creator user
- indexes cover each ownership key, pet, scheduled time, status, plus
  clinic/time and owner/time composites

## Derived clinic patients

**Purpose:** Backend query foundation for the Clinic Web Patients page.

There is no `patients` table. A clinic patient is derived as:

```text
unique appointments.pet_id
where appointments.clinic_profile_id = current clinic profile
and appointments.status <> 'cancelled'
```

The repository joins the derived pet IDs to `pets`, `owner_profiles`, and
`breeds` for safe patient list/detail responses. Cross-clinic access is hidden
as 404 at the service layer.

## `medical_records`

**Purpose:** Clinic-owned clinical notes for one patient visit.

Important fields:

- UUID `id`
- required `clinic_profile_id`, `pet_id`, and `created_by_user_id`
- nullable `appointment_id`
- required `visit_at` and `chief_complaint`
- nullable `clinical_findings`, `diagnosis`, `treatment_plan`, `medications`,
  `follow_up_instructions`, `next_follow_up_at`, `weight_kg`,
  `temperature_c`, and `notes`
- timestamps

Rules:

- `clinic_profile_id` references `clinic_profiles(id)`
- `pet_id` references `pets(id)`
- `appointment_id` references `appointments(id)` and is nullable
- `created_by_user_id` references `users(id)`
- partial unique index `idx_medical_records_appointment_id_unique` allows at
  most one medical record per non-null appointment
- indexes cover clinic, pet, clinic/visit date, and pet/visit date lookups
- `weight_kg` and `temperature_c` must be positive when supplied
- `next_follow_up_at` must not be earlier than `visit_at`

No `owner_profile_id` is stored directly because owner data is derived through:

```text
medical_records.pet_id -> pets.owner_profile_id -> owner_profiles.id
```
