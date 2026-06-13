package models

import (
	"backend/config"
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Email         string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PhoneNo       *string   `gorm:"size:30" json:"phone_no"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	UserGroupID   *uint     `json:"user_group_id"`
	UserGroup     *UserGroup `gorm:"foreignKey:UserGroupID" json:"user_group"`
	DepartementID *uint     `json:"departement_id"`
	Departement   *Departement `gorm:"foreignKey:DepartementID" json:"departement"`
	Avatar        *string   `gorm:"size:255" json:"avatar"`
	Sign          *string   `gorm:"size:255" json:"sign"`
	PropertyID    *int64    `json:"property_id"`
	Enable        bool      `gorm:"default:true" json:"enable"`
	Inisial       *string   `gorm:"size:5" json:"inisial"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserCreated   int64     `gorm:"not null" json:"user_created"`
	UserUpdate    int64     `gorm:"not null" json:"user_update"`
	RoleName      string    `gorm:"->;column:role_name" json:"role_name"`
	DeptName      string    `gorm:"->;column:dept_name" json:"dept_name"`
}

func (User) TableName() string {
	return "users"
}

type Property struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Code   string `gorm:"size:50;not null" json:"code"`
	Name   string `gorm:"size:255;not null" json:"name"`
	Enable bool   `gorm:"default:true" json:"enable"`
}

func (Property) TableName() string {
	return "master_properties"
}

type UserGroup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	PropertyID  *uint     `json:"property_id"`
	UserCreated *uint     `json:"user_created"`
	UserUpdate  *uint     `json:"user_update"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

type Departement struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	PropertyID  *uint     `json:"property_id"`
	UserCreated *uint     `json:"user_created"`
	UserUpdate  *uint     `json:"user_update"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Departement) TableName() string {
	return "master_departements"
}

type CustomerCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null;unique" json:"name"`
	PropertyID  *uint     `json:"property_id"`
	UserCreated *uint     `json:"user_created"`
	UserUpdate  *uint     `json:"user_update"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CustomerCategory) TableName() string {
	return "customer_category"
}

type Customer struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CategoryID    uint      `gorm:"not null" json:"category_id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Address       *string   `gorm:"type:text" json:"address"`
	Phone         *string   `gorm:"size:255" json:"phone"`
	Email         *string   `gorm:"size:255" json:"email"`
	NPWP          *string   `gorm:"size:255" json:"npwp"`
	Contact       *string   `gorm:"size:255" json:"contact"`
	Enable        bool      `gorm:"not null" json:"enable"`
	Notes         *string   `gorm:"type:text" json:"notes"`
	SalesID       *uint     `json:"sales_id"`
	PropertyID    uint      `gorm:"not null" json:"property_id"`
	UserCreated   *uint     `json:"user_created"`
	UserUpdate    *uint     `json:"user_update"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UsersID       uint      `gorm:"not null" json:"users_id"`
	AllowToVendor string    `gorm:"size:255;not null;default:'1'" json:"allow_to_vendor"`
	AccID         *uint     `json:"acc_id"`
	BeginDate     *time.Time `json:"begin_date"`
	Balance       *float64  `json:"balance"`
	Contacts      []CustomerContact `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"contacts"`
}

func (Customer) TableName() string {
	return "customer"
}

type CustomerContact struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	CustomerID uint    `gorm:"not null" json:"customer_id"`
	Name       string  `gorm:"size:255;not null" json:"name"`
	Phone      *string `gorm:"size:255" json:"phone"`
	Email      *string `gorm:"size:255" json:"email"`
	Position   *string `gorm:"size:255" json:"position"`
}

func (CustomerContact) TableName() string {
	return "customer_contact"
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type MasterTableAccess struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"size:255;not null" json:"name"`
	ParentID  *uint   `json:"parent_id"` // null jika menu utama
	MenuName  string  `gorm:"size:255;not null" json:"menu_name"`
	Path      *string `gorm:"size:255" json:"path"`
	Endpoint  *string `gorm:"size:255" json:"endpoint"`
	Icon      *string `gorm:"size:255" json:"icon"`
	SortOrder int     `gorm:"default:0" json:"sort_order"`
	IsActive  bool    `gorm:"default:true" json:"is_active"`
}

func (MasterTableAccess) TableName() string {
	return "menu_access"
}

type UserGroupPolicy struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	GroupID    uint      `gorm:"not null" json:"group_id"`
	TargetTableName string    `gorm:"column:table_name;size:255;not null" json:"table_name"`
	TableID    uint      `gorm:"column:table_id;not null" json:"table_id"`
	Action     string    `gorm:"size:255;not null" json:"action"`
	PropertyID *uint     `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	MasterTableAccess MasterTableAccess `gorm:"foreignKey:TableID;references:ID" json:"table_access"`
}

func (UserGroupPolicy) TableName() string {
	return "user_group_policies"
}

type GroupPolicy struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	GroupID    uint      `gorm:"not null" json:"group_id"`
	TargetTableName string    `gorm:"column:table_name;size:255;not null" json:"table_name"`
	TableID    uint      `gorm:"column:table_id;not null" json:"table_id"`
	Action     string    `gorm:"size:255;not null" json:"action"`
	PropertyID *uint     `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	MasterTableAccess MasterTableAccess `gorm:"foreignKey:TableID;references:ID" json:"table_access"`
}

func (GroupPolicy) TableName() string {
	return "group_policies"
}

type MenuItem struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uint      `json:"parent_id"`
	MenuName  string     `json:"menu_name"`
	Path      *string    `json:"path"`
	Icon      *string    `json:"icon"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
	CanRead   bool       `json:"can_read"`
	CanWrite  bool       `json:"can_write"`
	CanEdit   bool       `json:"can_edit"`
	CanDelete bool       `json:"can_delete"`
	Children  []MenuItem `json:"children"`
}

type PaymentTerm struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Day         int       `json:"day"`
	PropertyID  *uint     `json:"property_id"`
	UserCreated *uint     `json:"user_created"`
	UserUpdate  *uint     `json:"user_update"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PaymentTerm) TableName() string {
	return "payment_term"
}

type ProjectLevel struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

func (ProjectLevel) TableName() string {
	return "project_level"
}

type ProjectPriority struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

func (ProjectPriority) TableName() string {
	return "project_priority"
}

type QuotationProgress struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string  `gorm:"size:255;not null" json:"name"`
	Progress float64 `gorm:"not null" json:"progress"`
}

func (QuotationProgress) TableName() string {
	return "quotation_progress"
}

type QuotationStatus struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

func (QuotationStatus) TableName() string {
	return "quotation_status"
}

type Unit struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string    `gorm:"column:units_name;size:255;not null" json:"name"`
	PropertyID  *uint     `gorm:"-" json:"property_id,omitempty"`
	UserCreated *string   `gorm:"column:user_created" json:"user_created"`
	UserUpdate  *string   `gorm:"column:user_updated" json:"user_update"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updated_at"`
}

func (Unit) TableName() string {
	if config.AppConfig != nil && config.AppConfig.DBStockName != "" {
		return config.AppConfig.DBStockName + ".units"
	}
	return "units"
}

type Quotation struct {
	ID              string             `gorm:"primaryKey;column:id;size:20" json:"id"`
	PropertyID      int               `gorm:"column:property_id;not null" json:"property_id"`
	QuotationType   int               `gorm:"column:quotation_type;not null;default:1" json:"quotation_type"`
	QuotationID     string            `gorm:"column:quotation_id;size:255" json:"quotation_id"`
	QuotationDate   *time.Time        `gorm:"column:quotation_date" json:"quotation_date"`
	CustomerID     uint              `gorm:"column:customer_id;not null" json:"customer_id"`
	Customer       Customer          `gorm:"foreignKey:CustomerID" json:"customer"`
	ContactID      *uint             `gorm:"column:contact_id" json:"contact_id"`
	Contact        *CustomerContact  `gorm:"foreignKey:ContactID" json:"contact"`
	Subject        *string          `gorm:"column:subject" json:"subject"`
	Total          *float64         `gorm:"column:total" json:"total"`
	Tax            *float64         `gorm:"column:tax" json:"tax"`
	TaxValue       float64          `gorm:"column:tax_value;default:0" json:"tax_value"`
	PPh            float64          `gorm:"column:pph;default:0" json:"pph"`
	PPhValue       float64          `gorm:"column:pph_value;default:0" json:"pph_value"`
	Disc           *float64         `gorm:"column:disc" json:"disc"`
	GrandTotal     *float64         `gorm:"column:grand_total" json:"grand_total"`
	HppTotal       *float64         `gorm:"column:hpp_total" json:"hpp_total"`
	Profit        *float64         `gorm:"column:profit" json:"profit"`
	ProfitValue   *float64         `gorm:"column:profit_value" json:"profit_value"`
	PaymentTermID *uint            `gorm:"column:payment_term" json:"payment_term_id"`
	PaymentTerm   *PaymentTerm     `gorm:"foreignKey:PaymentTermID" json:"payment_term"`
	ValidUntil    *time.Time       `gorm:"column:valid_until" json:"valid_until"`
	Commision     float64          `gorm:"column:commision" json:"commision"`
	Notes         *string         `gorm:"column:notes;type:text" json:"notes"`
	Status        int              `gorm:"column:status;not null;default:1" json:"status"`
	StatusInfo    *QuotationStatus `gorm:"foreignKey:Status" json:"status_info"`
	StatusReview  *string          `gorm:"column:status_review;type:text" json:"status_review"`
	ProgressID    uint              `gorm:"column:progress;not null;default:2" json:"progress_id"`
	ProgressInfo  *QuotationProgress `gorm:"foreignKey:ProgressID;references:ID" json:"progress_info"`
	FollowupBy    *uint            `gorm:"column:followup_by" json:"followup_by"`
	FollowupDate *time.Time       `gorm:"column:followup_date" json:"followup_date"`
	NextFollowup  *time.Time       `gorm:"column:next_followup" json:"next_followup"`
	Folder        *string         `gorm:"column:folder" json:"folder"`
	SalesID       *uint            `gorm:"column:sales_id" json:"sales_id"`
	SalesPerson   *User            `gorm:"foreignKey:SalesID" json:"sales_person"`
	UserCreated   *uint            `gorm:"column:user_created" json:"user_created"`
	UserUpdate   *uint            `gorm:"column:user_update" json:"user_update"`
	CreatedAt    time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"column:updated_at" json:"updated_at"`
	Notif        bool             `gorm:"column:notif;default:false" json:"notif"`
	LevelID       *uint            `gorm:"column:level" json:"level_id"`
	Level         *ProjectLevel    `gorm:"foreignKey:LevelID" json:"level"`
	PriorityID   *uint            `gorm:"column:priority" json:"priority_id"`
	Priority     *ProjectPriority `gorm:"foreignKey:PriorityID" json:"priority"`
	ProjectStart *time.Time       `gorm:"column:project_start" json:"project_start"`
	ProjectEnd   *time.Time       `gorm:"column:project_end" json:"project_end"`
	PoNo          *string         `gorm:"column:po_no" json:"po_no"`
	PoDate        *time.Time       `gorm:"column:po_date" json:"po_date"`
	PoFile        *string         `gorm:"column:po_file" json:"po_file"`
	PoAssignTo    *string         `gorm:"column:po_assign_to" json:"po_assign_to"`
	Details       []QuotationDetail `gorm:"foreignKey:id;references:id" json:"details"`
	Subdetails    []QuotationSubdetail `gorm:"-" json:"subdetails"`
	Revisions     []QuotationMaster `gorm:"foreignKey:id;references:id" json:"revisions"`
}

func (Quotation) TableName() string {
	return "quotation"
}

type QuotationDetail struct {
	ID            string    `gorm:"primaryKey;column:id;size:20" json:"id"`
	RevID         int       `gorm:"primaryKey;column:rev_id;not null" json:"rev_id"`
	Line          int       `gorm:"primaryKey;column:line;not null" json:"line"`
	No            *int      `gorm:"column:no" json:"no"`
	ProductID     *uint     `gorm:"column:product_id" json:"product_id"`
	PartNo        *string   `gorm:"column:part_no" json:"part_no"`
	ProductType   int       `gorm:"column:product_type;not null;default:1" json:"product_type"`
	Description   *string   `gorm:"column:descriptions;type:text" json:"description"`
	Qty           *float64  `gorm:"column:qty" json:"qty"`
	UnitID        *uint     `gorm:"column:unit_id" json:"unit_id"`
	Unit          *Unit     `gorm:"foreignKey:UnitID" json:"unit"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	Total         float64   `gorm:"column:total;not null" json:"total"`
	OtherCost     float64   `gorm:"column:other_cost;not null;default:0" json:"other_cost"`
	Hpp           *float64  `gorm:"column:hpp" json:"hpp"`
	HppTotal      *float64  `gorm:"column:hpp_total" json:"hpp_total"`
}

func (QuotationDetail) TableName() string {
	return "quotation_detail"
}

type QuotationSubdetail struct {
	ID            string    `gorm:"primaryKey;column:id;size:20" json:"id"`
	RevID         int       `gorm:"primaryKey;column:rev_id;not null" json:"rev_id"`
	Line          int       `gorm:"primaryKey;column:line;not null" json:"line"`
	Subline       int       `gorm:"primaryKey;column:subline;not null" json:"subline"`
	No            *int      `gorm:"column:no" json:"no"`
	ProductID     *uint     `gorm:"column:product_id" json:"product_id"`
	PartNo        *string   `gorm:"column:part_no" json:"part_no"`
	ProductType   int       `gorm:"column:product_type;not null;default:1" json:"product_type"`
	Description   *string   `gorm:"column:descriptions;type:text" json:"description"`
	Qty           *float64  `gorm:"column:qty" json:"qty"`
	UnitID        *uint     `gorm:"column:unit_id" json:"unit_id"`
	Unit          *Unit     `gorm:"foreignKey:UnitID" json:"unit"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	Total         float64   `gorm:"column:total;not null" json:"total"`
	OtherCost     float64   `gorm:"column:other_cost;not null;default:0" json:"other_cost"`
	Hpp           *float64  `gorm:"column:hpp" json:"hpp"`
	HppTotal      *float64  `gorm:"column:hpp_total" json:"hpp_total"`
}

func (QuotationSubdetail) TableName() string {
	return "quotation_subdetail"
}

type QuotationMaster struct {
	ID            string     `gorm:"primaryKey;column:id;size:20" json:"id"`
	RevID         int        `gorm:"primaryKey;column:rev_id;not null" json:"rev_id"`
	QuotationDate *time.Time  `gorm:"column:quotation_date" json:"quotation_date"`
	Subject       *string    `gorm:"column:subject" json:"subject"`
	Total         *float64   `gorm:"column:total" json:"total"`
	Disc          *float64   `gorm:"column:disc" json:"disc"`
	Tax           *float64   `gorm:"column:tax" json:"tax"`
	TaxValue      float64    `gorm:"column:tax_value;default:0" json:"tax_value"`
	PPh           float64    `gorm:"column:pph;default:0" json:"pph"`
	PPhValue      float64    `gorm:"column:pph_value;default:0" json:"pph_value"`
	GrandTotal    *float64   `gorm:"column:grand_total" json:"grand_total"`
	HppTotal      *float64   `gorm:"column:hpp_total" json:"hpp_total"`
	Profit        float64    `gorm:"column:profit;not null;default:0" json:"profit"`
	ProfitValue  *float64   `gorm:"column:profit_value" json:"profit_value"`
	PaymentTermID *uint      `gorm:"column:payment_term" json:"payment_term_id"`
	ValidUntil    *time.Time `gorm:"column:valid_until" json:"valid_until"`
	Commision     float64   `gorm:"column:commision" json:"commision"`
	Notes        *string    `gorm:"column:notes;type:text" json:"notes"`
	UserCreated   *uint      `gorm:"column:user_created" json:"user_created"`
	SalesID       *uint      `gorm:"column:sales_id" json:"sales_id"`
	UserUpdate   *uint      `gorm:"column:user_update" json:"user_update"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DefaultQuot  bool      `gorm:"column:default_quot;default:false" json:"default_quot"`
	LevelID       *uint      `gorm:"column:level" json:"level_id"`
	PriorityID   *uint      `gorm:"column:priority" json:"priority_id"`
	ProjectStart *time.Time `gorm:"column:project_start" json:"project_start"`
	ProjectEnd   *time.Time `gorm:"column:project_end" json:"project_end"`
	PoNo          *string    `gorm:"column:po_no" json:"po_no"`
	PoDate        *time.Time `gorm:"column:po_date" json:"po_date"`
	PoFile        *string    `gorm:"column:po_file" json:"po_file"`
	PoAssignTo    *string    `gorm:"column:po_assign_to" json:"po_assign_to"`
}

func (QuotationMaster) TableName() string {
	return "quotation_master"
}

type QuotationFollowup struct {
	ID           string     `gorm:"column:id;primaryKey;size:20;not null" json:"line_id"`
	LineID       int        `gorm:"column:line_id;primaryKey;not null" json:"id"`
	PropertyID   int        `gorm:"column:property_id;not null" json:"property_id"`
	RevID        int        `gorm:"-" json:"rev_id"`
	Status       *uint      `gorm:"column:status" json:"status"`
	Progress     *uint      `gorm:"column:progress" json:"progress"`
	FollowupDate *time.Time `gorm:"column:followup_date" json:"followup_date"`
	FollowupBy   *uint      `gorm:"column:followup_by" json:"followup_by"`
	NextFollowup *time.Time `gorm:"column:next_followup" json:"next_followup"`
	Notes        *string    `gorm:"column:notes;type:text" json:"notes"`
	PoNo         *string    `gorm:"column:po_no" json:"po_no"`
	PoDate       *time.Time `gorm:"column:po_date" json:"po_date"`
	PoFile       *string    `gorm:"column:po_file" json:"po_file"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Preloaded Relations
	FollowupByUser *User              `gorm:"foreignKey:FollowupBy" json:"followup_by_user"`
	StatusInfo     *QuotationStatus   `gorm:"foreignKey:Status" json:"status_info"`
	ProgressInfo   *QuotationProgress `gorm:"foreignKey:Progress" json:"progress_info"`
}

func (QuotationFollowup) TableName() string {
	return "quotation_followup"
}

type CounterID struct {
	PropertyID   int    `gorm:"primaryKey;column:property_id;not null" json:"property_id"`
	CounterName string `gorm:"primaryKey;column:counter_name;size:255;not null" json:"counter_name"`
	Ym          string `gorm:"primaryKey;column:ym;size:6;not null" json:"ym"`
	Type        int    `gorm:"primaryKey;column:type;not null;default:0" json:"type"`
	Counter     int    `gorm:"column:counter;not null" json:"counter"`
}

func (CounterID) TableName() string {
	return "counter_id"
}

type Setting struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	PropertyID int    `gorm:"column:property_id;not null" json:"property_id"`
	Code       int    `gorm:"column:code;not null" json:"code"`
	Name       string `gorm:"column:name;size:255;not null" json:"name"`
	Value      string `gorm:"column:value;type:text;not null" json:"value"`
}

func (Setting) TableName() string {
	return "setting"
}
