package generatejwt

import (
	"fmt"
	"time"
	"user_service/internal/core/users"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(user *users.User, secretKey string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string, secretKey string) (*Claims, error) {
	claims := &Claims{}

	// parse the token, validate the signature, and unmarshal the claims into the Claims struct.
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // provides the key for validation.
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
