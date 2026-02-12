package handler_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	customValidator "todo_app/pkg/validator"
)

// TestMain registers custom validators (alphanumunder, nospaces, strongpassword)
// with gin's binding engine before any handler tests run.
// Without this, tests that bind DTOs using these tags will panic.
func TestMain(m *testing.M) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := customValidator.RegisterCustomValidators(v); err != nil {
			panic("failed to register custom validators: " + err.Error())
		}
	}

	os.Exit(m.Run())
}
