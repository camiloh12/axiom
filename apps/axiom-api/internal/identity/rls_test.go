package identity_test

import (
	"context"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRLSIsolation(t *testing.T) {
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	svc := identity.NewService(pool, issuer)
	ctx := context.Background()

	firm1, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Firm Alpha",
		AdminEmail: "alpha@test.com",
		AdminName:  "Alpha Admin",
		Password:   "pass1234",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	firm2, err := svc.RegisterFirm(ctx, identity.RegisterFirmInput{
		FirmName:   "Firm Beta",
		AdminEmail: "beta@test.com",
		AdminName:  "Beta Admin",
		Password:   "pass1234",
		Country:    "US",
		StaffCount: "1-10",
		AuditTypes: []string{"SOC2"},
	})
	require.NoError(t, err)

	_, err = svc.CreateClient(ctx, firm1.Firm.ID, identity.CreateClientInput{
		Name:     "Alpha Client",
		Industry: "Tech",
	})
	require.NoError(t, err)

	alphaClients, err := svc.ListClients(ctx, firm1.Firm.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, alphaClients, 1, "firm alpha should see its own client")

	betaClients, err := svc.ListClients(ctx, firm2.Firm.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, betaClients, 0, "firm beta must not see firm alpha's client")
}
