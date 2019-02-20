package models

import "time"

// Bin represents the container that holds incoming webhook payloads
type Bin struct {
	ID        string
	Title     string
	AccountID string
	UpdatedAt *time.Time
	CreatedAt *time.Time
}
