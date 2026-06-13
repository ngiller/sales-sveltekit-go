package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) FindByCode(code int, propertyID int) (*models.Setting, error) {
	var setting models.Setting
	err := r.db.Where("code = ? AND property_id = ?", code, propertyID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}
