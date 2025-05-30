package database

import (
	"fmt"
	"log"
	"realestate/did"
	"realestate/models"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbInitOnce sync.Once
)

// InitDB initializes the SQLite database and performs migrations
func InitDB() *gorm.DB {
	dbInitOnce.Do(func() {
		var err error
		dbInstance, err = gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ DB 연결 실패 (realestate.db):", err)
		}

		err = dbInstance.AutoMigrate(
			&models.User{},
			&models.Property{},
			&models.BrokerVC{},
			&models.DIDDocumentStore{},
		)
		if err != nil {
			log.Fatal("❌ 테이블 마이그레이션 실패 (realestate.db):", err)
		}
		log.Println("✅ SQLite DB 초기화 완료 (realestate.db)")
	})

	return dbInstance
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	if dbInstance == nil {
		return InitDB()
	}
	return dbInstance
}

// 💾 BrokerVC 저장
func StoreBrokerVC(vc *models.BrokerVC) error {
	return GetDB().Create(vc).Error
}

// 🔍 BrokerVC 조회
func GetBrokerVC(id string) (*models.BrokerVC, error) {
	var vc models.BrokerVC
	if err := GetDB().First(&vc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vc, nil
}

// 💾 DIDDocument 저장
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
	fmt.Println("💾 DID Document successfully stored/updated in DB (realestate.db) for DID:", agentDID)
	return nil
}

// 🔍 DIDDocument 조회
func GetDIDDocument(agentDID string) (*did.DIDDocument, error) {
	db := GetDB()
	var record models.DIDDocumentStore

	if err := db.First(&record, "did = ?", agentDID).Error; err != nil {
		return nil, fmt.Errorf("DID document not found for DID [%s] in DB (realestate.db): %w", agentDID, err)
	}

	doc, err := did.FromJson(record.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DID document from JSON for DID [%s] (realestate.db): %w", agentDID, err)
	}
	return &doc, nil
}
