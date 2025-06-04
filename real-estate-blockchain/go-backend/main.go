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
	// DB 초기화
	db := database.InitDB()

	// ✅ 커맨드라인 사용자 등록: go run main.go register TestUser9
	if len(os.Args) == 3 && os.Args[1] == "register" {
		username := os.Args[2]
		err := handler.RegisterUserCLI(username)
		if err != nil {
			log.Fatalf("❌ 사용자 등록 실패: %v\n", err)
		}
		fmt.Printf("✅ 사용자 등록 성공: %s\n", username)
		return
	}

	// ✅ Gin 서버 실행
	router := gin.Default()

	// ✅ CORS 설정
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")         // 개발 중에는 모든 출처를 허용합니다. 실제 배포 시에는 특정 도메인으로 제한하는 것이 좋습니다.
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 만약 쿠키나 인증 헤더를 사용한다면 필요할 수 있습니다.
		// 허용할 헤더 목록을 더 명시적으로 지정합니다. 프론트엔드에서 보내는 Content-Type 및 ngrok 관련 헤더를 포함합니다.
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning")
		// 허용할 HTTP 메소드 목록입니다. OPTIONS가 포함되어야 preflight 요청을 처리할 수 있습니다.
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 브라우저가 보내는 Preflight 요청 (OPTIONS 메소드)에 대한 처리입니다.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 204 No Content 응답으로 Preflight 요청 성공을 알립니다.
			return
		}
		c.Next()
	})

	// ✅ 정적 파일 경로 설정 (사진 접근용)
	router.Static("/uploads", "./uploads")

	// ✅ 공통 API 라우팅
	router.POST("/upload-photo", handler.UploadPhoto)
	router.POST("/add-property", handler.AddProperty)
	router.GET("/property/:id", handler.GetProperty)
	router.GET("/properties", handler.GetAllProperties)
	router.POST("/update-property", handler.UpdateProperty)
	router.POST("/reserve-property", handler.ReserveProperty)
	router.POST("/signup", handler.Signup)
	router.POST("/login", handler.Login)
	router.GET("/my-properties", handler.GetMyProperties)
	//	router.POST("/auth/login", handler.Login)

	// ✅ DID 기반 공인중개사 회원가입
	brokerApiGroup := router.Group("/api/brokers")
	{
		brokerApiGroup.POST("/register-with-did", handler.SignUpBrokerAndIssueDID)
	}

	// ✅ 역할 기반 사용자 등록 API 추가
	router.POST("/register-user", handler.RegisterUser(db))

	// ✅ 서버 실행
	log.Println("🚀 서버 실행 중: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
