package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type PaymentTermRepository struct {
	db *gorm.DB
}

func NewPaymentTermRepository(db *gorm.DB) *PaymentTermRepository {
	return &PaymentTermRepository{db: db}
}

func (r *PaymentTermRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.PaymentTerm, int64, error) {
	var items []models.PaymentTerm
	var total int64

	query := r.db.Model(&models.PaymentTerm{})
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

func (r *PaymentTermRepository) FindByID(id uint) (*models.PaymentTerm, error) {
	var item models.PaymentTerm
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PaymentTermRepository) Create(item *models.PaymentTerm) error {
	return r.db.Create(item).Error
}

func (r *PaymentTermRepository) Update(item *models.PaymentTerm) error {
	return r.db.Save(item).Error
}

func (r *PaymentTermRepository) Delete(id uint) error {
	return r.db.Delete(&models.PaymentTerm{}, id).Error
}

func (r *PaymentTermRepository) GetDB() *gorm.DB {
	return r.db
}
