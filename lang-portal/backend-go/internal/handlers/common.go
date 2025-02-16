package handlers

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// PaginationQuery represents common pagination parameters
type PaginationQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=100"`
}

// PaginationResponse represents the standard pagination response structure
type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// respondWithError sends a JSON error response
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: message})
}

// getPaginationParams extracts and validates pagination parameters
func getPaginationParams(c *gin.Context) PaginationQuery {
	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		query.Page = 1
		query.PageSize = 100
	}
	
	// Ensure page size doesn't exceed maximum
	if query.PageSize > 100 {
		query.PageSize = 100
	}
	
	// Ensure page number is positive
	if query.Page < 1 {
		query.Page = 1
	}
	
	return query
}

// calculatePagination creates a pagination response
func calculatePagination(currentPage, pageSize, totalItems int) PaginationResponse {
	totalPages := (totalItems + pageSize - 1) / pageSize
	
	return PaginationResponse{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: pageSize,
	}
}
