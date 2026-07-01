# Database migrations

PetNexus uses explicit PostgreSQL migrations and does not use GORM `AutoMigrate`.

Sprint 3 migrations:

1. `001_create_enums.sql` enables `pgcrypto` and creates the `user_role` enum.
2. `002_create_users.sql` creates the `users` table and its role index.

Apply them in numeric order. PowerShell commands are documented in the project `README.md`.

Do not create ad-hoc tables before reviewing `docs/database-plan.md`. Future migrations must preserve numeric ordering and stay inside the active Sprint scope. A dedicated runner such as `golang-migrate` can be introduced later.
