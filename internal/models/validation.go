package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom validation functions
	validate.RegisterValidation("username", validateUsername)
	validate.RegisterValidation("slug", validateSlug)
	validate.RegisterValidation("article_status", validateArticleStatus)
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return validate
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	// Username should contain only alphanumeric characters, underscores, and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

// validateSlug validates slug format
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if len(slug) == 0 {
		return false
	}
	// Slug should contain only lowercase letters, numbers, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return matched
}

// validateArticleStatus validates article status enum
func validateArticleStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return status == string(StatusDraft) || status == string(StatusPublished) || status == string(StatusArchived)
}

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// FormatValidationErrors converts validator errors to custom format
func FormatValidationErrors(err error) ValidationErrors {
	var validationErrors ValidationErrors
	
	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErrors {
			validationError := ValidationError{
				Field: fieldError.Field(),
				Value: fieldError.Value(),
			}
			
			switch fieldError.Tag() {
			case "required":
				validationError.Message = fieldError.Field() + " is required"
			case "email":
				validationError.Message = fieldError.Field() + " must be a valid email address"
			case "min":
				validationError.Message = fieldError.Field() + " must be at least " + fieldError.Param() + " characters long"
			case "max":
				validationError.Message = fieldError.Field() + " must be at most " + fieldError.Param() + " characters long"
			case "username":
				validationError.Message = fieldError.Field() + " must be 3-50 characters long and contain only letters, numbers, underscores, and hyphens"
			case "slug":
				validationError.Message = fieldError.Field() + " must contain only lowercase letters, numbers, and hyphens"
			case "article_status":
				validationError.Message = fieldError.Field() + " must be one of: draft, published, archived"
			default:
				validationError.Message = fieldError.Field() + " is invalid"
			}
			
			validationErrors = append(validationErrors, validationError)
		}
	}
	
	return validationErrors
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return FormatValidationErrors(err)
	}
	return nil
}