package handler

import (
	"encoding/json" // VC JSON 파싱을 위해 추가
	"fmt"
	"net/http"
	"time"

	"realestate/database"
	"realestate/models"
	"realestate/vc" // VC 검증을 위해 vc 패키지 임포트

	// JWT 토큰 생성을 위해 utils 패키지 임포트

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	//	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key") // 실제 서비스에서는 환경변수나 안전한 저장소 사용

// 회원가입 핸들러
func Signup(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다"})
		return
	}

	db := database.GetDB()

	// 중복 사용자 확인
	var existing models.User
	if err := db.First(&existing, "id = ?", req.ID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "이미 존재하는 사용자입니다"})
		return
	}

	// 비밀번호 해싱
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 해싱 실패"})
		return
	}

	// 사용자 DB 저장
	newUser := models.User{
		ID:        req.ID,
		Password:  string(hashed),
		Email:     req.Email,
		Enrolled:  false,
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 저장 실패"})
		return
	}
	fmt.Println("✅ 사용자 DB 저장 완료:", newUser.ID)

	// Wallet에 사용자 등록
	if err := RegisterUserCLI(req.ID); err != nil {
		fmt.Printf("❗ Wallet 등록 실패: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"message": "✅ 회원가입 완료 (단, wallet 등록 실패)"})
		return
	}

	// wallet 등록 성공 시 DB 상태 업데이트
	newUser.Enrolled = true
	db.Save(&newUser)

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ 회원가입 및 wallet 등록 완료",
	})
}

// 로그인 핸들러 (6/3 vc검증로직 추가함.)
func Login(c *gin.Context) {
	var req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		VC       string `json:"vc" binding:"required"` // VC를 JSON 문자열로 받음

	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다"})
		return
	}

	// 사용자 id,pw검증(기존db로직.)
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, "id = ?", req.ID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "존재하지 않는 사용자입니다"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "비밀번호가 일치하지 않습니다"})
		return
	}

	// 2.VC 검증 (6/3 추가)
	// vc 패키지에 정의된 VerifiableCredential 구조체를 사용합니다.
	var receivedVC vc.VerifiableCredential
	err := json.Unmarshal([]byte(req.VC), &receivedVC)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VC 파싱에 실패했습니다. 유효한 VC JSON을 제공해주세요: " + err.Error()})
		return
	}

	// 3. VC 유효성 검증
	// vc/validate.go의 ValidateVC 함수를 호출하여 VC를 검증합니다.
	// 이 함수는 서명, 발급자, 유효기간, 그리고 'fraudConvictionRecordStatus' 클레임을 검사합니다.
	err = vc.ValidateVC(receivedVC)
	if err != nil {
		// VC 검증 실패 시 구체적인 오류 메시지를 클라이언트에 반환합니다.
		fmt.Printf("VC 검증 실패: %v\n", err) // 서버 로그에 에러 출력
		c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 검증 실패: " + err.Error()})
		return
	}

	// 3. 모든 과정 통과하면 : ✅ JWT 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "토큰 생성 실패"})
		return
	}

	// ✅ 역할 함께 응답
	c.JSON(http.StatusOK, gin.H{
		"message": "✅ 로그인 성공",
		"token":   tokenString,
		"user":    user.ID,
		"role":    user.Role,
	})
}
