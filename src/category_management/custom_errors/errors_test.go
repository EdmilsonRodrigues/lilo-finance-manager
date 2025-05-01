//go:build unit

package customerrors_test

import (
	"net/http"
	"testing"
	"testing/quick"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
)

func TestInternalServerError(t *testing.T) {
	assertion := func(message string) bool {
		err := customerrors.InternalServerError(message)
		if err.Details.Message != "Internal Server Error: "+message {
			t.Errorf("Expected message to be 'Internal Server Error: %s', got '%s'", message, err.Details.Message)
			return false
		}
		if err.Details.Status != http.StatusInternalServerError {
			t.Errorf("Expected status to be 500, got %d", err.Details.Status)
			return false
		}
		return true
	}
	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 1000,
	}); err != nil {
		t.Error(err)
	}
}
