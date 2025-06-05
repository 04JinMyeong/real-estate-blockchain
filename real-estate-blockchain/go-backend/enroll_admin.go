package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func Enrollmain() {
	sdk, err := fabsdk.New(config.FromFile("./connection-org1.yaml"))
	if err != nil {
		log.Fatalf("❌ SDK 생성 실패: %v", err)
	}
	defer sdk.Close()

	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		log.Fatalf("❌ MSP 클라이언트 생성 실패: %v", err)
	}

	err = mspClient.Enroll("admin", msp.WithSecret("adminpw"))
	if err != nil {
		log.Fatalf("❌ Admin Enroll 실패: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		log.Fatalf("❌ Wallet 열기 실패: %v", err)
	}

	signingID, err := mspClient.GetSigningIdentity("admin")
	if err != nil {
		log.Fatalf("❌ SigningIdentity 불러오기 실패: %v", err)
	}

	cert := signingID.EnrollmentCertificate()
	key, err := signingID.PrivateKey().Bytes()
	if err != nil {
		log.Fatalf("❌ 개인키 가져오기 실패: %v", err)
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("admin", identity)
	if err != nil {
		log.Fatalf("❌ Wallet 저장 실패: %v", err)
	}

	fmt.Println("✅ Admin 등록 성공")
}
