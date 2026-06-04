package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanListRepository struct {
	db *gorm.DB
}

func NewKanbanListRepository(db *gorm.DB) *KanbanListRepository {
	return &KanbanListRepository{db: db}
}

func (r *KanbanListRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanListRepository) FindByBoardID(boardID uint) ([]models.KanbanList, error) {
	var lists []models.KanbanList
	err := r.db.Where("board_id = ?", boardID).Order("position ASC").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_archived = ?", false).Order("position ASC")
		}).
		Preload("Cards.Members").
		Preload("Cards.Labels").
		Find(&lists).Error
	return lists, err
}

func (r *KanbanListRepository) FindByID(id uint) (*models.KanbanList, error) {
	var list models.KanbanList
	err := r.db.First(&list, id).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *KanbanListRepository) Create(list *models.KanbanList) error {
	return r.db.Create(list).Error
}

func (r *KanbanListRepository) Update(list *models.KanbanList) error {
	return r.db.Save(list).Error
}

func (r *KanbanListRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanList{}, id).Error
}

func (r *KanbanListRepository) GetMaxPosition(boardID uint) (int, error) {
	var maxPos struct {
		Max int
	}
	err := r.db.Model(&models.KanbanList{}).Select("COALESCE(MAX(position), -1) as max").Where("board_id = ?", boardID).Scan(&maxPos).Error
	return maxPos.Max + 1, err
}

func (r *KanbanListRepository) Reorder(ids []uint) error {
	for i, id := range ids {
		if err := r.db.Model(&models.KanbanList{}).Where("id = ?", id).Update("position", i).Error; err != nil {
			return err
		}
	}
	return nil
}
