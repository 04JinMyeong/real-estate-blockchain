package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// 주소를 입력받아 고정된 해시값(ID)을 반환
func GeneratePropertyID(address string) string {
	hash := sha256.Sum256([]byte(address))
	return hex.EncodeToString(hash[:])
}
