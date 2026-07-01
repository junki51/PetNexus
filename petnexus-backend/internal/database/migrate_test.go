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

func allMigrationSQL() string {
	var builder strings.Builder
	for _, step := range migrationSteps {
		builder.WriteString(step.sql)
		builder.WriteByte('\n')
	}
	return builder.String()
}
