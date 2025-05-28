// package database

// import (
// 	"fmt" // 에러 포맷팅용
// 	"log"
// 	"realestate/did" // did 패키지 임포트 (did.DIDDocument 타입 사용)
// 	"sync"

// 	"realestate/models"

// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// var (

// 	realestateDB     *gorm.DB
//     realestateDBOnce sync.Once

// 	db   *gorm.DB
// 	once sync.Once
// )

// // GetDB returns a singleton GORM DB instance
// func GetDB() *gorm.DB {	//GetRealestateDB
// 	once.Do(func() {
// 		d, err := gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
// 		if err != nil {
// 			panic("failed to connect database: " + err.Error())
// 		}
// 		err = d.AutoMigrate(&models.BrokerVC{}, &models.User{}, &models.DIDDocumentStore{})
// 		if err != nil {
// 			panic("failed to auto migrate database: " + err.Error())
// 		}
// 		realestateDB = d
// 	})
// 	return realestateDB
// }

// // StoreBrokerVC saves a BrokerVC record to the database
// func StoreBrokerVC(vc *models.BrokerVC) error {
// 	return GetDB().Create(vc).Error
// }

// // GetBrokerVC retrieves a BrokerVC by its DID
// func GetBrokerVC(id string) (*models.BrokerVC, error) {
// 	var vc models.BrokerVC
// 	if err := GetDB().First(&vc, "id = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &vc, nil
// }

// // --- DID Document 저장 및 조회 함수 추가 ---

// func StoreDIDDocument(agentDID string, doc did.DIDDocument) error {
// 	db := GetDB()
// 	docJson, err := doc.ToJson()
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal DID document to JSON: %w", err)
// 	}

// 	record := models.DIDDocumentStore{
// 		DID:      agentDID,
// 		Document: docJson,
// 	}

// 	if err := db.Save(&record).Error; err != nil {
// 		return fmt.Errorf("failed to save DID document to DB: %w", err)
// 	}
// 	fmt.Println("💾 DID Document successfully stored/updated in DB for DID:", agentDID)
// 	return nil
// }

// func GetDIDDocument(agentDID string) (*did.DIDDocument, error) {
// 	db := GetDB()
// 	var record models.DIDDocumentStore

// 	if err := db.First(&record, "did = ?", agentDID).Error; err != nil {
// 		return nil, fmt.Errorf("DID document not found for DID [%s] in DB: %w", agentDID, err)
// 	}

// 	// did.FromJson 함수가 did 패키지 내에 구현되어 json.Unmarshal을 수행한다고 가정
// 	doc, err := did.FromJson(record.Document)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal DID document from JSON for DID [%s]: %w", agentDID, err)
// 	}
// 	return &doc, nil
// }

// func InitDB() {
// 	var err error
// 	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("❌ DB 연결 실패:", err)
// 	}

// 	err = DB.AutoMigrate(&models.User{})
// 	if err != nil {
// 		log.Fatal("❌ 사용자 테이블 생성 실패:", err)
// 	}

// 	log.Println("✅ SQLite DB 초기화 완료")
// }

// real-estate-blockchain-feature-did-vc/database/db.go
package database

import (
	"fmt" // 에러 포맷팅용
	"log"
	"realestate/did" // did 패키지 임포트 (did.DIDDocument 타입 사용)
	"sync"

	"realestate/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// "realestate.db" 용 GORM DB 인스턴스 및 초기화 제어 변수
	// DID, VC, 중개사 정보 등 핵심 애플리케이션 데이터용
	realestateDB     *gorm.DB
	realestateDBOnce sync.Once

	// "users.db" 용 GORM DB 인스턴스 (플랫폼 사용자 인증용)
	// handler/auth.go 에서 database.DB 로 참조됩니다.
	DB *gorm.DB
	// usersDBOnce sync.Once // InitDB가 main에서 한 번만 호출되므로, 여기서는 sync.Once가 필수는 아님
)

// GetDB returns a singleton GORM DB instance for "realestate.db"
// 이 함수는 애플리케이션의 주요 데이터 (DID 문서, VC, 사용자 정보 등)를 관리하는 DB에 접근합니다.
func GetDB() *gorm.DB {
	realestateDBOnce.Do(func() { // realestateDB 용 sync.Once 사용
		d, err := gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database (realestate.db): " + err.Error())
		}
		// User 모델도 여기서 마이그레이션 됩니다.
		// SignUpBrokerAndIssueDID 핸들러가 GetDB()를 사용하므로, User 정보가 realestate.db에 저장됩니다.
		err = d.AutoMigrate(&models.BrokerVC{}, &models.User{}, &models.DIDDocumentStore{})
		if err != nil {
			panic("failed to auto migrate (realestate.db): " + err.Error())
		}
		realestateDB = d
		log.Println("✅ SQLite DB 초기화 완료 (realestate.db)")
	})
	return realestateDB
}

// StoreBrokerVC saves a BrokerVC record to the "realestate.db"
func StoreBrokerVC(vc *models.BrokerVC) error {
	return GetDB().Create(vc).Error
}

// GetBrokerVC retrieves a BrokerVC by its DID from "realestate.db"
func GetBrokerVC(id string) (*models.BrokerVC, error) {
	var vc models.BrokerVC
	if err := GetDB().First(&vc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vc, nil
}

// StoreDIDDocument saves a DIDDocument to "realestate.db"
func StoreDIDDocument(agentDID string, doc did.DIDDocument) error {
	db := GetDB() // GetDB()는 realestateDB 인스턴스를 반환
	docJson, err := doc.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal DID document to JSON: %w", err)
	}

	record := models.DIDDocumentStore{
		DID:      agentDID,
		Document: docJson,
	}

	if err := db.Save(&record).Error; err != nil { // realestateDB에 저장
		return fmt.Errorf("failed to save DID document to DB: %w", err)
	}
	fmt.Println("💾 DID Document successfully stored/updated in DB (realestate.db) for DID:", agentDID)
	return nil
}

// GetDIDDocument retrieves a DIDDocument from "realestate.db"
func GetDIDDocument(agentDID string) (*did.DIDDocument, error) {
	db := GetDB() // realestateDB 인스턴스 사용
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

// InitDB initializes the GORM DB instance for "users.db"
// 이 함수는 주로 일반 사용자 가입/로그인(handler/auth.go)에서 사용될 DB를 초기화합니다.
// main.go에서 호출됩니다.
func InitDB() {
	var err error
	// 패키지 변수 DB (대문자) 에 할당
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ DB 연결 실패 (users.db):", err)
	}

	// "users.db"에 User 테이블만 마이그레이션
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ 사용자 테이블 생성 실패 (users.db):", err)
	}
	log.Println("✅ SQLite DB 초기화 완료 (users.db)")
}
