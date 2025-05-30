package models

import (
	"gorm.io/datatypes"
)

type Property struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Address      string         `json:"address"`
	CreatedBy    string         `json:"createdBy"`
	PriceHistory datatypes.JSON `json:"priceHistory"` // JSON 문자열로 저장
	OwnerHistory datatypes.JSON `json:"ownerHistory"` // JSON 문자열로 저장
	ReservedBy   string         `json:"reservedBy"`
	ReservedAt   string         `json:"reservedAt"`
	ExpiresAt    int64          `json:"expiresAt"`
}

type PriceEntry struct {
	Date  string `json:"date"`
	Price string `json:"price"`
}

type OwnerEntry struct {
	Date  string `json:"date"`
	Owner string `json:"owner"`
}
