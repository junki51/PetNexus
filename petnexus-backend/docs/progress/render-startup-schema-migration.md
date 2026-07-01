# Render: Safe Startup Schema Migration

Updated: 2026-07-01

## Problem

Render PostgreSQL connected successfully and health checks passed, but auth
register failed when the deployed database did not have the required Sprint 3
schema.

A later startup migration failed on existing databases because GORM
`AutoMigrate` tried to run:

```sql
ALTER TABLE "users" DROP CONSTRAINT "uni_users_email";
```

That constraint did not exist on the current `users` table, so PostgreSQL
returned `SQLSTATE 42704`.

## What changed

- Replaced startup `AutoMigrate` for `users` with explicit idempotent SQL.
- Added `database.RunMigrations(db *gorm.DB) error`.
- Startup migration now runs after PostgreSQL connection and before Gin routes.
- Migration failure still stops the server with a clear log message.
- Migration success logs `database migration completed successfully`.

## Current startup migration scope

Only Sprint 3 Auth/User schema is created:

- `pgcrypto` extension
- `user_role` enum
- `users` table
- `idx_users_email_unique` unique index
- `idx_users_role` role index

`user_role` keeps the current API role `clinic_staff` and also includes
`clinic` for compatibility with deployment instructions. Auth behavior was not
changed.

## Important safety behavior

- Does not call GORM `AutoMigrate` for `users`.
- Does not run `ALTER TABLE users DROP CONSTRAINT uni_users_email`.
- Does not fail if `users` already exists.
- Does not fail if `idx_users_email_unique` already exists.
- Does not log `DATABASE_URL`, database passwords, `JWT_SECRET`, or other
  secrets.

## Not implemented

- Pet
- Breed
- QR
- Clinic access
- Visit
- Timeline
- UI

## Verification

- `gofmt`
- `go test ./...`

