package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"realestate/vc"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("사용법: go run cmd/verify_test/main.go <vc.json 파일 경로>")
    }
    path := os.Args[1]

    data, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("VC 파일 읽기 실패: %v", err)
    }

    // 환경변수 ISSUER_DID, PUBLIC_KEY_PATH, PRIVATE_KEY_PATH 반드시 설정
    valid, err := vc.VerifyVC(string(data))
    if err != nil {
        log.Fatalf("VC 검증 오류: %v", err)
    }
    if valid {
        fmt.Println("✔ VC 검증 성공: 이 자격 증명서는 유효합니다.")
    } else {
        fmt.Println("✘ VC 검증 실패: 자격 증명서가 위조되었거나 변조되었습니다.")
    }
}
