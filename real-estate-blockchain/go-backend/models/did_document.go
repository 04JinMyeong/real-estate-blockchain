// real-estate-blockchain-feature-did-vc/models/did_document.go
package models

import "time"

// DIDDocumentStore는 데이터베이스에 저장될 DID Document의 스키마입니다.
type DIDDocumentStore struct {
	DID       string `gorm:"primaryKey"` // DID 문자열 (예: did:realestate:...)
	Document  string `gorm:"type:text"`  // DID Document 전체를 JSON 문자열로 저장 (TEXT 타입으로 충분한 공간 확보)
	CreatedAt time.Time
	UpdatedAt time.Time
}
