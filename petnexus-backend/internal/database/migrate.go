package database

import (
	"fmt"

	"gorm.io/gorm"
)

const createPGCryptoExtensionSQL = `CREATE EXTENSION IF NOT EXISTS pgcrypto;`

const createUserRoleEnumSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'user_role'
          AND n.nspname = current_schema()
    ) THEN
        CREATE TYPE user_role AS ENUM ('owner', 'clinic', 'clinic_staff', 'admin');
    END IF;
END
$$;`

const ensureUserRoleValuesSQL = `
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'user_role'
          AND n.nspname = current_schema()
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            JOIN pg_namespace n ON n.oid = t.typnamespace
            WHERE t.typname = 'user_role'
              AND n.nspname = current_schema()
              AND e.enumlabel = 'owner'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'owner';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            JOIN pg_namespace n ON n.oid = t.typnamespace
            WHERE t.typname = 'user_role'
              AND n.nspname = current_schema()
              AND e.enumlabel = 'clinic'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'clinic';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            JOIN pg_namespace n ON n.oid = t.typnamespace
            WHERE t.typname = 'user_role'
              AND n.nspname = current_schema()
              AND e.enumlabel = 'clinic_staff'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'clinic_staff';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            JOIN pg_namespace n ON n.oid = t.typnamespace
            WHERE t.typname = 'user_role'
              AND n.nspname = current_schema()
              AND e.enumlabel = 'admin'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'admin';
        END IF;
    END IF;
END
$$;`

const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    password_hash TEXT NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createUsersEmailUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(email);`

const createUsersRoleIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`

type migrationStep struct {
	name string
	sql  string
}

var migrationSteps = []migrationStep{
	{name: "ensure pgcrypto extension", sql: createPGCryptoExtensionSQL},
	{name: "ensure user_role enum", sql: createUserRoleEnumSQL},
	{name: "ensure user_role enum values", sql: ensureUserRoleValuesSQL},
	{name: "ensure users table", sql: createUsersTableSQL},
	{name: "ensure users email unique index", sql: createUsersEmailUniqueIndexSQL},
	{name: "ensure users role index", sql: createUsersRoleIndexSQL},
}

// RunMigrations creates only the schema required by features that are currently
// implemented. It intentionally avoids GORM AutoMigrate for users because
// AutoMigrate may try to alter or drop constraints on an existing database.
func RunMigrations(db *gorm.DB) error {
	for _, step := range migrationSteps {
		if err := db.Exec(step.sql).Error; err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

// Migrate preserves the previous public function name while delegating to the
// safe SQL migration runner.
func Migrate(db *gorm.DB) error {
	return RunMigrations(db)
}
