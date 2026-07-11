# Backend Overview

PetNexus is a Digital Pet Passport and Digital Pet Identity Platform for pet
owners and veterinary clinics. The backend is a Go/Gin REST API backed by
PostgreSQL and organized around owner-controlled pet identity.

## Current technology

- Go 1.22+
- Gin HTTP server
- PostgreSQL
- GORM PostgreSQL driver
- JWT access tokens
- bcrypt password hashing
- Docker Compose PostgreSQL for local development
- Render Web Service and Render Postgres for deployment

## Completed backend scope

Sprint 1-10 currently supports:

- application health check
- PostgreSQL health check
- owner and clinic user registration/login
- bcrypt password storage
- JWT authentication and role-based access control
- current authenticated user lookup
- owner profile create/get/partial update
- clinic profile create/get/partial update
- public dog/cat breed list and species filtering
- owner pet create/list/detail/partial update
- backend-generated permanent public pet IDs
- clinic pet lookup by public pet ID or exact owner phone
- owner and clinic appointment scheduling
- clinic calendar date/status/type filtering
- clinic appointment status and cancellation actions
- clinic patient list/detail derived from non-cancelled appointments
- clinic-owned medical record create/list/detail/update foundation
- database startup migration using guarded, idempotent SQL

## Supported owner flow

```text
Register/login as owner
→ Create or fetch owner profile
→ Select dog/cat and breed
→ Create pet
→ List/fetch/update own pets
→ Create/list/detail/cancel own appointments
```

The backend resolves the owner profile from the JWT user ID. Clients do not
choose `user_id` or `owner_profile_id`.

## Supported clinic flow

```text
Register/login as clinic
→ Create clinic profile
→ Fetch/update clinic profile settings
→ Look up limited pet identity by public pet ID or exact owner phone
→ Create/list appointments and manage calendar status
→ List/search patient pets and fetch patient detail summaries
→ Create/list/detail/update medical records for clinic patients
```

`clinic` is the canonical Clinic Web Dashboard role. Legacy `clinic_staff`
accounts remain accepted for compatibility. Clinic profile ownership is also
resolved from the JWT user ID.

## Not implemented yet

- QR sharing and QR scanning (a future QR may carry the public pet ID)
- clinic access requests
- owner approval/rejection
- authorized clinic/patient relationships
- verified visits and owner medical timeline
- pet timeline
- appointment overlap/capacity management and external calendar sync
- reports and analytics
- notifications
- file/image uploads, lab results, vaccination records, prescription tables
- complete Pet Passport backend
- clinic staff-member management

See [Backend roadmap](./roadmap.md) for the recommended next planning topic.
