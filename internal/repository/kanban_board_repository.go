package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanBoardRepository struct {
	db *gorm.DB
}

func NewKanbanBoardRepository(db *gorm.DB) *KanbanBoardRepository {
	return &KanbanBoardRepository{db: db}
}

func (r *KanbanBoardRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanBoardRepository) FindAll(search string, page, limit int) ([]models.KanbanBoard, int64, error) {
	var boards []models.KanbanBoard
	var total int64

	query := r.db.Model(&models.KanbanBoard{}).Where("is_archived = ?", false)
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&boards).Error
	return boards, total, err
}

func (r *KanbanBoardRepository) FindByID(id uint) (*models.KanbanBoard, error) {
	var board models.KanbanBoard
	err := r.db.Preload("Lists", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Preload("Lists.Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_archived = ?", false).Order("position ASC")
	}).Preload("Lists.Cards.Members").Preload("Lists.Cards.Labels").First(&board, id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *KanbanBoardRepository) Create(board *models.KanbanBoard) error {
	return r.db.Create(board).Error
}

func (r *KanbanBoardRepository) Update(board *models.KanbanBoard) error {
	return r.db.Save(board).Error
}

func (r *KanbanBoardRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanBoard{}, id).Error
}
