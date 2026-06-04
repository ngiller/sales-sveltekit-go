package models

import "time"

type Brand struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	BrandName   string    `gorm:"column:brand_name" json:"brand_name"`
	UserCreated string    `gorm:"column:user_created" json:"user_created"`
	UserUpdated string    `gorm:"column:user_updated" json:"user_updated"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Brand) TableName() string {
	return "brand"
}

type ProductCategory struct {
	ID                  uint      `gorm:"primaryKey;column:id" json:"id"`
	ProductCategoryName string    `gorm:"column:product_category_name" json:"product_category_name"`
	UserCreated         string    `gorm:"column:user_created" json:"user_created"`
	UserUpdated         string    `gorm:"column:user_updated" json:"user_updated"`
	CreatedAt           time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}

type Product struct {
	ID           uint             `gorm:"primaryKey;column:id" json:"id"`
	ProductName  string           `gorm:"column:product_name" json:"product_name"`
	ProductCode  string           `gorm:"column:product_code" json:"product_code"`
	Type         string           `gorm:"column:type" json:"type"`
	CategoryID   uint             `gorm:"column:category_id" json:"category_id"`
	VendorID     *uint            `gorm:"column:vendor_id" json:"vendor_id"`
	Descriptions string           `gorm:"column:descriptions" json:"descriptions"`
	UserCreated  string           `gorm:"column:user_created" json:"user_created"`
	UserUpdated  string           `gorm:"column:user_updated" json:"user_updated"`
	CreatedAt    time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	BrandID      uint             `gorm:"column:brand_id" json:"brand_id"`
	Brand        Brand            `gorm:"foreignKey:BrandID" json:"brand"`
	Category     ProductCategory  `gorm:"foreignKey:CategoryID" json:"category"`
	Stocks       []ProductStock   `gorm:"foreignKey:ProductID" json:"stocks"`
}

func (Product) TableName() string {
	return "product"
}

type ProductStock struct {
	ID               uint       `gorm:"primaryKey;column:id" json:"id"`
	ProductID        uint       `gorm:"column:product_id" json:"product_id"`
	UnitID           uint       `gorm:"column:unit_id" json:"unit_id"`
	StartStock       int        `gorm:"column:start_stock" json:"start_stock"`
	EndStock         int        `gorm:"column:end_stock" json:"end_stock"`
	MaxStock         int        `gorm:"column:max_stock" json:"max_stock"`
	MinStock         int        `gorm:"column:min_stock" json:"min_stock"`
	StartPrice       float64    `gorm:"column:start_price" json:"start_price"`
	LastBuyPrice     string     `gorm:"column:last_buy_price" json:"last_buy_price"`
	Enable           string     `gorm:"column:enable" json:"enable"`
	AllowBuy         string     `gorm:"column:allow_buy" json:"allow_buy"`
	AllowSell        string     `gorm:"column:allow_sell" json:"allow_sell"`
	UserCreated      string     `gorm:"column:user_created" json:"user_created"`
	UserUpdated      string     `gorm:"column:user_updated" json:"user_updated"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	LastBuyDate      *time.Time `gorm:"column:last_buy_date" json:"last_buy_date"`
	LastSellingPrice string     `gorm:"column:last_selling_price" json:"last_selling_price"`
}

func (ProductStock) TableName() string {
	return "product_stock"
}
