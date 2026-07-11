# Backend Roadmap

## Completed

Sprint 1-9 are implemented and tested:

- Backend and health foundation
- PostgreSQL connection and guarded startup migration
- Auth, bcrypt, JWT, and role middleware
- Owner Profile
- Breed catalog and basic Pet management
- Clinic Profile foundation
- Pet public ID and clinic pet lookup
- Appointment Calendar foundation
- Clinic Patient List backend

See [Sprint 1-9 summary](../sprints/sprint-1-to-9-summary.md).

## Not implemented

- QR sharing
- Clinic QR scanning
- Clinic access request
- Owner approval/rejection
- Authorized clinic relationships
- Medical records
- Verified visits
- Pet timeline
- Appointment overlap/capacity and staff scheduling
- Reports and analytics
- Notifications
- Real file/image upload
- Full Pet Passport backend
- Clinic staff-member management

Placeholder files for some future domains are not implemented APIs and must not
be treated as working features.

## Recommended next planning topic

The next backend sprint should design **Medical Records / Visit Records
Foundation** or **Clinic Access Request**, but should not combine both. If
Clinic Access Request is selected, QR should remain only an optional transport
for the public pet ID. The design should answer:

1. How a clinic requests access after finding a pet.
2. What pet data is visible before and after owner authorization.
3. Whether a QR carries only `public_pet_id` or a separate token.
4. If a token is used, its lifetime, revocation, and replay rules.
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
