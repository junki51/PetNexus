package services

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

const (
	maxClinicNameLength    = 200
	maxClinicPhoneLength   = 30
	maxClinicEmailLength   = 255
	maxClinicAddressLength = 1000
)

// ClinicProfileService owns clinic profile validation and business rules.
type ClinicProfileService interface {
	CreateClinicProfile(currentUserID string, req dto.CreateClinicProfileRequest) (*dto.ClinicProfileResponse, error)
	GetMyClinicProfile(currentUserID string) (*dto.ClinicProfileResponse, error)
	UpdateMyClinicProfile(currentUserID string, req dto.UpdateClinicProfileRequest) (*dto.ClinicProfileResponse, error)
}

type clinicProfileService struct {
	profileRepo repositories.ClinicProfileRepository
}

func NewClinicProfileService(profileRepo repositories.ClinicProfileRepository) ClinicProfileService {
	return &clinicProfileService{profileRepo: profileRepo}
}

func (s *clinicProfileService) CreateClinicProfile(currentUserID string, req dto.CreateClinicProfileRequest) (*dto.ClinicProfileResponse, error) {
	userID, err := parseClinicUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	profile, err := buildClinicProfile(userID, req)
	if err != nil {
		return nil, err
	}

	exists, err := s.profileRepo.ExistsByUserID(userID)
	if err != nil {
		return nil, internalServerError(err)
	}
	if exists {
		return nil, clinicProfileAlreadyExistsError()
	}

	if err := s.profileRepo.Create(profile); err != nil {
		if errors.Is(err, repositories.ErrClinicProfileAlreadyExists) {
			return nil, clinicProfileAlreadyExistsError()
		}
		return nil, internalServerError(err)
	}

	response := toClinicProfileResponse(profile)
	return &response, nil
}

func (s *clinicProfileService) GetMyClinicProfile(currentUserID string) (*dto.ClinicProfileResponse, error) {
	userID, err := parseClinicUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return nil, clinicProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toClinicProfileResponse(profile)
	return &response, nil
}

func (s *clinicProfileService) UpdateMyClinicProfile(currentUserID string, req dto.UpdateClinicProfileRequest) (*dto.ClinicProfileResponse, error) {
	userID, err := parseClinicUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	if !hasClinicProfileUpdate(req) {
		return nil, clinicProfileValidationError("Request body must contain at least one clinic profile field")
	}

	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return nil, clinicProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	if err := applyClinicProfileUpdate(profile, req); err != nil {
		return nil, err
	}
	if err := s.profileRepo.Update(profile); err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return nil, clinicProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toClinicProfileResponse(profile)
	return &response, nil
}

func buildClinicProfile(userID uuid.UUID, req dto.CreateClinicProfileRequest) (*models.ClinicProfile, error) {
	clinicName, err := normalizeRequiredClinicField("clinic_name", req.ClinicName, maxClinicNameLength)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := normalizeOptionalClinicField("phone_number", req.PhoneNumber, maxClinicPhoneLength)
	if err != nil {
		return nil, err
	}
	email, err := normalizeClinicEmail(req.Email)
	if err != nil {
		return nil, err
	}
	address, err := normalizeOptionalClinicField("address", req.Address, maxClinicAddressLength)
	if err != nil {
		return nil, err
	}

	return &models.ClinicProfile{
		UserID:      userID,
		ClinicName:  clinicName,
		PhoneNumber: phoneNumber,
		Email:       email,
		Address:     address,
	}, nil
}

func applyClinicProfileUpdate(profile *models.ClinicProfile, req dto.UpdateClinicProfileRequest) error {
	if req.ClinicName != nil {
		value, err := normalizeRequiredClinicField("clinic_name", *req.ClinicName, maxClinicNameLength)
		if err != nil {
			return err
		}
		profile.ClinicName = value
	}
	if req.PhoneNumber != nil {
		value, err := normalizeOptionalClinicField("phone_number", *req.PhoneNumber, maxClinicPhoneLength)
		if err != nil {
			return err
		}
		profile.PhoneNumber = value
	}
	if req.Email != nil {
		value, err := normalizeClinicEmail(*req.Email)
		if err != nil {
			return err
		}
		profile.Email = value
	}
	if req.Address != nil {
		value, err := normalizeOptionalClinicField("address", *req.Address, maxClinicAddressLength)
		if err != nil {
			return err
		}
		profile.Address = value
	}
	return nil
}

func hasClinicProfileUpdate(req dto.UpdateClinicProfileRequest) bool {
	return req.ClinicName != nil || req.PhoneNumber != nil || req.Email != nil || req.Address != nil
}

func normalizeRequiredClinicField(field, value string, maxLength int) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", clinicProfileValidationError(field + " is required")
	}
	if utf8.RuneCountInString(trimmed) > maxLength {
		return "", clinicProfileValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	return trimmed, nil
}

func normalizeOptionalClinicField(field, value string, maxLength int) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if utf8.RuneCountInString(trimmed) > maxLength {
		return nil, clinicProfileValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	if trimmed == "" {
		return nil, nil
	}
	return &trimmed, nil
}

func normalizeClinicEmail(value string) (*string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	if email == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(email) > maxClinicEmailLength {
		return nil, clinicProfileValidationError("email must not exceed 255 characters")
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return nil, clinicProfileValidationError("email format is invalid")
	}
	return &email, nil
}

func parseClinicUserID(value string) (uuid.UUID, error) {
	userID, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Unauthorized",
			"Invalid authenticated user",
			err,
		)
	}
	return userID, nil
}

func toClinicProfileResponse(profile *models.ClinicProfile) dto.ClinicProfileResponse {
	return dto.ClinicProfileResponse{
		ID:          profile.ID.String(),
		ClinicName:  profile.ClinicName,
		PhoneNumber: profile.PhoneNumber,
		Email:       profile.Email,
		Address:     profile.Address,
		CreatedAt:   profile.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   profile.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func clinicProfileValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func clinicProfileAlreadyExistsError() *utils.AppError {
	return utils.NewAppError(
		http.StatusConflict,
		"CLINIC_PROFILE_ALREADY_EXISTS",
		"Clinic profile already exists",
		"The authenticated clinic user already has a profile",
		nil,
	)
}

func clinicProfileNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"CLINIC_PROFILE_NOT_FOUND",
		"Clinic profile not found",
		"The authenticated clinic user does not have a profile",
		cause,
	)
}
