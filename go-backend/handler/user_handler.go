package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

// 📌 공통 등록 로직 함수
func register(username string) error {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	// Wallet 열기
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return fmt.Errorf("wallet 열기 실패: %w", err)
	}

	if wallet.Exists(username) {
		return nil // 이미 존재하면 그냥 성공으로 처리
	}

	// SDK 초기화
	sdk, err := fabsdk.New(config.FromFile(filepath.Clean(ccpPath)))
	if err != nil {
		return fmt.Errorf("SDK 초기화 실패: %w", err)
	}
	defer sdk.Close()

	// MSP 클라이언트 생성
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		return fmt.Errorf("MSP 클라이언트 생성 실패: %w", err)
	}

	// 사용자 등록 시도
	secret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Type:        "client",
		Affiliation: "org1.department1",
	})
	if err != nil {
		// 등록 실패 시 enroll 시도하지 말고 에러 반환
		return fmt.Errorf("사용자 등록 실패: %w", err)
	}

	// 사용자 Enroll
	err = mspClient.Enroll(username, msp.WithSecret(secret))
	if err != nil {
		return fmt.Errorf("사용자 Enroll 실패: %w", err)
	}

	// SigningIdentity에서 인증서와 개인키 추출
	signingIdentity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		return fmt.Errorf("SigningIdentity 조회 실패: %w", err)
	}

	cert := signingIdentity.EnrollmentCertificate()
	key, err := signingIdentity.PrivateKey().Bytes()
	if err != nil {
		return fmt.Errorf("개인키 추출 실패: %w", err)
	}

	// Wallet에 저장
	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	err = wallet.Put(username, identity)
	if err != nil {
		return fmt.Errorf("wallet 저장 실패: %w", err)
	}

	return nil
}

// ✅ HTTP API용 사용자 등록 (POST /register-user)
func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	err := register(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("✅ 사용자 '%s' 등록 완료!", req.Username)})
}

// ✅ CLI용 사용자 등록 (go run main.go register TestUser9)
func RegisterUserCLI(username string) error {
	return register(username)
}
