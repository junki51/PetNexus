package services

import (
	"errors"
	"net/http"
	"net/url"
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

const dateOnlyLayout = "2006-01-02"

var allowedOwnerGenders = map[string]struct{}{
	"male":              {},
	"female":            {},
	"prefer_not_to_say": {},
	"other":             {},
}

// OwnerProfileService owns owner profile validation and business rules.
type OwnerProfileService interface {
	CreateProfile(currentUserID string, req dto.CreateOwnerProfileRequest) (*dto.OwnerProfileResponse, error)
	GetProfile(currentUserID string) (*dto.OwnerProfileResponse, error)
	UpdateProfile(currentUserID string, req dto.UpdateOwnerProfileRequest) (*dto.OwnerProfileResponse, error)
}

type ownerProfileService struct {
	profileRepo repositories.OwnerProfileRepository
}

// NewOwnerProfileService creates an owner profile service.
func NewOwnerProfileService(profileRepo repositories.OwnerProfileRepository) OwnerProfileService {
	return &ownerProfileService{profileRepo: profileRepo}
}

func (s *ownerProfileService) CreateProfile(currentUserID string, req dto.CreateOwnerProfileRequest) (*dto.OwnerProfileResponse, error) {
	userID, err := parseOwnerUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	profile, err := buildOwnerProfile(userID, req, time.Now())
	if err != nil {
		return nil, err
	}

	exists, err := s.profileRepo.ExistsByUserID(userID)
	if err != nil {
		return nil, internalServerError(err)
	}
	if exists {
		return nil, ownerProfileAlreadyExistsError()
	}

	if err := s.profileRepo.Create(profile); err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileAlreadyExists) {
			return nil, ownerProfileAlreadyExistsError()
		}
		return nil, internalServerError(err)
	}

	response := toOwnerProfileResponse(profile)
	return &response, nil
}

func (s *ownerProfileService) GetProfile(currentUserID string) (*dto.OwnerProfileResponse, error) {
	userID, err := parseOwnerUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileNotFound) {
			return nil, ownerProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toOwnerProfileResponse(profile)
	return &response, nil
}

func (s *ownerProfileService) UpdateProfile(currentUserID string, req dto.UpdateOwnerProfileRequest) (*dto.OwnerProfileResponse, error) {
	userID, err := parseOwnerUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	if !hasOwnerProfileUpdate(req) {
		return nil, ownerValidationError("Request body must contain at least one profile field")
	}

	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileNotFound) {
			return nil, ownerProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	if err := applyOwnerProfileUpdate(profile, req, time.Now()); err != nil {
		return nil, err
	}
	if err := s.profileRepo.Update(profile); err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileNotFound) {
			return nil, ownerProfileNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toOwnerProfileResponse(profile)
	return &response, nil
}

func buildOwnerProfile(userID uuid.UUID, req dto.CreateOwnerProfileRequest, now time.Time) (*models.OwnerProfile, error) {
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)
	phoneNumber := strings.TrimSpace(req.PhoneNumber)
	if err := validateRequiredOwnerField("first_name", firstName, 100); err != nil {
		return nil, err
	}
	if err := validateRequiredOwnerField("last_name", lastName, 100); err != nil {
		return nil, err
	}
	if err := validateRequiredOwnerField("phone_number", phoneNumber, 30); err != nil {
		return nil, err
	}

	gender, err := normalizeGender(req.Gender)
	if err != nil {
		return nil, err
	}
	dateOfBirth, err := parseDateOfBirth(req.DateOfBirth, now)
	if err != nil {
		return nil, err
	}
	avatarURL, err := normalizeAvatarURL(req.AvatarURL)
	if err != nil {
		return nil, err
	}

	addressLine1, err := normalizeOptionalField("address_line1", req.AddressLine1, 255)
	if err != nil {
		return nil, err
	}
	addressLine2, err := normalizeOptionalField("address_line2", req.AddressLine2, 255)
	if err != nil {
		return nil, err
	}
	province, err := normalizeOptionalField("province", req.Province, 100)
	if err != nil {
		return nil, err
	}
	district, err := normalizeOptionalField("district", req.District, 100)
	if err != nil {
		return nil, err
	}
	subdistrict, err := normalizeOptionalField("subdistrict", req.Subdistrict, 100)
	if err != nil {
		return nil, err
	}
	postalCode, err := normalizeOptionalField("postal_code", req.PostalCode, 20)
	if err != nil {
		return nil, err
	}

	return &models.OwnerProfile{
		UserID:       userID,
		FirstName:    firstName,
		LastName:     lastName,
		Gender:       gender,
		DateOfBirth:  dateOfBirth,
		PhoneNumber:  phoneNumber,
		AvatarURL:    avatarURL,
		AddressLine1: addressLine1,
		AddressLine2: addressLine2,
		Province:     province,
		District:     district,
		Subdistrict:  subdistrict,
		PostalCode:   postalCode,
	}, nil
}

func applyOwnerProfileUpdate(profile *models.OwnerProfile, req dto.UpdateOwnerProfileRequest, now time.Time) error {
	if req.FirstName != nil {
		value := strings.TrimSpace(*req.FirstName)
		if err := validateRequiredOwnerField("first_name", value, 100); err != nil {
			return err
		}
		profile.FirstName = value
	}
	if req.LastName != nil {
		value := strings.TrimSpace(*req.LastName)
		if err := validateRequiredOwnerField("last_name", value, 100); err != nil {
			return err
		}
		profile.LastName = value
	}
	if req.PhoneNumber != nil {
		value := strings.TrimSpace(*req.PhoneNumber)
		if err := validateRequiredOwnerField("phone_number", value, 30); err != nil {
			return err
		}
		profile.PhoneNumber = value
	}
	if req.Gender != nil {
		value, err := normalizeGender(*req.Gender)
		if err != nil {
			return err
		}
		profile.Gender = value
	}
	if req.DateOfBirth != nil {
		value, err := parseDateOfBirth(*req.DateOfBirth, now)
		if err != nil {
			return err
		}
		profile.DateOfBirth = value
	}
	if req.AvatarURL != nil {
		value, err := normalizeAvatarURL(*req.AvatarURL)
		if err != nil {
			return err
		}
		profile.AvatarURL = value
	}

	optionalUpdates := []struct {
		input  *string
		field  string
		limit  int
		target **string
	}{
		{req.AddressLine1, "address_line1", 255, &profile.AddressLine1},
		{req.AddressLine2, "address_line2", 255, &profile.AddressLine2},
		{req.Province, "province", 100, &profile.Province},
		{req.District, "district", 100, &profile.District},
		{req.Subdistrict, "subdistrict", 100, &profile.Subdistrict},
		{req.PostalCode, "postal_code", 20, &profile.PostalCode},
	}
	for _, update := range optionalUpdates {
		if update.input == nil {
			continue
		}
		value, err := normalizeOptionalField(update.field, *update.input, update.limit)
		if err != nil {
			return err
		}
		*update.target = value
	}

	return nil
}

func hasOwnerProfileUpdate(req dto.UpdateOwnerProfileRequest) bool {
	return req.FirstName != nil || req.LastName != nil || req.Gender != nil ||
		req.DateOfBirth != nil || req.PhoneNumber != nil || req.AvatarURL != nil ||
		req.AddressLine1 != nil || req.AddressLine2 != nil || req.Province != nil ||
		req.District != nil || req.Subdistrict != nil || req.PostalCode != nil
}

func validateRequiredOwnerField(field, value string, maxLength int) error {
	if value == "" {
		return ownerValidationError(field + " is required")
	}
	if utf8.RuneCountInString(value) > maxLength {
		return ownerValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	return nil
}

func normalizeOptionalField(field, value string, maxLength int) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if utf8.RuneCountInString(trimmed) > maxLength {
		return nil, ownerValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	if trimmed == "" {
		return nil, nil
	}
	return &trimmed, nil
}

func normalizeGender(value string) (*string, error) {
	gender := strings.ToLower(strings.TrimSpace(value))
	if gender == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(gender) > 30 {
		return nil, ownerValidationError("gender must not exceed 30 characters")
	}
	if _, allowed := allowedOwnerGenders[gender]; !allowed {
		return nil, ownerValidationError("gender must be male, female, prefer_not_to_say, or other")
	}
	return &gender, nil
}

func parseDateOfBirth(value string, now time.Time) (*time.Time, error) {
	dateText := strings.TrimSpace(value)
	if dateText == "" {
		return nil, nil
	}
	date, err := time.Parse(dateOnlyLayout, dateText)
	if err != nil {
		return nil, ownerValidationError("date_of_birth must use YYYY-MM-DD format")
	}
	today, _ := time.Parse(dateOnlyLayout, now.UTC().Format(dateOnlyLayout))
	if date.After(today) {
		return nil, ownerValidationError("date_of_birth must not be in the future")
	}
	return &date, nil
}

func normalizeAvatarURL(value string) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, ownerValidationError("avatar_url must be a valid HTTP or HTTPS URL")
	}
	return &trimmed, nil
}

func parseOwnerUserID(value string) (uuid.UUID, error) {
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

func toOwnerProfileResponse(profile *models.OwnerProfile) dto.OwnerProfileResponse {
	var dateOfBirth *string
	if profile.DateOfBirth != nil {
		formatted := profile.DateOfBirth.Format(dateOnlyLayout)
		dateOfBirth = &formatted
	}

	return dto.OwnerProfileResponse{
		ID:           profile.ID.String(),
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		DisplayName:  strings.TrimSpace(profile.FirstName + " " + profile.LastName),
		Gender:       profile.Gender,
		DateOfBirth:  dateOfBirth,
		PhoneNumber:  profile.PhoneNumber,
		AvatarURL:    profile.AvatarURL,
		AddressLine1: profile.AddressLine1,
		AddressLine2: profile.AddressLine2,
		Province:     profile.Province,
		District:     profile.District,
		Subdistrict:  profile.Subdistrict,
		PostalCode:   profile.PostalCode,
		CreatedAt:    profile.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    profile.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ownerValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func ownerProfileAlreadyExistsError() *utils.AppError {
	return utils.NewAppError(
		http.StatusConflict,
		"OWNER_PROFILE_ALREADY_EXISTS",
		"Owner profile already exists",
		"The authenticated owner already has a profile",
		nil,
	)
}

func ownerProfileNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"OWNER_PROFILE_NOT_FOUND",
		"Owner profile not found",
		"The authenticated owner does not have a profile",
		cause,
	)
}
