package vc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"realestate/crypto"
)

// VerifyVC 검증 결과와 이유(에러 메시지)를 반환합니다.
// vcJSON: 검증할 VC 전문(JSON 문자열)
func VerifyVC(vcJSON string) (bool, error) {
	fmt.Println("===== [VerifyVC] 함수 시작 =====")
	// 입력된 VC JSON 문자열의 일부를 로그로 남깁니다. (너무 길 경우를 대비)
	// func min(a, b int) int { if a < b { return a; }; return b; } 와 같은 헬퍼 함수가 필요합니다.
	// 이 함수는 파일 하단이나 별도 유틸리티 파일에 추가해주세요.
	// fmt.Println("받은 VC JSON (일부):", string([]rune(vcJSON)[:min(100, len(vcJSON))]))

	// 1. VC 전체 JSON을 VerifiableCredential 구조체로 언마샬
	var vcData VerifiableCredential // issuer.go 에 정의된 구조체 사용
	if err := json.Unmarshal([]byte(vcJSON), &vcData); err != nil {
		fmt.Println("❌ [VerifyVC] VC JSON 파싱 실패:", err)
		return false, fmt.Errorf("VC JSON 파싱 실패: %w", err)
	}
	// CredentialSubject가 map[string]interface{}로 파싱되었는지 확인 후 주요 정보 로깅
	var subjectID string
	if csMap, ok := vcData.CredentialSubject.(map[string]interface{}); ok {
		if id, idOk := csMap["id"].(string); idOk {
			subjectID = id
		}
	}
	fmt.Println("✅ [VerifyVC] VC JSON 파싱 성공. Issuer:", vcData.Issuer, ", Subject ID (from VC):", subjectID)

	// 2. Proof 필드가 존재하는지 확인
	if vcData.Proof == nil {
		fmt.Println("❌ [VerifyVC] Proof 필드 없음")
		return false, errors.New("proof 필드가 없습니다")
	}
	proof := vcData.Proof // 이제 proof는 *vc.Proof 타입입니다.
	fmt.Println("✅ [VerifyVC] Proof 필드 확인됨. Type:", proof.Type)

	// 3. verificationMethod가 맞는지 확인 (ex: issuerDID#key-1)
	issuerDID := os.Getenv("ISSUER_DID")
	if issuerDID == "" {
		fmt.Println("❌ [VerifyVC] 환경 변수 ISSUER_DID 설정 안됨")
		return false, errors.New("서버 설정 오류: 발급자 DID가 지정되지 않았습니다")
	}
	// issuer.go에서 #key-1 로 생성했으므로 여기서도 동일하게 사용
	expectedVM := issuerDID + "#key-1"
	if proof.VerificationMethod != expectedVM {
		fmt.Printf("❌ [VerifyVC] verificationMethod 불일치: VC내 값(%s) vs 기대값(%s)\n", proof.VerificationMethod, expectedVM)
		return false, fmt.Errorf("verificationMethod 불일치: VC내 값(%s) vs 기대값(%s)", proof.VerificationMethod, expectedVM)
	}
	fmt.Println("✅ [VerifyVC] verificationMethod 일치 확인됨:", proof.VerificationMethod)

	// 4. payloadBytes 재생산: Proof 제외 상태의 VC 전체 JSON
	//    vcData는 이미 VerifiableCredential 구조체이므로, Proof 필드만 nil로 설정하고 마샬링합니다.
	vcForSigning := vcData   // 구조체 복사
	vcForSigning.Proof = nil // 서명 대상에서 Proof 제외
	payloadBytes, err := json.Marshal(vcForSigning)
	if err != nil {
		fmt.Println("❌ [VerifyVC] payload 재생성 실패:", err)
		return false, fmt.Errorf("VC 서명 대상 페이로드 재생성 실패: %w", err)
	}

	fmt.Println("===== [VerifyVC DEBUG] 검증 시 payloadBytes(JSON) =====")
	fmt.Println(string(payloadBytes))
	fmt.Println("============================================")

	// 5. Base64로 저장된 JWS(서명)를 디코딩하여 바이트 시그니처 얻기
	decodedSig, err := base64.StdEncoding.DecodeString(proof.Jws)
	if err != nil {
		fmt.Println("❌ [VerifyVC] JWS 서명 디코딩 실패:", err)
		return false, fmt.Errorf("JWS 서명 디코딩 실패: %w", err)
	}
	fmt.Println("✅ [VerifyVC] JWS 서명 디코딩 성공")

	// 6. 서명 검증 (crypto.Verify 호출)
	// crypto.Verify는 PUBLIC_KEY_PATH 환경 변수를 내부적으로 사용
	signatureValid, err := crypto.Verify(payloadBytes, decodedSig)
	if err != nil {
		fmt.Println("❌ [VerifyVC] 서명 검증 중 crypto.Verify 오류:", err)
		return false, fmt.Errorf("서명 검증 중 오류 발생: %w", err)
	}
	if !signatureValid {
		fmt.Println("❌ [VerifyVC] 서명 불일치! VC가 위조 또는 변조됨")
		return false, errors.New("서명 불일치: VC가 위조 또는 변조됨")
	}
	fmt.Println("✅ [VerifyVC] 서명 검증 통과!")

	// 7. fraudConvictionRecordStatus 클레임 확인 (핵심 추가 로직)
	// vcData.CredentialSubject는 interface{} 타입이므로, 실제 타입으로 변환(타입 단언)해야 합니다.
	// vc/issuer.go에서는 map[string]interface{}로 생성했으므로 동일하게 가정합니다.
	credSubject, ok := vcData.CredentialSubject.(map[string]interface{})
	if !ok {
		fmt.Println("❌ [VerifyVC] CredentialSubject 형식이 map[string]interface{}가 아님")
		return false, errors.New("VC의 CredentialSubject 형식이 올바르지 않습니다 (map[string]interface{} 기대)")
	}

	fraudStatusInterface, claimExists := credSubject["fraudConvictionRecordStatus"]
	if !claimExists {
		// 'fraudConvictionRecordStatus' 클레임 자체가 없는 경우
		// 데모 시나리오에서는 이 클레임이 항상 존재하고 유효한 값을 가진다고 가정할 수 있습니다.
		// 만약 없다면, 오류로 처리하거나 "문제 없음"으로 간주할지 정책 결정 필요.
		// 여기서는 "없으면 문제 없음"으로 간주하고 로그만 남깁니다.
		fmt.Println("ℹ️ [VerifyVC] 'fraudConvictionRecordStatus' 클레임이 VC에 존재하지 않음 (정상으로 간주)")
	} else {
		fraudStatus, typeOk := fraudStatusInterface.(string)
		if !typeOk {
			fmt.Println("❌ [VerifyVC] 'fraudConvictionRecordStatus' 클레임 값이 문자열이 아님. 실제 타입:", fmt.Sprintf("%T", fraudStatusInterface))
			return false, errors.New("'fraudConvictionRecordStatus' 클레임의 값이 문자열이 아닙니다")
		}

		fmt.Println("ℹ️ [VerifyVC] 'fraudConvictionRecordStatus' 클레임 값:", fraudStatus)
		// 대소문자 구분 없이 "Exists"와 비교
		if strings.EqualFold(fraudStatus, "Exists") {
			fmt.Println("❌ [VerifyVC] 전과 기록 확인됨! 로그인 차단 필요.")
			return false, errors.New("전과 기록이 확인되어 로그인이 제한됩니다.") // 이 오류 메시지가 Login 핸들러로 전달됨
		}
		fmt.Println("✅ [VerifyVC] 전과 기록 없음.")
	}

	// 모든 검증 통과
	fmt.Println("✅ [VerifyVC] 모든 검증 통과. VC 유효함.")
	return true, nil
}
