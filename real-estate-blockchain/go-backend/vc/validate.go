// File: go-backend/vc/validate.go

package vc

import (
	"encoding/base64" // base64 인코딩/디코딩을 위해 추가
	"encoding/json"   // JSON 마샬링/언마샬링을 위해 추가
	"errors"          // 에러 처리를 위해 추가
	"fmt"
	"os"   // 환경 변수 읽기를 위해 추가
	"time" // 시간 처리를 위해 추가

	"realestate/crypto" // crypto 패키지 임포트
)

// VerifiableCredential 및 Proof 구조체는 vc/issuer.go에 정의되어 있으므로,
// 여기서는 다시 정의하지 않습니다. 동일한 'vc' 패키지 내에서 자동으로 사용 가능합니다.

// ValidateVC VC의 유효성을 검증하는 함수입니다.
// 서명, 발급자, 유효기간, 그리고 특히 fraudConvictionRecordStatus 클레임을 검증합니다.
func ValidateVC(vcData VerifiableCredential) error { // vc/issuer.go에 정의된 VerifiableCredential 사용
	// 1. IssuanceDate 유효성 검사
	// VC가 너무 오래된 것은 아닌지 확인합니다. (옵션)
	issuanceTime, err := time.Parse(time.RFC3339, vcData.IssuanceDate)
	if err != nil {
		return fmt.Errorf("VC 발급일자 파싱 실패: %w", err)
	}
	// 예시: 발급된 지 1년 이상 된 VC는 유효하지 않다고 판단
	if time.Since(issuanceTime) > 24*365*time.Hour { // 1년
		return fmt.Errorf("VC 유효기간 만료: 발급일자가 너무 오래되었습니다 (발급일: %s)", vcData.IssuanceDate)
	}

	// 2. 발급자 (Issuer) DID 유효성 검사
	// VC를 발급한 주체(Issuer)가 우리가 기대하는 플랫폼(Issuer)인지 확인합니다.
	expectedIssuerDID := os.Getenv("ISSUER_DID")
	if expectedIssuerDID == "" {
		return fmt.Errorf("환경 변수 ISSUER_DID가 설정되지 않아 발급자 유효성 검증 불가")
	}
	if vcData.Issuer != expectedIssuerDID {
		return fmt.Errorf("VC 발급자 DID가 일치하지 않습니다: 기대값 %s, 실제값 %s", expectedIssuerDID, vcData.Issuer)
	}

	// 3. 서명 검증
	// VC의 무결성(위변조 여부)과 진위(정말 발급자가 서명했는지)를 확인합니다.
	// 서명 대상이 되는 페이로드(payload)는 Proof 필드를 제외한 VC 전체입니다.
	vcForVerification := vcData
	vcForVerification.Proof = nil // 서명 검증 시 Proof 필드는 제외

	payloadBytes, err := json.Marshal(vcForVerification) // 서명 대상 VC를 JSON 바이트로 변환
	if err != nil {
		return fmt.Errorf("VC 검증을 위한 페이로드 마샬링 실패: %w", err)
	}

	// Proof의 Jws 필드는 Base64 인코딩된 문자열이므로, 먼저 디코딩해야 합니다.
	if vcData.Proof == nil {
		return errors.New("VC에 Proof 필드가 누락되었습니다.")
	}
	jwsSignatureBytes, err := base64.StdEncoding.DecodeString(vcData.Proof.Jws)
	if err != nil {
		return fmt.Errorf("JWS 서명 디코딩 실패: %w", err)
	}

	// crypto.Verify 함수를 사용하여 서명 검증을 수행합니다.
	// 이 함수는 PUBLIC_KEY_PATH 환경 변수를 내부적으로 사용하여 공개키를 로드합니다.
	isValid, err := crypto.Verify(payloadBytes, jwsSignatureBytes)
	if err != nil {
		return fmt.Errorf("VC 서명 검증 중 오류 발생: %w", err)
	}
	if !isValid {
		return errors.New("VC 서명이 유효하지 않습니다.")
	}

	// 4. CredentialSubject 내부 클레임 유효성 검증 (핵심 로직)
	// 'credentialSubject'는 'interface{}' 타입이므로, 'map[string]interface{}'로 타입 단언해야 합니다.
	subjectMap, ok := vcData.CredentialSubject.(map[string]interface{})
	if !ok {
		return errors.New("VC CredentialSubject 형식이 올바르지 않습니다.")
	}

	// 'isLicensedBroker' 클레임 확인 (공인중개사 자격 여부)
	if isLicensed, exists := subjectMap["isLicensedBroker"].(bool); exists {
		if !isLicensed {
			return errors.New("로그인 제한: 해당 공인중개사는 유효한 라이센스를 가지고 있지 않습니다.")
		}
	} else {
		return errors.New("VC에 'isLicensedBroker' 클레임이 누락되었습니다.")
	}

	// 'fraudConvictionRecordStatus' 클레임 확인 (전과 기록 여부 - 가장 중요!)
	if status, exists := subjectMap["fraudConvictionRecordStatus"].(string); exists {
		if status == "Exists" {
			return errors.New("로그인 제한: 해당 공인중개사는 전과 기록이 확인되었습니다.")
		}
	} else {
		return errors.New("VC에 전과 기록 상태 클레임 ('fraudConvictionRecordStatus')이 누락되었습니다.")
	}

	return nil // 모든 검증 통과
}
