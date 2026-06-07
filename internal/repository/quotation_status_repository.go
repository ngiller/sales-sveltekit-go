package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type QuotationStatusRepository struct {
	db *gorm.DB
}

func NewQuotationStatusRepository(db *gorm.DB) *QuotationStatusRepository {
	return &QuotationStatusRepository{db: db}
}

func (r *QuotationStatusRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.QuotationStatus, int64, error) {
	var items []models.QuotationStatus
	var total int64

	query := r.db.Model(&models.QuotationStatus{})
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

func (r *QuotationStatusRepository) FindByID(id uint) (*models.QuotationStatus, error) {
	var item models.QuotationStatus
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *QuotationStatusRepository) Create(item *models.QuotationStatus) error {
	return r.db.Create(item).Error
}

func (r *QuotationStatusRepository) Update(item *models.QuotationStatus) error {
	return r.db.Save(item).Error
}

func (r *QuotationStatusRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuotationStatus{}, id).Error
}

func (r *QuotationStatusRepository) GetDB() *gorm.DB {
	return r.db
}
