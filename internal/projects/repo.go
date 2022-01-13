package projects

import (
	"github.com/fylerx/fyler/pkg/utils/randutils"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]*Project, error)
	GetByID(id uint32) (*Project, error)
	GetByAPIKey(api_key string) (*Project, error)
	Create(project *Project) (*Project, error)
	Update(id uint32, name string) (bool, error)
	Delete(id uint32) (bool, error)
}

type ProjectsRepo struct {
	projects *gorm.DB
}

func InitRepo(db *gorm.DB) Repository {
	return &ProjectsRepo{db.Model(&Project{}).Debug()}
}

func (repo *ProjectsRepo) GetAll() ([]*Project, error) {
	var projects []*Project
	err := repo.projects.Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (repo *ProjectsRepo) GetByID(id uint32) (*Project, error) {
	var project *Project
	err := repo.projects.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return project, err
}

func (repo *ProjectsRepo) GetByAPIKey(api_key string) (*Project, error) {
	var project *Project
	err := repo.projects.Where("api_key = ?", api_key).First(&project).Error
	if err != nil {
		return nil, err
	}
	return project, err
}

func (repo *ProjectsRepo) Create(project *Project) (*Project, error) {
	project.APIKey = randutils.RandString(32)
	err := repo.projects.Select("Name", "APIKey").Create(project).Error
	if err != nil {
		return nil, err
	}

	return project, err
}

func (repo *ProjectsRepo) Update(id uint32, name string) (bool, error) {
	_, err := repo.GetByID(id)
	if err != nil {
		return false, err
	}

	err = repo.projects.Where("id = ?", id).Update("name", name).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *ProjectsRepo) Delete(id uint32) (bool, error) {
	err := repo.projects.Delete(&Project{}, id).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
