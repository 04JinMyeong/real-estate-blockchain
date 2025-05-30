package vc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"realestate/crypto" // crypto 패키지 import 확인
)

// GenerateAndSignVC issues a Verifiable Credential for a real estate agent.
// It validates the agent's qualification, constructs the VC JSON, and attaches a digital signature.
// brokerStatus 파라미터 추가 (이전 안내 내용 반영)
func GenerateAndSignVC(did, name, licenseNum, phone string, brokerStatus string) (string, error) { // ◀◀◀ brokerStatus 파라미터 추가
	// 1. Validate broker qualification (현재는 licenseNum이 비어있지 않은지만 확인)
	if !validateBroker(name, licenseNum) { // validateBroker 함수는 vc/validate.go 에 정의되어 있어야 함
		return "", fmt.Errorf("broker validation failed for license: %s", licenseNum)
	}

	// 2. Read issuer DID from env
	issuerDID := os.Getenv("ISSUER_DID")
	if issuerDID == "" {
		return "", fmt.Errorf("ISSUER_DID environment variable is not set")
	}

	// 3. Build the credential document
	vc := map[string]interface{}{
		"@context":     []string{"https://www.w3.org/2018/credentials/v1"},
		"type":         []string{"VerifiableCredential", "RealEstateAgentVC"},
		"issuer":       issuerDID,
		"issuanceDate": time.Now().Format(time.RFC3339),
		"credentialSubject": map[string]interface{}{ // map[string]string 에서 map[string]interface{} 로 변경하여 유연성 확보
			"id":         did,
			"name":       name,
			"licenseNum": licenseNum,
			"phone":      phone,
			"status":     brokerStatus, // ◀◀◀ 파라미터로 받은 brokerStatus 사용
		},
	}

	// 4. Sign the credentialSubject
	subjectBytes, err := json.Marshal(vc["credentialSubject"])
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentialSubject: %w", err)
	}
	signature, err := crypto.Sign(subjectBytes) // crypto.Sign 함수는 crypto/keys.go 에 정의
	if err != nil {
		return "", fmt.Errorf("failed to sign VC: %w", err)
	}
	jws := base64.StdEncoding.EncodeToString(signature)

	// 5. Attach proof
	vc["proof"] = map[string]string{
		"type":               "Ed25519Signature2020",
		"created":            time.Now().Format(time.RFC3339),
		"proofPurpose":       "assertionMethod",
		"verificationMethod": issuerDID + "#key-1", // 발급자 DID 문서 내 검증키 ID (예시)
		"jws":                jws,
	}

	// 6. Return the final VC JSON
	vcBytes, err := json.MarshalIndent(vc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal final VC: %w", err)
	}
	return string(vcBytes), nil
}
