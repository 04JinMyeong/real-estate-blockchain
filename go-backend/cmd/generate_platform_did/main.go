package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil" // ioutil로 변경 (os.Open, bufio.Scanner 대신 간단하게 파일 전체 읽기)
	"log"
	"path/filepath"
	"strings" // TrimSpace 사용을 위해 추가

	// "realestate/did" 패키지의 정확한 import 경로를 사용해야 합니다.
	// 현재 프로젝트 구조를 기반으로 한 상대 경로나 모듈 경로를 사용하세요.
	// 예를 들어, go.mod 파일에 모듈 이름이 "realestate"로 되어 있다면 그대로 사용 가능합니다.
	"realestate/did"
)

func main() {
	// 1. issuer_public.key 파일 경로 정의
	// 현재 cmd/generate_platform_did/main.go 에서 실행되므로, 상대 경로는 ../../keystore 가 됩니다.
	// 또는, go run ./cmd/generate_platform_did/main.go 와 같이 go-backend 디렉토리에서 실행한다면 keystore/issuer_public.key 가 됩니다.
	// 여기서는 go-backend 디렉토리에서 실행하는 것을 기준으로 합니다.
	keyPath := filepath.Join("keystore", "issuer_public.key")

	// 2. Base64 인코딩된 공개키 파일 읽기
	pubKeyBase64Bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("공개키 파일 '%s'을(를) 여는 데 실패했습니다: %v", keyPath, err)
	}

	// 파일 끝의 개행 문자 등을 제거
	pubKeyBase64 := strings.TrimSpace(string(pubKeyBase64Bytes))

	if pubKeyBase64 == "" {
		log.Fatalf("공개키 파일 '%s'이(가) 비어 있거나 읽을 수 없습니다.", keyPath)
	}

	// 3. Base64 공개키 디코딩
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		log.Fatalf("Base64 공개키 디코딩에 실패했습니다: %v", err)
	}

	// (선택 사항) 공개키 길이 검증 - ed25519.PublicKeySize는 crypto/ed25519 패키지에 정의되어 있습니다.
	// import "crypto/ed25519" 추가 필요
	/*
		if len(pubKeyBytes) != ed25519.PublicKeySize {
			log.Fatalf("디코딩된 공개키의 크기가 올바르지 않습니다. 예상 크기: %d, 실제 크기: %d", ed25519.PublicKeySize, len(pubKeyBytes))
		}
	*/

	// 4. 공개키 바이트로부터 DID 생성
	platformDID := did.GenerateDIDFromPublicKey(pubKeyBytes) // did 패키지의 함수 사용

	fmt.Println("플랫폼 DID가 성공적으로 생성되었습니다.")
	fmt.Println("플랫폼 공개키 파일:", keyPath)
	fmt.Println("생성된 플랫폼 DID:", platformDID)
}
