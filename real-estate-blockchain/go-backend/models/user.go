package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Enrolled  bool      `json:"enrolled"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    // "user" 또는 "agent"
	DID       string    `json:"did,omitempty"` // 공인중개사의 DID 저장용
}
