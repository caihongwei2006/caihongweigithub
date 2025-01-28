package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User represents the user model for authentication
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Username         string    `json:"username" gorm:"unique;not null"`
	Email            string    `json:"email" gorm:"unique;not null"`
	Password         string    `json:"-" gorm:"not null"` // "-" means this field won't be shown in JSON
	Role             string    `json:"role" gorm:"default:'user'"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Cookie           string    `json:"cookie"`
	CookieExpiration string    `json:"cookie_expiration"`
}

func GenCookie(user User, db *gorm.DB) (error, string) {
	user.CookieExpiration = time.Now().Add(72 * time.Hour).Format(time.RFC3339) // Set the expiration to 72 hours from now
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	token := user.Username + user.Password + user.CreatedAt.String()
	user.Cookie = token
	//cookie 写入数据库
	if err := db.Create(&user).Error; err != nil {
		return err, ""
	}
	return nil, token
}

func CheckCookie(cookie string, db *gorm.DB) error {
	user := User{}
	if err := db.Where("cookie = ?", cookie).First(&user).Error; err != nil {
		return errors.New("invalid cookie")
	}
	expirationTime, err := time.Parse(time.RFC3339, user.CookieExpiration) // Use the format you stored the expiration in
	if err != nil {
		return errors.New("invalid expiration format")
	}

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		return errors.New("cookie expired")
	}
	return nil
}
