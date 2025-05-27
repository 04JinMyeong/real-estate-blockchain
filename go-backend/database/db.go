package database

import (
	"log"

	"realestate/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ DB 연결 실패:", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ 사용자 테이블 생성 실패:", err)
	}

	log.Println("✅ SQLite DB 초기화 완료")
}
