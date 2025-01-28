package main

import (
    "fmt"
    "log"

    "github.com/caihongwei2006/AI_API/database"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化数据库连接
    db, err := database.InitDB()
    if err != nil {
        log.Fatalf("数据库初始化失败: %v", err)
    }

    // 创建一个默认的 Gin 路由器
    router := gin.Default()

    // 注册路由和处理函数
    // 注册用户
    router.POST("/register", func(c *gin.Context) {
        database.PostRegisterHandler(c, db)
    })

    // 登录用户
    router.POST("/login", func(c *gin.Context) {
        database.PostLoginHandler(c, db)
    })

    // 示例路由：根路径
    router.GET("/", func(c *gin.Context) {
        c.String(200, "Hello, World!")
    })

    // 启动服务器
    address := ":8080"
    fmt.Printf("服务器正在运行在 %s\n", address)
    if err := router.Run(address); err != nil {
        log.Fatalf("服务器启动失败: %v", err)
    }
}