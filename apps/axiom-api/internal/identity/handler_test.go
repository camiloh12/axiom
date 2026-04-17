package identity_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHandler(t *testing.T) http.Handler {
	t.Helper()
	pool := platform.TestDB(t)
	privKey, pubKey := testKeyPair(t)
	issuer := identity.NewJWTIssuer(privKey, pubKey)
	svc := identity.NewService(pool, issuer)
	handler := identity.NewHandler(svc, issuer)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)
	return r
}

func TestHandler_RegisterAndLogin(t *testing.T) {
	router := setupHandler(t)

	regBody, _ := json.Marshal(map[string]any{
		"firm_name":           "Handler Test Firm",
		"admin_email":         "handler@test.com",
		"admin_name":          "Handler Admin",
		"password":            "testpass123",
		"country":             "US",
		"staff_count_range":   "1-10",
		"primary_audit_types": []string{"SOC2"},
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var regResp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &regResp))
	assert.NotEmpty(t, regResp["access_token"])

	loginBody, _ := json.Marshal(map[string]string{
		"email":    "handler@test.com",
		"password": "testpass123",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var loginResp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &loginResp))
	assert.NotEmpty(t, loginResp["access_token"])
}
