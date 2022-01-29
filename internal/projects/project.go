package projects

import (
	"time"
)

type Project struct {
	ID        uint32    `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"apikey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
