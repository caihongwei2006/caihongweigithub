package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a JWT token for the user.
func GenerateToken(c User) (string, error) {
	// Define expiration time.
	expirationTime := time.Now().Add(24 * time.Hour) //24小时有效期
	claims := &jwt.MapClaims{
		"username": c.Username,
		"exp":      expirationTime.Unix(),
		"ip":       c.IPv4,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token with secret key.
	return token.SignedString([]byte("chw"))
}

// ParseToken parses the JWT token.
func ParseToken(tokenString string) (*jwt.Token, error) {
	// Validate token using the secret key.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte("chw"), nil //用密钥”chw“来验证token合法性
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}
