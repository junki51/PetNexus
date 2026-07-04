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

The SQL files can still be applied manually in numeric order. PowerShell commands are documented in the project `README.md`.

Do not create ad-hoc tables before reviewing `docs/database-plan.md`. Future migrations must preserve numeric ordering and stay inside the active Sprint scope. A dedicated runner such as `golang-migrate` can be introduced later.
