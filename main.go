// real-estate-blockchain-feature-did-vc/main.go
package main

import (
	"log"
	"realestate/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() // Logger와 Recovery 미들웨어가 기본으로 포함됨

	// ─── CORS 설정 (gin-contrib/cors 사용) ───────────────────────────
	// 이 미들웨어는 다른 라우트 핸들러들보다 먼저 등록되어야 합니다.
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // 프론트엔드 개발 서버 주소
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "ngrok-skip-browser-warning"}, // 일반적으로 필요한 헤더들 추가
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 자격 증명(쿠키, 인증 헤더 등) 허용 여부
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	// ─────────────────────────────────────────────────────────────────

	// API 라우팅 설정
	// router.POST("/register-user", handler.RegisterUser) // 일반 사용자 Fabric 등록 (필요시 주석 해제)
	// router.POST("/add-property", handler.AddProperty)   // 매물 등록 (DID Auth 연동 전)
	// router.GET("/property/:id", handler.GetProperty)    // 매물 조회 (필요시 주석 해제)

	// 공인중개사 관련 API 그룹
	brokerApiGroup := router.Group("/api/brokers")
	{
		brokerApiGroup.POST("/register-with-did", handler.SignUpBrokerAndIssueDID)  // 공인중개사 회원가입 (DID 발급 포함). 프론트엔드 SignupAgent.js의 DID_SIGNUP_API_ENDPOINT와 일치해야 함

		// VC 검증 엔드포인트 (기존 유지)
		brokerApiGroup.POST("/verify", handler.VerifyBroker)
	}

	// (선택적) 일반 사용자 회원가입 API (만약 SignupUser.js에서 이 백엔드를 호출하도록 수정한다면)
	// userApiGroup := router.Group("/api/users")
	// {
	// 	userApiGroup.POST("/signup", handler.RegisterPlatformUser) // RegisterPlatformUser 핸들러 필요
	// }

	log.Println("🚀 서버 실행 중: http://localhost:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
