package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type QuotationProgressRepository struct {
	db *gorm.DB
}

func NewQuotationProgressRepository(db *gorm.DB) *QuotationProgressRepository {
	return &QuotationProgressRepository{db: db}
}

func (r *QuotationProgressRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.QuotationProgress, int64, error) {
	var items []models.QuotationProgress
	var total int64

	query := r.db.Model(&models.QuotationProgress{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "name"
	}
	validSort := map[string]bool{
		"id":       true,
		"name":     true,
		"progress": true,
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

func (r *QuotationProgressRepository) FindByID(id uint) (*models.QuotationProgress, error) {
	var item models.QuotationProgress
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *QuotationProgressRepository) Create(item *models.QuotationProgress) error {
	return r.db.Create(item).Error
}

func (r *QuotationProgressRepository) Update(item *models.QuotationProgress) error {
	return r.db.Save(item).Error
}

func (r *QuotationProgressRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuotationProgress{}, id).Error
}

func (r *QuotationProgressRepository) GetDB() *gorm.DB {
	return r.db
}
