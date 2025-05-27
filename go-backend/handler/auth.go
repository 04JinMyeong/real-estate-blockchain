package handler

import (
	"log"
	"net/http"
	"time"

	"realestate/database"
	"realestate/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 회원가입 핸들러 (자동 사용자 등록 포함)
func Signup(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다"})
		return
	}

	// 이미 존재하는 사용자 확인
	var existing models.User
	if err := database.DB.First(&existing, "id = ?", req.ID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "이미 존재하는 사용자입니다"})
		return
	}

	// 비밀번호 해시화
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 해싱 실패"})
		return
	}

	newUser := models.User{
		ID:        req.ID,
		Password:  string(hashed),
		Email:     req.Email,
		Enrolled:  false,
		CreatedAt: time.Now(),
	}

	// DB 저장
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 저장 실패"})
		return
	}

	// Fabric 사용자 등록 (req.ID 직접 전달)
	if err := RegisterUserCLI(req.ID); err != nil {
		log.Printf("❗ 사용자 등록 실패 (Fabric): %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "✅ 회원가입 완료 (단, 블록체인 사용자 등록 실패)"})
		return
	}

	// 등록 성공 시 DB에 반영
	newUser.Enrolled = true
	database.DB.Save(&newUser)

	c.JSON(http.StatusOK, gin.H{"message": "✅ 회원가입 및 블록체인 사용자 등록 완료"})
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

	var user models.User
	if err := database.DB.First(&user, "id = ?", req.ID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "존재하지 않는 사용자입니다"})
		return
	}

	// 비밀번호 대조
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "비밀번호가 일치하지 않습니다"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ 로그인 성공", "user": user.ID})
}
