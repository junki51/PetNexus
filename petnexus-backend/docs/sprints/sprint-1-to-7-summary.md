# Sprint 1–7 Backend Summary

## Sprint 1–6

Sprint 1–6 established the Go/Gin server, PostgreSQL/GORM connection, guarded
startup migrations, JWT/bcrypt authentication, owner profiles, breed catalog,
owner-managed pets, and clinic profiles. See the preserved
[Sprint 1–6 summary](./sprint-1-to-6-summary.md) for the detailed history.

## Sprint 7: Pet Public ID + Clinic Pet Lookup

**Added:** A permanent backend-generated public identifier for every pet,
including safe startup backfill for existing rows. IDs use
`PNX-PET-XXXXXX`, with a database unique index and collision retry.

**Main endpoint:**

- `GET /api/clinic/pet-lookup?pet_id=PNX-PET-XXXXXX`
- `GET /api/clinic/pet-lookup?owner_phone=<exact-phone>`

**Schema change:** `pets.public_pet_id VARCHAR(50) NOT NULL` and unique index
`idx_pets_public_pet_id_unique`.

**Access:** JWT + canonical `clinic` role or legacy-compatible
`clinic_staff`. Owner receives 403 and missing authentication receives 401.

**Privacy:** Lookup returns only limited pet identity, optional breed, owner
display name, and masked owner phone. Phone lookup is exact-only and no matches
return an empty list.

**Intentionally excluded:** QR generation/scanning, clinic access request,
owner approval, authorized patient relationships, medical records, visits,
timeline, staff, reports, notifications, and frontend changes.

## Current backend boundary

Clinics can identify a pet but do not gain access to it. A future sprint should
design and implement the explicit owner-controlled clinic access relationship;
a QR may later act as a shortcut carrying the public pet ID.
