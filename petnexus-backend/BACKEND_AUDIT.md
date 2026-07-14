# PetNexus Backend Audit

Audit date: 2026-07-14  
Scope: current repository code only, Sprint 1-10 backend implementation.  
Mode: inspection and reporting only. No application code, migrations, tests, configuration, routes, models, DTOs, handlers, services, repositories, middleware, or existing documentation were modified by this audit.

## 1. Executive summary

The PetNexus backend is clearly organized around the intended dependency direction:

```text
Handler -> Service -> Repository -> Database
```

The implemented Sprint 1-10 code has a solid learning/demo foundation. The strongest areas are:

- JWT middleware and role middleware are centralized.
- Most ownership identifiers come from JWT-derived profiles rather than request bodies.
- Owner pet access, owner appointments, clinic appointments, clinic patients, and medical records are scoped in services and repository queries.
- GORM `AutoMigrate` has been replaced with guarded SQL startup migrations.
- Unit and route tests cover many permission and validation boundaries.

The main production-readiness risks are not accidental syntax issues; they are product-security boundaries that have not been implemented yet:

1. The app can run with a known default `JWT_SECRET`.
2. Public clinic registration plus clinic-created appointments can establish a clinic-patient relationship for a pet without owner approval.
3. Once that clinic-patient relationship exists, clinic patient detail and medical record APIs expose more owner/pet data than a public deployment should allow without an authorization model.

No confirmed destructive migration behavior or password-hash exposure was found in the current code.

## 2. Current architecture overview

Startup is in `cmd/api/main.go`:

1. `config.Load()`
2. `database.ConnectPostgres(cfg)`
3. `database.RunMigrations(db)`
4. repository construction
5. service construction
6. handler construction
7. `gin.Default()`
8. `routes.Register(...)`
9. `router.Run(":" + cfg.Port)`

Layer responsibilities are generally separated well:

- `internal/handlers`: JSON binding, path/query extraction, response writing.
- `internal/services`: validation, ownership decisions, response mapping.
- `internal/repositories`: GORM queries and database error mapping.
- `internal/middleware`: JWT parsing and role checks.
- `internal/database`: PostgreSQL connection, ping, and guarded SQL migrations.
- `internal/models`: GORM table models.
- `internal/dto`: API request/response shapes.

Representative flows follow this pattern consistently:

- `POST /api/auth/login`: route -> `AuthHandler.Login` -> `authService.Login` -> `UserRepository.FindByEmail`.
- `POST /api/pets`: route owner middleware -> `PetHandler.CreatePet` -> `petService.CreatePet` -> owner profile, breed, pet repositories.
- `POST /api/owner/appointments`: owner middleware -> owner appointment handler/service -> owner profile, clinic, pet, appointment repositories.
- `GET /api/clinic/patients`: clinic middleware -> clinic patient handler/service -> clinic profile repository -> raw SQL derived from appointments.
- `POST /api/clinic/patients/:petId/medical-records`: clinic middleware -> medical record handler/service -> clinic profile and medical record repositories.

No circular import dependency was observed. Placeholder files for future modules exist but are not wired into routes.

## 3. What is implemented well

- Passwords are hashed with bcrypt in `internal/utils/password.go`; `models.User` comments explicitly warn not to serialize the model directly.
- JWT parsing in `internal/utils/jwt.go` requires HS256 and expiration.
- Auth middleware stores only `userID` and `userRole` in Gin context.
- Owner routes require `models.RoleOwner`.
- Clinic routes require `models.RoleClinic` or `models.RoleClinicStaff`.
- Pet DTOs do not accept `owner_profile_id` or `user_id`.
- Owner profile and clinic profile DTOs do not accept `user_id`.
- Medical record create DTO does not accept `clinic_profile_id`, `pet_id`, or `created_by_user_id`; those values come from URL/JWT.
- Cross-owner pet access uses `FindByIDAndOwnerProfileID`.
- Cross-clinic appointment and medical record access uses clinic-scoped repository methods.
- Startup migrations are explicit SQL and avoid GORM `AutoMigrate`.
- Breeds are seeded idempotently with `ON CONFLICT DO NOTHING`.
- Medical records use a partial unique index to prevent multiple records for the same non-null appointment.

## 4. Critical findings

### C-01: Known default JWT secret is accepted at runtime

- Severity: Critical
- Affected files: `internal/config/config.go`, `internal/utils/jwt.go`, `cmd/api/main.go`, `.env.example`
- Relevant functions/types: `config.Load`, `utils.GenerateAccessToken`, `utils.ParseAccessToken`, `main`
- Evidence: `config.Load` defaults `JWTSecret` to `change_me_in_production`. Startup does not reject this value when `APP_ENV` is production or when the value is too weak.
- Explanation: If a deployed environment forgets to set `JWT_SECRET`, tokens can be signed with a known public string.
- Realistic impact: Authentication bypass is possible in a misconfigured production or public demo environment because an attacker could forge valid HS256 tokens.
- Recommended correction: Add startup validation that fails when `APP_ENV != development` and `JWT_SECRET` is empty, default, or below an agreed entropy/length threshold. Document the required Render environment variable.
- Behavior change: Yes. Misconfigured production deployments would fail fast instead of serving traffic.

### C-02: A clinic can establish a patient relationship for an arbitrary pet without owner authorization

- Severity: Critical for public use; acceptable only as a known demo limitation.
- Affected files: `internal/services/auth_service.go`, `internal/routes/routes.go`, `internal/services/clinic_appointment_service.go`, `internal/repositories/pet_repository.go`, `internal/repositories/clinic_patient_repository.go`, `internal/services/medical_record_service.go`
- Relevant functions/types: `validateRegistration`, `RequireRole`, `CreateClinicAppointment`, `resolveAppointmentPet`, `FindByPublicPetID`, `FindByID`, `FindPatientsByClinicProfileID`, `PetHasNonCancelledAppointmentWithClinic`
- Evidence:
  - `validateRegistration` allows public `owner`, `clinic`, and `clinic_staff` registration.
  - `CreateClinicAppointment` accepts `pet_id` or `public_pet_id`.
  - `resolveAppointmentPet` loads the pet by ID or public pet ID without owner approval.
  - Clinic patients are derived from non-cancelled appointments for that clinic.
  - Medical record creation checks only that the pet has a non-cancelled appointment with the clinic.
- Explanation: Any clinic-side account can create an appointment for a pet if it knows the pet UUID or public pet ID. That appointment creates the relationship used by clinic patient and medical record flows.
- Realistic impact: In a public deployment, an unverified clinic user could create an appointment against another owner's pet and then access clinic patient detail or create clinic-owned medical records for that pet.
- Recommended correction: Before public use, add clinic verification and an owner authorization/access grant model. Clinic appointment creation, patient visibility, and medical record creation should require an approved relationship, QR/session authorization, or owner-confirmed appointment.
- Behavior change: Yes. Clinic-created appointments and patient access would become permission-gated.

## 5. High findings

### H-01: Public clinic and clinic_staff registration has no approval workflow

- Severity: High
- Affected files: `internal/services/auth_service.go`, `internal/routes/routes.go`
- Relevant functions/types: `validateRegistration`, `Register`, clinic route group in `routes.Register`
- Evidence: `validateRegistration` permits `models.RoleClinic` and `models.RoleClinicStaff` from the public register endpoint.
- Explanation: The role itself unlocks clinic-only endpoints. There is no admin approval, invitation, email verification, or clinic verification.
- Realistic impact: A user can self-select clinic-side privileges. This amplifies the impact of pet lookup, clinic appointment, clinic patients, and medical records.
- Recommended correction: Restrict public registration to owner only, or add an invitation/admin approval state before clinic role can access clinic endpoints.
- Behavior change: Yes. Clinic onboarding would change.

### H-02: Login/register endpoints have no brute-force or abuse protection

- Severity: High
- Affected files: `internal/handlers/auth_handler.go`, `internal/services/auth_service.go`, `cmd/api/main.go`
- Relevant functions/types: `AuthHandler.Login`, `authService.Login`, `AuthHandler.Register`
- Evidence: No rate limiter, IP throttling, account lockout, CAPTCHA, email verification, or failed-login tracking is present.
- Explanation: The login endpoint correctly returns generic invalid credentials, but nothing slows repeated guesses or automated account creation.
- Realistic impact: Credential stuffing, password guessing, and signup abuse are practical risks before public deployment.
- Recommended correction: Add rate limiting at the edge or app level, track failed attempts, consider email verification, and add monitoring.
- Behavior change: Yes. Some requests would be throttled or challenged.

### H-03: Medical record responses expose full owner phone and owner profile ID

- Severity: High
- Affected files: `internal/services/medical_record_service.go`, `internal/dto/medical_record_dto.go`
- Relevant functions/types: `toMedicalRecordOwnerSummary`, `MedicalRecordOwnerSummary`
- Evidence: `MedicalRecordOwnerSummary` returns `ID`, `FullName`, and `PhoneNumber`. `toMedicalRecordOwnerSummary` maps `pet.OwnerProfile.PhoneNumber` directly. In contrast, clinic lookup and clinic patient responses mask phone numbers.
- Explanation: The medical record response gives more owner PII than related clinic-facing endpoints.
- Realistic impact: If clinic relationship creation is not owner-authorized, full phone numbers can leak to a clinic-side account.
- Recommended correction: Mask phone by default or only expose full owner contact after explicit owner/clinic authorization. Revisit whether owner profile ID should be returned.
- Behavior change: Yes. API response shape or authorization requirements would change.

### H-04: No request size limit or HTTP server timeout configuration

- Severity: High before public deployment
- Affected files: `cmd/api/main.go`, all JSON handlers using `ShouldBindJSON`
- Relevant functions/types: `main`, `gin.Default`, `ShouldBindJSON`
- Evidence: `cmd/api/main.go` uses `gin.Default()` and `router.Run(...)`. No `http.Server` with read/write timeouts and no request body size limiter were found.
- Explanation: Large JSON bodies and slow clients are not bounded at the application level.
- Realistic impact: Resource exhaustion is possible, especially on endpoints accepting free-text fields such as medical records.
- Recommended correction: Use an explicit `http.Server` with timeouts and middleware such as `http.MaxBytesReader` or Gin body-size limiting. Coordinate with Render/proxy limits.
- Behavior change: Yes. Over-limit or slow requests would be rejected.

## 6. Medium findings

### M-01: Repository queries do not use request contexts

- Severity: Medium
- Affected files: all files under `internal/repositories`, selected database helpers
- Relevant functions/types: repository methods using `r.db`
- Evidence: Repository methods call `r.db.Where`, `r.db.Create`, `r.db.Raw`, etc. without `WithContext(c.Request.Context())`. Only `PingPostgres` uses context.
- Explanation: Client cancellation and deadlines do not propagate into database work.
- Realistic impact: Slow queries can continue after clients disconnect; timeouts are harder to enforce.
- Recommended correction: Pass `context.Context` from handlers to services and repositories, or attach context at repository call boundaries.
- Behavior change: Mostly no, except cancelled/timed-out requests would stop earlier.

### M-02: Some update repository methods rely on service scoping rather than scoped UPDATE statements

- Severity: Medium
- Affected files: `internal/repositories/owner_repository.go`, `internal/repositories/clinic_repository.go`, `internal/repositories/appointment_repository.go`
- Relevant functions/types: `OwnerProfileRepository.Update`, `ClinicProfileRepository.Update`, `AppointmentRepository.Update`
- Evidence:
  - Owner profile update uses `Where("id = ?", profile.ID)`.
  - Clinic profile update uses `Where("id = ?", profile.ID)`.
  - Appointment update uses `Where("id = ?", appointment.ID)`.
- Explanation: Services fetch scoped records first, so current behavior is mostly safe. However, the write query itself does not repeat the owner/clinic scope.
- Realistic impact: Lower defense-in-depth; future service changes or concurrent reassignment/manual data edits could make writes less robust.
- Recommended correction: Add scoped update methods or include `user_id`, `owner_profile_id`, or `clinic_profile_id` in update `WHERE` clauses where possible.
- Behavior change: No intended behavior change.

### M-03: Startup migrations can fail on existing dirty data when adding constraints

- Severity: Medium
- Affected files: `internal/database/migrate.go`, `migrations/*.sql`
- Relevant functions/types: `RunMigrations`, `ensureOwnerProfilesUserForeignKeySQL`, `ensurePets...`, `ensureAppointmentConstraintsSQL`, `ensureMedicalRecordConstraintsSQL`
- Evidence: Guarded `ALTER TABLE ... ADD CONSTRAINT` statements are idempotent, but they validate existing data immediately.
- Explanation: Existing Render data created outside the app or from earlier buggy versions could violate a new FK/check and stop startup.
- Realistic impact: A redeploy can fail until data is cleaned.
- Recommended correction: For future migrations, consider `NOT VALID` plus a deliberate validation step after data cleanup, or add preflight checks with clear error messages.
- Behavior change: No intended API behavior change.

### M-04: No migration version table or dedicated migration runner

- Severity: Medium
- Affected files: `internal/database/migrate.go`, `migrations/README.md`
- Relevant functions/types: `RunMigrations`
- Evidence: Startup executes every guarded migration step on each boot. `migrations/README.md` notes that a dedicated runner such as `golang-migrate` can be introduced later.
- Explanation: Guarded SQL works for the current scope, but there is no version history, rollback plan, or applied migration audit.
- Realistic impact: Harder to reason about schema state as the project grows.
- Recommended correction: Introduce a migration runner before schema complexity increases.
- Behavior change: Mostly operational; API behavior should not change.

### M-05: Medical record free-text fields have no application-level max length

- Severity: Medium
- Affected files: `internal/services/medical_record_service.go`, `internal/database/migrate.go`
- Relevant functions/types: `normalizeRequiredMedicalRecordText`, `normalizeOptionalMedicalRecordText`, `createMedicalRecordsTableSQL`
- Evidence: Medical record clinical fields are `TEXT`, and service normalization trims whitespace but does not limit length.
- Explanation: This is flexible for MVP notes, but dangerous without request body limits.
- Realistic impact: Large payloads can increase memory usage and database storage unexpectedly.
- Recommended correction: Define per-field limits and enforce request body size limits.
- Behavior change: Yes. Over-limit medical text would be rejected.

### M-06: Email uniqueness depends on service normalization rather than database normalization

- Severity: Medium
- Affected files: `internal/services/auth_service.go`, `internal/database/migrate.go`, `migrations/002_create_users.sql`
- Relevant functions/types: `normalizeEmail`, `createUsersEmailUniqueIndexSQL`
- Evidence: Service lowercases email before insert. Database unique index is `ON users(email)`, not `ON lower(email)`, and the column is not `citext`.
- Explanation: Current application paths normalize email, but direct DB writes or future code paths could insert case variants.
- Realistic impact: Duplicate logical accounts can appear if data bypasses the service.
- Recommended correction: Use a unique expression index on `lower(email)` or PostgreSQL `citext`.
- Behavior change: Potentially yes for existing mixed-case data and direct DB behavior.

### M-07: No real PostgreSQL integration tests for repositories or migrations

- Severity: Medium
- Affected files: `internal/repositories/*`, `internal/database/migrate_test.go`, `internal/services/*_test.go`, `internal/routes/*_test.go`
- Relevant functions/types: repository methods, migration SQL constants
- Evidence: Existing service tests use stubs/fakes, route tests use `httptest`, and migration tests assert SQL strings. No repository test spins up PostgreSQL.
- Explanation: Unit tests cover business rules, but they do not verify actual SQL behavior, constraints, joins, indexes, or transaction/race behavior.
- Realistic impact: SQL issues can slip to manual smoke tests or Render.
- Recommended correction: Add a small integration test suite against ephemeral PostgreSQL for migrations and high-risk repositories.
- Behavior change: No.

### M-08: Appointment status transitions are intentionally permissive

- Severity: Medium
- Affected files: `internal/services/clinic_appointment_service.go`, `internal/services/appointment_service_helpers.go`
- Relevant functions/types: `UpdateClinicAppointmentStatus`, `normalizeAppointmentStatus`, `setAppointmentStatus`
- Evidence: Any allowed status can be set by clinic routes; there is no transition matrix.
- Explanation: The docs state this is simple for the sprint, but public use usually needs transition rules and auditability.
- Realistic impact: A clinic could move appointments between states inconsistently.
- Recommended correction: Define a transition policy and add audit logging when status changes.
- Behavior change: Yes. Some transitions would be rejected.

## 7. Low findings

### L-01: Response JSON naming is inconsistent across modules

- Severity: Low
- Affected files: `internal/dto/auth_dto.go`, `internal/dto/pet_dto.go`, `internal/dto/medical_record_dto.go`, `docs/backend/api-reference.md`
- Relevant functions/types: `UserResponse`, `PetResponse`, `MedicalRecordDetailResponse`
- Evidence: Auth uses `createdAt`, pets use `created_at`, and medical records use `visitAt`/`publicPetId`.
- Explanation: The docs call out some inconsistency, but clients must handle mixed conventions.
- Realistic impact: Frontend integration friction and accidental mapping bugs.
- Recommended correction: Choose one API naming convention for future endpoints and plan versioned cleanup if needed.
- Behavior change: Yes if existing fields are renamed.

### L-02: Placeholder future modules can confuse repository readers

- Severity: Low
- Affected files: `internal/services/qr_service.go`, `internal/services/visit_service.go`, `internal/services/timeline_service.go`, `internal/models/authorization.go`, `internal/models/visit.go`, etc.
- Relevant functions/types: placeholder files only
- Evidence: Files exist for QR, authorization, visits, timeline, notifications, audit logs, and clinic staff but contain only comments and are not wired in routes.
- Explanation: This is harmless but can make the codebase look larger than the implemented feature set.
- Realistic impact: New contributors may assume features exist.
- Recommended correction: Keep placeholders documented or move future-only planning to docs until implementation starts.
- Behavior change: No.

### L-03: Public `/health/db` reveals database reachability

- Severity: Low
- Affected files: `internal/routes/routes.go`, `internal/handlers/health_handler.go`
- Relevant functions/types: `DatabaseHealth`
- Evidence: `GET /health/db` is registered before `/api` auth middleware and is public.
- Explanation: This is common for platform health checks, but it exposes whether PostgreSQL is reachable.
- Realistic impact: Minor information disclosure.
- Recommended correction: Keep public if Render needs it, or use a less revealing public health endpoint and protect detailed readiness checks.
- Behavior change: Maybe, depending on deployment health check setup.

### L-04: GORM model tags can drift from manual SQL migrations

- Severity: Low
- Affected files: `internal/models/*.go`, `internal/database/migrate.go`, `migrations/*.sql`
- Relevant functions/types: GORM models and SQL constants
- Evidence: AutoMigrate is intentionally not used, so model tags are documentation/helpers rather than the migration authority. Example: `MedicalRecord.AppointmentID` has a GORM `uniqueIndex` tag while SQL implements a partial unique index.
- Explanation: Future contributors might assume tags define schema.
- Realistic impact: Confusion or accidental AutoMigrate reintroduction.
- Recommended correction: Document that SQL migrations are authoritative and add schema alignment tests where practical.
- Behavior change: No.

## 8. Security observations

- No password hashes are returned through DTOs found in the audited flows.
- JWT parsing rejects non-HS256 algorithms and requires expiration.
- Missing or invalid auth returns 401; wrong role returns 403.
- Owner and clinic cross-scope reads generally return 404 inside service/repository flows.
- No raw SQL injection was found in repository query construction. Dynamic order clauses are chosen from service-validated sort keys.
- No committed `.env` file was found; `.gitignore` includes `.env`, `tmp/`, and `*.exe`.
- Local defaults are appropriate for Docker development but unsafe for production if used unchanged.
- No rate limiting, request-size limit, request timeout, or brute-force protection was found.
- Gin production mode/trusted proxy configuration is not set in application code.

## 9. Database and migration observations

Implemented tables:

- `users`
- `owner_profiles`
- `breeds`
- `pets`
- `clinic_profiles`
- `appointments`
- `medical_records`

Good database practices observed:

- `pgcrypto` is enabled idempotently.
- UUID primary keys use `gen_random_uuid()`.
- `user_role` enum is created and extended idempotently.
- Email, owner profile per user, clinic profile per user, breed species/name, pet public ID, and medical record appointment uniqueness are indexed.
- Foreign keys are guarded against duplicate creation.
- `appointments` and `medical_records` have useful lookup indexes.
- Startup migration avoids GORM AutoMigrate constraint rewriting.

Risks:

- No migration version table.
- Existing bad data can break a deploy when adding FKs/checks.
- Manual migrations and startup SQL must be kept in sync.
- Email uniqueness is not database-normalized.
- Medical record clinical text fields are unbounded.

## 10. Testing gaps

Packages with tests:

- `internal/config`
- `internal/database`
- `internal/routes`
- `internal/services`
- `internal/utils`

Packages without direct tests:

- `cmd/api`
- `internal/dto`
- `internal/handlers`
- `internal/middleware`
- `internal/models`
- `internal/repositories`

Important business rules currently tested:

- JWT creation/parsing and expired token rejection.
- bcrypt hash/check behavior.
- owner/clinic route role boundaries.
- owner profile create/get/patch validation.
- clinic profile create/get/patch validation.
- owner pet ownership and breed/species rules.
- owner and clinic appointment scoping and validation.
- clinic pet lookup route and limited response service behavior.
- clinic patient derivation and cross-clinic not-found behavior at service level.
- medical record patient validation, appointment validation, duplicate appointment conflict, clinic scoping, and patch immutability at service level.
- migration SQL string-level safety checks.

Important gaps:

- No real PostgreSQL repository integration tests.
- No migration execution test against a fresh and dirty PostgreSQL database.
- No concurrency/race test for duplicate medical record creation for one appointment.
- No brute-force/rate-limit tests because no feature exists.
- No tests for production startup rejecting default secrets because no guard exists.
- No full end-to-end HTTP + real DB integration suite in CI.
- No test coverage for server timeouts/body-size limits because no feature exists.

## 11. Documentation gaps

Current backend docs are useful and mostly aligned with implemented Sprints 1-10. Gaps:

- Older planning docs remain and can conflict with current code if read as current truth.
- The frontend integration guide is explicitly through Sprint 1-5 and does not cover Sprint 6-10.
- API naming inconsistency is documented but not resolved.
- Security limitations are mentioned in places, but the docs should more prominently warn that public clinic registration and owner approval are not production-safe yet.
- There is no single security hardening checklist for pre-public deployment.
- Migration docs describe startup migrations but do not include a versioned migration strategy.

## 12. Recommended action order

1. Fix C-01: fail startup outside development if `JWT_SECRET` is default/weak.
2. Fix C-02/H-01: define clinic onboarding and owner authorization before public clinic access.
3. Reduce H-03: mask owner phone in medical record responses or guard it behind explicit authorization.
4. Add H-02/H-04 hardening: rate limiting, request size limits, and HTTP timeouts.
5. Add PostgreSQL integration tests for migrations, appointments, clinic patients, and medical records.
6. Add context propagation to repository calls.
7. Add scoped update write queries for defense-in-depth.
8. Introduce a migration runner/version table before the next schema-heavy sprint.
9. Standardize API JSON naming for future endpoints.
10. Split current docs into "current implementation" and "historical planning" more clearly.
