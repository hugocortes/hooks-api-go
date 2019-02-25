package models

import "time"

// Bin represents the container that holds incoming webhook payloads
type Bin struct {
	ID        string `gorm:"primary_key;type:char(36)"`
	Title     string `gorm:"size:255;not null"`
	AccountID string `gorm:"type:char(36);not null;index:idx_account_id"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// Initialized validates if Bin is intialized
func (m *Bin) Initialized() bool {
	return m.ID != "" && m.CreatedAt != nil && m.UpdatedAt != nil
}
