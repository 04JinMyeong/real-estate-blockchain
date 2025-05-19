// real-estate-blockchain-feature-did-vc/database/db.go
package database

import (
	"fmt"            // 에러 포맷팅용
	"realestate/did" // did 패키지 임포트 (did.DIDDocument 타입 사용)
	"realestate/models"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns a singleton GORM DB instance
func GetDB() *gorm.DB {
	once.Do(func() {
		d, err := gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
		err = d.AutoMigrate(&models.BrokerVC{}, &models.User{}, &models.DIDDocumentStore{})
		if err != nil {
			panic("failed to auto migrate database: " + err.Error())
		}
		db = d
	})
	return db
}

// StoreBrokerVC saves a BrokerVC record to the database
func StoreBrokerVC(vc *models.BrokerVC) error {
	return GetDB().Create(vc).Error
}

// GetBrokerVC retrieves a BrokerVC by its DID
func GetBrokerVC(id string) (*models.BrokerVC, error) {
	var vc models.BrokerVC
	if err := GetDB().First(&vc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vc, nil
}

// --- DID Document 저장 및 조회 함수 추가 ---

func StoreDIDDocument(agentDID string, doc did.DIDDocument) error {
	db := GetDB()
	docJson, err := doc.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal DID document to JSON: %w", err)
	}

	record := models.DIDDocumentStore{
		DID:      agentDID,
		Document: docJson,
	}

	if err := db.Save(&record).Error; err != nil {
		return fmt.Errorf("failed to save DID document to DB: %w", err)
	}
	fmt.Println("💾 DID Document successfully stored/updated in DB for DID:", agentDID)
	return nil
}

func GetDIDDocument(agentDID string) (*did.DIDDocument, error) {
	db := GetDB()
	var record models.DIDDocumentStore

	if err := db.First(&record, "did = ?", agentDID).Error; err != nil {
		return nil, fmt.Errorf("DID document not found for DID [%s] in DB: %w", agentDID, err)
	}

	// did.FromJson 함수가 did 패키지 내에 구현되어 json.Unmarshal을 수행한다고 가정
	doc, err := did.FromJson(record.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DID document from JSON for DID [%s]: %w", agentDID, err)
	}
	return &doc, nil
}
