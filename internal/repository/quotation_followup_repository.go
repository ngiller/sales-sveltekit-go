package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type QuotationFollowupRepository struct {
	db *gorm.DB
}

func NewQuotationFollowupRepository(db *gorm.DB) *QuotationFollowupRepository {
	return &QuotationFollowupRepository{db: db}
}

func (r *QuotationFollowupRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *QuotationFollowupRepository) FindAllByQuotationID(quotationID string) ([]models.QuotationFollowup, error) {
	var items []models.QuotationFollowup
	err := r.db.Preload("FollowupByUser").Preload("StatusInfo").Preload("ProgressInfo").
		Where("id = ?", quotationID).
		Order("followup_date DESC").
		Find(&items).Error
	return items, err
}

func (r *QuotationFollowupRepository) FindByID(id string, lineID int) (*models.QuotationFollowup, error) {
	var item models.QuotationFollowup
	err := r.db.Preload("FollowupByUser").Preload("StatusInfo").Preload("ProgressInfo").
		Where("id = ? AND line_id = ?", id, lineID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *QuotationFollowupRepository) Create(item *models.QuotationFollowup) error {
	return r.db.Create(item).Error
}

func (r *QuotationFollowupRepository) Update(item *models.QuotationFollowup) error {
	return r.db.Save(item).Error
}

func (r *QuotationFollowupRepository) Delete(id string, lineID int) error {
	return r.db.Where("id = ? AND line_id = ?", id, lineID).Delete(&models.QuotationFollowup{}).Error
}
