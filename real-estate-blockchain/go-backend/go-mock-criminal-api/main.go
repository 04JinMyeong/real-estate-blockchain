package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	router.GET("/check", func(c *gin.Context) {
		userID := c.Query("id")
		fmt.Printf("✅ Mock API: Received criminal record check for User: %s\n", userID)

		// 데모 시나리오: 사용자 ID에 "bad"라는 단어가 포함되면 전과가 있는 것으로 간주
		if strings.Contains(userID, "bad") {
			fmt.Println("➡️ Mock API: User has criminal record. Responding with true.")
			c.JSON(http.StatusOK, gin.H{"hasCriminalRecord": true})
		} else {
			fmt.Println("➡️ Mock API: User has no criminal record. Responding with false.")
			c.JSON(http.StatusOK, gin.H{"hasCriminalRecord": false})
		}
	})

	fmt.Println("🚀 Mock Criminal Record API server is running on http://localhost:8082")
	router.Run(":8082")
}
