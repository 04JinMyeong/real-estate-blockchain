package blockchain

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func ReserveListing(user, id string, expiresAt int64) error {
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		return fmt.Errorf("wallet 로드 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return fmt.Errorf("user '%s' 가 wallet에 없음", user)
	}

	ccpPath := "./connection-org1.yaml"
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return fmt.Errorf("네트워크 접근 실패: %v", err)
	}

	contract := network.GetContract("realestate_v2")

	// ✅ 파라미터 순서: id, expiresAt, reservedBy
	_, err = contract.SubmitTransaction("ReserveListing", id, strconv.FormatInt(expiresAt, 10), user)
	if err != nil {
		return fmt.Errorf("체인코드 ReserveListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 ReserveListing 호출 성공")
	return nil
}
