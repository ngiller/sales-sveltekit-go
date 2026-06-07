package repository

import (
	"backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserGroup").Preload("Departement").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if user.UserGroup != nil {
		user.RoleName = user.UserGroup.Name
	}
	if user.Departement != nil {
		user.DeptName = user.Departement.Name
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserGroup").Preload("Departement").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if user.UserGroup != nil {
		user.RoleName = user.UserGroup.Name
	}
	if user.Departement != nil {
		user.DeptName = user.Departement.Name
	}

	return &user, nil
}
func (r *UserRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	// Use explicit JOIN to fetch the names directly into the User struct virtual fields
	query := r.db.Table("users").
		Select("users.*, user_groups.name as role_name, master_departements.name as dept_name").
		Joins("left join user_groups on user_groups.id = users.user_group_id").
		Joins("left join master_departements on master_departements.id = users.departement_id")

	if search != "" {
		searchStr := "%" + search + "%"
		query = query.Where("users.name LIKE ? OR users.email LIKE ?", searchStr, searchStr)
	}

	query.Count(&total)

	// Default sort
	if sortBy == "" {
		sortBy = "users.name"
	} else if sortBy == "name" {
		sortBy = "users.name"
	} else if sortBy == "email" {
		sortBy = "users.email"
	} else if sortBy == "role_name" {
		sortBy = "user_groups.name"
	} else if sortBy == "dept_name" {
		sortBy = "master_departements.name"
	} else {
		sortBy = "users.name"
	}

	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Scan(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
