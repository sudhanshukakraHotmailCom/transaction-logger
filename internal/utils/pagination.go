package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Pagination contains pagination request data
type Pagination struct {
	Page     int
	PageSize int
}

// NewPagination creates a new pagination instance from gin context
func NewPagination(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if pageSize < 1 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset returns the offset for SQL query
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
