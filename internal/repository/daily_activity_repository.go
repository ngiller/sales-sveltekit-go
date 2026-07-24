package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type DailyActivityRepository struct {
	db *gorm.DB
}

func NewDailyActivityRepository(db *gorm.DB) *DailyActivityRepository {
	return &DailyActivityRepository{db: db}
}

func (r *DailyActivityRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *DailyActivityRepository) FindAll(userIDFilter uint, search string, page, limit int, sortBy, sortDir string, startDate, endDate string) ([]models.DailyActivity, int64, error) {
	var items []models.DailyActivity
	var total int64

	query := r.db.Table("daily_activities").
		Select("daily_activities.*, users.name as user_name").
		Joins("left join users on users.id = daily_activities.user_id")

	if userIDFilter > 0 {
		query = query.Where("daily_activities.user_id = ?", userIDFilter)
	}

	if startDate != "" && endDate != "" {
		query = query.Where("daily_activities.activity_date BETWEEN ? AND ?", startDate, endDate)
	}

	if search != "" {
		searchStr := "%" + search + "%"
		query = query.Where("(daily_activities.activity LIKE ? OR daily_activities.process LIKE ? OR daily_activities.issues LIKE ? OR daily_activities.result LIKE ? OR daily_activities.notes LIKE ? OR users.name LIKE ?)",
			searchStr, searchStr, searchStr, searchStr, searchStr, searchStr)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "daily_activities.activity_date"
	} else if sortBy == "activity_date" {
		sortBy = "daily_activities.activity_date"
	} else if sortBy == "user_name" {
		sortBy = "users.name"
	} else if sortBy == "activity" {
		sortBy = "daily_activities.activity"
	} else if sortBy == "id" {
		sortBy = "daily_activities.id"
	} else {
		sortBy = "daily_activities.activity_date"
	}

	if sortDir == "" {
		sortDir = "desc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *DailyActivityRepository) FindByID(id uint) (*models.DailyActivity, error) {
	var item models.DailyActivity
	err := r.db.Preload("User").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DailyActivityRepository) Create(item *models.DailyActivity) error {
	return r.db.Create(item).Error
}

func (r *DailyActivityRepository) Update(item *models.DailyActivity) error {
	return r.db.Save(item).Error
}

func (r *DailyActivityRepository) Delete(id uint) error {
	return r.db.Delete(&models.DailyActivity{}, id).Error
}
