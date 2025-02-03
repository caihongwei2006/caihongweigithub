package database

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRegisterHandler(c *gin.Context, db *gorm.DB) {
	if c.Request.Method != "POST" {
		c.JSON(405, gin.H{
			"error": "invalid method",
		})
		return
	}

	var requestData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	user := User{
		Username: requestData.Username,
		Password: requestData.Password,
	}

	cookie, err := GenCookie(user, db, c.Writer)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate cookie",
		})
		return
	}

	// SetCookie parameters:
	// name, value, maxAge (in seconds), path, domain, secure, httpOnly
	c.SetCookie("Cookie", cookie, 72*3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": "registration successful",
	})
}

func PostLoginHandler(c *gin.Context, db *gorm.DB) {
	// 检查是否已有有效 cookie
	cookie, err := c.Request.Cookie("Cookie") // 指定 cookie 名
	if err == nil {
		// 如果能获取到 session_cookie，则根据你的业务逻辑做校验（可选）
		if err := CheckCookie(cookie.Value, db); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful (via cookie)",
			})
			return
		}
		// 如果 cookie 无效，则继续执行后续的用户名密码校验
	}

	// 从 JSON 中获取用户名密码
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 在数据库中查找用户
	var user User
	if err := db.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 简单的明文密码比较（不安全，仅作演示）
	if user.Password != loginReq.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 如果用户名和密码正确，就生成新的 cookie
	newCookie, err := GenCookie(user, db, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate cookie"})
		return
	}

	// 设置 session_cookie，以便后续请求中可使用该 cookie
	c.SetCookie("Cookie", newCookie, 72*3600, "/", "localhost", false, true)

	// 返回登录成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"cookie":  newCookie, // 可选：把生成的 cookie 返回给前端做调试
	})
}

func PostQuerryHandler(w *gin.Context) {

}
