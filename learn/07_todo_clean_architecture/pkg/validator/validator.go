package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Custom validation tags:
// - nospaces: Ensures string has no spaces
// - alphanumunder: Only alphanumeric and underscores
// - strongpassword: At least one uppercase, lowercase, number, special char

// RegisterCustomValidators registers all custom validation rules with the validator
func RegisterCustomValidators(v *validator.Validate) error {
	// Register 'nospaces' validator - ensures no whitespace
	if err := v.RegisterValidation("nospaces", noSpaces); err != nil {
		return err
	}

	// Register 'alphanumunder' validator - only letters, numbers, underscores
	if err := v.RegisterValidation("alphanumunder", alphaNumericUnderscore); err != nil {
		return err
	}

	// Register 'strongpassword' validator - password strength rules
	if err := v.RegisterValidation("strongpassword", strongPassword); err != nil {
		return err
	}

	return nil
}

// noSpaces validates that a string contains no whitespace characters
func noSpaces(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !strings.Contains(value, " ") && !strings.ContainsAny(value, "\t\n\r")
}

// alphaNumericUnderscore validates that a string only contains letters, numbers, and underscores
func alphaNumericUnderscore(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Regex: start with letter, then letters/numbers/underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, value)
	return matched
}

// strongPassword validates password strength:
// - At least 8 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one number
// - At least one special character
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Length is handled by 'min' tag, but double-check here
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// GetValidationErrors converts validator errors to user-friendly messages
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())

			switch fieldError.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "min":
				errors[field] = field + " must be at least " + fieldError.Param() + " characters"
			case "max":
				errors[field] = field + " must not exceed " + fieldError.Param() + " characters"
			case "email":
				errors[field] = "invalid email format"
			case "nospaces":
				errors[field] = field + " must not contain spaces"
			case "alphanumunder":
				errors[field] = field + " must start with a letter and contain only letters, numbers, and underscores"
			case "strongpassword":
				errors[field] = "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}
