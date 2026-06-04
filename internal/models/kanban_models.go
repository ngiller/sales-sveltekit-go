package models

import "time"

type KanbanBoard struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Background  *string   `gorm:"size:255" json:"background,omitempty"`
	IsArchived  bool      `gorm:"default:false" json:"is_archived"`
	PropertyID  *uint     `json:"property_id,omitempty"`
	UserCreated *uint     `json:"user_created,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Lists       []KanbanList `gorm:"foreignKey:BoardID" json:"lists,omitempty"`
}

func (KanbanBoard) TableName() string { return "kanban_boards" }

type KanbanList struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	BoardID  uint   `gorm:"not null;index" json:"board_id"`
	Name     string `gorm:"size:255;not null" json:"name"`
	Position int    `gorm:"default:0" json:"position"`
	Color    *string `gorm:"size:50" json:"color,omitempty"`
	Cards    []KanbanCard `gorm:"foreignKey:ListID" json:"cards,omitempty"`
}

func (KanbanList) TableName() string { return "kanban_lists" }

type KanbanCard struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ListID      uint      `gorm:"not null;index" json:"list_id"`
	BoardID     uint      `gorm:"not null;index" json:"board_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Position    int       `gorm:"default:0" json:"position"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	CoverImage  *string   `gorm:"size:255" json:"cover_image,omitempty"`
	IsArchived  bool      `gorm:"default:false" json:"is_archived"`
	UserCreated *uint     `json:"user_created,omitempty"`
	UserUpdated *uint     `json:"user_updated,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Members     []User             `gorm:"many2many:kanban_card_members;" json:"members,omitempty"`
	Labels      []KanbanLabel      `gorm:"many2many:kanban_card_labels;" json:"labels,omitempty"`
	Checklists  []KanbanChecklist  `gorm:"foreignKey:CardID" json:"checklists,omitempty"`
	Attachments []KanbanAttachment `gorm:"foreignKey:CardID" json:"attachments,omitempty"`
	Comments    []KanbanComment    `gorm:"foreignKey:CardID" json:"comments,omitempty"`
}

func (KanbanCard) TableName() string { return "kanban_cards" }

type KanbanLabel struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	BoardID uint   `gorm:"not null;index" json:"board_id"`
	Name    string `gorm:"size:100;not null" json:"name"`
	Color   string `gorm:"size:50;not null" json:"color"`
}

func (KanbanLabel) TableName() string { return "kanban_labels" }

type KanbanChecklist struct {
	ID       uint                  `gorm:"primaryKey" json:"id"`
	CardID   uint                  `gorm:"not null;index" json:"card_id"`
	Name     string                `gorm:"size:255;not null" json:"name"`
	Position int                   `gorm:"default:0" json:"position"`
	Items    []KanbanChecklistItem `gorm:"foreignKey:ChecklistID" json:"items,omitempty"`
}

func (KanbanChecklist) TableName() string { return "kanban_checklists" }

type KanbanChecklistItem struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ChecklistID uint   `gorm:"not null;index" json:"checklist_id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	IsChecked   bool   `gorm:"default:false" json:"is_checked"`
	Position    int    `gorm:"default:0" json:"position"`
	AssigneeID  *uint  `json:"assignee_id,omitempty"`
}

func (KanbanChecklistItem) TableName() string { return "kanban_checklist_items" }

type KanbanAttachment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CardID     uint      `gorm:"not null;index" json:"card_id"`
	FileName   string    `gorm:"size:255;not null" json:"file_name"`
	FilePath   string    `gorm:"size:255;not null" json:"file_path"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `gorm:"size:100" json:"mime_type"`
	UploadedBy *uint     `json:"uploaded_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (KanbanAttachment) TableName() string { return "kanban_attachments" }

type KanbanComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null;index" json:"card_id"`
	UserID    *uint     `json:"user_id,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (KanbanComment) TableName() string { return "kanban_comments" }
