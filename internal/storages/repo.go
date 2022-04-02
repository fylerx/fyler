package storages

import (
	"gorm.io/gorm"
)

type IStorage interface {
	CreateStorage(input *Storage) error
	DeleteStorage(id uint32) error
}

type StorageRepo struct {
	storages *gorm.DB
}

func InitRepo(db *gorm.DB) IStorage {
	return &StorageRepo{db.Model(&Storage{}).Debug()}
}

func (repo *StorageRepo) CreateStorage(input *Storage) error {
	return repo.storages.Create(input).Error
}

func (repo *StorageRepo) DeleteStorage(id uint32) error {
	return repo.storages.Delete(&Storage{}, id).Error
}
