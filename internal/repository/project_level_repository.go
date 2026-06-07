package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type ProjectLevelRepository struct {
	db *gorm.DB
}

func NewProjectLevelRepository(db *gorm.DB) *ProjectLevelRepository {
	return &ProjectLevelRepository{db: db}
}

func (r *ProjectLevelRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.ProjectLevel, int64, error) {
	var items []models.ProjectLevel
	var total int64

	query := r.db.Model(&models.ProjectLevel{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "name"
	}
	validSort := map[string]bool{
		"id":   true,
		"name": true,
	}
	if !validSort[sortBy] {
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

func (r *ProjectLevelRepository) FindByID(id uint) (*models.ProjectLevel, error) {
	var item models.ProjectLevel
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ProjectLevelRepository) Create(item *models.ProjectLevel) error {
	return r.db.Create(item).Error
}

func (r *ProjectLevelRepository) Update(item *models.ProjectLevel) error {
	return r.db.Save(item).Error
}

func (r *ProjectLevelRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProjectLevel{}, id).Error
}

func (r *ProjectLevelRepository) GetDB() *gorm.DB {
	return r.db
}
