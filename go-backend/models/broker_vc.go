package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// VC 구조 정의 클래스
type BrokerVC struct {
    // DID as primary key
    ID         string         `gorm:"primaryKey;column:id" json:"id"`
    Name       string         `gorm:"column:name" json:"name"`
    LicenseNo  string         `gorm:"column:license_no" json:"license_no"`
    Issuer     string         `gorm:"column:issuer" json:"issuer"`
    IssuedAt   time.Time      `gorm:"column:issued_at" json:"issued_at"`
    Signature  string         `gorm:"column:signature" json:"signature"`
    CreatedAt  time.Time      `gorm:"autoCreateTime"`
    UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook generates DID if not set
func (vc *BrokerVC) BeforeCreate(tx *gorm.DB) (err error) {
    if vc.ID == "" {
        vc.ID = "did:realestate:" + uuid.NewString()
    }
    return
}