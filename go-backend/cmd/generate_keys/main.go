// File: go-backend/crypto/generate_keys.go
// 백엔드 서버의 키 쌍을 생성하여 파일로 저장하는 로직.

package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	// 1) ed25519 키 쌍 생성
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("키 생성 실패: %v", err)
	}

	// 2) Base64 인코딩
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	privB64 := base64.StdEncoding.EncodeToString(priv)

	// 3) 파일에 쓰기 (keystore 디렉토리에 issuer_ 접두어로 저장)
	ksDir := filepath.Join(".", "keystore")
	if err := ioutil.WriteFile(filepath.Join(ksDir, "issuer_public.key"), []byte(pubB64), 0600); err != nil {
		log.Fatalf("공개키 저장 실패: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(ksDir, "issuer_private.key"), []byte(privB64), 0600); err != nil {
		log.Fatalf("개인키 저장 실패: %v", err)
	}

	fmt.Println("키 쌍 생성 및 저장 완료:")
	fmt.Println(" - issuer_public.key")
	fmt.Println(" - issuer_private.key")
}
