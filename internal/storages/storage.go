package storages

import (
	"time"

	gormcrypto "github.com/pkasila/gorm-crypto"
)

type Storage struct {
	ProjectID       uint32                    `gorm:"primarykey" json:"project_id" mapstructure:"project_id"`
	AccessKeyID     gormcrypto.EncryptedValue `json:"access_key_id" mapstructure:"access_key"`
	SecretAccessKey gormcrypto.EncryptedValue `json:"secret_access_key" mapstructure:"secret_key"`
	Bucket          string                    `json:"bucket"`
	Endpoint        string                    `json:"endpoint"`
	Region          string                    `json:"region"`
	DisableSSL      bool                      `json:"disable_ssl" mapstructure:"disable_ssl"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// type StorageInput struct {
// 	ProjectID       uint32                    `mapstructure:"project_id"`
// 	AccessKeyID     gormcrypto.EncryptedValue `mapstructure:"access_key"`
// 	SecretAccessKey gormcrypto.EncryptedValue `mapstructure:"secret_key"`
// 	Bucket          string
// 	Endpoint        string
// 	Region          string
// 	DisableSSL      bool `mapstructure:"disable_ssl"`
// }
