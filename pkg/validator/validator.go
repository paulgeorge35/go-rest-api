package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required" example:"reset-token-123"`
	Password string `json:"password" binding:"required,min=8" example:"newPassword123"`
}

type MagicLinkRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	// Add custom validations here
	_ = v.RegisterValidation("password", validatePassword)

	return &Validator{
		validate: v,
	}
}

func (v *Validator) ValidateRegisterRequest(req *RegisterRequest) error {
	return v.validate.Struct(req)
}

func (v *Validator) ValidateLoginRequest(req *LoginRequest) error {
	return v.validate.Struct(req)
}

// Custom password validation
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// At least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// At least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// At least one number
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	// At least one special character
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return false
	}

	return true
}

func (v *Validator) ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	return v.validate.Var(email, "required,email") == nil
}
