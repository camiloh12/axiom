package identity

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	FirmID uuid.UUID `json:"firm_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	privKey         *rsa.PrivateKey
	pubKey          *rsa.PublicKey
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTIssuer(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) *JWTIssuer {
	return &JWTIssuer{
		privKey:         privKey,
		pubKey:          pubKey,
		accessDuration:  15 * time.Minute,
		refreshDuration: 7 * 24 * time.Hour,
	}
}

func NewJWTIssuerWithDurations(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, access, refresh time.Duration) *JWTIssuer {
	return &JWTIssuer{
		privKey:         privKey,
		pubKey:          pubKey,
		accessDuration:  access,
		refreshDuration: refresh,
	}
}

func (j *JWTIssuer) Issue(userID, firmID uuid.UUID, role string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := &Claims{
		UserID: userID,
		FirmID: firmID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			Issuer:    "axiom",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessStr, err := accessToken.SignedString(j.privKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID: userID,
		FirmID: firmID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			Issuer:    "axiom",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(j.privKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}

func (j *JWTIssuer) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func (j *JWTIssuer) Refresh(refreshTokenStr string) (*TokenPair, error) {
	claims, err := j.Verify(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	return j.Issue(claims.UserID, claims.FirmID, claims.Role)
}
