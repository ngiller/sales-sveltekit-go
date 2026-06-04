package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.Unit, int64, error) {
	var items []models.Unit
	var total int64

	query := r.db.Model(&models.Unit{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("units_name LIKE ?", searchTerm)
	}

	query.Count(&total)

	if sortBy == "" || sortBy == "name" {
		sortBy = "units_name"
	}
	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *UnitRepository) FindByID(id uint) (*models.Unit, error) {
	var item models.Unit
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *UnitRepository) Create(item *models.Unit) error {
	return r.db.Create(item).Error
}

func (r *UnitRepository) Update(item *models.Unit) error {
	return r.db.Save(item).Error
}

func (r *UnitRepository) Delete(id uint) error {
	return r.db.Delete(&models.Unit{}, id).Error
}

func (r *UnitRepository) GetDB() *gorm.DB {
	return r.db
}
