package repository

import (
	"testing"

	"backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildMenuTreeEmpty(t *testing.T) {
	tree := buildMenuTree([]models.GroupPolicy{})
	assert.Empty(t, tree)
	assert.IsType(t, []models.MenuItem{}, tree)
}

func TestBuildMenuTreeSingleRoot(t *testing.T) {
	policies := []models.GroupPolicy{
		{
			GroupID: 1,
			Action:  "read",
			TableID: 1,
			MasterTableAccess: models.MasterTableAccess{
				ID:        1,
				Name:      "Dashboard",
				MenuName:  "Dashboard",
				ParentID:  nil,
				Path:      toStrPtr("/dashboard"),
				Icon:      toStrPtr("home"),
				SortOrder: 1,
				IsActive:  true,
			},
		},
	}

	tree := buildMenuTree(policies)
	require.Len(t, tree, 1)
	assert.Equal(t, uint(1), tree[0].ID)
	assert.Equal(t, "Dashboard", tree[0].Name)
	assert.Equal(t, "/dashboard", *tree[0].Path)
	assert.Equal(t, "home", *tree[0].Icon)
	assert.True(t, tree[0].CanRead)
	assert.False(t, tree[0].CanWrite)
	assert.False(t, tree[0].CanEdit)
	assert.False(t, tree[0].CanDelete)
	assert.True(t, tree[0].IsActive)
	assert.Empty(t, tree[0].Children)
}

func TestBuildMenuTreeWithChildren(t *testing.T) {
	policies := []models.GroupPolicy{
		{
			GroupID: 1,
			Action:  "read",
			TableID: 1,
			MasterTableAccess: models.MasterTableAccess{
				ID:        1,
				Name:      "Master Data",
				MenuName:  "Master Data",
				ParentID:  nil,
				SortOrder: 1,
				IsActive:  true,
			},
		},
		{
			GroupID: 1,
			Action:  "read",
			TableID: 2,
			MasterTableAccess: models.MasterTableAccess{
				ID:        2,
				Name:      "Users",
				MenuName:  "Users",
				ParentID:  toUintPtr(1),
				SortOrder: 1,
				IsActive:  true,
			},
		},
		{
			GroupID: 1,
			Action:  "create",
			TableID: 2,
			MasterTableAccess: models.MasterTableAccess{
				ID:        2,
				Name:      "Users",
				MenuName:  "Users",
				ParentID:  toUintPtr(1),
				SortOrder: 1,
				IsActive:  true,
			},
		},
		{
			GroupID: 1,
			Action:  "read",
			TableID: 3,
			MasterTableAccess: models.MasterTableAccess{
				ID:        3,
				Name:      "Roles",
				MenuName:  "Roles",
				ParentID:  toUintPtr(1),
				SortOrder: 2,
				IsActive:  true,
			},
		},
	}

	tree := buildMenuTree(policies)

	require.Len(t, tree, 1)
	assert.Equal(t, "Master Data", tree[0].Name)
	require.Len(t, tree[0].Children, 2)

	// Children should be sorted by SortOrder
	assert.Equal(t, "Users", tree[0].Children[0].Name)
	assert.Equal(t, "Roles", tree[0].Children[1].Name)

	// Check Users permissions
	users := tree[0].Children[0]
	assert.True(t, users.CanRead)
	assert.True(t, users.CanWrite)  // create maps to Write
	assert.False(t, users.CanEdit)
	assert.False(t, users.CanDelete)

	// Check Roles permissions
	roles := tree[0].Children[1]
	assert.True(t, roles.CanRead)
	assert.False(t, roles.CanWrite)
	assert.False(t, roles.CanEdit)
	assert.False(t, roles.CanDelete)
}

func TestBuildMenuTreeNestedMultipleLevels(t *testing.T) {
	policies := []models.GroupPolicy{
		{GroupID: 1, Action: "read", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Root", MenuName: "Root", ParentID: nil, SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 2, MasterTableAccess: models.MasterTableAccess{ID: 2, Name: "Child", MenuName: "Child", ParentID: toUintPtr(1), SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 3, MasterTableAccess: models.MasterTableAccess{ID: 3, Name: "Grandchild", MenuName: "Grandchild", ParentID: toUintPtr(2), SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 4, MasterTableAccess: models.MasterTableAccess{ID: 4, Name: "Sibling", MenuName: "Sibling", ParentID: nil, SortOrder: 2, IsActive: true}},
	}

	tree := buildMenuTree(policies)

	require.Len(t, tree, 2)
	assert.Equal(t, "Root", tree[0].Name)
	require.Len(t, tree[0].Children, 1)
	assert.Equal(t, "Child", tree[0].Children[0].Name)
	require.Len(t, tree[0].Children[0].Children, 1)
	assert.Equal(t, "Grandchild", tree[0].Children[0].Children[0].Name)
	assert.Equal(t, "Sibling", tree[1].Name)
}

func TestBuildMenuTreeInactiveMenuExcluded(t *testing.T) {
	policies := []models.GroupPolicy{
		{
			GroupID: 1,
			Action:  "read",
			TableID: 1,
			MasterTableAccess: models.MasterTableAccess{
				ID:        1,
				Name:      "Inactive Menu",
				MenuName:  "Inactive",
				ParentID:  nil,
				SortOrder: 1,
				IsActive:  false, // Inactive
			},
		},
		{
			GroupID: 1,
			Action:  "read",
			TableID: 2,
			MasterTableAccess: models.MasterTableAccess{
				ID:        2,
				Name:      "Active Menu",
				MenuName:  "Active",
				ParentID:  nil,
				SortOrder: 2,
				IsActive:  true,
			},
		},
	}

	tree := buildMenuTree(policies)
	require.Len(t, tree, 1)
	assert.Equal(t, "Active Menu", tree[0].Name)
}

func TestBuildMenuTreeDuplicatePolicies(t *testing.T) {
	policies := []models.GroupPolicy{
		{GroupID: 1, Action: "read", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Menu", MenuName: "Menu", ParentID: nil, SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Menu", MenuName: "Menu", ParentID: nil, SortOrder: 1, IsActive: true}}, // Duplicate
		{GroupID: 1, Action: "create", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Menu", MenuName: "Menu", ParentID: nil, SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "update", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Menu", MenuName: "Menu", ParentID: nil, SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "delete", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Menu", MenuName: "Menu", ParentID: nil, SortOrder: 1, IsActive: true}},
	}

	tree := buildMenuTree(policies)
	require.Len(t, tree, 1)
	assert.True(t, tree[0].CanRead)
	assert.True(t, tree[0].CanWrite)
	assert.True(t, tree[0].CanEdit)
	assert.True(t, tree[0].CanDelete)
}

func TestBuildMenuTreeMultipleRoots(t *testing.T) {
	policies := []models.GroupPolicy{
		{GroupID: 1, Action: "read", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "B Menu", MenuName: "B", ParentID: nil, SortOrder: 2, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 2, MasterTableAccess: models.MasterTableAccess{ID: 2, Name: "A Menu", MenuName: "A", ParentID: nil, SortOrder: 1, IsActive: true}},
		{GroupID: 1, Action: "read", TableID: 3, MasterTableAccess: models.MasterTableAccess{ID: 3, Name: "C Menu", MenuName: "C", ParentID: nil, SortOrder: 3, IsActive: true}},
	}

	tree := buildMenuTree(policies)
	require.Len(t, tree, 3)

	// Should be sorted by SortOrder
	assert.Equal(t, "A Menu", tree[0].Name)
	assert.Equal(t, "B Menu", tree[1].Name)
	assert.Equal(t, "C Menu", tree[2].Name)
}

func TestBuildMenuTreeEmptyNameFallback(t *testing.T) {
	policies := []models.GroupPolicy{
		{GroupID: 1, Action: "read", TableID: 1, MasterTableAccess: models.MasterTableAccess{ID: 1, Name: "Fallback Name", MenuName: "", ParentID: nil, SortOrder: 1, IsActive: true}},
	}

	tree := buildMenuTree(policies)
	require.Len(t, tree, 1)
	// When MenuName is empty, it falls back to Name
	assert.Equal(t, "Fallback Name", tree[0].MenuName)
}

func toStrPtr(s string) *string { return &s }
func toUintPtr(u uint) *uint   { return &u }
