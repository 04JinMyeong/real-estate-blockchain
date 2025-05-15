package blockchain

import (
	"fmt"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func SubmitUpdateListing(user, id, newOwner, newPrice string) error {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return fmt.Errorf("wallet 열기 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return fmt.Errorf("사용자 '%s'가 wallet에 없습니다", user)
	}

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

	contract := network.GetContract("realEstate")

	_, err = contract.SubmitTransaction("UpdateListing", id, newOwner, newPrice)
	if err != nil {
		return fmt.Errorf("체인코드 UpdateListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 UpdateListing 호출 성공")
	return nil
}
