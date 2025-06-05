package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"realestate/vc"
)

func main() {
	// 0. 환경 변수 확인 (ISSUER_DID, PRIVATE_KEY_PATH 필수)
	issuerDID := os.Getenv("ISSUER_DID")
	if issuerDID == "" {
		log.Fatal("환경 변수 ISSUER_DID가 설정되어 있지 않습니다. (예: did:realestate:platformIssuer001)")
	}
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		log.Fatal("환경 변수 PRIVATE_KEY_PATH가 설정되어 있지 않습니다. (예: keystore/issuer_private.key)")
	}
	// PUBLIC_KEY_PATH는 vc.GenerateAndSignVC에서는 직접 사용하지 않지만, VC 검증 시 필요하므로 설정되어 있는지 확인하는 것이 좋습니다.
	// if os.Getenv("PUBLIC_KEY_PATH") == "" {
	// 	log.Fatal("환경 변수 PUBLIC_KEY_PATH가 설정되어 있지 않습니다. (예: keystore/issuer_public.key)")
	// }

	// --- 시나리오 1: '전과 기록 없음' VC 생성 시작 ---
	fmt.Println("--- '전과 기록 없음' VC 생성 시작 ---")
	normalAgentDID := "did:realestate:3d7c5c838186c1ac13501ee94f76386f3525a2f3964d4fc4b00bfcd07b88930c" // 'TRUE' 사용자의 실제 DID
	normalAgentName := "김정상"                                                                            // CSV 파일 기준

	// additionalClaims 맵 생성 (여기서 brokerLicenseNumber, isLicensedBroker도 함께 전달)
	normalAgentClaims := map[string]interface{}{
		"licenseHolderName":           normalAgentName, // vc.GenerateAndSignVC의 name 인자와 중복되지만, additionalClaims를 통해 일괄 관리 가능
		"licenseNumber":               "110-2025-00001",
		"phone":                       "010-1111-1111",
		"isLicensedBroker":            true,   // 라이센스 유효 여부
		"fraudConvictionRecordStatus": "None", // 전과 없음
	}

	// vc.GenerateAndSignVC 호출 (agentID는 직접 전달, 나머지 클레임은 additionalClaims 맵으로 전달)
	vcNormalJSON, err := vc.GenerateAndSignVC(normalAgentDID, normalAgentName, normalAgentClaims) // name은 agentID로 사용
	if err != nil {
		log.Fatalf("'TRUE' 사용자 (김정상) VC 생성 실패: %v", err)
	}
	fmt.Printf("--- 'TRUE' 사용자 (김정상, DID: %s) VC ---\n", normalAgentDID)
	fmt.Println(vcNormalJSON)
	err = ioutil.WriteFile("TRUE_normal_agent_vc.json", []byte(vcNormalJSON), 0644)
	if err != nil {
		log.Printf("VC 저장 실패: %v", err)
	}
	fmt.Println("✅ 'TRUE' 사용자 VC를 TRUE_normal_agent_vc.json 파일로 저장했습니다.\n")

	// --- 시나리오 2: '전과 기록 있음' VC 생성 시작 ---
	fmt.Println("--- '전과 기록 있음' VC 생성 시작 ---")
	criminalAgentDID := "did:realestate:2141f50f9b48f987f920518789eb16d0e831de9087fcb9bcc193e5def9ad5d27" // 'FALSE' 사용자의 실제 DID
	criminalAgentName := "박사기"                                                                            // CSV 파일 기준

	// additionalClaims 맵 생성
	criminalAgentClaims := map[string]interface{}{
		"licenseHolderName":           criminalAgentName,
		"licenseNumber":               "110-2025-00002",
		"phone":                       "010-2222-2222",
		"isLicensedBroker":            true,
		"fraudConvictionRecordStatus": "Exists", // 전과 있음
	}

	// vc.GenerateAndSignVC 호출
	vcCriminalJSON, err := vc.GenerateAndSignVC(criminalAgentDID, criminalAgentName, criminalAgentClaims)
	if err != nil {
		log.Fatalf("'FALSE' 사용자 (박사기) VC 생성 실패: %v", err)
	}
	fmt.Printf("--- 'FALSE' 사용자 (박사기, DID: %s) VC ---\n", criminalAgentDID)
	fmt.Println(vcCriminalJSON)
	err = ioutil.WriteFile("FALSE_criminal_agent_vc.json", []byte(vcCriminalJSON), 0644)
	if err != nil {
		log.Printf("VC 저장 실패: %v", err)
	}
	fmt.Println("✅ 'FALSE' 사용자 VC를 FALSE_criminal_agent_vc.json 파일로 저장했습니다.\n")

	fmt.Println("🎉 모든 VC 생성이 완료되었습니다.")
	fmt.Println("ℹ️ 이 스크립트를 실행한 터미널에 ISSUER_DID와 PRIVATE_KEY_PATH 환경 변수가 올바르게 설정되어 있어야 합니다.")
}
