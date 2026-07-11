package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/utils"
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

const addPetsPublicPetIDColumnSQL = `
ALTER TABLE pets ADD COLUMN IF NOT EXISTS public_pet_id VARCHAR(50);`

const normalizeEmptyPetsPublicPetIDSQL = `
UPDATE pets SET public_pet_id = NULL
WHERE public_pet_id IS NOT NULL AND BTRIM(public_pet_id) = '';`

const createPetsPublicPetIDUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_pets_public_pet_id_unique
ON pets(public_pet_id);`

const ensurePetsPublicPetIDNotNullSQL = `
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'pets'
          AND column_name = 'public_pet_id'
          AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE pets ALTER COLUMN public_pet_id SET NOT NULL;
    END IF;
END
$$;`

const createClinicProfilesTableSQL = `
CREATE TABLE IF NOT EXISTS clinic_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    clinic_name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(30),
    email VARCHAR(255),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createClinicProfilesUserIDUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_clinic_profiles_user_id_unique
ON clinic_profiles(user_id);`

const ensureClinicProfilesUserForeignKeySQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_clinic_profiles_user'
          AND conrelid = 'clinic_profiles'::regclass
    ) THEN
        ALTER TABLE clinic_profiles
        ADD CONSTRAINT fk_clinic_profiles_user
        FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END
$$;`

const createAppointmentsTableSQL = `
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_profile_id UUID NOT NULL,
    clinic_profile_id UUID NOT NULL,
    pet_id UUID NOT NULL,
    title VARCHAR(150),
    appointment_type VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    note TEXT,
    created_by_user_id UUID,
    created_by_role VARCHAR(20) NOT NULL,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createAppointmentIndexesSQL = `
CREATE INDEX IF NOT EXISTS idx_appointments_owner_profile_id
ON appointments(owner_profile_id);
CREATE INDEX IF NOT EXISTS idx_appointments_clinic_profile_id
ON appointments(clinic_profile_id);
CREATE INDEX IF NOT EXISTS idx_appointments_pet_id
ON appointments(pet_id);
CREATE INDEX IF NOT EXISTS idx_appointments_scheduled_at
ON appointments(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_appointments_status
ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_clinic_scheduled_at
ON appointments(clinic_profile_id, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_appointments_owner_scheduled_at
ON appointments(owner_profile_id, scheduled_at);`

const ensureAppointmentConstraintsSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_appointments_owner_profile'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT fk_appointments_owner_profile
        FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_appointments_clinic_profile'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT fk_appointments_clinic_profile
        FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_appointments_pet'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT fk_appointments_pet
        FOREIGN KEY (pet_id) REFERENCES pets(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_appointments_created_by_user'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT fk_appointments_created_by_user
        FOREIGN KEY (created_by_user_id) REFERENCES users(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_appointments_type'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT chk_appointments_type CHECK (
            appointment_type IN (
                'checkup', 'vaccination', 'consultation', 'follow_up',
                'grooming', 'emergency', 'other'
            )
        );
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_appointments_status'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT chk_appointments_status CHECK (
            status IN ('requested', 'scheduled', 'checked_in', 'completed', 'cancelled')
        );
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_appointments_created_by_role'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT chk_appointments_created_by_role
        CHECK (created_by_role IN ('owner', 'clinic'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_appointments_duration'
          AND conrelid = 'appointments'::regclass
    ) THEN
        ALTER TABLE appointments
        ADD CONSTRAINT chk_appointments_duration
        CHECK (duration_minutes BETWEEN 5 AND 480);
    END IF;
END
$$;`

const createMedicalRecordsTableSQL = `
CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_profile_id UUID NOT NULL,
    pet_id UUID NOT NULL,
    appointment_id UUID,
    created_by_user_id UUID NOT NULL,
    visit_at TIMESTAMPTZ NOT NULL,
    chief_complaint TEXT NOT NULL,
    clinical_findings TEXT,
    diagnosis TEXT,
    treatment_plan TEXT,
    medications TEXT,
    follow_up_instructions TEXT,
    next_follow_up_at TIMESTAMPTZ,
    weight_kg NUMERIC(6,2),
    temperature_c NUMERIC(5,2),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

const createMedicalRecordIndexesSQL = `
CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_profile_id
ON medical_records(clinic_profile_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id
ON medical_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_visit_at
ON medical_records(clinic_profile_id, visit_at DESC);
CREATE INDEX IF NOT EXISTS idx_medical_records_pet_visit_at
ON medical_records(pet_id, visit_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_medical_records_appointment_id_unique
ON medical_records(appointment_id)
WHERE appointment_id IS NOT NULL;`

const ensureMedicalRecordConstraintsSQL = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_clinic_profile'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_clinic_profile
        FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_pet'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_pet
        FOREIGN KEY (pet_id) REFERENCES pets(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_appointment'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_appointment
        FOREIGN KEY (appointment_id) REFERENCES appointments(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_created_by_user'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_created_by_user
        FOREIGN KEY (created_by_user_id) REFERENCES users(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_weight_positive'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_weight_positive
        CHECK (weight_kg IS NULL OR weight_kg > 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_temperature_positive'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_temperature_positive
        CHECK (temperature_c IS NULL OR temperature_c > 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_next_follow_up_after_visit'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_next_follow_up_after_visit
        CHECK (next_follow_up_at IS NULL OR next_follow_up_at >= visit_at);
    END IF;
END
$$;`

type migrationStep struct {
	name string
	sql  string
	run  func(*gorm.DB) error
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
	{name: "ensure pets public_pet_id column", sql: addPetsPublicPetIDColumnSQL},
	{name: "normalize empty pets public_pet_id", sql: normalizeEmptyPetsPublicPetIDSQL},
	{name: "ensure pets public_pet_id unique index", sql: createPetsPublicPetIDUniqueIndexSQL},
	{name: "backfill pets public_pet_id", run: backfillMissingPublicPetIDs},
	{name: "ensure pets public_pet_id not null", sql: ensurePetsPublicPetIDNotNullSQL},
	{name: "ensure clinic_profiles table", sql: createClinicProfilesTableSQL},
	{name: "ensure clinic_profiles user unique index", sql: createClinicProfilesUserIDUniqueIndexSQL},
	{name: "ensure clinic_profiles user foreign key", sql: ensureClinicProfilesUserForeignKeySQL},
	{name: "ensure appointments table", sql: createAppointmentsTableSQL},
	{name: "ensure appointment indexes", sql: createAppointmentIndexesSQL},
	{name: "ensure appointment constraints", sql: ensureAppointmentConstraintsSQL},
	{name: "ensure medical_records table", sql: createMedicalRecordsTableSQL},
	{name: "ensure medical_records indexes", sql: createMedicalRecordIndexesSQL},
	{name: "ensure medical_records constraints", sql: ensureMedicalRecordConstraintsSQL},
}

// RunMigrations creates only the schema required by features that are currently
// implemented. It intentionally avoids GORM AutoMigrate because AutoMigrate
// may try to alter or drop constraints on an existing database.
func RunMigrations(db *gorm.DB) error {
	for _, step := range migrationSteps {
		if step.run != nil {
			if err := step.run(db); err != nil {
				return fmt.Errorf("%s: %w", step.name, err)
			}
			continue
		}
		if err := db.Exec(step.sql).Error; err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func backfillMissingPublicPetIDs(db *gorm.DB) error {
	var petIDs []string
	if err := db.Table("pets").
		Where("public_pet_id IS NULL OR BTRIM(public_pet_id) = ''").
		Pluck("id", &petIDs).Error; err != nil {
		return fmt.Errorf("list pets missing public pet ID: %w", err)
	}

	for _, petID := range petIDs {
		assigned := false
		for attempt := 0; attempt < 20; attempt++ {
			publicPetID, err := utils.GeneratePublicPetID()
			if err != nil {
				return err
			}
			result := db.Exec(
				"UPDATE pets SET public_pet_id = ? WHERE id = ? AND (public_pet_id IS NULL OR BTRIM(public_pet_id) = '')",
				publicPetID,
				petID,
			)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
					continue
				}
				return fmt.Errorf("backfill pet %s public pet ID: %w", petID, result.Error)
			}
			assigned = result.RowsAffected == 1
			if result.RowsAffected == 0 {
				assigned = true
			}
			break
		}
		if !assigned {
			return fmt.Errorf("backfill pet %s public pet ID: exhausted unique ID attempts", petID)
		}
	}
	return nil
}

// Migrate preserves the previous public function name while delegating to the
// safe SQL migration runner.
func Migrate(db *gorm.DB) error {
	return RunMigrations(db)
}
