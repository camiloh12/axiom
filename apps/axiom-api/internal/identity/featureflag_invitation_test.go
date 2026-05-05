package identity_test

// This file holds in-memory tests for the launch-posture feature flag in the
// identity service that don't require a live Postgres connection. The flag
// check executes before any DB call, so we can exercise it with a nil pool.

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform/featureflag"
)

func TestCreateInvitation_ClientHubDisabled_NoDB(t *testing.T) {
	// Build a Service with no pool and the flag off. The flag check should
	// short-circuit before any DB access happens.
	svc := identity.NewServiceWithFlags(nil, nil, featureflag.New(false))

	for _, role := range []string{"ClientAdmin", "ClientUser"} {
		t.Run(role, func(t *testing.T) {
			_, err := svc.CreateInvitation(context.Background(), uuid.New(), uuid.New(), identity.CreateInvitationInput{
				Email:        "client@test.com",
				AssignedRole: role,
			})
			require.Error(t, err)
			var appErr *platform.AppError
			require.True(t, errors.As(err, &appErr), "expected *platform.AppError, got %T: %v", err, err)
			assert.Equal(t, 422, appErr.Code)
			assert.Equal(t, "CLIENT_HUB_DISABLED", appErr.Message)
		})
	}
}
