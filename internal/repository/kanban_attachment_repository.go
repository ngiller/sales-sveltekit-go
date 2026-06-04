package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanAttachmentRepository struct {
	db *gorm.DB
}

func NewKanbanAttachmentRepository(db *gorm.DB) *KanbanAttachmentRepository {
	return &KanbanAttachmentRepository{db: db}
}

func (r *KanbanAttachmentRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanAttachmentRepository) FindByCardID(cardID uint) ([]models.KanbanAttachment, error) {
	var attachments []models.KanbanAttachment
	err := r.db.Where("card_id = ?", cardID).Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}

func (r *KanbanAttachmentRepository) FindByID(id uint) (*models.KanbanAttachment, error) {
	var a models.KanbanAttachment
	err := r.db.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *KanbanAttachmentRepository) Create(a *models.KanbanAttachment) error {
	return r.db.Create(a).Error
}

func (r *KanbanAttachmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanAttachment{}, id).Error
}
