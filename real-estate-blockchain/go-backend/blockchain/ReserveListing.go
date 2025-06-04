// blockchain/ReserveListing.go
package blockchain

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// ReserveListing submits a transaction to reserve a listing until `expiresAt` (Unix timestamp).
// Parameters:
//   - user: the identity in the wallet performing the reservation
//   - id: the listing ID
//   - expiresAt: Unix timestamp when the reservation expires
func ReserveListing(user, id string, expiresAt int64) error {
	// Load wallet and check identity
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		return fmt.Errorf("wallet 로드 실패: %v", err)
	}
	if !wallet.Exists(user) {
		return fmt.Errorf("user '%s'가 wallet에 없습니다", user)
	}

	// Load connection profile
	ccpPath := "./connection-org1.yaml"
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	// Get network and contract
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return fmt.Errorf("네트워크 접근 실패: %v", err)
	}
	contract := network.GetContract("realEstate")

	// Submit transaction: ReserveListing(id, expiresAt, user)
	_, err = contract.SubmitTransaction("ReserveListing", id, strconv.FormatInt(expiresAt, 10), user)
	if err != nil {
		return fmt.Errorf("체인코드 ReserveListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 ReserveListing 호출 성공")
	return nil
}

// ReleaseListing submits a transaction to clear reservation fields for a listing.
// Parameters:
//   - user: the identity in the wallet performing the release (typically 'admin')
//   - id: the listing ID to release
func ReleaseListing(user, id string) error {
	// Load wallet and check identity
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		return fmt.Errorf("wallet 불러오기 실패: %v", err)
	}
	if !wallet.Exists(user) {
		return fmt.Errorf("사용자 '%s'가 wallet에 없습니다", user)
	}

	// Load connection profile
	ccpPath := "./connection-org1.yaml"
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	// Get network and contract
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return fmt.Errorf("네트워크 접근 실패: %v", err)
	}
	contract := network.GetContract("realEstate")

	// Submit transaction: ReleaseListing(id)
	_, err = contract.SubmitTransaction("ReleaseListing", id)
	if err != nil {
		return fmt.Errorf("체인코드 ReleaseListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 ReleaseListing 호출 성공")
	return nil
}
