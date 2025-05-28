package main

import (
	"fmt"
	"log"
	"os"
	"realestate/database"
	"realestate/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	// 커맨드라인 사용자 등록: 예) go run main.go register TestUser9
	if len(os.Args) == 3 && os.Args[1] == "register" {
		username := os.Args[2]
		err := handler.RegisterUserCLI(username)
		if err != nil {
			log.Fatalf("❌ 사용자 등록 실패: %v\n", err)
		}
		fmt.Printf("✅ 사용자 등록 성공: %s\n", username)
		return
	}

	// Gin 서버 실행
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

	// ✅ 정적 파일 경로 설정 (사진 접근용)
	router.Static("/uploads", "./uploads")

	// ✅ 사진 업로드 API
	router.POST("/upload-photo", handler.UploadPhoto)

	// 기존 라우팅
	router.POST("/register-with-did", handler.SignUpBrokerAndIssueDID)
	router.POST("/add-property", handler.AddProperty)
	router.GET("/property/:id", handler.GetProperty)
	router.GET("/properties", handler.GetAllProperties)
	router.POST("/update-property", handler.UpdateProperty)
	router.POST("/reserve-property", handler.ReserveProperty)
	router.POST("/signup", handler.Signup)
	router.POST("/login", handler.Login)
	router.GET("/my-properties", handler.GetMyProperties)
	router.POST("/auth/login", handler.Login)

	// --- 현 프로젝트의 DID 발급 관련 라우트 추가 ---
	// SignupAgent.js가 호출하는 경로와 일치해야 합니다.
	brokerApiGroup := router.Group("/api/brokers") // 사용자님의 기능은 /api/brokers 그룹 하위에 있었음
	{
		// 공인중개사 회원가입 (DID 발급 포함)
		brokerApiGroup.POST("/register-with-did", handler.SignUpBrokerAndIssueDID) //
	}

	log.Println("🚀 서버 실행 중: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
