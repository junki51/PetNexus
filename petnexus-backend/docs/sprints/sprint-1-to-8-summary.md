# Sprint 1–8 Backend Summary

## Sprint 1–7

Sprint 1–7 established the Go/Gin and PostgreSQL foundation, authentication,
owner and clinic profiles, owner-managed pets, public pet IDs, and
privacy-limited clinic pet lookup. See the preserved
[Sprint 1–7 summary](./sprint-1-to-7-summary.md).

## Sprint 8: Appointment Calendar Foundation

**Added:** Appointment persistence and owner/clinic appointment APIs suitable
for the first real Clinic Web Calendar integration.

**Owner access:** JWT role `owner`; create only for an owned pet and
list/get/cancel only appointments belonging to the JWT-derived owner profile.

**Clinic access:** JWT role `clinic` or legacy `clinic_staff`; create,
list/get, update status, and cancel only under the JWT-derived clinic profile.

**Calendar:** Optional exact UTC day, inclusive UTC date range, status, and
appointment-type filters. Results are ordered by scheduled time ascending.

**Privacy:** Responses expose useful pet, masked owner, and clinic summaries,
but omit internal profile/user ownership IDs and authentication data.

**Intentionally excluded:** Medical records, dashboard aggregates, reports,
notifications, staff schedules, payments, Google Calendar sync, full QR/access
workflow, and frontend implementation.
