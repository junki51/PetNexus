package database

import (
	"strings"
	"testing"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

func TestUserRoleEnumMatchesAuthRoles(t *testing.T) {
	for _, role := range []string{models.RoleOwner, models.RoleClinic, models.RoleClinicStaff, models.RoleAdmin} {
		if !strings.Contains(allMigrationSQL(), "'"+role+"'") {
			t.Fatalf("user_role enum SQL does not include %q", role)
		}
	}
}

func TestUserRoleEnumIncludesCanonicalClinicRole(t *testing.T) {
	if !strings.Contains(allMigrationSQL(), "'clinic'") {
		t.Fatal("user_role enum SQL should include clinic for Render/local compatibility")
	}
}

func TestUsersTableSQLMatchesCurrentUserModel(t *testing.T) {
	for _, fragment := range []string{
		"id UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"email VARCHAR(255) NOT NULL",
		"phone VARCHAR(30)",
		"password_hash TEXT NOT NULL",
		"role user_role NOT NULL",
		"created_at TIMESTAMPTZ NOT NULL DEFAULT now()",
		"updated_at TIMESTAMPTZ NOT NULL DEFAULT now()",
	} {
		if !strings.Contains(createUsersTableSQL, fragment) {
			t.Fatalf("users table SQL does not include %q", fragment)
		}
	}
}

func TestUsersEmailUniqueIndexIsIdempotent(t *testing.T) {
	const want = "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(email);"
	if !strings.Contains(createUsersEmailUniqueIndexSQL, want) {
		t.Fatalf("users email unique index SQL = %q, want it to include %q", createUsersEmailUniqueIndexSQL, want)
	}
}

func TestMigrationSQLDoesNotDropMissingUserEmailConstraint(t *testing.T) {
	migrationSQL := strings.ToUpper(allMigrationSQL())

	if strings.Contains(migrationSQL, "DROP CONSTRAINT") {
		t.Fatal("migration SQL must not drop constraints")
	}
	if strings.Contains(migrationSQL, "UNI_USERS_EMAIL") {
		t.Fatal("migration SQL must not reference the old GORM constraint name uni_users_email")
	}
}

func TestOwnerProfilesTableSQLMatchesCurrentModel(t *testing.T) {
	for _, fragment := range []string{
		"id UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"user_id UUID NOT NULL",
		"first_name VARCHAR(100) NOT NULL",
		"last_name VARCHAR(100) NOT NULL",
		"gender VARCHAR(30)",
		"date_of_birth DATE",
		"phone_number VARCHAR(30) NOT NULL",
		"avatar_url TEXT",
		"address_line1 VARCHAR(255)",
		"address_line2 VARCHAR(255)",
		"province VARCHAR(100)",
		"district VARCHAR(100)",
		"subdistrict VARCHAR(100)",
		"postal_code VARCHAR(20)",
		"created_at TIMESTAMPTZ NOT NULL DEFAULT now()",
		"updated_at TIMESTAMPTZ NOT NULL DEFAULT now()",
	} {
		if !strings.Contains(createOwnerProfilesTableSQL, fragment) {
			t.Fatalf("owner_profiles table SQL does not include %q", fragment)
		}
	}
}

func TestOwnerProfilesConstraintsAreIdempotent(t *testing.T) {
	if !strings.Contains(createOwnerProfilesUserIDUniqueIndexSQL, "CREATE UNIQUE INDEX IF NOT EXISTS idx_owner_profiles_user_id_unique") {
		t.Fatal("owner_profiles user_id unique index must be idempotent")
	}
	if !strings.Contains(ensureOwnerProfilesUserForeignKeySQL, "IF NOT EXISTS") ||
		!strings.Contains(ensureOwnerProfilesUserForeignKeySQL, "FOREIGN KEY (user_id) REFERENCES users(id)") {
		t.Fatal("owner_profiles foreign key must be guarded and reference users(id)")
	}
}

func TestBreedsMigrationIsSafeAndSeedsExpectedBreeds(t *testing.T) {
	for _, fragment := range []string{
		"CREATE TABLE IF NOT EXISTS breeds",
		"species VARCHAR(20) NOT NULL",
		"name VARCHAR(100) NOT NULL",
		"name_th VARCHAR(100)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_breeds_species_name_unique",
		"CREATE INDEX IF NOT EXISTS idx_breeds_species",
		"CHECK (species IN ('dog', 'cat'))",
		"ON CONFLICT (species, name) DO NOTHING",
	} {
		if !strings.Contains(allMigrationSQL(), fragment) {
			t.Fatalf("breed migration SQL does not include %q", fragment)
		}
	}
	if strings.Count(seedBreedsSQL, "('dog',") != 8 || strings.Count(seedBreedsSQL, "('cat',") != 8 {
		t.Fatal("breed seed must include exactly 8 dog and 8 cat breeds")
	}
}

func TestPetsMigrationMatchesModelAndUsesGuardedConstraints(t *testing.T) {
	for _, fragment := range []string{
		"CREATE TABLE IF NOT EXISTS pets",
		"owner_profile_id UUID NOT NULL",
		"breed_id UUID",
		"species VARCHAR(20) NOT NULL",
		"name VARCHAR(100) NOT NULL",
		"gender VARCHAR(30)",
		"date_of_birth DATE",
		"weight_kg NUMERIC(6,2)",
		"microchip_id VARCHAR(100)",
		"avatar_url TEXT",
		"color VARCHAR(100)",
		"distinctive_marks TEXT",
		"is_neutered BOOLEAN",
	} {
		if !strings.Contains(createPetsTableSQL, fragment) {
			t.Fatalf("pets table SQL does not include %q", fragment)
		}
	}
	migrationSQL := allMigrationSQL()
	for _, fragment := range []string{
		"CREATE INDEX IF NOT EXISTS idx_pets_owner_profile_id",
		"CREATE INDEX IF NOT EXISTS idx_pets_breed_id",
		"CREATE INDEX IF NOT EXISTS idx_pets_species",
		"FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id)",
		"FOREIGN KEY (breed_id) REFERENCES breeds(id)",
		"CHECK (species IN ('dog', 'cat'))",
		"gender IN ('male', 'female', 'unknown')",
	} {
		if !strings.Contains(migrationSQL, fragment) {
			t.Fatalf("pet migration SQL does not include %q", fragment)
		}
	}
}

func TestPetsPublicPetIDMigrationIsSafeAndIdempotent(t *testing.T) {
	for _, fragment := range []string{
		"ALTER TABLE pets ADD COLUMN IF NOT EXISTS public_pet_id VARCHAR(50)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_pets_public_pet_id_unique",
		"ON pets(public_pet_id)",
		"ALTER TABLE pets ALTER COLUMN public_pet_id SET NOT NULL",
	} {
		if !strings.Contains(allMigrationSQL(), fragment) {
			t.Fatalf("public pet ID migration SQL does not include %q", fragment)
		}
	}

	foundBackfill := false
	for _, step := range migrationSteps {
		if step.name == "backfill pets public_pet_id" && step.run != nil {
			foundBackfill = true
			break
		}
	}
	if !foundBackfill {
		t.Fatal("startup migration must include application-level public pet ID backfill")
	}
}

func TestClinicProfilesMigrationMatchesModelAndUsesGuardedConstraints(t *testing.T) {
	for _, fragment := range []string{
		"CREATE TABLE IF NOT EXISTS clinic_profiles",
		"id UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"user_id UUID NOT NULL",
		"clinic_name VARCHAR(200) NOT NULL",
		"phone_number VARCHAR(30)",
		"email VARCHAR(255)",
		"address TEXT",
		"created_at TIMESTAMPTZ NOT NULL DEFAULT now()",
		"updated_at TIMESTAMPTZ NOT NULL DEFAULT now()",
	} {
		if !strings.Contains(createClinicProfilesTableSQL, fragment) {
			t.Fatalf("clinic_profiles table SQL does not include %q", fragment)
		}
	}
	if !strings.Contains(createClinicProfilesUserIDUniqueIndexSQL, "CREATE UNIQUE INDEX IF NOT EXISTS idx_clinic_profiles_user_id_unique") {
		t.Fatal("clinic_profiles user_id unique index must be idempotent")
	}
	if !strings.Contains(ensureClinicProfilesUserForeignKeySQL, "IF NOT EXISTS") ||
		!strings.Contains(ensureClinicProfilesUserForeignKeySQL, "FOREIGN KEY (user_id) REFERENCES users(id)") {
		t.Fatal("clinic_profiles foreign key must be guarded and reference users(id)")
	}
}

func TestAppointmentsMigrationIsSafeAndMatchesSprint8Model(t *testing.T) {
	migrationSQL := allMigrationSQL()
	for _, fragment := range []string{
		"CREATE TABLE IF NOT EXISTS appointments",
		"owner_profile_id UUID NOT NULL",
		"clinic_profile_id UUID NOT NULL",
		"pet_id UUID NOT NULL",
		"title VARCHAR(150)",
		"appointment_type VARCHAR(50) NOT NULL",
		"scheduled_at TIMESTAMPTZ NOT NULL",
		"duration_minutes INTEGER NOT NULL",
		"status VARCHAR(50) NOT NULL",
		"created_by_user_id UUID",
		"created_by_role VARCHAR(20) NOT NULL",
		"cancelled_at TIMESTAMPTZ",
		"CREATE INDEX IF NOT EXISTS idx_appointments_clinic_scheduled_at",
		"CREATE INDEX IF NOT EXISTS idx_appointments_owner_scheduled_at",
		"FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id)",
		"FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id)",
		"FOREIGN KEY (pet_id) REFERENCES pets(id)",
		"FOREIGN KEY (created_by_user_id) REFERENCES users(id)",
		"'checkup', 'vaccination', 'consultation', 'follow_up'",
		"'requested', 'scheduled', 'checked_in', 'completed', 'cancelled'",
		"created_by_role IN ('owner', 'clinic')",
		"duration_minutes BETWEEN 5 AND 480",
	} {
		if !strings.Contains(migrationSQL, fragment) {
			t.Fatalf("appointment migration SQL does not include %q", fragment)
		}
	}
	if strings.Contains(strings.ToUpper(createAppointmentsTableSQL+ensureAppointmentConstraintsSQL), "DROP ") {
		t.Fatal("appointment migration must not drop existing schema objects")
	}
}

func TestMedicalRecordsMigrationIsSafeAndMatchesSprint10Model(t *testing.T) {
	migrationSQL := allMigrationSQL()
	for _, fragment := range []string{
		"CREATE TABLE IF NOT EXISTS medical_records",
		"clinic_profile_id UUID NOT NULL",
		"pet_id UUID NOT NULL",
		"appointment_id UUID",
		"created_by_user_id UUID NOT NULL",
		"visit_at TIMESTAMPTZ NOT NULL",
		"chief_complaint TEXT NOT NULL",
		"clinical_findings TEXT",
		"diagnosis TEXT",
		"treatment_plan TEXT",
		"medications TEXT",
		"follow_up_instructions TEXT",
		"next_follow_up_at TIMESTAMPTZ",
		"weight_kg NUMERIC(6,2)",
		"temperature_c NUMERIC(5,2)",
		"CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_profile_id",
		"CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id",
		"CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_visit_at",
		"CREATE INDEX IF NOT EXISTS idx_medical_records_pet_visit_at",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_medical_records_appointment_id_unique",
		"WHERE appointment_id IS NOT NULL",
		"FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id)",
		"FOREIGN KEY (pet_id) REFERENCES pets(id)",
		"FOREIGN KEY (appointment_id) REFERENCES appointments(id)",
		"FOREIGN KEY (created_by_user_id) REFERENCES users(id)",
		"CHECK (weight_kg IS NULL OR weight_kg > 0)",
		"CHECK (temperature_c IS NULL OR temperature_c > 0)",
		"CHECK (next_follow_up_at IS NULL OR next_follow_up_at >= visit_at)",
	} {
		if !strings.Contains(migrationSQL, fragment) {
			t.Fatalf("medical record migration SQL does not include %q", fragment)
		}
	}
	if strings.Contains(strings.ToUpper(createMedicalRecordsTableSQL+ensureMedicalRecordConstraintsSQL), "DROP ") {
		t.Fatal("medical record migration must not drop existing schema objects")
	}
}

func allMigrationSQL() string {
	var builder strings.Builder
	for _, step := range migrationSteps {
		builder.WriteString(step.sql)
		builder.WriteByte('\n')
	}
	return builder.String()
}
