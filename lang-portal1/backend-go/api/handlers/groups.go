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

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		group := &models.Group{
			Name:        req.Name,
			Description: req.Description,
		}

		if err := models.CreateGroup(db, group); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create group")
			return
		}

		c.JSON(http.StatusCreated, group)
	}
}

func UpdateGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		var req CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		group := &models.Group{
			ID:   id,
			Name: req.Name,
		}

		if err := models.UpdateGroup(db, group); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to update group")
			return
		}

		c.JSON(http.StatusOK, group)
	}
}

func DeleteGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		if err := models.DeleteGroup(db, id); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete group")
			return
		}

		c.Status(http.StatusOK)
	}
}

func AddWordToGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupIDStr := c.Param("id")
		wordIDStr := c.Param("wordId")

		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		wordID, err := strconv.Atoi(wordIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		if err := models.AddWordToGroup(db, groupID, wordID); err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Group or word not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to add word to group")
			return
		}

		c.Status(http.StatusOK)
	}
}

func RemoveWordFromGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupIDStr := c.Param("id")
		wordIDStr := c.Param("wordId")

		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		wordID, err := strconv.Atoi(wordIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		if err := models.RemoveWordFromGroup(db, groupID, wordID); err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Group or word not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to remove word from group")
			return
		}

		c.Status(http.StatusOK)
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
