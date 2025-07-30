package model

const (
	PaginationDefaultLimit = 20
	PaginationDefaultPage  = 1
	PaginationMaxLimit     = 1000
)

// Pagination represents for a response struct
type Pagination struct {
	TotalCount  int64 `json:"totalCount"`
	CurrentPage int   `json:"currentPage"`
	Limit       int   `json:"limit"`
}
