package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type DepartementRepository struct {
	db *gorm.DB
}

func NewDepartementRepository(db *gorm.DB) *DepartementRepository {
	return &DepartementRepository{db: db}
}

func (r *DepartementRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.Departement, int64, error) {
	var depts []models.Departement
	var total int64

	query := r.db.Model(&models.Departement{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	query.Count(&total)

	// Default sort
	if sortBy == "" {
		sortBy = "name"
	}
	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&depts).Error
	return depts, total, err
}

func (r *DepartementRepository) FindByID(id uint) (*models.Departement, error) {
	var dept models.Departement
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartementRepository) Create(dept *models.Departement) error {
	return r.db.Create(dept).Error
}

func (r *DepartementRepository) Update(dept *models.Departement) error {
	return r.db.Save(dept).Error
}

func (r *DepartementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Departement{}, id).Error
}

func (r *DepartementRepository) GetDB() *gorm.DB {
	return r.db
}
