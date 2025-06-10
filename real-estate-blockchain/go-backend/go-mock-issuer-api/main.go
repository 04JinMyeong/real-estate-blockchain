package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// VC, Proof 구조체 정의
type VerifiableCredential struct {
	Context           []string    `json:"@context"`
	ID                string      `json:"id"`
	Type              []string    `json:"type"`
	Issuer            string      `json:"issuer"`
	IssuanceDate      string      `json:"issuanceDate"`
	CredentialSubject interface{} `json:"credentialSubject"`
	Proof             *Proof      `json:"proof,omitempty"`
}

type Proof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	VerificationMethod string `json:"verificationMethod"`
	ProofPurpose       string `json:"proofPurpose"`
	Jws                string `json:"jws"`
}

var issuerPrivateKey *ecdsa.PrivateKey
var issuerPublicKey *ecdsa.PublicKey
var issuerDID = "did:example:mock-issuer-123"

func init() {
	var err error
	issuerPrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate issuer private key: %v", err)
	}
	issuerPublicKey = &issuerPrivateKey.PublicKey
	fmt.Println("✅ Mock Issuer: 서명을 위한 키 쌍 생성 완료.")
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	router.Use(cors.New(config))

	// --- ▼▼▼ 핸들러들을 여기에 독립적으로 정의합니다 ▼▼▼ ---

	// 1. VC 발급을 처리하는 API 핸들러
	router.POST("/issue-vc", func(c *gin.Context) {
		var req struct{ Name, ID, DID string }
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, id, did 필드가 필요합니다."})
			return
		}

		hasCriminalRecord := strings.Contains(req.ID, "bad")
		credentialSubject := map[string]interface{}{
			"id":              req.DID,
			"name":            req.Name,
			"license_active":  true,
			"criminal_record": hasCriminalRecord,
		}

		vcToSign := VerifiableCredential{
			Context:           []string{"https://www.w3.org/2018/credentials/v1"},
			ID:                "http://example.edu/credentials/3732",
			Type:              []string{"VerifiableCredential", "AgentCredential"},
			Issuer:            issuerDID,
			IssuanceDate:      time.Now().UTC().Format(time.RFC3339),
			CredentialSubject: credentialSubject,
		}

		payloadBytes, _ := json.Marshal(vcToSign)
		hash := sha256.Sum256(payloadBytes)
		r, s, err := ecdsa.Sign(rand.Reader, issuerPrivateKey, hash[:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "VC 서명 실패"})
			return
		}
		signatureBytes := append(r.Bytes(), s.Bytes()...)
		jwsSignature := base64.RawURLEncoding.EncodeToString(signatureBytes)

		vcToSign.Proof = &Proof{
			Type:               "EcdsaSecp256k1Signature2019",
			Created:            time.Now().UTC().Format(time.RFC3339),
			VerificationMethod: issuerDID + "#key-1",
			ProofPurpose:       "assertionMethod",
			Jws:                jwsSignature,
		}

		fmt.Println("➡️ Mock Issuer: VC Generated and Signed. Responding to client.")
		c.JSON(http.StatusOK, gin.H{"vc": vcToSign})
	})

	// 2. 공개키를 반환하는 API 핸들러
	router.GET("/public-key", func(c *gin.Context) {
		publicKeyBytes := elliptic.Marshal(issuerPublicKey.Curve, issuerPublicKey.X, issuerPublicKey.Y)
		c.JSON(http.StatusOK, gin.H{
			"publicKey": base64.StdEncoding.EncodeToString(publicKeyBytes),
		})
	})
	// --- ▲▲▲ 핸들러 정의 끝 ▲▲▲ ---

	fmt.Println("🚀 Mock VC Issuer API server is running on http://localhost:8083")
	router.Run(":8083")
}
