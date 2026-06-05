package domain

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p *PaginationParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
}
