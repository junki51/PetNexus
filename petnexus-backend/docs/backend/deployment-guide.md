# Deployment Guide

## Local backend

Start PostgreSQL:

```powershell
docker compose up -d
docker compose ps
```

Run the backend:

```powershell
go run ./cmd/api
```

Default local API URL:

```text
http://localhost:8080
```

The backend connects to PostgreSQL and runs guarded startup migrations before
Gin begins serving routes. Startup fails immediately if connection or migration
fails.

## Local PostgreSQL

Docker Compose currently provides one PostgreSQL service with database
`petnexus` on port 5432. Local defaults come from config and `.env.example`.
Do not reuse local credentials as production secrets.

## Render deployment

- Backend runtime: Render Web Service
- Production database: Render Postgres
- Current API base URL: `https://petnexus-api.onrender.com`
- Production database connection: `DATABASE_URL`

When `DATABASE_URL` is non-empty, it takes precedence over separate local
`DB_*` fields. The application does not log the DSN or secrets.

## Environment names

The current config reads:

- `APP_ENV`
- `PORT`
- `DATABASE_URL`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`
- `JWT_SECRET`
- `JWT_EXPIRES_IN`

Set real values in the deployment platform. Never document or commit production
values for passwords, database URL, or JWT secret.

## Startup migration behavior

Startup safely ensures:

- `pgcrypto`
- `user_role` enum values
- `users`
- `owner_profiles`
- `breeds` and seed data
- `pets`
- unique public pet IDs and backfill for existing pets
- `clinic_profiles`
- `appointments`, its calendar indexes, ownership foreign keys, and checks
- required indexes, checks, and foreign keys

Migration SQL uses `IF NOT EXISTS` or guarded PostgreSQL blocks. See
[Database schema](./database-schema.md) and `migrations/README.md`.

## Deploy checklist

1. Run `gofmt` for changed Go files.
2. Run `go test ./...` and preferably `go vet ./...`.
3. Review migration safety and `git diff`.
4. Push the intended commit.
5. Let Render deploy/redeploy the Web Service.
6. Confirm logs show database connection and migration success.
7. Verify:

```powershell
$baseUrl = "https://petnexus-api.onrender.com"
Invoke-RestMethod -Method GET "$baseUrl/health"
Invoke-RestMethod -Method GET "$baseUrl/health/db"
```

8. Run the functional smoke flow in [Testing guide](./testing-guide.md).

## Troubleshooting

- `/health` fails: check service deployment/start command and `PORT`.
- `/health` works but `/health/db` fails: check Render Postgres state and
  `DATABASE_URL` binding.
- Startup exits during migration: inspect the named migration step in logs;
  never work around it by dropping production constraints blindly.
- Auth token failures after deployment: confirm consistent `JWT_SECRET` and
  expiration configuration without printing their values.
