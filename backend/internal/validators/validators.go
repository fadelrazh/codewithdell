package validators

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidator represents custom validation functions
type CustomValidator struct {
	validate *validator.Validate
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	
	// Register custom validation functions
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("slug", validateSlug)
	v.RegisterValidation("hexcolor", validateHexColor)
	v.RegisterValidation("url", validateURL)
	v.RegisterValidation("content", validateContent)
	
	return &CustomValidator{
		validate: v,
	}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	
	// Username must be 3-30 characters long
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	
	// Username can only contain letters, numbers, underscores, and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return false
	}
	
	// Username cannot start or end with underscore or hyphen
	if strings.HasPrefix(username, "_") || strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "_") || strings.HasSuffix(username, "-") {
		return false
	}
	
	// Username cannot contain consecutive underscores or hyphens
	if strings.Contains(username, "__") || strings.Contains(username, "--") ||
		strings.Contains(username, "_-") || strings.Contains(username, "-_") {
		return false
	}
	
	return true
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Password must be at least 8 characters long
	if len(password) < 8 {
		return false
	}
	
	// Password must contain at least one uppercase letter
	hasUpper := false
	// Password must contain at least one lowercase letter
	hasLower := false
	// Password must contain at least one digit
	hasDigit := false
	// Password must contain at least one special character
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasDigit && hasSpecial
}

// validateSlug validates slug format
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	
	// Slug must be 3-100 characters long
	if len(slug) < 3 || len(slug) > 100 {
		return false
	}
	
	// Slug can only contain lowercase letters, numbers, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	if !matched {
		return false
	}
	
	// Slug cannot start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}
	
	// Slug cannot contain consecutive hyphens
	if strings.Contains(slug, "--") {
		return false
	}
	
	return true
}

// validateHexColor validates hex color format
func validateHexColor(fl validator.FieldLevel) bool {
	color := fl.Field().String()
	
	// Hex color must be 3 or 6 characters long (with or without #)
	if len(color) == 0 {
		return true // Allow empty color
	}
	
	// Remove # if present
	if strings.HasPrefix(color, "#") {
		color = color[1:]
	}
	
	// Check if it's a valid hex color
	matched, _ := regexp.MatchString(`^[0-9A-Fa-f]{3}$|^[0-9A-Fa-f]{6}$`, color)
	return matched
}

// validateURL validates URL format
func validateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	
	if len(url) == 0 {
		return true // Allow empty URL
	}
	
	// Basic URL validation
	matched, _ := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, url)
	return matched
}

// validateContent validates content length and format
func validateContent(fl validator.FieldLevel) bool {
	content := fl.Field().String()
	
	// Content must be at least 10 characters long
	if len(content) < 10 {
		return false
	}
	
	// Content must not exceed 50,000 characters
	if len(content) > 50000 {
		return false
	}
	
	// Content must not contain only whitespace
	if strings.TrimSpace(content) == "" {
		return false
	}
	
	return true
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// GetValidationErrors converts validator errors to custom format
func GetValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			error := ValidationError{
				Field: e.Field(),
				Tag:   e.Tag(),
				Value: fmt.Sprintf("%v", e.Value()),
			}
			
			// Custom error messages
			switch e.Tag() {
			case "required":
				error.Message = fmt.Sprintf("%s is required", e.Field())
			case "email":
				error.Message = fmt.Sprintf("%s must be a valid email address", e.Field())
			case "min":
				error.Message = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
			case "max":
				error.Message = fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param())
			case "username":
				error.Message = fmt.Sprintf("%s must be 3-30 characters long and contain only letters, numbers, underscores, and hyphens", e.Field())
			case "password":
				error.Message = fmt.Sprintf("%s must be at least 8 characters long and contain uppercase, lowercase, digit, and special character", e.Field())
			case "slug":
				error.Message = fmt.Sprintf("%s must be 3-100 characters long and contain only lowercase letters, numbers, and hyphens", e.Field())
			case "hexcolor":
				error.Message = fmt.Sprintf("%s must be a valid hex color code", e.Field())
			case "url":
				error.Message = fmt.Sprintf("%s must be a valid URL", e.Field())
			case "content":
				error.Message = fmt.Sprintf("%s must be 10-50,000 characters long", e.Field())
			default:
				error.Message = fmt.Sprintf("%s is invalid", e.Field())
			}
			
			errors = append(errors, error)
		}
	}
	
	return errors
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	input = re.ReplaceAllString(input, "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return phoneRegex
}

// ValidateDate validates date format (YYYY-MM-DD)
func ValidateDate(date string) bool {
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(date) {
		return false
	}
	
	// Additional validation for valid date
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// ValidateTime validates time format (HH:MM:SS)
func ValidateTime(timeStr string) bool {
	timeRegex := regexp.MatchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`, timeStr)
	return timeRegex
} 