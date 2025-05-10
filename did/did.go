package did

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"

    "realestate/crypto"
)

// GenerateDID 로딩된 공개키로부터 DID를 생성합니다.
// 환경 변수 PUBLIC_KEY_PATH 에서 공개키 파일을 읽어 옵니다.
func GenerateDID() (string, error) {
    pubKey, err := crypto.LoadPublicKey("PUBLIC_KEY_PATH")
    if err != nil {
        return "", fmt.Errorf("공개키 로드 실패: %w", err)
    }
    return GenerateDIDFromPublicKey(pubKey), nil
}

// GenerateDIDFromPublicKey 주어진 공개키 바이트를 해시하여 DID 문자열을 반환합니다.
func GenerateDIDFromPublicKey(pubKey []byte) string {
    hash := sha256.Sum256(pubKey)
    // hex 인코딩된 해시값을 DID suffix 로 사용
    return fmt.Sprintf("did:realestate:%s", hex.EncodeToString(hash[:]))
}
