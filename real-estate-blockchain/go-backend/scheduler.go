// scheduler.go
package main

import (
	"encoding/json"
	"log"
	"realestate/blockchain"
	"time"
	// 실제 모듈 경로에 맞게 수정하세요
)

// ListingInfo는 GetAllListings 결과를 언마샬링할 때 사용할 최소 구조체입니다.
// JSON 필드명("id", "reservedBy", "expiresAt")과 정확히 일치해야 합니다.
type ListingInfo struct {
	ID         string  `json:"id"`
	ReservedBy string  `json:"reservedBy"`
	ExpiresAt  float64 `json:"expiresAt"` // JSON에서 number 타입으로 내려오므로 float64 사용
}

func StartReservationReleaser() {
	// 테스트용으로 1분마다 실행하려면 time.Minute으로 변경하세요. (운영 시에는 30*time.Minute 권장)
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for range ticker.C {
			log.Println("🔄 예약 만료 스케줄러 실행 중…")
			now := time.Now().Unix()

			// 1) 모든 매물 조회 (admin 계정)
			respStr, err := blockchain.QueryAllProperties("admin")
			if err != nil {
				log.Printf("❌ 전체 매물 조회 실패: %v\n", err)
				continue
			}

			// 2) JSON 파싱
			var allListings []ListingInfo
			if err := json.Unmarshal([]byte(respStr), &allListings); err != nil {
				log.Printf("❌ JSON 파싱 실패: %v\n", err)
				continue
			}

			// 3) 만료된 예약 자동 해제
			for _, l := range allListings {
				if l.ReservedBy != "" && now > int64(l.ExpiresAt) {
					log.Printf("⏰ 만료 감지: ID=%s, expiresAt=%.0f (현재=%d) → 자동 해제\n", l.ID, l.ExpiresAt, now)
					if err := blockchain.ReleaseListing("admin", l.ID); err != nil {
						log.Printf("⚠️ 자동 해제 실패 (ID=%s): %v\n", l.ID, err)
					} else {
						log.Printf("✅ 자동 해제 완료 (ID=%s)\n", l.ID)
					}
				}
			}
		}
	}()
}
