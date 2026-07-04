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

const createOwnerProfilesTableSQL = `
CREATE TABLE IF NOT EXISTS owner_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    gender VARCHAR(30),
    date_of_birth DATE,
    phone_number VARCHAR(30) NOT NULL,
    avatar_url TEXT,
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    province VARCHAR(100),
    district VARCHAR(100),
    subdistrict VARCHAR(100),
    postal_code VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createOwnerProfilesUserIDUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_owner_profiles_user_id_unique
ON owner_profiles(user_id);`

const ensureOwnerProfilesUserForeignKeySQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_owner_profiles_user'
          AND conrelid = 'owner_profiles'::regclass
    ) THEN
        ALTER TABLE owner_profiles
        ADD CONSTRAINT fk_owner_profiles_user
        FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END
$$;`

const createBreedsTableSQL = `
CREATE TABLE IF NOT EXISTS breeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    species VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    name_th VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createBreedsSpeciesNameUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_breeds_species_name_unique
ON breeds(species, name);`

const createBreedsSpeciesIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_breeds_species ON breeds(species);`

const ensureBreedsSpeciesCheckSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_breeds_species'
          AND conrelid = 'breeds'::regclass
    ) THEN
        ALTER TABLE breeds
        ADD CONSTRAINT chk_breeds_species CHECK (species IN ('dog', 'cat'));
    END IF;
END
$$;`

const seedBreedsSQL = `
INSERT INTO breeds (species, name, name_th) VALUES
    ('dog', 'Golden Retriever', NULL),
    ('dog', 'Labrador Retriever', NULL),
    ('dog', 'Poodle', NULL),
    ('dog', 'Shiba Inu', NULL),
    ('dog', 'Siberian Husky', NULL),
    ('dog', 'Chihuahua', NULL),
    ('dog', 'Pomeranian', NULL),
    ('dog', 'Thai Bangkaew', NULL),
    ('cat', 'Persian', NULL),
    ('cat', 'Scottish Fold', NULL),
    ('cat', 'British Shorthair', NULL),
    ('cat', 'Siamese', NULL),
    ('cat', 'Maine Coon', NULL),
    ('cat', 'Ragdoll', NULL),
    ('cat', 'Sphynx', NULL),
    ('cat', 'Domestic Shorthair', NULL)
ON CONFLICT (species, name) DO NOTHING;`

const createPetsTableSQL = `
CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_profile_id UUID NOT NULL,
    breed_id UUID,
    species VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(30),
    date_of_birth DATE,
    weight_kg NUMERIC(6,2),
    microchip_id VARCHAR(100),
    avatar_url TEXT,
    color VARCHAR(100),
    distinctive_marks TEXT,
    is_neutered BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createPetsOwnerProfileIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_pets_owner_profile_id ON pets(owner_profile_id);`

const createPetsBreedIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_pets_breed_id ON pets(breed_id);`

const createPetsSpeciesIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_pets_species ON pets(species);`

const ensurePetsOwnerProfileForeignKeySQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_pets_owner_profile'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT fk_pets_owner_profile
        FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id);
    END IF;
END
$$;`

const ensurePetsBreedForeignKeySQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_pets_breed'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT fk_pets_breed
        FOREIGN KEY (breed_id) REFERENCES breeds(id);
    END IF;
END
$$;`

const ensurePetsSpeciesCheckSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_pets_species'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT chk_pets_species CHECK (species IN ('dog', 'cat'));
    END IF;
END
$$;`

const ensurePetsGenderCheckSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_pets_gender'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT chk_pets_gender
        CHECK (gender IS NULL OR gender IN ('male', 'female', 'unknown'));
    END IF;
END
$$;`

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
	{name: "ensure owner_profiles table", sql: createOwnerProfilesTableSQL},
	{name: "ensure owner_profiles user unique index", sql: createOwnerProfilesUserIDUniqueIndexSQL},
	{name: "ensure owner_profiles user foreign key", sql: ensureOwnerProfilesUserForeignKeySQL},
	{name: "ensure breeds table", sql: createBreedsTableSQL},
	{name: "ensure breeds species-name unique index", sql: createBreedsSpeciesNameUniqueIndexSQL},
	{name: "ensure breeds species index", sql: createBreedsSpeciesIndexSQL},
	{name: "ensure breeds species check", sql: ensureBreedsSpeciesCheckSQL},
	{name: "seed breeds", sql: seedBreedsSQL},
	{name: "ensure pets table", sql: createPetsTableSQL},
	{name: "ensure pets owner profile index", sql: createPetsOwnerProfileIndexSQL},
	{name: "ensure pets breed index", sql: createPetsBreedIndexSQL},
	{name: "ensure pets species index", sql: createPetsSpeciesIndexSQL},
	{name: "ensure pets owner profile foreign key", sql: ensurePetsOwnerProfileForeignKeySQL},
	{name: "ensure pets breed foreign key", sql: ensurePetsBreedForeignKeySQL},
	{name: "ensure pets species check", sql: ensurePetsSpeciesCheckSQL},
	{name: "ensure pets gender check", sql: ensurePetsGenderCheckSQL},
}

// RunMigrations creates only the schema required by features that are currently
// implemented. It intentionally avoids GORM AutoMigrate because AutoMigrate
// may try to alter or drop constraints on an existing database.
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
