package identity_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform/featureflag"
)

// setupService constructs a Service for tests with CLIENT_HUB_ENABLED=true so
// the existing test corpus continues to validate the full both-sided feature
// surface. Tests that exercise the launch-posture flag-off behavior should
// call setupServiceWithFlags directly.
func setupService(t *testing.T) *identity.Service {
	t.Helper()
	return setupServiceWithFlags(t, featureflag.New(true))
}

func setupServiceWithFlags(t *testing.T, flags *featureflag.Flags) *identity.Service {
	t.Helper()
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	return identity.NewServiceWithFlags(pool, issuer, flags)
}

func TestRegisterFirm(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	result, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Acme CPAs",
		AdminEmail: "admin@acme.com",
		AdminName:  "Alice Admin",
		Password:   "securepassword123",
		Country:    "US",
		StaffCount: "21-40",
		AuditTypes: []string{"ISO27001", "SOC2"},
	})
	require.NoError(t, err)
	assert.Equal(t, "Acme CPAs", result.Firm.Name)
	assert.Equal(t, "admin@acme.com", result.User.Email)
	assert.Equal(t, "FirmAdmin", result.User.Role)
	assert.NotEmpty(t, result.Tokens.AccessToken)
}

func TestRegisterFirm_DuplicateEmail(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	input := identity.RegisterFirmInput{
		FirmName:   "Firm A",
		AdminEmail: "dup@example.com",
		AdminName:  "User A",
		Password:   "password123",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	}

	_, err := svc.RegisterFirm(ctx, input)
	require.NoError(t, err)

	_, err = svc.RegisterFirm(ctx, input)
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	_, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Test Firm",
		AdminEmail: "login@test.com",
		AdminName:  "Login User",
		Password:   "mypassword",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	result, err := svc.Login(ctx, "login@test.com", "mypassword")
	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	_, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Test Firm",
		AdminEmail: "wrong@test.com",
		AdminName:  "User",
		Password:   "correctpassword",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	_, err = svc.Login(ctx, "wrong@test.com", "wrongpassword")
	assert.Error(t, err)
}
