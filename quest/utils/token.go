package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// 将密码哈希话处理一下
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 解码对比密码
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateToken(username string) (string, error) { //生成token
	// Define expiration time.
	expirationTime := time.Now().Add(24 * time.Hour) //24小数有效期
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token with secret key.
	return token.SignedString([]byte("chw")) //签名：chw
}

// ParseToken parses and validates the JWT token.
func ParseToken(tokenString string) (*jwt.Token, error) {
	// Validate token using the secret key.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Your token is invalid..........")
		}
		return []byte("chw"), nil //用密钥验证。
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}
