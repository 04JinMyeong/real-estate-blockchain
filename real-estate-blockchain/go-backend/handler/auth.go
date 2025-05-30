package handler

import (
	"fmt"
	"net/http"
	"time"

	"realestate/database"
	"realestate/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

// 로그인 핸들러
func Login(c *gin.Context) {
	var req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다"})
		return
	}

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

	// ✅ JWT 생성
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
