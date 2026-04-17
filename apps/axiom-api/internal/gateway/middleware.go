package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

type Middleware struct {
	jwtIssuer *identity.JWTIssuer
}

func NewMiddleware(jwtIssuer *identity.JWTIssuer) *Middleware {
	return &Middleware{jwtIssuer: jwtIssuer}
}

// GetClaims is retained for existing callers; prefer identity.ClaimsFromContext.
func GetClaims(ctx context.Context) *identity.Claims {
	return identity.ClaimsFromContext(ctx)
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			platform.WriteError(w, platform.ErrUnauthorized("missing authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			platform.WriteError(w, platform.ErrUnauthorized("invalid authorization header format"))
			return
		}

		claims, err := m.jwtIssuer.Verify(parts[1])
		if err != nil {
			platform.WriteError(w, platform.ErrUnauthorized("invalid or expired token"))
			return
		}

		ctx := identity.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) WithRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := identity.ClaimsFromContext(r.Context())
			if claims == nil {
				platform.WriteError(w, platform.ErrUnauthorized("no claims in context"))
				return
			}
			if !allowed[claims.Role] {
				platform.WriteError(w, platform.ErrForbidden("insufficient permissions"))
				return
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
