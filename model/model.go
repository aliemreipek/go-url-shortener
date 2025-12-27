package model

import "time"

// Url is GormModel definition (Customized to use JSON tags effectively)
type Url struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Redirect  string    `gorm:"not null" json:"redirect"`   // The original long URL
	URL       string    `gorm:"unique;not null" json:"url"` // The short code (e.g., "aX9d")
	Clicked   uint64    `gorm:"default:0" json:"clicked"`   // Analytics: Click count
	Random    bool      `gorm:"default:true" json:"random"` // Is it auto-generated or custom?
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
