package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanCommentRepository struct {
	db *gorm.DB
}

func NewKanbanCommentRepository(db *gorm.DB) *KanbanCommentRepository {
	return &KanbanCommentRepository{db: db}
}

func (r *KanbanCommentRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanCommentRepository) FindByCardID(cardID uint) ([]models.KanbanComment, error) {
	var comments []models.KanbanComment
	err := r.db.Where("card_id = ?", cardID).Order("created_at ASC").
		Preload("User").Find(&comments).Error
	return comments, err
}

func (r *KanbanCommentRepository) FindByID(id uint) (*models.KanbanComment, error) {
	var c models.KanbanComment
	err := r.db.Preload("User").First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *KanbanCommentRepository) Create(comment *models.KanbanComment) error {
	return r.db.Create(comment).Error
}

func (r *KanbanCommentRepository) Update(comment *models.KanbanComment) error {
	return r.db.Save(comment).Error
}

func (r *KanbanCommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanComment{}, id).Error
}
