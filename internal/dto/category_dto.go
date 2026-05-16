package dtos

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryResponse struct {
	ID  string `json:"id"`
	Name string `json:"name"`
}
