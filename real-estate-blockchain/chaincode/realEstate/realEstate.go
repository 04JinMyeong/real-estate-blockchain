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
	CreatedBy    string       `json:"createdBy"`
	PriceHistory []PriceEntry `json:"priceHistory"`
	OwnerHistory []OwnerEntry `json:"ownerHistory"`
	ReservedBy   string       `json:"reservedBy"`
	ReservedAt   string       `json:"reservedAt"`
	ExpiresAt    int64        `json:"expiresAt"`
	TxID         string       `json:"txID"`
	Timestamp    string       `json:"timestamp"`
}

type PriceEntry struct {
	Price string `json:"price"`
	Date  string `json:"date"`
}

type OwnerEntry struct {
	Owner string `json:"owner"`
	Date  string `json:"date"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) AddListing(
	ctx contractapi.TransactionContextInterface,
	id, address, owner, price, createdBy string,
) error {
	listingJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("데이터 검색 중 오류: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	if listingJSON != nil {
		var existing Listing
		if err := json.Unmarshal(listingJSON, &existing); err != nil {
			return fmt.Errorf("Unmarshal 실패: %v", err)
		}
		existing.PriceHistory = append(existing.PriceHistory, PriceEntry{Price: price, Date: currentDate})
		existing.OwnerHistory = append(existing.OwnerHistory, OwnerEntry{Owner: owner, Date: currentDate})

		updated, _ := json.Marshal(existing)
		return ctx.GetStub().PutState(id, updated)
	}

	newListing := Listing{
		ID:        id,
		Address:   address,
		CreatedBy: createdBy,
		PriceHistory: []PriceEntry{
			{Price: price, Date: currentDate},
		},
		OwnerHistory: []OwnerEntry{
			{Owner: owner, Date: currentDate},
		},
		ReservedBy: "",
		ReservedAt: "",
		ExpiresAt:  0,
		TxID:       "",
		Timestamp:  "",
	}

	bytes, _ := json.Marshal(newListing)
	return ctx.GetStub().PutState(id, bytes)
}

func (s *SmartContract) GetListing(ctx contractapi.TransactionContextInterface, id string) (*Listing, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil || data == nil {
		return nil, fmt.Errorf("매물 정보 없음: %s", id)
	}
	var listing Listing
	if err := json.Unmarshal(data, &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}

func (s *SmartContract) ReserveListing(
	ctx contractapi.TransactionContextInterface,
	id string,
	expiresAt int64,
) error {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("매물 조회 실패: %v", err)
	}
	if data == nil {
		return fmt.Errorf("매물이 존재하지 않습니다: %s", id)
	}

	var listing Listing
	if err := json.Unmarshal(data, &listing); err != nil {
		return fmt.Errorf("JSON 파싱 실패: %v", err)
	}

	// 이미 예약되어 있고 아직 만료 전인 경우
	if listing.ReservedBy != "" && time.Now().Unix() < listing.ExpiresAt {
		return fmt.Errorf("이미 예약된 매물입니다 (예약자: %s)", listing.ReservedBy)
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("클라이언트 ID 조회 실패: %v", err)
	}

	listing.ReservedBy = clientID
	listing.ReservedAt = time.Now().Format("2006-01-02 15:04:05")
	listing.ExpiresAt = expiresAt

	updated, err := json.Marshal(listing)
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %v", err)
	}

	return ctx.GetStub().PutState(id, updated)
}

func (s *SmartContract) ReleaseListing(
	ctx contractapi.TransactionContextInterface,
	id string,
) error {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("매물 조회 실패: %v", err)
	}
	if data == nil {
		return fmt.Errorf("매물이 존재하지 않습니다: %s", id)
	}

	var listing Listing
	if err := json.Unmarshal(data, &listing); err != nil {
		return fmt.Errorf("JSON 파싱 실패: %v", err)
	}

	// 예약된 상태가 아니면 그대로
	if listing.ReservedBy == "" {
		return nil
	}

	// 현재 시각이 ExpiresAt을 넘으면 해제
	if time.Now().Unix() > listing.ExpiresAt {
		listing.ReservedBy = ""
		listing.ReservedAt = ""
		listing.ExpiresAt = 0

		updated, _ := json.Marshal(listing)
		return ctx.GetStub().PutState(id, updated)
	}
	return nil
}

func (s *SmartContract) GetAllListings(ctx contractapi.TransactionContextInterface) ([]*Listing, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer it.Close()

	var result []*Listing
	for it.HasNext() {
		resp, _ := it.Next()
		var l Listing
		_ = json.Unmarshal(resp.Value, &l)
		// 빈 TxID, Timestamp 필드 초기화
		l.TxID = ""
		l.Timestamp = ""
		result = append(result, &l)
	}
	return result, nil
}

func (s *SmartContract) GetListingHistory(ctx contractapi.TransactionContextInterface, id string) ([]*Listing, error) {
	iter, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, fmt.Errorf("이력 조회 실패: %v", err)
	}
	defer iter.Close()

	var history []*Listing
	for iter.HasNext() {
		mod, err := iter.Next()
		if err != nil {
			return nil, err
		}

		var listing Listing
		if err := json.Unmarshal(mod.Value, &listing); err != nil {
			continue
		}

		listing.TxID = mod.TxId
		listing.Timestamp = mod.Timestamp.AsTime().Format("2006-01-02 15:04:05")
		history = append(history, &listing)
	}

	return history, nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("체인코드 생성 실패: %v", err))
	}
	if err := cc.Start(); err != nil {
		panic(fmt.Sprintf("체인코드 시작 실패: %v", err))
	}
}
