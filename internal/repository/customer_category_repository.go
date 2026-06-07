package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type CustomerCategoryRepository struct {
	db *gorm.DB
}

func NewCustomerCategoryRepository(db *gorm.DB) *CustomerCategoryRepository {
	return &CustomerCategoryRepository{db: db}
}

func (r *CustomerCategoryRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *CustomerCategoryRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.CustomerCategory, int64, error) {
	var categories []models.CustomerCategory
	var total int64

	query := r.db.Model(&models.CustomerCategory{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	// Default sort
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
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&categories).Error
	return categories, total, err
}

func (r *CustomerCategoryRepository) FindByID(id uint) (*models.CustomerCategory, error) {
	var category models.CustomerCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CustomerCategoryRepository) Create(category *models.CustomerCategory) error {
	return r.db.Create(category).Error
}

func (r *CustomerCategoryRepository) Update(category *models.CustomerCategory) error {
	return r.db.Save(category).Error
}

func (r *CustomerCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.CustomerCategory{}, id).Error
}
