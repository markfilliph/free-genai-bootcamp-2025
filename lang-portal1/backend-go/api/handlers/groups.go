package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/models"
)

func GetGroups(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, perPage := getPaginationParams(c)

		groups, total, err := models.GetGroups(db, page, perPage)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get groups")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(groups, page, total, perPage))
	}
}

func GetGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		group, err := models.GetGroup(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Group not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get group")
			return
		}

		c.JSON(http.StatusOK, group)
	}
}

func GetGroupWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		words, err := models.GetGroupWords(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Group not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get group words")
			return
		}

		c.JSON(http.StatusOK, words)
	}
}

func GetGroupStudySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		page, perPage := getPaginationParams(c)

		sessions, total, err := models.GetStudySessions(db, id, page, perPage)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get group study sessions")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(sessions, page, total, perPage))
	}
}
