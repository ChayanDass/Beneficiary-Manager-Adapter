package models

// PaginationMeta contains metadata for paginated responses
type PaginationMeta struct {
	ResourceCount int    `json:"resource_count" example:"200"`
	TotalPages    int64  `json:"total_pages,omitempty" example:"20"`
	Page          int64  `json:"page,omitempty" example:"10"`
	Limit         int64  `json:"limit,omitempty" example:"10"`
	Next          string `json:"next,omitempty" example:"/api/v1/schemes?limit=10&page=11"`
	Previous      string `json:"previous,omitempty" example:"/api/v1/schemes?limit=10&page=9"`
}

// PaginationInput is the input model for pagination
type PaginationInput struct {
	Page  int64 `json:"page" example:"10"`
	Limit int64 `json:"limit" example:"10"`
}

// PaginationParse defines behavior for pagination inputs
type PaginationParse interface {
	GetOffset() int64
	GetLimit() int64
}

// GetOffset returns offset value
func (p PaginationInput) GetOffset() int64 {
	return (p.Page - 1) * p.Limit
}

// GetLimit returns limit value
func (p PaginationInput) GetLimit() int64 {
	return p.Limit
}

// ------------------ API Responses ------------------

// ErrorResponse for API error output
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse for success output
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SchemeResponse for paginated schemes
type SchemeResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}
