package database

import (
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
        // 여기서는 예시로 sqlite 사용.
        // 실제로는 Postgres/MySQL 등 선택 가능.
        d, err := gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
        if err != nil {
            panic("failed to connect database: " + err.Error())
        }
        // 자동 마이그레이션
        d.AutoMigrate(&models.BrokerVC{})
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