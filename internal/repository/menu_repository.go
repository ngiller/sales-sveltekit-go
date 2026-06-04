package repository

import (
	"backend/internal/models"
	"sort"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMenusByUserGroupID(userGroupID uint) ([]models.MenuItem, error) {
	var policies []models.GroupPolicy

	// Ensure the relation is loaded correctly using the new string-based join
	err := r.db.Preload("MasterTableAccess").
		Where("group_id = ?", userGroupID).
		Find(&policies).Error

	if err != nil {
		return nil, err
	}

	return buildMenuTree(policies), nil
}

func buildMenuTree(policies []models.GroupPolicy) []models.MenuItem {
	var nodes []models.MenuItem
	
	type menuData struct {
		Access models.MasterTableAccess
		Read   bool
		Write  bool
		Edit   bool
		Delete bool
	}

	items := make(map[uint]*menuData)
	for _, p := range policies {
		if p.MasterTableAccess.IsActive {
			if _, ok := items[p.MasterTableAccess.ID]; !ok {
				items[p.MasterTableAccess.ID] = &menuData{
					Access: p.MasterTableAccess,
				}
			}
			
			switch p.Action {
			case "read":
				items[p.MasterTableAccess.ID].Read = true
			case "create":
				items[p.MasterTableAccess.ID].Write = true
			case "update":
				items[p.MasterTableAccess.ID].Edit = true
			case "delete":
				items[p.MasterTableAccess.ID].Delete = true
			}
		}
	}

	var findChildren func(parentID uint) []models.MenuItem
	findChildren = func(parentID uint) []models.MenuItem {
		var children []models.MenuItem
		for _, v := range items {
			if v.Access.ParentID != nil && *v.Access.ParentID == parentID {
				mName := v.Access.MenuName
				if mName == "" { mName = v.Access.Name }

				child := models.MenuItem{
					ID:        v.Access.ID,
					Name:      v.Access.Name,
					ParentID:  v.Access.ParentID,
					MenuName:  mName,
					Path:      v.Access.Path,
					Icon:      v.Access.Icon,
					SortOrder: v.Access.SortOrder,
					IsActive:  v.Access.IsActive,
					CanRead:   v.Read,
					CanWrite:  v.Write,
					CanEdit:   v.Edit,
					CanDelete: v.Delete,
					Children:  findChildren(v.Access.ID),
				}
				children = append(children, child)
			}
		}
		
		sort.Slice(children, func(i, j int) bool {
			return items[children[i].ID].Access.SortOrder < items[children[j].ID].Access.SortOrder
		})
		
		if len(children) == 0 {
			return []models.MenuItem{} // Return empty array instead of nil for JSON consistency
		}
		return children
	}

	for _, v := range items {
		if v.Access.ParentID == nil || *v.Access.ParentID == 0 {
			mName := v.Access.MenuName
			if mName == "" { mName = v.Access.Name }

			root := models.MenuItem{
				ID:        v.Access.ID,
				Name:      v.Access.Name,
				ParentID:  v.Access.ParentID,
				MenuName:  mName,
				Path:      v.Access.Path,
				Icon:      v.Access.Icon,
				SortOrder: v.Access.SortOrder,
				IsActive:  v.Access.IsActive,
				CanRead:   v.Read,
				CanWrite:  v.Write,
				CanEdit:   v.Edit,
				CanDelete: v.Delete,
				Children:  findChildren(v.Access.ID),
			}
			nodes = append(nodes, root)
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return items[nodes[i].ID].Access.SortOrder < items[nodes[j].ID].Access.SortOrder
	})

	if len(nodes) == 0 {
		return []models.MenuItem{}
	}

	return nodes
}
