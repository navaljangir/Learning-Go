package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTUtil provides JWT token generation and validation
type JWTUtil struct {
	secret      string
	expiryHours int
	issuer      string
}

// NewJWTUtil creates a new JWT utility
func NewJWTUtil(secret string, expiryHours int, issuer string) *JWTUtil {
	return &JWTUtil{
		secret:      secret,
		expiryHours: expiryHours,
		issuer:      issuer,
	}
}

// GenerateToken generates a new JWT token for the given user
func (j *JWTUtil) GenerateToken(userID, username string) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * time.Duration(j.expiryHours))

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); 
		!ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token expired")
		}

		// Validate issuer
		if claims.Issuer != j.issuer {
			return nil, errors.New("invalid token issuer")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
