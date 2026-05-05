// Package featureflag provides typed access to runtime feature flags.
//
// Flags are read once from environment variables at startup and exposed as
// typed accessors. The launch posture (per docs/superpowers/specs/implementation-plan-design.md
// §2.1) gates the auditee-side surface behind CLIENT_HUB_ENABLED, defaulting
// to false so client-facing endpoints, invitations, and UI are disabled until
// the flag is flipped on.
package featureflag

import (
	"os"
	"strconv"
	"sync"
)

type Flags struct {
	clientHubEnabled bool
}

var (
	once   sync.Once
	loaded *Flags
)

// Load reads flag values from the environment. Subsequent calls are no-ops —
// flags are intentionally immutable for the lifetime of the process. Tests
// that need to flip a flag should construct a Flags value directly with New().
func Load() *Flags {
	once.Do(func() {
		loaded = &Flags{
			clientHubEnabled: parseBool(os.Getenv("CLIENT_HUB_ENABLED"), false),
		}
	})
	return loaded
}

// New constructs an explicit Flags value for tests.
func New(clientHubEnabled bool) *Flags {
	return &Flags{clientHubEnabled: clientHubEnabled}
}

// ClientHubEnabled reports whether the auditee-side Client Hub surface is
// active. When false (default), client-facing endpoints return 410 Gone,
// ClientAdmin/ClientUser invitations are rejected, and the SPA hides
// client-side navigation.
func (f *Flags) ClientHubEnabled() bool {
	if f == nil {
		return false
	}
	return f.clientHubEnabled
}

// ClientHubEnabled is a convenience accessor that reads the loaded flag set.
// Callers that already hold a *Flags should prefer the method form for
// testability.
func ClientHubEnabled() bool {
	return Load().ClientHubEnabled()
}

func parseBool(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return v
}
