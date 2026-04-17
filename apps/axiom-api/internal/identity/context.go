package identity

import "context"

type contextKey string

const claimsKey contextKey = "claims"

// ContextWithClaims returns a new context carrying the given Claims.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves the Claims stored by the auth middleware, or nil
// if none are present.
func ClaimsFromContext(ctx context.Context) *Claims {
	claims, _ := ctx.Value(claimsKey).(*Claims)
	return claims
}
