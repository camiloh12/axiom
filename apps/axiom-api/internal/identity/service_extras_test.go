package identity_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform/featureflag"
)

func TestGetFirm(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	reg, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName: "Get Firm", AdminEmail: "gf@test.com", AdminName: "A",
		Password: "pwpwpwpw", Country: "US", StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	firm, err := svc.GetFirm(ctx, reg.Firm.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Firm", firm.Name)
}

func TestClientCRUD(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	reg, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName: "Client Firm", AdminEmail: "cf@test.com", AdminName: "A",
		Password: "pwpwpwpw", Country: "US", StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	client, err := svc.CreateClient(ctx, reg.Firm.ID, identity.CreateClientInput{
		Name: "Acme Client", Industry: "Manufacturing",
		PrimaryContactEmail: "ops@acme.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "Acme Client", client.Name)
	assert.Equal(t, reg.Firm.ID, client.FirmID)

	list, err := svc.ListClients(ctx, reg.Firm.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Acme Client", list[0].Name)

	got, err := svc.GetClient(ctx, client.ID)
	require.NoError(t, err)
	assert.Equal(t, "Manufacturing", got.Industry)
}

func TestInvitationFlow(t *testing.T) {
	svc := setupService(t)
	ctx := context.Background()

	reg, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName: "Inv Firm", AdminEmail: "inv@test.com", AdminName: "Admin",
		Password: "pwpwpwpw", Country: "US", StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	inv, err := svc.CreateInvitation(ctx, reg.Firm.ID, reg.User.ID, identity.CreateInvitationInput{
		Email: "staff@test.com", AssignedRole: "Staff",
	})
	require.NoError(t, err)
	require.NotEmpty(t, inv.Token, "token should be returned at creation")

	validated, err := svc.ValidateInvitationToken(ctx, inv.Token)
	require.NoError(t, err)
	assert.Equal(t, "staff@test.com", validated.Email)

	result, err := svc.AcceptInvitation(ctx, inv.Token, "Staffer", "staffpass123")
	require.NoError(t, err)
	assert.Equal(t, "Staff", result.User.Role)
	assert.NotEmpty(t, result.Tokens.AccessToken)

	// Re-accepting the same token must fail.
	_, err = svc.AcceptInvitation(ctx, inv.Token, "Other", "pwpwpwpw")
	assert.Error(t, err)
}

// TestCreateInvitation_ClientHubDisabled verifies the launch-posture flag
// rejects ClientAdmin / ClientUser invitations when CLIENT_HUB_ENABLED=false.
// See docs/superpowers/specs/implementation-plan-design.md §2.1.
func TestCreateInvitation_ClientHubDisabled(t *testing.T) {
	svc := setupServiceWithFlags(t, featureflag.New(false))
	ctx := context.Background()

	reg, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName: "Disabled Hub Firm", AdminEmail: "dh@test.com", AdminName: "Admin",
		Password: "pwpwpwpw", Country: "US", StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	for _, role := range []string{"ClientAdmin", "ClientUser"} {
		t.Run(role, func(t *testing.T) {
			_, err := svc.CreateInvitation(ctx, reg.Firm.ID, reg.User.ID, identity.CreateInvitationInput{
				Email: "client@test.com", AssignedRole: role,
			})
			require.Error(t, err)
			var appErr *platform.AppError
			require.True(t, errors.As(err, &appErr), "expected *platform.AppError, got %T", err)
			assert.Equal(t, 422, appErr.Code)
			assert.Equal(t, "CLIENT_HUB_DISABLED", appErr.Message)
		})
	}

	// Firm-side roles still succeed even when the flag is off.
	inv, err := svc.CreateInvitation(ctx, reg.Firm.ID, reg.User.ID, identity.CreateInvitationInput{
		Email: "staff@test.com", AssignedRole: "Staff",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, inv.Token)
}

// TestCreateInvitation_ClientHubEnabled verifies that with the flag flipped
// on, client-role invitations succeed.
func TestCreateInvitation_ClientHubEnabled(t *testing.T) {
	svc := setupServiceWithFlags(t, featureflag.New(true))
	ctx := context.Background()

	reg, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName: "Enabled Hub Firm", AdminEmail: "eh@test.com", AdminName: "Admin",
		Password: "pwpwpwpw", Country: "US", StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	inv, err := svc.CreateInvitation(ctx, reg.Firm.ID, reg.User.ID, identity.CreateInvitationInput{
		Email: "clientadmin@test.com", AssignedRole: "ClientAdmin",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, inv.Token)
	assert.Equal(t, "ClientAdmin", inv.AssignedRole)
}
