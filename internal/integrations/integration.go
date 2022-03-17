package integrations

import (
	"time"

	"github.com/fylerx/fyler/internal/enum"
	gormcrypto "github.com/pkasila/gorm-crypto"
)

type Integration struct {
	ID          uint32                    `gorm:"primarykey" json:"id"`
	ProjectID   uint32                    `json:"project_id"`
	TaskType    enum.Service              `gorm:"type:service" json:"service"`
	Credentials gormcrypto.EncryptedValue `json:"-"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}
