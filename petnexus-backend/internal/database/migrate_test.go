package database

import (
	"strings"
	"testing"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

func TestUserRoleEnumMatchesAuthRoles(t *testing.T) {
	for _, role := range []string{models.RoleOwner, models.RoleClinicStaff, models.RoleAdmin} {
		if !strings.Contains(allMigrationSQL(), "'"+role+"'") {
			t.Fatalf("user_role enum SQL does not include %q", role)
		}
	}
}

func TestUserRoleEnumIncludesClinicAlias(t *testing.T) {
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

func allMigrationSQL() string {
	var builder strings.Builder
	for _, step := range migrationSteps {
		builder.WriteString(step.sql)
		builder.WriteByte('\n')
	}
	return builder.String()
}
