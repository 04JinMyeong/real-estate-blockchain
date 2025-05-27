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

// Reserved 필드를 항상 포함시키기 위해 omitempty 제거
type Listing struct {
	ID           string       `json:"id"`
	Address      string       `json:"address"`
	CreatedBy    string       `json:"createdBy"`
	PriceHistory []PriceEntry `json:"priceHistory"`
	OwnerHistory []OwnerEntry `json:"ownerHistory"`
	ReservedBy   string       `json:"reservedBy"` // 빈 문자열이라도 반드시 존재
	ReservedAt   string       `json:"reservedAt"` // 빈 문자열이라도 반드시 존재
	ExpiresAt    int64        `json:"expiresAt"`  // 0이라도 반드시 존재
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
		// 기존 매물에 히스토리만 추가, reserved 필드는 그대로 유지
		var existing Listing
		if err := json.Unmarshal(listingJSON, &existing); err != nil {
			return fmt.Errorf("Unmarshal 실패: %v", err)
		}
		existing.PriceHistory = append(existing.PriceHistory, PriceEntry{Price: price, Date: currentDate})
		existing.OwnerHistory = append(existing.OwnerHistory, OwnerEntry{Owner: owner, Date: currentDate})

		updated, _ := json.Marshal(existing)
		return ctx.GetStub().PutState(id, updated)
	}

	// 새 매물 생성 시 reserved 필드를 빈값/0으로 초기화
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
	}

	bytes, _ := json.Marshal(newListing)
	return ctx.GetStub().PutState(id, bytes)
}

func (s *SmartContract) GetListing(
	ctx contractapi.TransactionContextInterface,
	id string,
) (*Listing, error) {
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

	// 이미 예약된 경우
	if listing.ReservedBy != "" && time.Now().Unix() < listing.ExpiresAt {
		return fmt.Errorf("이미 예약된 매물입니다 (예약자: %s)", listing.ReservedBy)
	}

	// 예약 정보 갱신
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("클라이언트 ID 조회 실패: %v", err)
	}

	listing.ReservedBy = clientID
	listing.ReservedAt = time.Now().Format("2006-01-02")
	listing.ExpiresAt = expiresAt

	updated, err := json.Marshal(listing)
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %v", err)
	}

	return ctx.GetStub().PutState(id, updated)
}

func (s *SmartContract) GetAllListings(
	ctx contractapi.TransactionContextInterface,
) ([]*Listing, error) {
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
		result = append(result, &l)
	}
	return result, nil
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
