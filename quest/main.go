package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"quest/database"
	"quest/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func ShowFrontPage(c *gin.Context) {
	c.HTML(200, "frontpage.html", nil)
}
func PostForm() {
	data := make(url.Values)
	data.Add("name", "chw")
	data.Add("password", "123")
	payload := data.Encode()
	r, _ := http.Post(
		"http://localhost:8080/signup",
		"application/x-www-form-urlencoded",
		strings.NewReader(payload), //本质上是字符串
	)
	defer func() { _ = r.Body.Close() }()
	content, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s", content)
}

func Signup(username, password string) {
	// Create form data
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	// Send POST request to the backend
	resp, err := http.Post("http://localhost:8080/signup", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func main() {
	// Connect to the database
	database.ConnectDB()

	router := gin.Default()
	router.LoadHTMLGlob("pages/*") //一开始由template生产的html文件，加载一下

	// Routes
	router.GET("/login", models.GetLogin)     // 处理Get请求
	router.POST("/login", models.PostLogin)   // 登录Post请求
	router.GET("/signup", models.GetSignup)   // Get signup
	router.POST("/signup", models.PostSignup) // 注册Post请求
	router.GET("/", ShowFrontPage)            // 主页，没什么用

	router.Run(":8080")

}
