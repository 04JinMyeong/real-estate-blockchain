// go-backend/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"realestate/database"
	"realestate/handler"

	"net/http" // ◀◀◀ 이 라인이 있는지 확인하고, 없다면 추가합니다.

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // .env 파일 로드를 위해 유지
)

func main() {
	// .env 파일 로드 (애플리케이션 시작 시) - 사용자님 코드에서 가져옴
	errEnv := godotenv.Load() // 변수명 변경 (아래 err과의 충돌 방지)
	if errEnv != nil {
		log.Println("환경 변수 파일을 찾을 수 없습니다 (.env). 수동으로 설정된 환경 변수를 사용합니다.")
	}

	// DB 초기화 - 팀원 코드 방식 적용
	_ = database.InitDB() // InitDB()가 *gorm.DB를 반환하고, 팀원 버전의 db.go는 realestate.db를 초기화
	// 만약 InitDB()가 아무것도 반환하지 않는 이전 버전이라면, 다음 라인으로 GetDB()를 호출하여 초기화 유도
	// _ = database.GetDB() // 또는 이 방식 사용

	// 커맨드라인 사용자 등록: 예) go run main.go register TestUser9
	if len(os.Args) == 3 && os.Args[1] == "register" {
		username := os.Args[2]
		// RegisterUserCLI는 이제 DB 인스턴스를 받지 않도록 유지 (내부에서 GetDB() 사용 가정)
		// 또는, 팀원의 RegisterUser(db)와 유사한 형태로 변경되었을 수 있음 - user_handler.go 확인 필요
		err := handler.RegisterUserCLI(username)
		if err != nil {
			log.Fatalf("❌ 사용자 등록 실패: %v\n", err)
		}
		fmt.Printf("✅ 사용자 등록 성공: %s\n", username)
		return
	}

	// Gin 서버 실행
	router := gin.Default()

	// CORS 설정 - 팀원 코드 방식 적용
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning") // 필요한 헤더 추가
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")                         // 허용할 메소드 명시
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}
		c.Next()
	})

	// 정적 파일 경로 설정 (사진 접근용)
	router.Static("/uploads", "./uploads")

	// 사진 업로드 API
	router.POST("/upload-photo", handler.UploadPhoto)

	// 매물 관련 라우팅
	router.POST("/add-property", handler.AddProperty)
	router.GET("/property/:id", handler.GetProperty)
	router.GET("/properties", handler.GetAllProperties)
	router.POST("/update-property", handler.UpdateProperty)
	router.POST("/reserve-property", handler.ReserveProperty)
	router.GET("/my-properties", handler.GetMyProperties) // 내 매물 조회

	// 사용자 인증/인가 관련 라우팅
	router.POST("/signup", handler.Signup) // 일반 사용자 가입
	router.POST("/login", handler.Login)   // 일반 사용자 로그인
	// router.POST("/auth/login", handler.Login) // /login과 중복되므로 하나만 사용 권장

	// DID 기반 공인중개사 관련 라우팅
	brokerApiGroup := router.Group("/api/brokers")
	{
		brokerApiGroup.POST("/register-with-did", handler.SignUpBrokerAndIssueDID) // DID 발급 공인중개사 가입
		// brokerApiGroup.POST("/verify", handler.VerifyBroker) // VerifyBroker 핸들러가 property_handler.go에 있다면 여기에 라우트 필요
	}

	// 역할 기반 사용자 등록 API 추가 (팀원 코드) - handler.RegisterUser 함수가 정의되어 있어야 함
	// 만약 handler.RegisterUser 함수가 아직 없다면, 이 라인은 주석 처리하거나 팀원에게 해당 함수 코드 요청
	// router.POST("/register-user", handler.RegisterUser(db)) // db 변수 전달

	log.Println("🚀 서버 실행 중: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
