package database

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PostRegisterHandler(w *gin.Context, db *gorm.DB) {
	if w.Request.Method != "POST" {
		w.JSON(405, gin.H{
			"error": "invalid method",
		})
		return
	}
	username := w.PostForm("username")
	password := w.PostForm("password")
	user := User{
		Username: username,
		Password: password,
	}
	err, cookie := GenCookie(user, db)
	if err != nil {
		w.JSON(500, gin.H{
			"error": "failed to generate cookie",
		})
	}
	w.SetCookie("Cookie", cookie, 72*3600, "/", "localhost", false, true)

}
func PostLoginHandler(w *gin.Context, db *gorm.DB) {
	if w.Request.Method != "POST" {
		w.JSON(405, gin.H{
			"error": "invalid method",
		})
		return
	}
	cookie := w.Request.Header.Get("Cookie")
	if cookie != "" { //若有cookie，验证
		err := CheckCookie(cookie, db)
		if err != nil {
			w.JSON(405, gin.H{
				"error": "invalid method",
			})
			return
		}
		w.JSON(200, gin.H{
			"message": "Login successful",
		})
	}
	username := w.PostForm("username")
	password := w.PostForm("password")

	user := User{}
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		w.JSON(404, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.JSON(404, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	err, cookie := GenCookie(user, db)
	if err != nil {
		w.JSON(500, gin.H{
			"error": "failed to generate cookie",
		})
		if cookie == "" {
			w.JSON(500, gin.H{
				"error": "failed to generate cookie",
			})
			return
		}

		w.SetCookie("Cookie", cookie, 72*3600, "/", "localhost", false, true) //浏览器保存cookie

		w.JSON(200, gin.H{
			"message": "Login successful",
			"cookie":  cookie,
		})
	}
}

func PostQuerryHandler(w *gin.Context) {

}
