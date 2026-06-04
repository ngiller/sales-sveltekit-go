package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type PolicyRepository struct {
	db *gorm.DB
}

func NewPolicyRepository(db *gorm.DB) *PolicyRepository {
	return &PolicyRepository{db: db}
}

func (r *PolicyRepository) FindAll() ([]models.GroupPolicy, error) {
	var policies []models.GroupPolicy
	err := r.db.Preload("MasterTableAccess").Find(&policies).Error
	return policies, err
}

func (r *PolicyRepository) FindByID(id uint) (*models.GroupPolicy, error) {
	var policy models.GroupPolicy
	err := r.db.Preload("MasterTableAccess").First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *PolicyRepository) Create(policy *models.GroupPolicy) error {
	return r.db.Create(policy).Error
}

func (r *PolicyRepository) Update(policy *models.GroupPolicy) error {
	return r.db.Save(policy).Error
}

func (r *PolicyRepository) Delete(id uint) error {
	return r.db.Delete(&models.GroupPolicy{}, id).Error
}

func (r *PolicyRepository) FindByGroupID(groupID uint) ([]models.GroupPolicy, error) {
	var policies []models.GroupPolicy
	err := r.db.Preload("MasterTableAccess").Where("group_id = ?", groupID).Find(&policies).Error
	return policies, err
}

func (r *PolicyRepository) FindByGroupIDAndTableAccessID(groupID uint, tableAccessID uint) (*models.GroupPolicy, error) {
	var policies []models.GroupPolicy
	err := r.db.Where("group_id = ? AND table_id = ?", groupID, tableAccessID).Find(&policies).Error
	if err != nil {
		return nil, err
	}
	if len(policies) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &policies[0], nil
}

func (r *PolicyRepository) DeleteByGroupID(groupID uint) error {
	return r.db.Where("group_id = ?", groupID).Delete(&models.GroupPolicy{}).Error
}

func (r *PolicyRepository) BulkCreate(policies []models.GroupPolicy) error {
	if len(policies) == 0 {
		return nil
	}
	return r.db.CreateInBatches(policies, 100).Error
}

func (r *PolicyRepository) FindByGroupIDWithTree(groupID uint) ([]models.MenuItem, error) {
	//return r.FindAllForTree(groupID)
	return nil, nil
}
