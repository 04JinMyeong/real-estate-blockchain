package handler

import (
	"encoding/json" // VC JSON 파싱을 위해 추가
	"fmt"
	"net/http"
	"time"

	"realestate/database"
	"realestate/models"
	"realestate/vc" // vc.VerifiableCredential, vc.VerifyVC 사용

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
		VC       string `json:"vc" binding:"required"` // VC를 JSON 문자열로 받음
	}

	// 요청 바인딩 및 기본 검증
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("[Login Handler] Error binding JSON:", err.Error()) // 바인딩 에러 시 로그 추가
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다: " + err.Error()})
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

	// VC가 비어있는 경우의 처리 (binding:"required"가 이미 처리하지만, 방어적으로 추가)
	if req.VC == "" { // 이 부분은 binding:"required"에 의해 사실상 도달하기 어려움
		fmt.Println("[Login Handler] VC string is empty, though binding was required.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "VC가 제출되지 않았습니다."})
		return
	}

	// VC 문자열을 vc.VerifiableCredential 구조체로 언마샬링
	var receivedUserVC vc.VerifiableCredential             // vc 패키지에 정의된 구조체 사용
	err := json.Unmarshal([]byte(req.VC), &receivedUserVC) // err 변수 새로 선언
	if err != nil {
		fmt.Println("[Login Handler] Error unmarshalling VC:", err.Error()) // VC 파싱 에러 로그
		c.JSON(http.StatusBadRequest, gin.H{"error": "제출된 VC의 JSON 형식이 올바르지 않습니다: " + err.Error()})
		return
	}
	fmt.Println("[Login Handler] VC unmarshalled successfully. VC ID:", receivedUserVC.ID)

	// VC 소유권 확인 (공인중개사 역할일 때)
	if user.Role == "agent" {
		if user.DID == "" {
			fmt.Println("[Login Handler] User is agent but DID is missing in DB for user:", user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자의 DID 정보가 DB에 없습니다. VC 검증 불가."})
			return
		}

		credSubject, ok := receivedUserVC.CredentialSubject.(map[string]interface{})
		if !ok {
			fmt.Println("[Login Handler] VC CredentialSubject is not map[string]interface{}")
			c.JSON(http.StatusBadRequest, gin.H{"error": "VC의 CredentialSubject 형식이 올바르지 않습니다."})
			return
		}
		vcOwnerDID, ok := credSubject["id"].(string)
		if !ok {
			fmt.Println("[Login Handler] VC CredentialSubject.id is missing or not a string")
			c.JSON(http.StatusBadRequest, gin.H{"error": "VC의 CredentialSubject에 id(DID) 필드가 없거나 문자열이 아닙니다."})
			return
		}

		if user.DID != vcOwnerDID {
			fmt.Printf("[Login Handler] VC Ownership Mismatch: User DID (%s) != VC DID (%s)\n", user.DID, vcOwnerDID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 소유권 검증 실패: 제출된 VC는 현재 사용자의 것이 아닙니다."})
			return
		}
		fmt.Printf("[Login Handler] VC ownership verified for user %s\n", user.ID)
	}

	// VC 유효성 및 클레임 검증 (vc.VerifyVC 함수 사용)
	// vc.VerifyVC는 (bool, error)를 반환한다고 가정합니다.
	isValid, verificationErr := vc.VerifyVC(req.VC)
	if verificationErr != nil {
		fmt.Printf("[Login Handler] VC verification error: %v\n", verificationErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 유효성 검증 오류: " + verificationErr.Error()})
		return
	}
	if !isValid {
		fmt.Println("[Login Handler] VC is not valid (isValid is false)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "제출된 VC가 유효하지 않습니다."})
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
