// // VC를 발급하는 핵심 로직들을 구현한 파일입니다.

// package vc

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"time"

// 	"realestate/crypto"
// )

// // VerifiableCredential 구조체 정의(W3C VC Data Model 형식으로 구성함.)
// type VerifiableCredential struct {
// 	Context           []string    `json:"@context"`
// 	ID                string      `json:"id"` // VC의 고유 ID (예: "urn:uuid:...")
// 	Type              []string    `json:"type"`
// 	Issuer            string      `json:"issuer"` // 발급자 DID (환경 변수 ISSUER_DID 사용)
// 	IssuanceDate      string      `json:"issuanceDate"`
// 	ExpirationDate    string      `json:"expirationDate,omitempty"` // 선택 사항
// 	CredentialSubject interface{} `json:"credentialSubject"`        // 다양한 형태의 주제를 담을 수 있도록 interface{} 사용
// 	Proof             *Proof      `json:"proof,omitempty"`
// }

// type Proof struct {
// 	Type               string `json:"type"`
// 	Created            string `json:"created"`
// 	VerificationMethod string `json:"verificationMethod"` // 발급자의 DID#key-id
// 	ProofPurpose       string `json:"proofPurpose"`
// 	Jws                string `json:"jws"`
// }

// // VC를 생성하고 서명한 후 JSON 문자열로 반환하는 함수.(이 파일의 핵심함수임)
// func GenerateAndSignVC(agentDID, name, licenseNum, phone string, additionalClaims map[string]interface{}) (string, error) {
// 	// 3.1. 발급자(Issuer) DID 및 개인키 경로 환경 변수에서 읽기
// 	issuerDID := os.Getenv("ISSUER_DID")
// 	if issuerDID == "" {
// 		return "", fmt.Errorf("ISSUER_DID 환경 변수가 설정되어 있지 않습니다")
// 	}

// 	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
// 	if privateKeyPath == "" {
// 		return "", fmt.Errorf("PRIVATE_KEY_PATH 환경 변수가 설정되어 있지 않습니다")
// 	}

// 	// 3.2. Credential Subject 구성
// 	credentialSubject := map[string]interface{}{ // 다양한 정보를 담기 위해 map 사용
// 		"id":                agentDID, // VC를 받는 주체(공인중개사)의 DID
// 		"licenseHolderName": name,
// 		"licenseNumber":     licenseNum,
// 		"phone":             phone,
// 	}

// 	// additionalClaims가 nil이 아니고 내용이 있다면 credentialSubject에 병합
// 	if additionalClaims != nil {
// 		for key, value := range additionalClaims {
// 			credentialSubject[key] = value // 여기에 fraudConvictionRecordStatus 등이 추가됨
// 		}
// 	}

// 	// 3.3. VC 기본 구조체 인스턴스 생성 (VerifiableCredential 구조체 사용)
// 	vc := VerifiableCredential{
// 		Context: []string{ // JSON-LD 컨텍스트
// 			"https://www.w3.org/2018/credentials/v1",
// 			// 필요시 프로젝트별 컨텍스트 추가
// 		},
// 		ID:           fmt.Sprintf("urn:uuid:%s-%s", agentDID, time.Now().Format("20060102150405")), // 예시: 더 고유한 VC ID 생성
// 		Type:         []string{"VerifiableCredential", "RealEstateAgentLicenseCredential"},         // VC의 타입
// 		Issuer:       issuerDID,                                                                    // 발급자 DID
// 		IssuanceDate: time.Now().UTC().Format(time.RFC3339),                                        // VC 발급 시간 (UTC 기준)
// 		// ExpirationDate: time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC3339), // 예: 1년 후 만료 (필요시 추가)
// 		CredentialSubject: credentialSubject, // 위에서 구성한 credentialSubject
// 		// Proof는 아래에서 서명 후 채워짐
// 	}

// 	// 3.4. 서명 대상(페이로드) 준비: VC에서 Proof 필드를 제외한 부분
// 	// JWS 표준에 따라, VC의 `proof` 필드를 제외한 나머지 부분을 서명 대상으로 합니다.
// 	vcForSigning := vc       // VC 복사
// 	vcForSigning.Proof = nil // 서명할 때는 Proof 필드를 비워둠

// 	payloadBytes, err := json.Marshal(vcForSigning) // 서명 대상 VC를 JSON 바이트로 변환
// 	if err != nil {
// 		return "", fmt.Errorf("VC 페이로드 JSON 마샬링 실패: %w", err)
// 	}

// 	// 3.5. JWS 디지털 서명 생성
// 	// crypto.Sign 함수는 privateKeyPath의 개인키를 사용하여 payloadBytes를 서명합니다.
// 	jwsSignature, err := crypto.Sign(privateKeyPath, payloadBytes)
// 	if err != nil {
// 		return "", fmt.Errorf("JWS 서명 생성 실패: %w", err)
// 	}

// 	// 3.6. Proof 객체 생성 및 VC에 추가 (Proof 구조체 사용)
// 	vc.Proof = &Proof{ // 포인터 타입이므로 & 사용
// 		Type:               "Ed25519Signature2018",                // 사용하는 서명 방식 (Ed25519 기반 가정)
// 		Created:            time.Now().UTC().Format(time.RFC3339), // 서명 생성 시간
// 		VerificationMethod: fmt.Sprintf("%s#keys-1", issuerDID),   // 서명 검증에 사용할 공개키 식별자 (예시)
// 		ProofPurpose:       "assertionMethod",                     // 증명의 목적
// 		Jws:                jwsSignature,                          // 생성된 JWS 서명 값
// 	}

// 	// 3.7. 최종 VC(Proof 포함)를 JSON 문자열로 변환하여 반환
// 	finalVCBytes, err := json.MarshalIndent(vc, "", "  ") // 사람이 읽기 좋도록 들여쓰기 적용
// 	if err != nil {
// 		return "", fmt.Errorf("최종 VC JSON 마샬링 실패: %w", err)
// 	}

// 	return string(finalVCBytes), nil
// }

// File: real-estate-blockchain/go-backend/vc/issuer.go
// VC를 발급하는 핵심 로직들을 구현한 파일입니다.

package vc

import (
	"encoding/base64" // Base64 인코딩을 위해 추가
	"encoding/json"
	"fmt"
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
		Context: []string{ // JSON-LD 컨텍스트
			"https://www.w3.org/2018/credentials/v1",
			"https://example.org/schemas/broker-v1", // 프로젝트별 컨텍스트 추가
		},
		ID:                fmt.Sprintf("urn:uuid:%s-%s", agentDID, time.Now().Format("20060102150405")), // 예시: 더 고유한 VC ID 생성
		Type:              []string{"VerifiableCredential", "RealEstateAgentLicenseCredential"},         // VC의 타입
		Issuer:            issuerDID,                                                                    // 발급자 DID
		IssuanceDate:      time.Now().UTC().Format(time.RFC3339),                                        // VC 발급 시간 (UTC 기준)
		CredentialSubject: credentialSubject,                                                            // 위에서 구성한 credentialSubject
	}

	// 3.4. 서명 대상(페이로드) 준비: VC에서 Proof 필드를 제외한 부분
	vcForSigning := vc       // VC 복사
	vcForSigning.Proof = nil // 서명할 때는 Proof 필드를 비워둠

	payloadBytes, err := json.Marshal(vcForSigning) // 서명 대상 VC를 JSON 바이트로 변환
	if err != nil {
		return "", fmt.Errorf("VC 페이로드 JSON 마샬링 실패: %w", err)
	}

	// 3.5. JWS 디지털 서명 생성
	// crypto.Sign 함수는 PRIVATE_KEY_PATH 환경 변수에서 개인키를 로드하여 payloadBytes를 서명합니다.
	// 따라서 인자를 payloadBytes만 넘깁니다. (첫 번째 오류 수정)
	jwsSignatureBytes, err := crypto.Sign(payloadBytes)
	if err != nil {
		return "", fmt.Errorf("JWS 서명 생성 실패: %w", err)
	}

	// 3.6. Proof 객체 생성 및 VC에 추가 (Proof 구조체 사용)
	vc.Proof = &Proof{ // 포인터 타입이므로 & 사용
		Type:               "Ed25519Signature2018",                // 사용하는 서명 방식 (Ed25519 기반 가정)
		Created:            time.Now().UTC().Format(time.RFC3339), // 서명 생성 시간
		VerificationMethod: fmt.Sprintf("%s#keys-1", issuerDID),   // 서명 검증에 사용할 공개키 식별자 (예시)
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
