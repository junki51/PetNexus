# Backend Roadmap

## Completed

Sprint 1–6 are implemented and tested:

- Backend and health foundation
- PostgreSQL connection and guarded startup migration
- Auth, bcrypt, JWT, and role middleware
- Owner Profile
- Breed catalog and basic Pet management
- Clinic Profile foundation

See [Sprint 1–6 summary](../sprints/sprint-1-to-6-summary.md).

## Not implemented

- QR sharing
- Clinic QR scanning
- Clinic access request
- Owner approval/rejection
- Authorized clinic relationships
- Clinic patient list
- Medical records
- Verified visits
- Pet timeline
- Calendar and appointments
- Reports and analytics
- Notifications
- Real file/image upload
- Full Pet Passport backend
- Clinic staff-member management

Placeholder files for some future domains are not implemented APIs and must not
be treated as working features.

## Recommended Sprint 7 planning topic

Sprint 7 should design **QR Sharing + Clinic Access Request**, not immediately
implement it. The design should answer:

1. What a QR token identifies and how long it remains valid.
2. Whether QR is single-use, reusable, or revocable.
3. What pet data is visible before owner authorization.
4. How a clinic requests access after scanning.
5. How owner approval/rejection is represented.
6. Authorization scope, expiry, revocation, and audit trail.
7. Replay protection and token hashing/storage.
8. Whether clinic profile or staff identity is the authorization principal.
9. Exact status transitions and idempotency rules.
10. Which events later appear in notifications and timeline.

Only after these decisions are documented should schema, migration, DTO,
service, and endpoint implementation begin.

## Scope guard

Do not combine QR, access authorization, visits, medical records, and timeline
in one sprint. They have different permission and audit requirements and should
be delivered in small, testable foundations.
