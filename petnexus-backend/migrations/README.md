# Database migrations

PetNexus will use PostgreSQL. Versioned migrations will be introduced in the database foundation sprint with `golang-migrate`.

Do not create ad-hoc tables or migrations before reviewing `docs/database-plan.md`. Migration order and schema must follow that plan.
