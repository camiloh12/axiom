package platform_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON_SetsStatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteJSON(rr, http.StatusCreated, map[string]string{"ok": "yes"})

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "yes", body["ok"])
}

func TestWriteError_AppErrorRendersCode(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteError(rr, platform.ErrNotFound("firm"))

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "firm", body["error"])
}

func TestWriteError_GenericErrorRenders500(t *testing.T) {
	rr := httptest.NewRecorder()
	platform.WriteError(rr, errors.New("boom"))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
