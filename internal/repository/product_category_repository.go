package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type ProductCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

func (r *ProductCategoryRepository) FindAll() ([]models.ProductCategory, error) {
	var categories []models.ProductCategory
	err := r.db.Find(&categories).Error
	return categories, err
}
