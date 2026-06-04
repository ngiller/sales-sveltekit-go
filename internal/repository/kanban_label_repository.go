package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanLabelRepository struct {
	db *gorm.DB
}

func NewKanbanLabelRepository(db *gorm.DB) *KanbanLabelRepository {
	return &KanbanLabelRepository{db: db}
}

func (r *KanbanLabelRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanLabelRepository) FindByBoardID(boardID uint) ([]models.KanbanLabel, error) {
	var labels []models.KanbanLabel
	err := r.db.Where("board_id = ?", boardID).Order("name ASC").Find(&labels).Error
	return labels, err
}

func (r *KanbanLabelRepository) FindByID(id uint) (*models.KanbanLabel, error) {
	var label models.KanbanLabel
	err := r.db.First(&label, id).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *KanbanLabelRepository) Create(label *models.KanbanLabel) error {
	return r.db.Create(label).Error
}

func (r *KanbanLabelRepository) Update(label *models.KanbanLabel) error {
	return r.db.Save(label).Error
}

func (r *KanbanLabelRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanLabel{}, id).Error
}
