// 📄 go-backend/vc/verify.go
package vc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings" // ◀◀◀ strings 패키지 import 추가

	"realestate/crypto"
)

// VerifiableCredential 구조체가 이 파일 또는 vc 패키지 내 다른 파일(예: issuer.go)에 정의되어 있어야 합니다.
// type VerifiableCredential struct { ... }
// type Proof struct { ... }
// 만약 issuer.go에만 있다면, 여기서도 접근 가능합니다 (같은 vc 패키지이므로).

// VerifyVC는 VC(JSON 문자열)를 받아서 서명(Signature)이 유효한지, 그리고 특정 클레임을 검사합니다.
func VerifyVC(vcJSON string) (bool, error) {
	// 1) VC 전체 JSON을 struct로 언마샬
	var vcData VerifiableCredential // vc 패키지 내에 VerifiableCredential 구조체가 정의되어 있다고 가정
	if err := json.Unmarshal([]byte(vcJSON), &vcData); err != nil {
		return false, fmt.Errorf("VC JSON 파싱 실패: %w", err)
	}

	// 2) Proof 필드가 존재하는지 확인
	if vcData.Proof == nil {
		return false, errors.New("proof 필드가 없습니다")
	}
	proof := vcData.Proof

	// 3) verificationMethod가 맞는지 확인 (ex: issuerDID#key-1)
	issuerDID := os.Getenv("ISSUER_DID")
	expectedVM := issuerDID + "#key-1" // 또는 #keys-1 등 vc/issuer.go 와 일치하는 식별자
	if proof.VerificationMethod != expectedVM {
		return false, fmt.Errorf("verificationMethod 불일치: VC(%s) vs 기대값(%s)", proof.VerificationMethod, expectedVM)
	}

	// 4) payloadBytes 재생산: Proof 제외 상태의 VC 전체 JSON
	vcCopy := vcData
	vcCopy.Proof = nil // proof를 삭제하여, 발급 시점과 동일한 상태로 만들기
	payloadBytes, err := json.Marshal(vcCopy)
	if err != nil {
		return false, fmt.Errorf("payload 재생성 실패: %w", err)
	}

	// 디버그 로그
	fmt.Println("===== [DEBUG] 검증 시 payloadBytes(JSON) =====")
	fmt.Println(string(payloadBytes))
	fmt.Println("============================================")

	// 5) Base64로 저장된 JWS(서명)를 디코딩하여 바이트 시그니처 얻기
	decodedSig, err := base64.StdEncoding.DecodeString(proof.Jws) // 변수명을 decodedSig로 변경 (sig는 crypto 패키지에서 사용 가능성)
	if err != nil {
		return false, fmt.Errorf("jws Base64 디코드 실패: %w", err)
	}

	// 6) ed25519.Verify를 통해 payloadBytes와 decodedSig 검증
	// crypto.Verify는 PUBLIC_KEY_PATH 환경 변수를 내부적으로 사용
	signatureValid, err := crypto.Verify(payloadBytes, decodedSig)
	if err != nil {
		// crypto.Verify 내부에서 공개키 로드 실패 등의 오류 발생 가능성
		return false, fmt.Errorf("서명 검증 중 오류 발생: %w", err)
	}
	if !signatureValid {
		return false, errors.New("서명 불일치: VC가 위조 또는 변조됨")
	}
	fmt.Println("✅ [VerifyVC] 서명 검증 통과") // 성공 로그 추가

	// 🔽🔽🔽 7) fraudConvictionRecordStatus 클레임 확인 (핵심 추가/수정 로직) 🔽🔽🔽
	// CredentialSubject가 interface{}이므로, 실제 타입으로 변환해야 합니다.
	// vc/issuer.go 에서는 map[string]interface{}로 생성했으므로 동일하게 가정합니다.
	credSubject, ok := vcData.CredentialSubject.(map[string]interface{})
	if !ok {
		// CredentialSubject가 예상한 map 형태가 아닐 경우 처리
		return false, errors.New("VC의 CredentialSubject 형식이 올바르지 않습니다 (map[string]interface{} 기대)")
	}

	fraudStatusInterface, ok := credSubject["fraudConvictionRecordStatus"]
	if !ok {
		// 'fraudConvictionRecordStatus' 클레임 자체가 없는 경우
		// 데모 시나리오에서는 이 클레임이 항상 존재한다고 가정하고, 없다면 오류로 처리하거나,
		// 또는 클레임이 없으면 "정상"으로 간주할 수도 있습니다. 여기서는 오류로 간주하지 않고 통과시킵니다 (선택적).
		// 만약 필수로 검사해야 한다면 아래 주석 해제:
		// return false, errors.New("'fraudConvictionRecordStatus' 클레임이 VC에 존재하지 않습니다.")
		fmt.Println("ℹ️ [VerifyVC] 'fraudConvictionRecordStatus' 클레임이 VC에 존재하지 않습니다. (검증 통과로 간주)")
	} else {
		fraudStatus, ok := fraudStatusInterface.(string)
		if !ok {
			// 클레임은 있지만 문자열 타입이 아닌 경우
			return false, errors.New("'fraudConvictionRecordStatus' 클레임의 값이 문자열이 아닙니다.")
		}

		// 대소문자 구분 없이 "Exists"와 비교 (더 안전한 비교)
		if strings.EqualFold(fraudStatus, "Exists") {
			fmt.Println("❌ [VerifyVC] 전과 기록 확인됨! fraudConvictionRecordStatus:", fraudStatus)
			return false, errors.New("전과 기록이 확인되어 로그인이 제한됩니다.") // 명확한 오류 메시지 반환
		}
		fmt.Println("✅ [VerifyVC] 전과 기록 없음. fraudConvictionRecordStatus:", fraudStatus) // 정상 로그 추가
	}

	// 모든 검증 통과
	return true, nil
}
