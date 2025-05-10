package main

import (
	"fmt"
	"log"
	"os"

	"realestate/did"
	"realestate/vc"
)

func main() {
    // 1) ISSUER_DID 환경변수가 반드시 설정되어 있어야 합니다.
    issuerDID := os.Getenv("ISSUER_DID")
    if issuerDID == "" {
        log.Fatal("ISSUER_DID 환경 변수가 설정되어 있지 않습니다")
    }

    // 2) 사용할 DID를 생성하거나(선택사항) 직접 지정할 수 있습니다.
    //    인자로 DID를 넘기면 그 값을 쓰고, 없으면 did.GenerateDID()로 새로 생성합니다.
    var agentDID string
    if len(os.Args) > 1 {
        agentDID = os.Args[1]
    } else {
        var err error
        agentDID, err = did.GenerateDID()
        if err != nil {
            log.Fatalf("DID 생성 실패: %v", err)
        }
    }

    // 3) 테스트용 공인중개사 정보
    name := "홍길동"
    licenseNum := "서울-123456"
    phone := "010-1234-5678"

    // 4) VC 발급 호출
    vcJSON, err := vc.GenerateAndSignVC(agentDID, name, licenseNum, phone)
    if err != nil {
        log.Fatalf("VC 발급 실패: %v", err)
    }

    // 5) 결과 출력
    fmt.Println("발급된 VC:")
    fmt.Println(vcJSON)
}
