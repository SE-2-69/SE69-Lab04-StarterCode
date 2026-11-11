package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer   = "se69-lab04"
	Audience = "equipment-api"
)

// payload of the token: role + standard claims (iss, aud, sub, iat, exp)
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// --- create token (signed with HS256) ---

func CreateToken(username, role, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Audience:  jwt.ClaimStrings{Audience},
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// --- verify token: algorithm, signature, issuer, audience, expiry ---

func VerifyToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{"HS256"}), // reject "none" and other algorithms
		jwt.WithIssuer(Issuer),
		jwt.WithAudience(Audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
