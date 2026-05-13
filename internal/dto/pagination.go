package dtos

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationRequest struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"limit" query:"limit"`
}

func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.Limit < 1 || p.Limit > MaxLimit {
		p.Limit = DefaultLimit
	}
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
