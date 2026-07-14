# PetNexus Backend Technical Debt

Audit date: 2026-07-14  
Scope: current Sprint 1-10 backend code only.

| ID | Severity | Category | Description | Affected files | Recommended fix | Complexity | Suggested timing |
| --- | --- | --- | --- | --- | --- | --- | --- |
| TD-001 | Critical | Security / config | App accepts the known default `JWT_SECRET=change_me_in_production` if env is missing. | `internal/config/config.go`, `cmd/api/main.go`, `internal/utils/jwt.go`, `.env.example` | Fail startup outside development when JWT secret is empty, default, or weak. Add Render setup docs/checklist. | Small | Fix now |
| TD-002 | Critical | Authorization model | Clinic accounts can create appointments for pets by `pet_id` or `public_pet_id` without owner approval; this creates the appointment-derived patient relationship. | `internal/services/auth_service.go`, `internal/services/clinic_appointment_service.go`, `internal/repositories/pet_repository.go`, `internal/repositories/clinic_patient_repository.go`, `internal/services/medical_record_service.go` | Add verified clinic onboarding plus owner authorization/access grants before a clinic can create patient relationships or access expanded patient data. | Large | Before public demo |
| TD-003 | High | Role management | Public registration allows `clinic` and `clinic_staff` roles. | `internal/services/auth_service.go`, `internal/routes/routes.go` | Restrict public registration to owner or add invite/admin approval for clinic roles. | Medium | Before public demo |
| TD-004 | High | Auth abuse protection | Login/register have no rate limit, lockout, email verification, or abuse controls. | `internal/handlers/auth_handler.go`, `internal/services/auth_service.go`, `cmd/api/main.go` | Add edge/app rate limiting, failed-attempt tracking, monitoring, and optional verification. | Medium | Before public demo |
| TD-005 | High | API privacy | Medical record responses return full owner phone and owner profile ID, while other clinic-facing flows mask phone. | `internal/services/medical_record_service.go`, `internal/dto/medical_record_dto.go` | Mask phone by default or require explicit owner authorization before exposing full contact info. | Small | Before public demo |
| TD-006 | High | HTTP hardening | Server uses `gin.Default()` and `router.Run()` without explicit request size limits or read/write timeouts. | `cmd/api/main.go`, all JSON handlers | Use explicit `http.Server` with timeouts and body-size limiting middleware. | Medium | Before public demo |
| TD-007 | Medium | Database context | Repository queries do not receive request context. | `internal/repositories/*`, service interfaces, handlers | Thread `context.Context` from handlers to services/repositories and use `db.WithContext(ctx)`. | Large | Before next feature |
| TD-008 | Medium | Defense-in-depth | Some update queries rely on prior service scoping and do not repeat ownership scope in the `UPDATE` statement. | `internal/repositories/owner_repository.go`, `internal/repositories/clinic_repository.go`, `internal/repositories/appointment_repository.go` | Add scoped update methods or include `user_id`, `owner_profile_id`, or `clinic_profile_id` in `WHERE` clauses. | Small | Before next feature |
| TD-009 | Medium | Migrations | Guarded `ALTER TABLE ADD CONSTRAINT` can fail startup on existing dirty data. | `internal/database/migrate.go`, `migrations/*.sql` | Add preflight checks or use `NOT VALID` plus later validation for future constraints. | Medium | Before next schema sprint |
| TD-010 | Medium | Migration operations | No migration version table or dedicated migration runner. | `internal/database/migrate.go`, `migrations/README.md` | Introduce `golang-migrate` or equivalent with versioned applied-state tracking. | Medium | Before next schema sprint |
| TD-011 | Medium | Data integrity | Email uniqueness is normalized in service, but database unique index is case-sensitive. | `internal/services/auth_service.go`, `internal/database/migrate.go`, `migrations/002_create_users.sql` | Use `citext` or a unique index on `lower(email)` after checking/backfilling existing data. | Medium | Before public demo |
| TD-012 | Medium | Testing | No repository or migration integration tests against real PostgreSQL. | `internal/repositories/*`, `internal/database/migrate_test.go` | Add ephemeral PostgreSQL integration tests for migrations, appointments, clinic patients, and medical records. | Large | Before next feature |
| TD-013 | Medium | Validation / resource usage | Medical record free-text fields are unbounded and request bodies are not size-limited. | `internal/services/medical_record_service.go`, `internal/database/migrate.go` | Define max lengths for clinical fields and add body-size limits. | Small | Before public demo |
| TD-014 | Medium | Appointment workflow | Clinic appointment status updates allow any supported status with no transition policy. | `internal/services/clinic_appointment_service.go`, `internal/services/appointment_service_helpers.go` | Define a status transition matrix and add audit logging for status changes. | Medium | Later |
| TD-015 | Low | API consistency | JSON naming is mixed: auth uses camelCase, most older DTOs use snake_case, medical records use camelCase. | `internal/dto/*.go`, `docs/backend/api-reference.md` | Choose one convention for future endpoints and plan versioned cleanup. | Medium | Later |
| TD-016 | Low | Documentation hygiene | Historical planning docs coexist with current docs and can be mistaken for implemented behavior. | `docs/*`, `README.md` | Add stronger "historical" labels or move old plans under an archive folder. | Small | Later |
| TD-017 | Low | Repository clarity | Placeholder files for QR, visits, timeline, notifications, audit logs, authorization, and clinic staff can imply features exist. | `internal/services/*`, `internal/handlers/*`, `internal/models/*`, `internal/repositories/*` placeholders | Keep placeholders documented or move future-only placeholders to docs until implementation starts. | Small | Later |
| TD-018 | Low | Health endpoint disclosure | Public `/health/db` exposes database reachability. | `internal/routes/routes.go`, `internal/handlers/health_handler.go` | Keep if Render needs it, or split public liveness from protected/internal readiness. | Small | Later |
| TD-019 | Low | Schema ownership | GORM tags may drift from manual SQL migrations; SQL is authoritative but this is easy to misunderstand. | `internal/models/*.go`, `internal/database/migrate.go`, `migrations/*.sql` | Document SQL authority clearly and add model/schema alignment tests. | Small | Later |
| TD-020 | Low | Production mode | Application code does not set Gin mode or trusted proxy policy based on config. | `cmd/api/main.go` | Set production mode explicitly in deployment or startup code and configure trusted proxies intentionally. | Small | Before public demo |

## Suggested execution order

1. TD-001
2. TD-002 and TD-003 together
3. TD-005
4. TD-004 and TD-006
5. TD-012
6. TD-007 and TD-008
7. TD-009 and TD-010
8. TD-011 and TD-013
9. TD-014 through TD-020 as cleanup/hardening
