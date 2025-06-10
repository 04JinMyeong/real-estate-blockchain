package handler

import (
	"encoding/json" // VC JSON 파싱을 위해 추가
	"fmt"
	"net/http"
	"strings"
	"time"

	"realestate/database"
	"realestate/models"
	"realestate/vc" // vc.VerifiableCredential, vc.VerifyVC 사용

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key") // 실제 서비스에서는 환경변수나 안전한 저장소 사용

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware는 토큰을 검증하고 사용자 정보를 다음 핸들러로 전달하는 미들웨어입니다.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 요청 헤더에서 'Authorization' 값을 가져옵니다.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "인증 헤더가 없습니다."})
			return
		}

		// 2. 토큰은 보통 "Bearer <token>" 형식이므로, "Bearer " 부분을 제거합니다.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // "Bearer "가 없는 경우
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer 토큰 형식이 아닙니다."})
			return
		}

		// 3. 토큰을 파싱하고 유효성을 검증합니다.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 서명 방식이 HMAC인지 확인하고, 우리가 사용하는 비밀키(jwtKey)를 반환합니다.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil // jwtKey는 파일 상단에 정의된 것과 동일한 키
		})

		// 4. 파싱 중 에러가 발생했거나, 토큰이 유효하지 않은 경우 요청을 차단합니다.
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			return
		}

		// 5. 토큰이 유효하면, 다음 핸들러(IssueVC)에서 사용할 수 있도록
		//    Context에 사용자 정보를 저장합니다.
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		// 6. 모든 검증을 통과했으므로, 다음 핸들러로 요청을 전달합니다.
		c.Next()
	}
}

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
	// --- ▼ 1단계: VC를 선택적으로 받도록 구조체 수정 ---
	var req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		VC       string `json:"vc"` // binding:"required" 제거
	}

	// 요청 바인딩 및 기본 검증
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값이 올바르지 않습니다: " + err.Error()})
		return
	}

	// DB에서 사용자 정보 조회 (기존과 동일)
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, "id = ?", req.ID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "존재하지 않는 사용자입니다"})
		return
	}

	// 비밀번호 검증 (기존과 동일)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "비밀번호가 일치하지 않습니다"})
		return
	}

	// --- ▼ 2단계: 역할(Role)에 따라 VC 검증 로직 분기 ---
	if user.Role == "agent" {
		// 역할이 'agent'인 경우, VC 검증 절차를 시작합니다.
		fmt.Printf("[Login Handler] User %s is an agent. Starting VC verification.\n", user.ID)

		// 1. 공인중개사가 VC를 제출했는지 확인
		if req.VC == "" {
			fmt.Printf("[Login Handler] Agent %s did not submit a VC.\n", user.ID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "공인중개사 로그인을 위해서는 VC를 반드시 제출해야 합니다."})
			return
		}

		// 2. VC 소유권 검증 (기존 코드를 그대로 활용)
		var receivedUserVC vc.VerifiableCredential
		if err := json.Unmarshal([]byte(req.VC), &receivedUserVC); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "제출된 VC의 JSON 형식이 올바르지 않습니다: " + err.Error()})
			return
		}
		credSubject, _ := receivedUserVC.CredentialSubject.(map[string]interface{})
		vcOwnerDID, _ := credSubject["id"].(string)
		if user.DID != vcOwnerDID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 소유권 검증 실패: 제출된 VC는 현재 사용자의 것이 아닙니다."})
			return
		}

		// 3. VC 유효성 및 클레임 검증 (기존 코드를 그대로 활용)
		isValid, verificationErr := vc.VerifyVC(req.VC) // 이 함수가 Claim 검증까지 내부적으로 처리한다고 가정
		if verificationErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 유효성 검증 오류: " + verificationErr.Error()})
			return
		}
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "제출된 VC가 유효하지 않거나 자격이 맞지 않습니다."})
			return
		}
		fmt.Printf("[Login Handler] Agent %s successfully verified VC.\n", user.ID)
	}
	// 역할이 'user'인 경우, 위 'if user.Role == "agent"' 블록을 모두 건너뛰고 바로 여기로 오게 됩니다.

	fmt.Printf("[Login Handler] All checks passed for user %s with role %s. Generating JWT.\n", user.ID, user.Role)

	// ✅ JWT 생성 (기존과 동일)
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

	// ✅ 역할 함께 응답 (기존과 동일)
	c.JSON(http.StatusOK, gin.H{
		"message": "✅ 로그인 성공",
		"token":   tokenString,
		"user":    user.ID,
		"role":    user.Role,
	})
}
