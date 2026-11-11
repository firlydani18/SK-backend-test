package data

import "time"

// Applicant is the persistence model for an applicant record.
type Applicant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Email     string    `gorm:"size:160;not null;uniqueIndex" json:"email"`
	Position  string    `gorm:"size:120;not null" json:"position"`
	Stage     string    `gorm:"size:30;not null;index" json:"stage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
