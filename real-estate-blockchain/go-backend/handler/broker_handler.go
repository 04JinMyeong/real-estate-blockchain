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
	"realestate/did"     // 사용자 정의 did 패키지
	"realestate/models"  // 모델
	"realestate/service" // API 호출을 위한 service 패키지
	"realestate/vc"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 공인중개사 회원가입(DID 발급 포함) 요청 시 사용될 구조체
type SignUpBrokerWithDIDRequest struct {
	PlatformUsername   string `json:"platform_username" binding:"required"`
	PlatformPassword   string `json:"platform_password" binding:"required"`
	AgentPublicKey     string `json:"agent_public_key" binding:"required"` // Base64 인코딩된 공개키
	AgentName          string `json:"agentName"`                           // 대표자 성명 필드 추가
	RegistrationNumber string `json:"registrationNumber"`                  // 중개사무소 등록번호 필드 추가
}

func SignUpBrokerAndIssueDID(c *gin.Context) {
	var req SignUpBrokerWithDIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// --- [추가] 🚀 자격증 검증 로직 (가장 먼저 수행) ---
	fmt.Println("[Broker Signup] 중개사 자격증 번호를 검증합니다:", req.AgentName)
	isValid, err := service.VerifyAgentLicense(req.AgentName, req.RegistrationNumber)
	if err != nil {
		// 목업 API 서버가 꺼져있거나 네트워크 오류가 발생한 경우
		c.JSON(http.StatusInternalServerError, gin.H{"error": "자격 검증 시스템 오류: " + err.Error()})
		return
	}
	if !isValid {
		// API가 '자격 없음'이라고 응답한 경우
		c.JSON(http.StatusForbidden, gin.H{"error": "유효한 공인중개사 정보가 아닙니다."})
		return
	}
	fmt.Println("[Broker Signup] Agent license successfully verified.")
	// --- [추가] 🚀 검증 로직 끝 ---

	// --- 여기서부터는 검증을 통과한 경우에만 실행되는 기존 로직입니다 ---

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
		ID:                 req.PlatformUsername,
		Password:           string(hashedPassword),
		Enrolled:           false,
		CreatedAt:          time.Now(),
		Role:               "agent",        // 중개인임을 명시
		DID:                agentDIDString, // <<< 여기에 생성된 DID 저장 (중요!)
		AgentName:          req.AgentName,
		RegistrationNumber: req.RegistrationNumber,
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

// IssueVC 함수는 인증된 공인중개사에게 VC를 발급합니다.
func IssueVC(c *gin.Context) {
	// 1. 미들웨어로부터 사용자 정보 받기
	//    이전 단계에서 만든 AuthMiddleware가 c.Get()으로 정보를 조회할 수 있도록
	//    Context에 'userID'와 'userRole'을 설정해준다고 가정합니다.
	userID_interface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증 정보를 찾을 수 없습니다."})
		return
	}
	userID := userID_interface.(string)

	userRole_interface, _ := c.Get("userRole")
	userRole := userRole_interface.(string)

	// 2. 역할 검증 (공인중개사만 발급 가능)
	if userRole != "agent" {
		c.JSON(http.StatusForbidden, gin.H{"error": "VC 발급 권한이 없습니다. 공인중개사만 가능합니다."})
		return
	}

	// 3. DB에서 사용자 정보 조회 (DID, 이름 등)
	var user models.User
	db := database.GetDB()
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다."})
		return
	}
	if user.DID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "이 사용자에게는 DID가 발급되지 않았습니다."})
		return
	}

	// 4. [목업]전과기록 API 호출
	//    (실제로는 이 로직도 service 패키지로 분리하는 것이 더 좋습니다)
	resp, err := http.Get("http://localhost:8082/check?id=" + user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "전과기록 조회 시스템 오류: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	var result struct {
		HasCriminalRecord bool `json:"hasCriminalRecord"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "전과기록 API 응답 처리 오류"})
		return
	}

	// 5. VC에 담을 정보(Claim) 생성
	vcClaims := map[string]interface{}{
		"license_active":  true,                     // 회원가입 시 자격이 검증되었음을 의미
		"criminal_record": result.HasCriminalRecord, // API 조회 결과를 Claim에 반영
	}

	// 6. vc/issuer.go의 함수를 호출하여 VC 생성
	//    CreateVC 함수는 (발급자DID, 사용자DID, 사용자이름, 클레임) 등을 인자로 받는다고 가정합니다.
	vcJSON, err := vc.GenerateAndSignVC(user.DID, user.AgentName, vcClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "VC 생성 실패: " + err.Error()})
		return
	}

	// 7. 생성된 VC를 사용자의 DB 레코드에 저장
	if err := db.Model(&user).Update("vc", vcJSON).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB에 VC 저장 실패: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VC 발급 성공", "vc": vcJSON})
}
