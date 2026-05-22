package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category  Category  `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	SKU string    `json:"sku" gorm:"unique;not null"`
	Name string   `json:"name" gorm:"not null"`
	Price int64  `json:"price" gorm:"not null"`
	Stock int	 `json:"stock" gorm:"not null"`
	Image string `json:"image" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
