package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Listing struct {
	ID           string       `json:"id"`
	Address      string       `json:"address"`
	PriceHistory []PriceEntry `json:"priceHistory"`
	OwnerHistory []OwnerEntry `json:"ownerHistory"`
}

type PriceEntry struct {
	Price string `json:"price"`
	Date  string `json:"date"`
}

type OwnerEntry struct {
	Owner string `json:"owner"`
	Date  string `json:"date"`
}

// InitLedger 함수 (초기화용)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// 매물 추가 또는 기존 매물에 이력 추가
func (s *SmartContract) AddListing(ctx contractapi.TransactionContextInterface, id string, address string, owner string, price string) error {
	listingJSON, err := ctx.GetStub().GetState(id)
	currentDate := time.Now().Format("2006-01-02")

	if err != nil {
		return fmt.Errorf("데이터 검색 중 오류: %v", err)
	}

	if listingJSON != nil {
		// 기존 매물일 경우 → 이력 추가
		var existing Listing
		_ = json.Unmarshal(listingJSON, &existing)

		existing.PriceHistory = append(existing.PriceHistory, PriceEntry{Price: price, Date: currentDate})
		existing.OwnerHistory = append(existing.OwnerHistory, OwnerEntry{Owner: owner, Date: currentDate})

		updatedJSON, _ := json.Marshal(existing)
		return ctx.GetStub().PutState(id, updatedJSON)
	}

	// 새로운 매물일 경우 → 새 블록
	newListing := Listing{
		ID:      id,
		Address: address,
		PriceHistory: []PriceEntry{
			{Price: price, Date: currentDate},
		},
		OwnerHistory: []OwnerEntry{
			{Owner: owner, Date: currentDate},
		},
	}

	listingBytes, _ := json.Marshal(newListing)
	return ctx.GetStub().PutState(id, listingBytes)
}

// 매물 상세 조회 (실제 구현)
func (s *SmartContract) GetListingById(ctx contractapi.TransactionContextInterface, id string) (*Listing, error) {
	listingJSON, err := ctx.GetStub().GetState(id)
	if err != nil || listingJSON == nil {
		return nil, fmt.Errorf("매물 정보가 존재하지 않습니다: %v", id)
	}

	var listing Listing
	_ = json.Unmarshal(listingJSON, &listing)
	return &listing, nil
}

// ❗ GetListing → GetListingById 별칭으로 연결
func (s *SmartContract) GetListing(ctx contractapi.TransactionContextInterface, id string) (*Listing, error) {
	return s.GetListingById(ctx, id)
}

// 전체 매물 목록 조회
func (s *SmartContract) GetAllListings(ctx contractapi.TransactionContextInterface) ([]*Listing, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var listings []*Listing
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		var listing Listing
		_ = json.Unmarshal(queryResponse.Value, &listing)
		listings = append(listings, &listing)
	}

	return listings, nil
}

// 매물 수정
func (s *SmartContract) UpdateListing(ctx contractapi.TransactionContextInterface, id string, newOwner string, newPrice string) error {
	listingJSON, err := ctx.GetStub().GetState(id)
	if err != nil || listingJSON == nil {
		return fmt.Errorf("매물 %s 이 존재하지 않습니다", id)
	}

	var listing Listing
	err = json.Unmarshal(listingJSON, &listing)
	if err != nil {
		return fmt.Errorf("데이터 언마샬 실패: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")

	if newPrice != "" {
		listing.PriceHistory = append(listing.PriceHistory, PriceEntry{
			Price: newPrice,
			Date:  currentDate,
		})
	}
	if newOwner != "" {
		listing.OwnerHistory = append(listing.OwnerHistory, OwnerEntry{
			Owner: newOwner,
			Date:  currentDate,
		})
	}

	updatedJSON, err := json.Marshal(listing)
	if err != nil {
		return fmt.Errorf("JSON 변환 실패: %v", err)
	}

	return ctx.GetStub().PutState(id, updatedJSON)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("체인코드 생성 실패: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("체인코드 시작 실패: %v", err))
	}
}
