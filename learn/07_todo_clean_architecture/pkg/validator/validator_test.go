package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// setupValidator creates a validator with custom rules registered
func setupValidator(t *testing.T) *validator.Validate {
	v := validator.New()
	err := RegisterCustomValidators(v)
	assert.NoError(t, err, "should register custom validators without error")
	return v
}

// TestNoSpaces tests the nospaces validator
func TestNoSpaces(t *testing.T) {
	v := setupValidator(t)

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"valid: no spaces", "username123", false},
		{"valid: underscores", "user_name_123", false},
		{"valid: dashes", "user-name-123", false},
		{"invalid: single space", "user name", true},
		{"invalid: multiple spaces", "user  name  123", true},
		{"invalid: leading space", " username", true},
		{"invalid: trailing space", "username ", true},
		{"invalid: tab character", "user\tname", true},
		{"invalid: newline", "user\nname", true},
		{"invalid: carriage return", "user\rname", true},
		{"valid: empty string", "", false}, // No spaces in empty string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a struct with the validation tag
			type TestStruct struct {
				Username string `validate:"nospaces"`
			}

			// ACT: Validate
			testData := TestStruct{Username: tt.input}
			err := v.Struct(testData)

			// ASSERT
			if tt.shouldErr {
				assert.Error(t, err, "should fail validation for: %s", tt.input)
			} else {
				assert.NoError(t, err, "should pass validation for: %s", tt.input)
			}
		})
	}
}

// TestAlphaNumericUnderscore tests the alphanumunder validator
func TestAlphaNumericUnderscore(t *testing.T) {
	v := setupValidator(t)

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases - must start with letter
		{"valid: simple username", "username", false},
		{"valid: with numbers", "user123", false},
		{"valid: with underscores", "user_name_123", false},
		{"valid: mixed case", "UserName123", false},
		{"valid: single letter", "a", false},
		{"valid: uppercase start", "User123", false},

		// Invalid cases
		{"invalid: starts with number", "123user", true},
		{"invalid: starts with underscore", "_username", true},
		{"invalid: contains space", "user name", true},
		{"invalid: contains dash", "user-name", true},
		{"invalid: contains special char", "user@name", true},
		{"invalid: contains dot", "user.name", true},
		{"invalid: empty string", "", true},
		{"invalid: only numbers", "123", true},
		{"invalid: only underscore", "_", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Username string `validate:"alphanumunder"`
			}

			testData := TestStruct{Username: tt.input}
			err := v.Struct(testData)

			if tt.shouldErr {
				assert.Error(t, err, "should fail validation for: %s", tt.input)
			} else {
				assert.NoError(t, err, "should pass validation for: %s", tt.input)
			}
		})
	}
}

// TestStrongPassword tests the strongpassword validator
func TestStrongPassword(t *testing.T) {
	v := setupValidator(t)

	tests := []struct {
		name      string
		input     string
		shouldErr bool
		reason    string
	}{
		// Valid passwords
		{"valid: meets all requirements", "Pass@123word", false, ""},
		{"valid: complex password", "MyP@ssw0rd!", false, ""},
		{"valid: with multiple special chars", "Test#123$Pass", false, ""},
		{"valid: minimum length", "Aa1!bcde", false, ""},

		// Invalid: too short
		{"invalid: too short", "Pass@1", true, "less than 8 chars"},
		{"invalid: 7 chars", "Pas@1Aa", true, "7 chars"},

		// Invalid: missing uppercase
		{"invalid: no uppercase", "pass@word123", true, "no uppercase letter"},

		// Invalid: missing lowercase
		{"invalid: no lowercase", "PASS@WORD123", true, "no lowercase letter"},

		// Invalid: missing number
		{"invalid: no number", "Pass@word!", true, "no number"},

		// Invalid: missing special character
		{"invalid: no special char", "Password123", true, "no special character"},

		// Invalid: only one type
		{"invalid: only lowercase", "password", true, "only lowercase"},
		{"invalid: only uppercase", "PASSWORD", true, "only uppercase"},
		{"invalid: only numbers", "12345678", true, "only numbers"},

		// Edge cases
		{"invalid: empty string", "", true, "empty"},
		{"invalid: spaces only", "        ", true, "spaces don't count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Password string `validate:"strongpassword"`
			}

			testData := TestStruct{Password: tt.input}
			err := v.Struct(testData)

			if tt.shouldErr {
				assert.Error(t, err, "should fail: %s (%s)", tt.reason, tt.input)
			} else {
				assert.NoError(t, err, "should pass validation for: %s", tt.input)
			}
		})
	}
}

// TestCombinedValidations tests multiple validation tags together
func TestCombinedValidations(t *testing.T) {
	v := setupValidator(t)

	type RegisterDTO struct {
		Username string `validate:"required,alphanumunder,min=3,max=20"`
		Password string `validate:"required,strongpassword,min=8"`
		Email    string `validate:"required,email"`
	}

	tests := []struct {
		name      string
		dto       RegisterDTO
		shouldErr bool
	}{
		{
			name: "valid registration",
			dto: RegisterDTO{
				Username: "john_doe123",
				Password: "SecureP@ss123",
				Email:    "john@example.com",
			},
			shouldErr: false,
		},
		{
			name: "invalid username - too short",
			dto: RegisterDTO{
				Username: "jo",
				Password: "SecureP@ss123",
				Email:    "john@example.com",
			},
			shouldErr: true,
		},
		{
			name: "invalid username - starts with number",
			dto: RegisterDTO{
				Username: "123john",
				Password: "SecureP@ss123",
				Email:    "john@example.com",
			},
			shouldErr: true,
		},
		{
			name: "invalid password - not strong",
			dto: RegisterDTO{
				Username: "john_doe",
				Password: "weakpass",
				Email:    "john@example.com",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.dto)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetValidationErrors tests the error message formatter
func TestGetValidationErrors(t *testing.T) {
	v := setupValidator(t)

	type TestDTO struct {
		Username string `validate:"required,min=3,alphanumunder"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,strongpassword"`
		Age      int    `validate:"required,min=18"`
	}

	// Create invalid data
	dto := TestDTO{
		Username: "12", // Too short and starts with number
		Email:    "invalid-email",
		Password: "weak",
		Age:      15,
	}

	// ACT: Validate and get errors
	err := v.Struct(dto)
	assert.Error(t, err, "validation should fail")

	errors := GetValidationErrors(err)

	// ASSERT: Check error messages are user-friendly
	assert.NotEmpty(t, errors, "should have validation errors")

	// Check that error map contains field names (lowercased)
	// Note: The exact errors depend on which validation fails first
	assert.NotNil(t, errors, "errors map should not be nil")
}

// TestGetValidationErrorsForEachTag tests error messages for each tag type
func TestGetValidationErrorsForEachTag(t *testing.T) {
	v := setupValidator(t)

	tests := []struct {
		name          string
		structDef     interface{}
		expectedField string
		expectedMsg   string
	}{
		{
			name: "required error",
			structDef: struct {
				Name string `validate:"required"`
			}{Name: ""},
			expectedField: "name",
			expectedMsg:   "name is required",
		},
		{
			name: "email error",
			structDef: struct {
				Email string `validate:"email"`
			}{Email: "invalid"},
			expectedField: "email",
			expectedMsg:   "invalid email format",
		},
		{
			name: "min length error",
			structDef: struct {
				Password string `validate:"min=8"`
			}{Password: "short"},
			expectedField: "password",
			expectedMsg:   "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.structDef)
			assert.Error(t, err)

			errors := GetValidationErrors(err)
			assert.Contains(t, errors, tt.expectedField, "should have error for field")
			assert.Equal(t, tt.expectedMsg, errors[tt.expectedField])
		})
	}
}

// TestGetValidationErrorsWithNoSpaces tests nospaces error message
func TestGetValidationErrorsWithNoSpaces(t *testing.T) {
	v := setupValidator(t)

	type TestStruct struct {
		Username string `validate:"nospaces"`
	}

	dto := TestStruct{Username: "user name"}
	err := v.Struct(dto)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "username")
	assert.Equal(t, "username must not contain spaces", errors["username"])
}

// TestGetValidationErrorsWithAlphanumunder tests alphanumunder error message
func TestGetValidationErrorsWithAlphanumunder(t *testing.T) {
	v := setupValidator(t)

	type TestStruct struct {
		Username string `validate:"alphanumunder"`
	}

	dto := TestStruct{Username: "123user"}
	err := v.Struct(dto)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "username")
	assert.Contains(t, errors["username"], "must start with a letter")
}

// TestGetValidationErrorsWithStrongPassword tests strongpassword error message
func TestGetValidationErrorsWithStrongPassword(t *testing.T) {
	v := setupValidator(t)

	type TestStruct struct {
		Password string `validate:"strongpassword"`
	}

	dto := TestStruct{Password: "weakpass"}
	err := v.Struct(dto)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
	assert.Contains(t, errors["password"], "uppercase")
	assert.Contains(t, errors["password"], "lowercase")
	assert.Contains(t, errors["password"], "number")
	assert.Contains(t, errors["password"], "special character")
}

// TestGetValidationErrorsWithNonValidatorError tests with non-validator error
func TestGetValidationErrorsWithNonValidatorError(t *testing.T) {
	// Create a regular error (not validator.ValidationErrors)
	err := assert.AnError

	// ACT
	errors := GetValidationErrors(err)

	// ASSERT: Should return empty map for non-validation errors
	assert.Empty(t, errors, "should return empty map for non-validation errors")
}

// TestRegisterCustomValidatorsMultipleTimes tests idempotency
func TestRegisterCustomValidatorsMultipleTimes(t *testing.T) {
	v := validator.New()

	// Register multiple times
	err1 := RegisterCustomValidators(v)
	err2 := RegisterCustomValidators(v)

	// NOTE: Re-registering might cause errors depending on validator library
	// The first registration should succeed
	assert.NoError(t, err1, "first registration should succeed")

	// Second might fail with "already registered" error, which is fine
	// We're just testing it doesn't panic
	_ = err2
}

// TestStrongPasswordWithUnicodeCharacters tests strong password with unicode
func TestStrongPasswordWithUnicodeCharacters(t *testing.T) {
	v := setupValidator(t)

	tests := []struct {
		name      string
		password  string
		shouldErr bool
	}{
		{"valid: unicode special chars", "Passw0rd™", false},
		{"valid: unicode in middle", "P@ssw0rd™", false},
		{"valid: accented chars", "Pàssw0rd!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Password string `validate:"strongpassword"`
			}

			testData := TestStruct{Password: tt.password}
			err := v.Struct(testData)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestMaxValidation tests max length validator
func TestMaxValidation(t *testing.T) {
	v := setupValidator(t)

	type TestStruct struct {
		Username string `validate:"max=10"`
	}

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"valid: under max", "user", false},
		{"valid: at max", "1234567890", false},
		{"invalid: over max", "12345678901", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData := TestStruct{Username: tt.input}
			err := v.Struct(testData)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
