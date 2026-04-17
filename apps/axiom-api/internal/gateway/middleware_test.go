package gateway_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/gateway"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
)

func setup(t *testing.T) (*identity.JWTIssuer, *gateway.Middleware) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	issuer := identity.NewJWTIssuer(privKey, &privKey.PublicKey)
	mw := gateway.NewMiddleware(issuer)
	return issuer, mw
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	issuer, mw := setup(t)
	userID := uuid.New()
	firmID := uuid.New()
	pair, err := issuer.Issue(userID, firmID, "Staff")
	require.NoError(t, err)

	handler := mw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := gateway.GetClaims(r.Context())
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, firmID, claims.FirmID)
		assert.Equal(t, "Staff", claims.Role)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	_, mw := setup(t)
	handler := mw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestWithRole_Allowed(t *testing.T) {
	issuer, mw := setup(t)
	pair, _ := issuer.Issue(uuid.New(), uuid.New(), "Partner")

	handler := mw.Auth(mw.WithRole("Partner", "FirmAdmin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWithRole_Denied(t *testing.T) {
	issuer, mw := setup(t)
	pair, _ := issuer.Issue(uuid.New(), uuid.New(), "Staff")

	handler := mw.Auth(mw.WithRole("FirmAdmin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}
