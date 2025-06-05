package vc

import (
	"encoding/base64" // Base64 인코딩을 위해 추가
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"realestate/crypto" // crypto 패키지 임포트
)

// VerifiableCredential 구조체 정의(W3C VC Data Model 형식으로 구성함.)
type VerifiableCredential struct {
	Context           []string    `json:"@context"`
	ID                string      `json:"id"` // VC의 고유 ID (예: "urn:uuid:...")
	Type              []string    `json:"type"`
	Issuer            string      `json:"issuer"` // 발급자 DID (환경 변수 ISSUER_DID 사용)
	IssuanceDate      string      `json:"issuanceDate"`
	ExpirationDate    string      `json:"expirationDate,omitempty"` // 선택 사항
	CredentialSubject interface{} `json:"credentialSubject"`        // 다양한 형태의 주제를 담을 수 있도록 interface{} 사용
	Proof             *Proof      `json:"proof,omitempty"`
}

type Proof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	VerificationMethod string `json:"verificationMethod"` // 발급자의 DID#key-id
	ProofPurpose       string `json:"proofPurpose"`
	Jws                string `json:"jws"`
}

// VC를 생성하고 서명한 후 JSON 문자열로 반환하는 함수.(이 파일의 핵심함수임)
// licenseNum, phone 인자는 additionalClaims 맵으로 통합되었습니다.
func GenerateAndSignVC(agentDID, name string, additionalClaims map[string]interface{}) (string, error) {
	// 3.1. 발급자(Issuer) DID 및 개인키 경로 환경 변수에서 읽기 (crypto.Sign에서 직접 사용하므로 여기서는 확인만)
	issuerDID := os.Getenv("ISSUER_DID")
	if issuerDID == "" {
		return "", fmt.Errorf("ISSUER_DID 환경 변수가 설정되어 있지 않습니다")
	}

	// 3.2. Credential Subject 구성
	credentialSubject := map[string]interface{}{ // 다양한 정보를 담기 위해 map 사용
		"id":                agentDID, // VC를 받는 주체(공인중개사)의 DID
		"licenseHolderName": name,     // 기본 클레임으로 포함
	}

	// additionalClaims가 nil이 아니고 내용이 있다면 credentialSubject에 병합
	if additionalClaims != nil {
		for key, value := range additionalClaims {
			credentialSubject[key] = value // 여기에 fraudConvictionRecordStatus, licenseNumber, phone 등이 추가됨
		}
	}

	// 3.3. VC 기본 구조체 인스턴스 생성 (VerifiableCredential 구조체 사용)
	vc := VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://example.org/schemas/broker-v1",
		},
		ID:                fmt.Sprintf("urn:uuid:%s-%s", agentDID, time.Now().Format("20060102150405")),
		Type:              []string{"VerifiableCredential", "RealEstateAgentLicenseCredential"},
		Issuer:            issuerDID,
		IssuanceDate:      time.Now().UTC().Format(time.RFC3339),
		CredentialSubject: credentialSubject,
		// Proof는 nil 상태
	}

	// 이 아래가 “서명용 JSON(payloadBytes)” 생성 부분
	payloadBytes, err := json.Marshal(vc)
	if err != nil {
		log.Fatalf("payload marshal error: %v", err)
	}

	// **디버깅용 로그 추가**: payloadBytes를 문자열로 출력
	fmt.Println("===== [DEBUG] 발급 시 payloadBytes(JSON) =====")
	fmt.Println(string(payloadBytes))
	fmt.Println("=============================================")

	// 실제 서명(Sign) 수행
	jwsSignatureBytes, err := crypto.Sign(payloadBytes)
	if err != nil {
		log.Fatalf("Sign error: %v", err)
	}

	// 3.6. Proof 객체 생성 및 VC에 추가 (Proof 구조체 사용)
	vc.Proof = &Proof{ // 포인터 타입이므로 & 사용
		Type:               "Ed25519Signature2018",                // 사용하는 서명 방식 (Ed25519 기반 가정)
		Created:            time.Now().UTC().Format(time.RFC3339), // 서명 생성 시간
		VerificationMethod: fmt.Sprintf("%s#key-1", issuerDID),    // 서명 검증에 사용할 공개키 식별자 (예시)
		ProofPurpose:       "assertionMethod",
		// jwsSignatureBytes는 []byte 타입이므로, string으로 변환하여 할당합니다. (두 번째 오류 수정)
		Jws: base64.StdEncoding.EncodeToString(jwsSignatureBytes),
	}

	// 3.7. 최종 VC(Proof 포함)를 JSON 문자열로 변환하여 반환
	finalVCBytes, err := json.MarshalIndent(vc, "", "  ") // 사람이 읽기 좋도록 들여쓰기 적용
	if err != nil {
		return "", fmt.Errorf("최종 VC JSON 마샬링 실패: %w", err)
	}

	return string(finalVCBytes), nil
}
