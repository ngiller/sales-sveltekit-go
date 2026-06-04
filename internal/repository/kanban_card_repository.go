package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanCardRepository struct {
	db *gorm.DB
}

func NewKanbanCardRepository(db *gorm.DB) *KanbanCardRepository {
	return &KanbanCardRepository{db: db}
}

func (r *KanbanCardRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanCardRepository) FindByListID(listID uint) ([]models.KanbanCard, error) {
	var cards []models.KanbanCard
	err := r.db.Where("list_id = ? AND is_archived = ?", listID, false).
		Order("position ASC").
		Preload("Members").Preload("Labels").
		Preload("Checklists", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Checklists.Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Attachments").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Comments.User").
		Find(&cards).Error
	return cards, err
}

func (r *KanbanCardRepository) FindByBoardID(boardID uint) ([]models.KanbanCard, error) {
	var cards []models.KanbanCard
	err := r.db.Where("board_id = ? AND is_archived = ?", boardID, false).
		Order("position ASC").
		Preload("Members").Preload("Labels").
		Find(&cards).Error
	return cards, err
}

func (r *KanbanCardRepository) FindByID(id uint) (*models.KanbanCard, error) {
	var card models.KanbanCard
	err := r.db.Preload("Members").Preload("Labels").
		Preload("Checklists", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Checklists.Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Attachments").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Comments.User").
		First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *KanbanCardRepository) Create(card *models.KanbanCard) error {
	return r.db.Create(card).Error
}

func (r *KanbanCardRepository) Update(card *models.KanbanCard) error {
	return r.db.Save(card).Error
}

func (r *KanbanCardRepository) Delete(id uint) error {
	return r.db.Delete(&models.KanbanCard{}, id).Error
}

func (r *KanbanCardRepository) GetMaxPosition(listID uint) (int, error) {
	var maxPos struct {
		Max int
	}
	err := r.db.Model(&models.KanbanCard{}).Select("COALESCE(MAX(position), -1) as max").Where("list_id = ? AND is_archived = ?", listID, false).Scan(&maxPos).Error
	return maxPos.Max + 1, err
}

func (r *KanbanCardRepository) MoveCard(cardID uint, toListID uint, newPosition int) error {
	return r.db.Model(&models.KanbanCard{}).Where("id = ?", cardID).Updates(map[string]interface{}{
		"list_id": toListID,
		"position": newPosition,
	}).Error
}

func (r *KanbanCardRepository) Reorder(ids []uint, listID uint) error {
	for i, id := range ids {
		if err := r.db.Model(&models.KanbanCard{}).Where("id = ?", id).Update("position", i).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *KanbanCardRepository) SyncMembers(cardID uint, userIDs []uint) error {
	var card models.KanbanCard
	if err := r.db.First(&card, cardID).Error; err != nil {
		return err
	}
	var users []models.User
	if len(userIDs) > 0 {
		if err := r.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			return err
		}
	}
	return r.db.Model(&card).Association("Members").Replace(users)
}

func (r *KanbanCardRepository) SyncLabels(cardID uint, labelIDs []uint) error {
	var card models.KanbanCard
	if err := r.db.First(&card, cardID).Error; err != nil {
		return err
	}
	var labels []models.KanbanLabel
	if len(labelIDs) > 0 {
		if err := r.db.Where("id IN ?", labelIDs).Find(&labels).Error; err != nil {
			return err
		}
	}
	return r.db.Model(&card).Association("Labels").Replace(labels)
}
