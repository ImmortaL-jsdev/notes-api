package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, duration time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("failed to sign token : %w", err)
	}

	return signedToken, nil

}

func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token :%w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
