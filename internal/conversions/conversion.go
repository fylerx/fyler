package conversions

import (
	"time"
)

type Conversion struct {
	ID           uint64    `gorm:"primarykey" json:"id"`
	TaskID       uint64    `json:"task_id"`
	DownloadTime int       `json:"download_time"`
	UploadTime   int       `json:"upload_time"`
	FileSize     int64     `json:"file_size"`
	ResultPath   string    `json:"result_path"`
	TimeSpent    int       `json:"time_spent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
