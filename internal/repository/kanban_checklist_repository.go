package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type KanbanChecklistRepository struct {
	db *gorm.DB
}

func NewKanbanChecklistRepository(db *gorm.DB) *KanbanChecklistRepository {
	return &KanbanChecklistRepository{db: db}
}

func (r *KanbanChecklistRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *KanbanChecklistRepository) FindByCardID(cardID uint) ([]models.KanbanChecklist, error) {
	var checklists []models.KanbanChecklist
	err := r.db.Where("card_id = ?", cardID).Order("position ASC").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).Find(&checklists).Error
	return checklists, err
}

func (r *KanbanChecklistRepository) FindChecklistByID(id uint) (*models.KanbanChecklist, error) {
	var cl models.KanbanChecklist
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).First(&cl, id).Error
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *KanbanChecklistRepository) CreateChecklist(cl *models.KanbanChecklist) error {
	return r.db.Create(cl).Error
}

func (r *KanbanChecklistRepository) UpdateChecklist(cl *models.KanbanChecklist) error {
	return r.db.Save(cl).Error
}

func (r *KanbanChecklistRepository) DeleteChecklist(id uint) error {
	return r.db.Delete(&models.KanbanChecklist{}, id).Error
}

func (r *KanbanChecklistRepository) GetMaxChecklistPosition(cardID uint) (int, error) {
	var maxPos struct { Max int }
	err := r.db.Model(&models.KanbanChecklist{}).Select("COALESCE(MAX(position), -1) as max").Where("card_id = ?", cardID).Scan(&maxPos).Error
	return maxPos.Max + 1, err
}

// Checklist Items
func (r *KanbanChecklistRepository) CreateItem(item *models.KanbanChecklistItem) error {
	return r.db.Create(item).Error
}

func (r *KanbanChecklistRepository) UpdateItem(item *models.KanbanChecklistItem) error {
	return r.db.Save(item).Error
}

func (r *KanbanChecklistRepository) DeleteItem(id uint) error {
	return r.db.Delete(&models.KanbanChecklistItem{}, id).Error
}

func (r *KanbanChecklistRepository) FindItemByID(id uint) (*models.KanbanChecklistItem, error) {
	var item models.KanbanChecklistItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *KanbanChecklistRepository) GetMaxItemPosition(checklistID uint) (int, error) {
	var maxPos struct { Max int }
	err := r.db.Model(&models.KanbanChecklistItem{}).Select("COALESCE(MAX(position), -1) as max").Where("checklist_id = ?", checklistID).Scan(&maxPos).Error
	return maxPos.Max + 1, err
}
