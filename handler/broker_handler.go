package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	//"time"

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
	Email            string `json:"email" binding:"required"`
	FullName         string `json:"full_name" binding:"required"`
	LicenseNumber    string `json:"license_number" binding:"required"`
	OfficeAddress    string `json:"office_address,omitempty"`
	AgentPublicKey   string `json:"agent_public_key" binding:"required"` // 프론트에서 전달하는 Base64 인코딩된 공개키
}

func SignUpBrokerAndIssueDID(c *gin.Context) {
	var req SignUpBrokerWithDIDRequest // 요청 DTO를 models 패키지로 옮기거나 현재 위치에 정의
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
	fmt.Printf("✅ Generated DID Document: %+v\n", didDoc) // 개발 중 확인용

	// --- 4. DID Document DB 저장 ---
	err = database.StoreDIDDocument(agentDIDString, didDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store DID Document: " + err.Error()})
		return
	}
	fmt.Println("✅ DID Document stored for:", agentDIDString)

	
	// --- 5. 플랫폼 사용자 계정 생성 및 정보 저장 ---
	// (주의: 이 부분은 실제 애플리케이션의 요구사항에 맞게 확장/보안 강화 필요)
	db := database.GetDB()

	// 플랫폼 사용자 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PlatformPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password: " + err.Error()})
		return
	}

	// 플랫폼 사용자 정보 (models.User - GORM 모델)
	platformUser := models.User{
		// ID는 GORM이 자동 생성하도록 둘 수 있음 (또는 별도 규칙)
		Username: req.PlatformUsername,
		Password: string(hashedPassword),
		// Email:    req.Email, // User 모델에 Email 필드가 있다면 추가
		// AgentDID: agentDIDString, // User 모델에 DID를 직접 저장하거나, 별도 매핑 테이블 사용
	}
	if err := db.Create(&platformUser).Error; err != nil {
		// 사용자 이름 중복 등 실제 에러 처리 필요
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform user: " + err.Error()})
		return
	}
	fmt.Println("✅ Platform user created:", platformUser.Username)

	/*
	(선택적) 공인중개사 프로필 정보 저장 (models.BrokerProfile 같은 새 모델 생성 고려)
	이 모델에는 FullName, LicenseNumber, OfficeAddress, 그리고 platformUser.ID (FK), agentDIDString 등을 저장
	예시:
	brokerProfile := models.BrokerProfile{
	 PlatformUserID: platformUser.ID,
	 DID: agentDIDString,
	 FullName: req.FullName,
	 LicenseNumber: req.LicenseNumber,
	 OfficeAddress: req.OfficeAddress,
	 Email: req.Email,
	}
	if err := db.Create(&brokerProfile).Error; err != nil {
	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create broker profile: " + err.Error()})
	 return
	}
	fmt.Println("✅ Broker profile created for DID:", agentDIDString)
	*/

	/*
	--- 6. (선택적) Fabric 네트워크용 사용자(Identity) 등록 ---
	이 부분은 필요하다면 `blockchain.RegisterAndEnrollUser` 등을 호출합니다.
	fabricUsername := agentDIDString // 또는 platformUser.Username 사용
	sdk, errSdk := fabsdk.New(config.FromFile(filepath.Clean("./connection-org1.yaml")))
	if errSdk == nil {
	   defer sdk.Close()
	   mspClient, _ := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	   // 에러 처리 및 상세 로직 필요
	   blockchain.RegisterAndEnrollUser(sdk, mspClient, fabricUsername)
	}
	*/

	// --- 7. 성공 응답 ---
	c.JSON(http.StatusOK, gin.H{
		"message": "Agent registered and DID issued successfully.",
		"did":     agentDIDString,
	})
}

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
