package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

func TestRegistrationAllowsOwnerAndClinicRoles(t *testing.T) {
	for _, role := range []string{models.RoleOwner, models.RoleClinic, models.RoleClinicStaff} {
		if err := validateRegistration("user@example.com", "", "password123", role); err != nil {
			t.Fatalf("validateRegistration() role %q error = %v", role, err)
		}
	}
}

func TestRegistrationStillRejectsPublicAdmin(t *testing.T) {
	err := validateRegistration("admin@example.com", "", "password123", models.RoleAdmin)
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != http.StatusForbidden || appErr.Code != "FORBIDDEN_ROLE" {
		t.Fatalf("status/code = %d/%s, want 403/FORBIDDEN_ROLE", appErr.HTTPStatus, appErr.Code)
	}
}
