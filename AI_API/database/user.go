package database

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// User represents the user model for authentication
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Username         string    `json:"username" gorm:"unique;not null"`
	Email            string    `json:"email"`
	Password         string    `json:"-" gorm:"not null"` // "-" means this field won't be shown in JSON
	Role             string    `json:"role" gorm:"default:'user'"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Cookie           string    `json:"cookie"`
	CookieExpiration time.Time `json:"cookie_expiration"`
}

func GenCookie(user User, db *gorm.DB) (string, error) {
	// Set CookieExpiration as time.Time instead of formatted string
	user.CookieExpiration = time.Now().Add(72 * time.Hour)

	// Set CreatedAt and UpdatedAt
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Generate a secure token (consider using UUID or similar)
	token := user.Username + user.Password + user.CreatedAt.String()
	user.Cookie = token

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		// Log the error for debugging
		log.Printf("Error creating user: %v", err)
		return "", err
	}
	return token, nil
}

func CheckCookie(cookie string, db *gorm.DB) error {
	user := User{}
	if err := db.Where("cookie = ?", cookie).First(&user).Error; err != nil {
		return errors.New("invalid cookie")
	}
	expirationTime := user.CookieExpiration

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		return errors.New("cookie expired")
	}
	return nil
}
