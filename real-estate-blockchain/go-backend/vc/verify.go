package vc

import (
	"crypto/ecdsa"    // ◀◀◀ ECDSA 암호화 관련 패키지
	"crypto/elliptic" // ◀◀◀ 타원 곡선 암호화 관련 패키지
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil" // ◀◀◀ HTTP 응답을 읽기 위한 패키지
	"math/big"  // ◀◀◀ 큰 정수 연산을 위한 패키지
	"net/http"  // ◀◀◀ HTTP 요청을 위한 패키지
)

// VC 구조체와 Proof 구조체는 vc 패키지 내 다른 파일(예: issuer.go)에 이미 정의되어 있어야 합니다.
// 만약 정의되어 있지 않다면, 아래 주석을 해제하여 추가해주세요.
/*
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
*/

// getIssuerPublicKey 함수는 목업 서버로부터 발급자의 공개키를 가져옵니다.
func getIssuerPublicKey(issuerDID string) (*ecdsa.PublicKey, error) {
	// 실제 서비스에서는 issuer의 DID Document를 리졸브해서 공개키를 찾아야 합니다.
	// 데모에서는 간단히 목업 서버의 API를 호출해서 공개키를 가져옵니다.
	resp, err := http.Get("http://localhost:8083/public-key")
	if err != nil {
		return nil, fmt.Errorf("목업 API에서 공개키를 가져오는 데 실패: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res struct {
		PublicKey string `json:"publicKey"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("공개키 API 응답 파싱 실패: %w", err)
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(res.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("공개키 디코딩 실패: %w", err)
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), pubKeyBytes)
	if x == nil {
		return nil, errors.New("잘못된 공개키 형식")
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// VerifyVC 함수는 VC의 서명과 내용을 검증합니다.
func VerifyVC(vcJSON string) (bool, error) {
	fmt.Println("===== [VerifyVC] 새로운 검증 로직 시작 =====")
	var vcData VerifiableCredential
	if err := json.Unmarshal([]byte(vcJSON), &vcData); err != nil {
		return false, fmt.Errorf("VC JSON 파싱 실패: %w", err)
	}

	// 1. Proof 필드 및 JWS 확인
	if vcData.Proof == nil || vcData.Proof.Jws == "" {
		fmt.Println("❌ [VerifyVC] Proof 필드 또는 JWS 서명이 없음")
		return false, errors.New("VC에 유효한 Proof(서명)가 없습니다")
	}
	fmt.Println("✅ [VerifyVC] Proof 필드 확인됨. Type:", vcData.Proof.Type)

	// 2. 서명(JWS) 값을 디코딩
	signatureBytes, err := base64.RawURLEncoding.DecodeString(vcData.Proof.Jws)
	if err != nil {
		fmt.Println("❌ [VerifyVC] JWS 서명 디코딩 실패:", err)
		return false, fmt.Errorf("VC 서명(JWS) 디코딩 실패: %w", err)
	}
	// P-256 서명은 64바이트 (r, s 각각 32바이트)
	if len(signatureBytes) != 64 {
		return false, errors.New("잘못된 ECDSA 서명 길이")
	}
	// r, s 값 분리
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// 3. 검증을 위해 원본 payload 재생성 (Proof를 제외한 나머지 부분)
	originalVCData := vcData
	originalVCData.Proof = nil
	payloadBytes, err := json.Marshal(originalVCData)
	if err != nil {
		return false, fmt.Errorf("검증용 VC 데이터 생성 실패: %w", err)
	}
	// 원본 데이터의 해시 계산
	hash := sha256.Sum256(payloadBytes)

	// 4. 발급자의 공개키 가져오기
	issuerPubKey, err := getIssuerPublicKey(vcData.Issuer)
	if err != nil {
		return false, err
	}

	// 5. 공개키로 서명 검증
	isValid := ecdsa.Verify(issuerPubKey, hash[:], r, s)
	if !isValid {
		fmt.Println("❌ [VerifyVC] ECDSA 서명 검증 실패")
		return false, errors.New("VC 서명이 유효하지 않습니다 (위조 또는 변조됨)")
	}
	fmt.Println("✅ [VerifyVC] 서명 검증 성공!")

	// 6. VC 내용(Claim) 검증
	credSubject, ok := vcData.CredentialSubject.(map[string]interface{})
	if !ok {
		return false, errors.New("CredentialSubject 형식 오류")
	}
	hasCriminalRecord, ok := credSubject["criminal_record"].(bool)
	if !ok {
		return false, errors.New("'criminal_record' 클레임을 찾을 수 없거나 형식이 다릅니다")
	}
	if hasCriminalRecord {
		return false, errors.New("전과 기록이 확인되어 로그인이 제한됩니다")
	}
	fmt.Println("✅ [VerifyVC] 전과 기록 없음 확인.")

	fmt.Println("✅ [VerifyVC] 모든 검증 통과! VC 유효함.")
	return true, nil
}
