package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.UserGroup, int64, error) {
	var roles []models.UserGroup
	var total int64

	query := r.db.Model(&models.UserGroup{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	// Default sort
	if sortBy == "" {
		sortBy = "name"
	}
	validSort := map[string]bool{
		"id":   true,
		"name": true,
	}
	if !validSort[sortBy] {
		sortBy = "name"
	}
	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *RoleRepository) FindByID(id uint) (*models.UserGroup, error) {
	var role models.UserGroup
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(role *models.UserGroup) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *models.UserGroup) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserGroup{}, id).Error
}
