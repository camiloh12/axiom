package platform_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// t.Setenv registers cleanup to restore prior state; os.Unsetenv then
	// clears the var so envconfig applies its `default` tag.
	t.Setenv("PORT", "x")
	t.Setenv("DATABASE_URL", "x")
	t.Setenv("ENVIRONMENT", "x")
	require.NoError(t, os.Unsetenv("PORT"))
	require.NoError(t, os.Unsetenv("DATABASE_URL"))
	require.NoError(t, os.Unsetenv("ENVIRONMENT"))

	cfg, err := platform.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
	assert.Contains(t, cfg.DatabaseURL, "axiom_db")
	assert.Equal(t, "development", cfg.Environment)
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://custom/db")
	t.Setenv("ENVIRONMENT", "production")

	cfg, err := platform.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://custom/db", cfg.DatabaseURL)
	assert.Equal(t, "production", cfg.Environment)
}
