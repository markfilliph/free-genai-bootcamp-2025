package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	// Custom validation tags
	customValidators = map[string]validator.Func{
		"name": validateName,
		"word": validateWord,
	}
)

// ValidationMiddleware adds custom validators and handles validation errors
func ValidationMiddleware() gin.HandlerFunc {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators
		for tag, fn := range customValidators {
			_ = v.RegisterValidation(tag, fn)
		}
	}

	return func(c *gin.Context) {
		c.Next()
	}
}

// ValidatePathID validates path parameters that should be positive integers
func ValidatePathID(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(paramName)
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Missing required path parameter: %s", paramName),
			})
			return
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil || idInt <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid %s: must be a positive integer", paramName),
			})
			return
		}

		c.Next()
	}
}

// ValidatePagination validates pagination query parameters
func ValidatePagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")

		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page: must be a positive integer",
			})
			return
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page_size: must be between 1 and 100",
			})
			return
		}

		c.Next()
	}
}

// ValidateJSON validates the request body against the provided struct
func ValidateJSON(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			var errorMessages []string
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, e := range ve {
					errorMessages = append(errorMessages, formatValidationError(e))
				}
			} else {
				errorMessages = append(errorMessages, "Invalid request body")
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errorMessages,
			})
			return
		}
		c.Next()
	}
}

// Custom validators
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	// Name should be 2-50 characters, alphanumeric with spaces and hyphens
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9\\s-]+$", name)
	return match
}

func validateWord(fl validator.FieldLevel) bool {
	word := fl.Field().String()
	// Word should be 1-100 characters, no special characters except spaces and hyphens
	if len(word) < 1 || len(word) > 100 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9\\s-]+$", word)
	return match
}

// Helper functions
func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "name":
		return fmt.Sprintf("%s must be 2-50 characters long and contain only letters, numbers, spaces, and hyphens", field)
	case "word":
		return fmt.Sprintf("%s must be 1-100 characters long and contain only letters, numbers, spaces, and hyphens", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
