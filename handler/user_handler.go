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

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	username := req.Username
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	// Wallet 열기
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wallet 열기 실패", "detail": err.Error()})
		return
	}

	if wallet.Exists(username) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("⚠️ 사용자 '%s' 이미 wallet에 있음", username)})
		return
	}

	// SDK 초기화
	sdk, err := fabsdk.New(config.FromFile(filepath.Clean(ccpPath)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SDK 초기화 실패", "detail": err.Error()})
		return
	}
	defer sdk.Close()

	// MSP 클라이언트 생성
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MSP 클라이언트 생성 실패", "detail": err.Error()})
		return
	}

	// 사용자 등록 시도
	secret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Type:        "client",
		Affiliation: "org1.department1",
	})
	if err != nil {
		// 이미 등록된 경우를 대비한 기본 비밀번호 사용
		secret = "userpw"
	}

	// 사용자 Enroll
	err = mspClient.Enroll(username, msp.WithSecret(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 Enroll 실패", "detail": err.Error()})
		return
	}

	// SigningIdentity에서 인증서와 개인키 직접 추출
	signingIdentity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SigningIdentity 조회 실패", "detail": err.Error()})
		return
	}

	cert := signingIdentity.EnrollmentCertificate()
	key, err := signingIdentity.PrivateKey().Bytes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "개인키 추출 실패", "detail": err.Error()})
		return
	}

	// Wallet에 저장
	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	err = wallet.Put(username, identity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wallet 저장 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("✅ 사용자 '%s' 등록 완료!", username)})
}
