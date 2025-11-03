package utils

type Pagination struct {
	Page  int `query:"page" validate:"gte=1"`
	Limit int `query:"limit" validate:"gte=1"`
}

const (
	defaultPage  int = 1
	defaultLimit int = 10
)

func NewPagination() Pagination {
	return Pagination{
		Page:  defaultPage,
		Limit: defaultLimit,
	}
}

func (p *Pagination) GettOffset() int {
	return (p.Page - 1) * p.Limit
}
