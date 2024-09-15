package tokens

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	accessSecret  = []byte("access_token_secret")
	refreshSecret = []byte("refresh_token_secret")
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new JWT access token.
func GenerateAccessToken(userID string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString(accessSecret)
}

// GenerateRefreshToken generates a new JWT refresh token.
func GenerateRefreshToken(userID string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(tokenStr string, isRefreshToken bool) (*Claims, error) {
	var secret []byte
	if isRefreshToken {
		secret = refreshSecret
	} else {
		secret = accessSecret
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshAccessToken takes a refresh token and generates a new access token.
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken, true)
	if err != nil {
		return "", err
	}

	return GenerateAccessToken(claims.UserID, time.Hour)
}
