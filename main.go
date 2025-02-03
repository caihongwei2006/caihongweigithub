package main

import (
	"fmt"
	"log"

	"github.com/caihongwei2006/AI_API/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Initialize GPT client (using environment variable for API key)
	if err := database.InitializeGPT("sk-zk23491c78f06488209ddbbbbb14622a69e673dfb23936e2", "https://api.zhizengzeng.com/v1"); err != nil {
		log.Fatalf("Failed to initialize GPT: %v", err)
	}

	// 3. Create Gin router
	r := gin.Default()

	// 4. Setup routes with database injection
	api := r.Group("/api")
	{
		// GPT endpoint
		api.POST("/ask", func(c *gin.Context) {
			database.PostGPTHandler(c.Writer, c.Request, db)
		})

		// Auth endpoints
		api.POST("/register", func(c *gin.Context) {
			database.PostRegisterHandler(c, db)
		})

		api.POST("/login", func(c *gin.Context) {
			database.PostLoginHandler(c, db)
		})
	}

	// 5. Start server
	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
