package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.Table("customer").
		Select("customer.*, customer_category.name as category_name").
		Joins("left join customer_category on customer_category.id = customer.category_id")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("customer.name LIKE ? OR customer.email LIKE ? OR customer.phone LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	// Default sort mapping
	if sortBy == "" || sortBy == "name" {
		sortBy = "customer.name"
	} else if sortBy == "category_name" {
		sortBy = "customer_category.name"
	} else if sortBy == "email" {
		sortBy = "customer.email"
	} else if sortBy == "phone" {
		sortBy = "customer.phone"
	}

	if sortDir == "" {
		sortDir = "asc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Preload("Contacts").Order(orderClause).Offset(offset).Limit(limit).Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) FindByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("Contacts").First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

func (r *CustomerRepository) GetDB() *gorm.DB {
	return r.db
}
