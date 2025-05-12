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
    // 1. JSON unmarshal
    var vc map[string]interface{}
    if err := json.Unmarshal([]byte(vcJSON), &vc); err != nil {
        return false, fmt.Errorf("VC JSON 파싱 실패: %w", err)
    }

    // 2. proof 필드 확인
    proofI, ok := vc["proof"]
    if !ok {
        return false, errors.New("proof 필드가 없습니다")
    }
    proof, ok := proofI.(map[string]interface{})
    if !ok {
        return false, errors.New("proof 형식이 잘못되었습니다")
    }

    // 3. verificationMethod 검증 (ISSUER_DID#key-1 과 일치)
    vmI, ok := proof["verificationMethod"]
    if !ok {
        return false, errors.New("verificationMethod가 없습니다")
    }
    vm, ok := vmI.(string)
    if !ok {
        return false, errors.New("verificationMethod 타입 오류")
    }
    issuerDID := os.Getenv("ISSUER_DID")
    expectedVM := issuerDID + "#key-1"
    if !strings.EqualFold(vm, expectedVM) {
        return false, fmt.Errorf("verificationMethod 불일치: %s vs %s", vm, expectedVM)
    }

    // 4. credentialSubject 부분만 다시 Marshal
    csI, ok := vc["credentialSubject"]
    if !ok {
        return false, errors.New("credentialSubject가 없습니다")
    }
    csBytes, err := json.Marshal(csI)
    if err != nil {
        return false, fmt.Errorf("credentialSubject 재Marshal 실패: %w", err)
    }

    // 5. jws(Base64 서명) 디코딩
    jwsI, ok := proof["jws"]
    if !ok {
        return false, errors.New("jws가 없습니다")
    }
    jwsStr, ok := jwsI.(string)
    if !ok {
        return false, errors.New("jws 타입 오류")
    }
    sig, err := base64.StdEncoding.DecodeString(jwsStr)
    if err != nil {
        return false, fmt.Errorf("jws Base64 디코드 실패: %w", err)
    }

    // 6. 서명 검증
    valid, err := crypto.Verify(csBytes, sig)
    if err != nil {
        return false, fmt.Errorf("서명 검증 오류: %w", err)
    }
    if !valid {
        return false, errors.New("서명 불일치: VC가 위조 또는 변조됨")
    }

    return true, nil
}
