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
	"realestate/did" // 사용자 정의 did 패키지
	"realestate/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // 플랫폼 사용자 비밀번호 해싱용
	// "github.com/google/uuid"
)

// 공인중개사 회원가입(DID 발급 포함) 요청 시 사용될 구조체로, 프론트엔드에서 보내는 JSON의 필드명과 일치해야 합니다.
type SignUpBrokerWithDIDRequest struct {
	PlatformUsername string `json:"platform_username" binding:"required"`
	PlatformPassword string `json:"platform_password" binding:"required"`
	// Email, FullName, LicenseNumber, OfficeAddress 필드 제거
	AgentPublicKey string `json:"agent_public_key" binding:"required"` // 프론트에서 전달하는 Base64 인코딩된 공개키
}

func SignUpBrokerAndIssueDID(c *gin.Context) {
	var req SignUpBrokerWithDIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// --- 1. 공개키 처리 ---
	agentPubKeyBytes, err := base64.StdEncoding.DecodeString(req.AgentPublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid public key encoding: " + err.Error()})
		return
	}

	// --- 2. DID 생성 ---
	// 공개키로부터 DID 생성 (did 패키지의 함수 사용)
	agentDIDString := did.GenerateDIDFromPublicKey(agentPubKeyBytes)
	fmt.Println("✅ Generated Agent DID:", agentDIDString)

	// --- 3. DID Document 생성 ---
	// 프론트엔드에서 Ed25519 키를 생성했으므로, 관련 키 타입을 명시합니다.
	keyTypeForDoc := "Ed25519VerificationKey2020" // 또는 "Ed25519VerificationKey2018"
	didDoc, err := did.CreateAgentDIDDocument(agentDIDString, agentPubKeyBytes, keyTypeForDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DID Document: " + err.Error()})
		return
	}
	// fmt.Printf("✅ Generated DID Document: %+v\n", didDoc) // 개발 중 확인용 로그, 필요시 주석 해제

	// --- 4. DID Document DB 저장 ---
	// database 패키지의 함수를 사용하여 DID Document를 DB에 저장
	err = database.StoreDIDDocument(agentDIDString, didDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store DID Document: " + err.Error()})
		return
	}
	fmt.Println("✅ DID Document stored for:", agentDIDString)

	// --- 5. 플랫폼 사용자 계정 생성 ---(주의: 이 부분은 실제 애플리케이션의 요구사항에 맞게 확장/보안 강화 필요)
	db := database.GetDB() // GORM DB 인스턴스 가져오기

	// 플랫폼 사용자 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PlatformPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password: " + err.Error()})
		return
	}

	// platformUser 초기화 시 Username 대신 ID 필드 사용
	platformUser := models.User{
		ID:        req.PlatformUsername, // req.PlatformUsername 값을 User 구조체의 ID 필드에 할당
		Password:  string(hashedPassword),
		Enrolled:  false,
		CreatedAt: time.Now(), // time 패키지 import 필요
	}

	// 사용자 이름 중복 확인 (선택적이지만 중요)
	var existingUser models.User
	if err := db.First(&existingUser, "id = ?", req.PlatformUsername).Error; err == nil {
		// err == nil 이면 사용자가 이미 존재함
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists with this username: " + req.PlatformUsername})
		return
	} // err != nil 이고 errors.Is(err, gorm.ErrRecordNotFound) 이면 사용자가 없어 생성 가능

	// 플랫폼 사용자 정보 DB에 생성
	if err := db.Create(&platformUser).Error; err != nil {
		// 실제 에러 처리 (예: DB 연결 오류 등)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform user: " + err.Error()})
		return
	}
	fmt.Println("✅ Platform user created:", platformUser.ID)

	// --- 최종 성공 응답 ---
	// 발급된 DID 정보와 함께 성공 메시지를 클라이언트에 전달
	c.JSON(http.StatusOK, gin.H{
		"message": "Agent registered and DID issued successfully.",
		"did":     agentDIDString, // 생성된 공인중개사의 DID
	})
}

// VerifyBroker 및 signVC 함수는 VC 검증 단계에서 사용되므로 현재 단계에서는 수정하지 않습니다.
// ... (VerifyBroker, signVC 함수 기존 코드 유지) ...

// VerifyBrokerRequest payload for VC verification
type VerifyBrokerRequest struct {
	ID string `json:"id"` // Broker DID
}

// VerifyBroker checks the validity of a broker's VC
func VerifyBroker(c *gin.Context) {
	var req VerifyBrokerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	// Retrieve stored VC
	vc, err := database.GetBrokerVC(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VC를 찾을 수 없습니다"})
		return
	}

	// Recalculate signature
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

// signVC computes a SHA-256 hash of the VC content as a signature
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
