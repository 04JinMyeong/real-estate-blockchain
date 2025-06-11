package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. 환경 변수 설정 (TLS 인증서 경로)
	os.Setenv("FABRIC_CA_CLIENT_TLS_CERTFILES", "./organizations/fabric-ca/org1/tls-cert.pem")

	// 2. Fabric SDK 초기화
	sdk, err := fabsdk.New(config.FromFile("./connection-org1.yaml"))
	if err != nil {
		log.Fatalf("❌ SDK 생성 실패: %v", err)
	}
	defer sdk.Close()

	// 3. MSP 클라이언트 생성
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		log.Fatalf("❌ MSP 클라이언트 생성 실패: %v", err)
	}

	// 4. Wallet 객체 생성
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		log.Fatalf("❌ Wallet 열기 실패: %v", err)
	}

	// 5. DB 열기
	db, err := sql.Open("sqlite3", "./realestate.db")
	if err != nil {
		log.Fatalf("❌ DB 연결 실패: %v", err)
	}
	defer db.Close()

	// 6. 사용자 목록 조회 (중개사만)
	rows, err := db.Query("SELECT id, password FROM users WHERE role = 'agent'")
	if err != nil {
		log.Fatalf("❌ DB 쿼리 실패: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, password string
		if err := rows.Scan(&id, &password); err != nil {
			log.Printf("❗ 사용자 정보 읽기 실패: %v", err)
			continue
		}

		// Wallet에 이미 있는지 확인
		exists := wallet.Exists(id)
		if err != nil {
			log.Printf("❗ Wallet 확인 실패 (%s): %v", id, err)
			continue
		}
		if exists {
			log.Printf("✅ 이미 Wallet에 등록됨: %s", id)
			continue
		}

		// Enroll 시도
		err = mspClient.Enroll(id, msp.WithSecret(password))
		if err != nil {
			log.Printf("❌ Enroll 실패 (%s): %v", id, err)
			continue
		}

		signingID, err := mspClient.GetSigningIdentity(id)
		if err != nil {
			log.Printf("❌ SigningIdentity 불러오기 실패 (%s): %v", id, err)
			continue
		}

		cert := signingID.EnrollmentCertificate()
		key, err := signingID.PrivateKey().Bytes()
		if err != nil {
			log.Printf("❌ 개인키 가져오기 실패 (%s): %v", id, err)
			continue
		}

		identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
		err = wallet.Put(id, identity)
		if err != nil {
			log.Printf("❌ Wallet 저장 실패 (%s): %v", id, err)
		} else {
			log.Printf("✅ Wallet 등록 완료: %s", id)
		}
	}
}
