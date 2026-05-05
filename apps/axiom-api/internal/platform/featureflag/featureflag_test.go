package featureflag_test

import (
	"testing"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform/featureflag"
)

func TestNew_ClientHubEnabled(t *testing.T) {
	tests := []struct {
		name string
		in   bool
		want bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := featureflag.New(tt.in)
			if got := f.ClientHubEnabled(); got != tt.want {
				t.Errorf("ClientHubEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilFlags_DefaultDisabled(t *testing.T) {
	var f *featureflag.Flags
	if f.ClientHubEnabled() {
		t.Error("nil Flags should report ClientHubEnabled() = false")
	}
}
