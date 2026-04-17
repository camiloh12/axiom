package platform_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func TestAppError_Error(t *testing.T) {
	err := platform.ErrNotFound("user")
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "user")
}

func TestErrorConstructors_SetCorrectStatusCodes(t *testing.T) {
	cases := []struct {
		name     string
		err      *platform.AppError
		expected int
	}{
		{"NotFound", platform.ErrNotFound("x"), http.StatusNotFound},
		{"Unauthorized", platform.ErrUnauthorized("x"), http.StatusUnauthorized},
		{"Forbidden", platform.ErrForbidden("x"), http.StatusForbidden},
		{"BadRequest", platform.ErrBadRequest("x"), http.StatusBadRequest},
		{"Conflict", platform.ErrConflict("x"), http.StatusConflict},
		{"Validation", platform.ErrValidation("x", "d"), http.StatusUnprocessableEntity},
		{"Internal", platform.ErrInternal("x"), http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.err.Code)
		})
	}
}

func TestErrValidation_CarriesDetail(t *testing.T) {
	err := platform.ErrValidation("missing fields", "email is required")
	assert.Equal(t, "missing fields", err.Message)
	assert.Equal(t, "email is required", err.Detail)
}
