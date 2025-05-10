package main

import (
	"log"
	"realestate/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS 설정
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 라우팅 설정
	router.POST("/register-user", handler.RegisterUser)
	router.POST("/add-property", handler.AddProperty)
	router.GET("/property/:id", handler.GetProperty) // ✅ 매물 조회 API 추가

	log.Println("🚀 서버 실행 중: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
