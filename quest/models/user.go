package models

import (
	"net"
	"net/http"
	"quest/database"
	"quest/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 待传入数据库
type User struct {
	Username string `gorm:"primaryKey"` //可以换成JSON格式，目前是encode模式
	Password string
	IPv4     string
	IPv6     string
	Token    string
}

// 处理signinget请求
func GetSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

// 处理post signin
func PostSignup(c *gin.Context) {
	var user User
	user.Username = c.PostForm("username")
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	if ip == "" {
		c.JSON(400, gin.H{"error": "Unable to detect IP address"})
		return
	}
	user.IPv4 = ip // Store detected IP (IPv6 can be handled similarly if needed)

	// Hash password before storing.
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Generate JWT token for the user.
	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	user.Token = token

	// Save user to the database.
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful", "token": token})
}

// 处理login get请求
func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// 处理 login post请求
func PostLogin(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	//将绑定的json数据传入loginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//检查token是否为空
	if loginData.Token != "" {
		// Validate the token
		_, err := ParseToken(loginData.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Token is valid, no need to check password.
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}

	//如果没有token，检查用户名是否存在
	var user User
	if err := database.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// 验证密码
	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
