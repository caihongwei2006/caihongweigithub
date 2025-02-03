package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// User represents the user model for authentication
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Username         string    `json:"username" gorm:"unique;not null"`
	Password         string    `json:"-" gorm:"not null"` // "-" means this field won't be shown in JSON
	Role             string    `json:"role" gorm:"default:'user'"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Cookie           string    `json:"cookie"`
	CookieExpiration time.Time `json:"cookie_expiration"`
}

func GenCookie(user User, db *gorm.DB) (string, error) {
	// 设置 Cookie 过期时间
	user.CookieExpiration = time.Now().Add(72 * time.Hour)

	// 更新 UpdatedAt 时间
	user.UpdatedAt = time.Now()

	// 生成一个简单的 Token（仅用于本地测试）
	token := "token_" + user.Username + "_" + time.Now().Format("20060102150405")
	user.Cookie = token

	// 使用 Save 方法更新已有的用户记录
	if err := db.Save(&user).Error; err != nil {
		// 记录错误日志以便调试
		log.Printf("Error updating user: %v", err)
		fmt.Printf("Error updating user: %v", err)
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
