# Database migrations

PetNexus keeps explicit PostgreSQL migration files for local/manual use. The
deployed API also runs a guarded startup migration before serving routes.

Startup migration behavior:

- creates `pgcrypto` when needed
- creates the `user_role` enum and missing enum values when needed
- creates the `users` table with raw SQL when needed
- creates `idx_users_email_unique` with `CREATE UNIQUE INDEX IF NOT EXISTS`
- creates the `owner_profiles` table with a guarded foreign key to `users(id)`
- creates the unique `owner_profiles.user_id` index idempotently
- creates `breeds` and `pets` with guarded checks, indexes, and foreign keys
- seeds 8 dog and 8 cat breeds with `ON CONFLICT DO NOTHING`
- creates `clinic_profiles` with an idempotent unique user index and guarded
  foreign key to `users(id)`
- adds unique `pets.public_pet_id`, backfills existing pets, and enforces
  non-null IDs after backfill
- creates `appointments` with guarded ownership foreign keys, checks, and
  calendar lookup indexes
- avoids GORM `AutoMigrate` so startup will not try to drop or
  rewrite missing constraints on an existing database
- stops application startup if migration fails

Current migrations:

1. `001_create_enums.sql` enables `pgcrypto` and creates the `user_role` enum.
2. `002_create_users.sql` creates the `users` table, its email unique index,
   and its role index.
3. `003_create_owner_profiles.sql` creates the Sprint 4 owner profile table,
   its unique user index, and its guarded user foreign key.
4. `004_create_breeds_and_pets.sql` creates the Sprint 5 breed catalog and pet
   tables, adds guarded integrity constraints, and seeds the initial breeds.
5. `005_create_clinic_profiles.sql` creates the Sprint 6 clinic profile table,
   unique user index, and guarded user foreign key.
6. `006_add_public_pet_id.sql` adds the Sprint 7 public pet identifier, safely
   backfills existing pets, creates its unique index, and enforces not-null.
7. `007_create_appointments.sql` creates the Sprint 8 appointment calendar
   foundation, indexes, foreign keys, and allowed-value checks.

The SQL files can still be applied manually in numeric order. PowerShell commands are documented in the project `README.md`.

Do not create ad-hoc tables before reviewing `docs/database-plan.md`. Future migrations must preserve numeric ordering and stay inside the active Sprint scope. A dedicated runner such as `golang-migrate` can be introduced later.
