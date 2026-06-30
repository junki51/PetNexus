# Database migrations

PetNexus uses PostgreSQL through GORM. Sprint 2 only establishes and verifies the database connection.

No tables, SQL migrations, or GORM AutoMigrate calls belong in this sprint. Versioned migrations will be introduced later with `golang-migrate`.

Do not create ad-hoc tables before reviewing `docs/database-plan.md`. Migration order and schema must follow that plan.
