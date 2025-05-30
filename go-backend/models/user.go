package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Enrolled  bool      `json:"enrolled"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role,omitempty"` // ◀◀◀ Role 필드 추가 (문자열 타입, 선택적 JSON)

	// ID       uint   `json:"id"`
	// Username string `json:"username"`
	// Password string `json:"password"`
}
