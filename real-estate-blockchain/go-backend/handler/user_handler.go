package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"realestate/models" // ✅ models 폴더 경로에 맞게 수정 완료

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"` // "user" or "agent"
}

// 📌 공통 등록 로직 함수
func register(username string) error {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return fmt.Errorf("wallet 열기 실패: %w", err)
	}

	if wallet.Exists(username) {
		return nil
	}

	sdk, err := fabsdk.New(config.FromFile(filepath.Clean(ccpPath)))
	if err != nil {
		return fmt.Errorf("SDK 초기화 실패: %w", err)
	}
	defer sdk.Close()

	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		return fmt.Errorf("MSP 클라이언트 생성 실패: %w", err)
	}

	secret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Type:        "client",
		Affiliation: "org1.department1",
	})
	if err != nil {
		return fmt.Errorf("사용자 등록 실패: %w", err)
	}

	err = mspClient.Enroll(username, msp.WithSecret(secret))
	if err != nil {
		return fmt.Errorf("사용자 Enroll 실패: %w", err)
	}

	signingIdentity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		return fmt.Errorf("SigningIdentity 조회 실패: %w", err)
	}

	cert := signingIdentity.EnrollmentCertificate()
	key, err := signingIdentity.PrivateKey().Bytes()
	if err != nil {
		return fmt.Errorf("개인키 추출 실패: %w", err)
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	err = wallet.Put(username, identity)
	if err != nil {
		return fmt.Errorf("wallet 저장 실패: %w", err)
	}

	return nil
}

// ✅ HTTP 사용자 등록 API (POST /register-user)
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "요청이 올바르지 않습니다"})
			return
		}

		// Fabric Wallet 등록
		if err := register(req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// DB 저장
		user := models.User{
			ID:       req.Username,
			Password: req.Password,
			Email:    req.Name,
			Enrolled: true,
			Role:     req.Role,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 사용자 저장 실패", "detail": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("✅ 사용자 '%s' 등록 완료", req.Username)})
	}
}

// ✅ CLI 사용자 등록
func RegisterUserCLI(username string) error {
	return register(username)
}
