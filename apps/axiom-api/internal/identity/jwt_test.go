package identity_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key, &key.PublicKey
}

func TestJWTIssueAndVerify(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)

	userID := uuid.New()
	firmID := uuid.New()

	pair, err := issuer.Issue(userID, firmID, "FirmAdmin")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)

	claims, err := issuer.Verify(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, firmID, claims.FirmID)
	assert.Equal(t, "FirmAdmin", claims.Role)
}

func TestJWTExpiredToken(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuerWithDurations(privKey, pubKey, -1*time.Second, 7*24*time.Hour)

	pair, err := issuer.Issue(uuid.New(), uuid.New(), "Staff")
	require.NoError(t, err)

	_, err = issuer.Verify(pair.AccessToken)
	assert.Error(t, err)
}

func TestJWTRefresh(t *testing.T) {
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)

	userID := uuid.New()
	firmID := uuid.New()

	pair, err := issuer.Issue(userID, firmID, "Manager")
	require.NoError(t, err)

	// Refresh tokens issued in the same second would yield identical access tokens
	// since exp/iat are both at one-second resolution. Pause briefly so the new
	// token has a later IssuedAt timestamp.
	time.Sleep(1100 * time.Millisecond)

	newPair, err := issuer.Refresh(pair.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newPair.AccessToken)
	assert.NotEqual(t, pair.AccessToken, newPair.AccessToken)
}
