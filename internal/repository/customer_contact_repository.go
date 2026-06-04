package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type CustomerContactRepository struct {
	db *gorm.DB
}

func NewCustomerContactRepository(db *gorm.DB) *CustomerContactRepository {
	return &CustomerContactRepository{db: db}
}

func (r *CustomerContactRepository) FindAllByCustomerID(customerID uint) ([]models.CustomerContact, error) {
	var contacts []models.CustomerContact
	err := r.db.Where("customer_id = ?", customerID).Find(&contacts).Error
	return contacts, err
}

func (r *CustomerContactRepository) FindByID(id uint) (*models.CustomerContact, error) {
	var contact models.CustomerContact
	err := r.db.First(&contact, id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *CustomerContactRepository) Create(contact *models.CustomerContact) error {
	return r.db.Create(contact).Error
}

func (r *CustomerContactRepository) Update(contact *models.CustomerContact) error {
	return r.db.Save(contact).Error
}

func (r *CustomerContactRepository) Delete(id uint) error {
	return r.db.Delete(&models.CustomerContact{}, id).Error
}
