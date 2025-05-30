// cmd/vc_test/main.go
package main

import (
	"fmt"
	"log"
	"os"

	// "realestate/did" // 이 예제에서는 특정 DID 문자열을 직접 사용합니다.
	"realestate/vc"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("환경 변수 파일을 찾을 수 없습니다 (.env). 수동으로 설정된 환경 변수를 사용합니다.")
	}

	// --- 시나리오 1: 정상 공인중개사 ---
	fmt.Println("--- 시나리오 1: 정상 공인중개사 VC 생성 중 ---")
	normalAgentDID := "did:realestate:demobroker_normal_001" // 데모용 정상 공인중개사 DID
	normalName := "홍길동 (정상)"
	normalLicenseNum := "서울-12345-정상"
	normalPhone := "010-1111-1111"
	normalBrokerStatus := "valid" // 또는 "licenseStatus:Active,fraudRecord:None" 등을 표현할 수 있는 복합 상태 문자열 또는 구조체 사용 고려 가능

	fmt.Printf("VC 발급 정보: DID=%s, Name=%s, License=%s, Phone=%s, Status=%s\n",
		normalAgentDID, normalName, normalLicenseNum, normalPhone, normalBrokerStatus)

	normalVcJSON, err := vc.GenerateAndSignVC(normalAgentDID, normalName, normalLicenseNum, normalPhone, normalBrokerStatus)
	if err != nil {
		log.Fatalf("정상 공인중개사 VC 발급 실패: %v", err)
	}

	fmt.Println("발급된 VC (정상 공인중개사):")
	fmt.Println(normalVcJSON)

	normalFileName := "normal_broker_vc.json"
	err = os.WriteFile(normalFileName, []byte(normalVcJSON), 0644)
	if err != nil {
		log.Fatalf("VC 파일 저장 실패 ('%s'): %v", normalFileName, err)
	}
	fmt.Printf("✅ 정상 공인중개사 VC가 '%s' 파일로 저장되었습니다.\n\n", normalFileName)

	// --- 시나리오 2: 전과 기록 보유 공인중개사 ---
	fmt.Println("--- 시나리오 2: 전과 기록 보유 공인중개사 VC 생성 중 ---")
	fraudAgentDID := "did:realestate:demobroker_fraud_002" // 데모용 전과 기록 보유 공인중개사 DID
	fraudName := "조마피아 (전과)"
	fraudLicenseNum := "부산-67890-전과"
	fraudPhone := "010-2222-2222"
	// 이 "status" 값을 어떻게 해석하여 매물 등록을 차단할지는 4단계의 VC 검증 로직에서 결정됩니다.
	// 예를 들어, "status" 필드 값에 "fraudRecord:Exists" 와 같은 정보를 포함시키거나,
	// 단순히 "invalid_fraud_record" 와 같은 특정 문자열을 사용할 수 있습니다.
	// GenerateAndSignVC 함수가 현재는 status 문자열 하나만 받으므로, 문자열로 구분 가능한 값을 사용합니다.
	fraudBrokerStatus := "fraudRecord_Exists" // 또는 "invalid_due_to_fraud_record" 등

	fmt.Printf("VC 발급 정보: DID=%s, Name=%s, License=%s, Phone=%s, Status=%s\n",
		fraudAgentDID, fraudName, fraudLicenseNum, fraudPhone, fraudBrokerStatus)

	fraudVcJSON, err := vc.GenerateAndSignVC(fraudAgentDID, fraudName, fraudLicenseNum, fraudPhone, fraudBrokerStatus)
	if err != nil {
		log.Fatalf("전과 기록 보유 공인중개사 VC 발급 실패: %v", err)
	}

	fmt.Println("발급된 VC (전과 기록 보유 공인중개사):")
	fmt.Println(fraudVcJSON)

	fraudFileName := "fraud_record_broker_vc.json"
	err = os.WriteFile(fraudFileName, []byte(fraudVcJSON), 0644)
	if err != nil {
		log.Fatalf("VC 파일 저장 실패 ('%s'): %v", fraudFileName, err)
	}
	fmt.Printf("✅ 전과 기록 보유 공인중개사 VC가 '%s' 파일로 저장되었습니다.\n", fraudFileName)
}
