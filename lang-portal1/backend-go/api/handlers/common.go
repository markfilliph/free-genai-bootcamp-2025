package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginatedResponse struct {
	Items         interface{} `json:"items"`
	CurrentPage   int        `json:"current_page"`
	TotalPages    int        `json:"total_pages"`
	TotalItems    int        `json:"total_items"`
	ItemsPerPage  int        `json:"items_per_page"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

const defaultPerPage = 100

func getPaginationParams(c *gin.Context) (page, perPage int) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", strconv.Itoa(defaultPerPage))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err = strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = defaultPerPage
	}

	return page, perPage
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: message})
}

func newPaginatedResponse(items interface{}, currentPage, totalItems, perPage int) PaginatedResponse {
	totalPages := (totalItems + perPage - 1) / perPage
	return PaginatedResponse{
		Items:         items,
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		TotalItems:    totalItems,
		ItemsPerPage:  perPage,
	}
}
