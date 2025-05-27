package models

type Property struct {
	ID           string       `json:"id"`
	Address      string       `json:"address"`
	CreatedBy    string       `json:"createdBy"`
	PriceHistory []PriceEntry `json:"priceHistory"`
	OwnerHistory []OwnerEntry `json:"ownerHistory"`
	ReservedBy   string       `json:"reservedBy,omitempty"`
	ReservedAt   string       `json:"reservedAt,omitempty"`
	ExpiresAt    int64        `json:"expiresAt,omitempty"`
}

type PriceEntry struct {
	Price string `json:"price"`
	Date  string `json:"date"`
}

type OwnerEntry struct {
	Owner string `json:"owner"`
	Date  string `json:"date"`
}
