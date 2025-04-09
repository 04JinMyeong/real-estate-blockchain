package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func EnrollAdmin(sdk *fabsdk.FabricSDK) error {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		return fmt.Errorf("MSP 클라이언트 생성 실패: %v", err)
	}

	adminID := "admin"
	adminSecret := "adminpw"

	_, err = mspClient.GetSigningIdentity(adminID)
	if err == nil {
		fmt.Println("⚠️ Admin 이미 등록되어 있음")
		return nil
	}

	err = mspClient.Enroll(adminID, msp.WithSecret(adminSecret))
	if err != nil {
		return fmt.Errorf("Admin 등록 실패: %v", err)
	}

	_, err = mspClient.GetSigningIdentity(adminID)
	if err != nil {
		return fmt.Errorf("Admin SigningIdentity 저장 실패: %v", err)
	}

	fmt.Println("🗂️ admin SigningIdentity 저장됨:", adminID)
	return nil
}

func RegisterAndEnrollUser(sdk *fabsdk.FabricSDK, mspClient *msp.Client, userID string) error {
	_, err := mspClient.GetSigningIdentity(userID)
	if err == nil {
		fmt.Println("⚠️ 사용자 이미 등록되어 있음")
		return nil
	}

	secret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        userID,
		Type:        "client",
		Affiliation: "",
	})
	if err != nil {
		return fmt.Errorf("사용자 등록 실패: %v", err)
	}
	fmt.Println("✅ 사용자 등록 성공:", secret)

	err = mspClient.Enroll(userID, msp.WithSecret(secret))
	if err != nil {
		return fmt.Errorf("사용자 인증서 발급 실패: %v", err)
	}

	fmt.Println("✅ 사용자 등록 완료")
	return nil
}
