package blockchain

import (
	"fmt"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// 매물 등록
func SubmitAddListing(user, id, address, owner, price, photoUrl string) error {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return fmt.Errorf("wallet 불러오기 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return fmt.Errorf("사용자 '%s'가 wallet에 존재하지 않습니다", user)
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

	contract := network.GetContract("realEstate") // ⬅️ 최신 체인코드 네임

	_, err = contract.SubmitTransaction("AddListing", id, address, owner, price, user, photoUrl)
	if err != nil {
		return fmt.Errorf("체인코드 AddListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 AddListing 호출 성공")
	return nil
}

// 개별 매물 조회
func QueryProperty(user, id string) (string, error) {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return "", fmt.Errorf("wallet 불러오기 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return "", fmt.Errorf("사용자 '%s'가 wallet에 존재하지 않습니다", user)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return "", fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return "", fmt.Errorf("네트워크 접근 실패: %v", err)
	}

	contract := network.GetContract("realEstate") // ⬅️ 최신 체인코드 네임

	result, err := contract.EvaluateTransaction("GetListing", id)
	if err != nil {
		return "", fmt.Errorf("체인코드 GetListing 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 GetListing 조회 성공")
	return string(result), nil
}

// 전체 매물 조회
func QueryAllProperties(user string) (string, error) {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return "", fmt.Errorf("wallet 불러오기 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return "", fmt.Errorf("사용자 '%s'가 wallet에 존재하지 않습니다", user)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return "", fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return "", fmt.Errorf("네트워크 접근 실패: %v", err)
	}

	contract := network.GetContract("realEstate") // ⬅️ 최신 체인코드 네임

	result, err := contract.EvaluateTransaction("GetAllListings")
	if err != nil {
		return "", fmt.Errorf("체인코드 GetAllListings 호출 실패: %v", err)
	}

	return string(result), nil
}

// 매물 이력 조회 (History)
func QueryPropertyHistory(user, id string) (string, error) {
	walletPath := "./wallet"
	ccpPath := "./connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return "", fmt.Errorf("wallet 불러오기 실패: %v", err)
	}

	if !wallet.Exists(user) {
		return "", fmt.Errorf("사용자 '%s'가 wallet에 존재하지 않습니다", user)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return "", fmt.Errorf("Gateway 연결 실패: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return "", fmt.Errorf("네트워크 접근 실패: %v", err)
	}

	contract := network.GetContract("realEstate") // ⬅️ 최신 체인코드 네임

	result, err := contract.EvaluateTransaction("GetListingHistory", id)
	if err != nil {
		return "", fmt.Errorf("체인코드 GetListingHistory 호출 실패: %v", err)
	}

	fmt.Println("✅ 체인코드 GetListingHistory 호출 성공")
	return string(result), nil
}
