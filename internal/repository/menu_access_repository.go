package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type MenuAccessRepository struct {
	db *gorm.DB
}

func NewMenuAccessRepository(db *gorm.DB) *MenuAccessRepository {
	return &MenuAccessRepository{db: db}
}

func (r *MenuAccessRepository) FindAll() ([]models.MasterTableAccess, error) {
	var items []models.MasterTableAccess
	err := r.db.Order("sort_order ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *MenuAccessRepository) FindByID(id uint) (*models.MasterTableAccess, error) {
	var item models.MasterTableAccess
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MenuAccessRepository) Create(item *models.MasterTableAccess) error {
	return r.db.Create(item).Error
}

func (r *MenuAccessRepository) Update(item *models.MasterTableAccess) error {
	return r.db.Save(item).Error
}

func (r *MenuAccessRepository) Delete(id uint) error {
	return r.db.Delete(&models.MasterTableAccess{}, id).Error
}
