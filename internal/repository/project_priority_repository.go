package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type ProjectPriorityRepository struct {
	db *gorm.DB
}

func NewProjectPriorityRepository(db *gorm.DB) *ProjectPriorityRepository {
	return &ProjectPriorityRepository{db: db}
}

func (r *ProjectPriorityRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.ProjectPriority, int64, error) {
	var items []models.ProjectPriority
	var total int64

	query := r.db.Model(&models.ProjectPriority{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "name"
	}
	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ProjectPriorityRepository) FindByID(id uint) (*models.ProjectPriority, error) {
	var item models.ProjectPriority
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ProjectPriorityRepository) Create(item *models.ProjectPriority) error {
	return r.db.Create(item).Error
}

func (r *ProjectPriorityRepository) Update(item *models.ProjectPriority) error {
	return r.db.Save(item).Error
}

func (r *ProjectPriorityRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProjectPriority{}, id).Error
}

func (r *ProjectPriorityRepository) GetDB() *gorm.DB {
	return r.db
}
