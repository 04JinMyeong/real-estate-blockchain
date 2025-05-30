package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"realestate/database"
	"realestate/did"    // 사용자 정의 did 패키지
	"realestate/models" // 모델

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 공인중개사 회원가입(DID 발급 포함) 요청 시 사용될 구조체
type SignUpBrokerWithDIDRequest struct {
	PlatformUsername string `json:"platform_username" binding:"required"`
	PlatformPassword string `json:"platform_password" binding:"required"`
	AgentPublicKey   string `json:"agent_public_key" binding:"required"` // Base64 인코딩된 공개키
}

func SignUpBrokerAndIssueDID(c *gin.Context) {
	var req SignUpBrokerWithDIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 1. 공개키 디코딩
	agentPubKeyBytes, err := base64.StdEncoding.DecodeString(req.AgentPublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid public key encoding: " + err.Error()})
		return
	}

	// 2. DID 생성
	agentDIDString := did.GenerateDIDFromPublicKey(agentPubKeyBytes)
	fmt.Println("✅ Generated Agent DID:", agentDIDString)

	// 3. DID Document 생성
	keyTypeForDoc := "Ed25519VerificationKey2020"
	didDoc, err := did.CreateAgentDIDDocument(agentDIDString, agentPubKeyBytes, keyTypeForDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DID Document: " + err.Error()})
		return
	}

	// 4. DID Document DB 저장
	if err := database.StoreDIDDocument(agentDIDString, didDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store DID Document: " + err.Error()})
		return
	}
	fmt.Println("✅ DID Document stored for:", agentDIDString)

	// 5. 플랫폼 사용자 DB 계정 생성
	db := database.GetDB()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PlatformPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password: " + err.Error()})
		return
	}

	// 중복 확인
	var existingUser models.User
	if err := db.First(&existingUser, "id = ?", req.PlatformUsername).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists with this username: " + req.PlatformUsername})
		return
	}

	platformUser := models.User{
		ID:        req.PlatformUsername,
		Password:  string(hashedPassword),
		Enrolled:  false,
		CreatedAt: time.Now(),
		Role:      "agent", // 중개인임을 명시
	}

	if err := db.Create(&platformUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform user: " + err.Error()})
		return
	}
	fmt.Println("✅ Platform user created:", platformUser.ID)

	// 6. Fabric Wallet 등록
	if err := RegisterUserCLI(req.PlatformUsername); err != nil {
		fmt.Printf("❗ Wallet 등록 실패: %v\n", err)
		// 실패해도 프로세스는 계속 진행 (에러 메시지 포함 응답)
		c.JSON(http.StatusOK, gin.H{
			"message": "Agent registered and DID issued successfully, but wallet registration failed.",
			"did":     agentDIDString,
		})
		return
	}

	// 등록 성공 시 DB 업데이트
	platformUser.Enrolled = true
	db.Save(&platformUser)

	// 7. 최종 응답
	c.JSON(http.StatusOK, gin.H{
		"message": "Agent registered and DID issued successfully.",
		"did":     agentDIDString,
	})
}

// VC 검증용 요청 구조체
type VerifyBrokerRequest struct {
	ID string `json:"id"` // Broker DID
}

// VC 서명 검증용 내부 함수
func signVC(vc models.BrokerVC) (string, error) {
	temp := vc
	temp.Signature = ""
	data, err := json.Marshal(temp)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// VC 검증 핸들러
func VerifyBroker(c *gin.Context) {
	var req VerifyBrokerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	vc, err := database.GetBrokerVC(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VC를 찾을 수 없습니다"})
		return
	}

	expectedSig, err := signVC(*vc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서명 검증 실패", "detail": err.Error()})
		return
	}
	if vc.Signature != expectedSig {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 검증 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VC 검증 성공", "vc": vc})
}
