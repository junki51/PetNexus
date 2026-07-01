package services

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

const (
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

// AuthService owns authentication validation and business rules.
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	GetCurrentUser(userID string) (*dto.UserResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      config.Config
}

// NewAuthService creates an authentication service with explicit dependencies.
func NewAuthService(userRepo repositories.UserRepository, cfg config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	phone := strings.TrimSpace(req.Phone)
	role := strings.ToLower(strings.TrimSpace(req.Role))

	if err := validateRegistration(email, phone, req.Password, role); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, internalServerError(err)
	}
	if exists {
		return nil, emailAlreadyExistsError()
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, internalServerError(err)
	}

	user := &models.User{
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role:         role,
	}
	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repositories.ErrEmailAlreadyExists) {
			return nil, emailAlreadyExistsError()
		}
		return nil, internalServerError(err)
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiresIn,
	)
	if err != nil {
		return nil, internalServerError(err)
	}

	return &dto.AuthResponse{
		User:        toUserResponse(user),
		AccessToken: accessToken,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	if email == "" {
		return nil, validationError("Email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, validationError("Password is required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, invalidCredentialsError()
		}
		return nil, internalServerError(err)
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, invalidCredentialsError()
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiresIn,
	)
	if err != nil {
		return nil, internalServerError(err)
	}

	return &dto.AuthResponse{
		User:        toUserResponse(user),
		AccessToken: accessToken,
	}, nil
}

func (s *authService) GetCurrentUser(userID string) (*dto.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Unauthorized",
			"Invalid authenticated user",
			err,
		)
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, utils.NewAppError(
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"User not found",
				"Authenticated user no longer exists",
				err,
			)
		}
		return nil, internalServerError(err)
	}

	response := toUserResponse(user)
	return &response, nil
}

func validateRegistration(email, phone, password, role string) error {
	if email == "" {
		return validationError("Email is required")
	}
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return validationError("Email format is invalid")
	}
	if strings.TrimSpace(password) == "" {
		return validationError("Password is required")
	}
	if utf8.RuneCountInString(password) < minPasswordLength {
		return validationError("Password must be at least 8 characters")
	}
	if len([]byte(password)) > maxPasswordBytes {
		return validationError("Password must not exceed 72 bytes")
	}
	if len(phone) > 30 {
		return validationError("Phone must not exceed 30 characters")
	}
	if role == "" {
		return validationError("Role is required")
	}
	if role != models.RoleOwner && role != models.RoleClinicStaff {
		return utils.NewAppError(
			http.StatusForbidden,
			"FORBIDDEN_ROLE",
			"Role is not allowed",
			"Public registration allows only owner or clinic_staff",
			nil,
		)
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func validationError(details string) *utils.AppError {
	return utils.NewAppError(
		http.StatusUnprocessableEntity,
		"VALIDATION_ERROR",
		"Validation failed",
		details,
		nil,
	)
}

func emailAlreadyExistsError() *utils.AppError {
	return utils.NewAppError(
		http.StatusConflict,
		"EMAIL_ALREADY_EXISTS",
		"Email already exists",
		"An account with this email already exists",
		nil,
	)
}

func invalidCredentialsError() *utils.AppError {
	return utils.NewAppError(
		http.StatusUnauthorized,
		"INVALID_CREDENTIALS",
		"Invalid credentials",
		"Email or password is incorrect",
		nil,
	)
}

func internalServerError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"Something went wrong",
		"An internal server error occurred",
		cause,
	)
}
