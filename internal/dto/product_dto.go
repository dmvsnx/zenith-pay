package dtos

type ProductRequest struct {
	CategoryID string `json:"category_id" validate:"required,uuid"`
	SKU        string `json:"sku"`
	Name       string `json:"name" validate:"required"`
	Price      int64  `json:"price" validate:"required,gt=0"`
	Stock      int    `json:"stock" validate:"required,gte=0"`
	Image      string `json:"-"`
}

type ProductResponse struct {
	ID           string `json:"id"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	Stock        int    `json:"stock"`
	Image        string `json:"image"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ProductUpdateRequest struct {
	CategoryID *string `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Name       *string `json:"name,omitempty"`
	Price      *int64  `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock      *int    `json:"stock,omitempty" validate:"omitempty,gte=0"`
	Image      *string `json:"-"`
}
