package storages

import (
	"time"

	"github.com/fylerx/fyler/internal/storage"
	gormcrypto "github.com/pkasila/gorm-crypto"
)

type Storage struct {
	ID              uint32                    `gorm:"primarykey" json:"id"`
	ProjectID       uint32                    `json:"project_id" mapstructure:"project_id"`
	AccessKeyID     gormcrypto.EncryptedValue `json:"access_key_id" mapstructure:"access_key"`
	SecretAccessKey gormcrypto.EncryptedValue `json:"secret_access_key" mapstructure:"secret_key"`
	Bucket          string                    `json:"bucket"`
	Endpoint        string                    `json:"endpoint"`
	Region          string                    `json:"region"`
	DisableSSL      bool                      `json:"disable_ssl" mapstructure:"disable_ssl"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

func (s *Storage) Config() storage.Config {
	return storage.Config{
		AccessKeyID:     s.AccessKeyID.Raw.(string),
		SecretAccessKey: s.SecretAccessKey.Raw.(string),
		Bucket:          s.Bucket,
		Endpoint:        s.Endpoint,
		Region:          s.Region,
		DisableSSL:      s.DisableSSL,
	}
}
