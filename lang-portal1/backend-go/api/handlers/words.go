package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/models"
)

func GetWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, perPage := getPaginationParams(c)

		words, total, err := models.GetWords(db, page, perPage)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get words")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(words, page, total, perPage))
	}
}

func GetWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		word, err := models.GetWord(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Word not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get word")
			return
		}

		c.JSON(http.StatusOK, word)
	}
}

type CreateWordRequest struct {
	Japanese string          `json:"japanese" binding:"required"`
	Romaji   string          `json:"romaji" binding:"required"`
	English  string          `json:"english" binding:"required"`
	Parts    json.RawMessage `json:"parts"`
}

func CreateWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateWordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		word := &models.Word{
			Japanese: req.Japanese,
			Romaji:   req.Romaji,
			English:  req.English,
			Parts:    req.Parts,
		}

		if err := models.CreateWord(db, word); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create word")
			return
		}

		c.JSON(http.StatusCreated, word)
	}
}

func UpdateWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		var req CreateWordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		word := &models.Word{
			ID:       id,
			Japanese: req.Japanese,
			Romaji:   req.Romaji,
			English:  req.English,
			Parts:    req.Parts,
		}

		if err := models.UpdateWord(db, word); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to update word")
			return
		}

		c.JSON(http.StatusOK, word)
	}
}

func DeleteWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		if err := models.DeleteWord(db, id); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete word")
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func AddWordToGroupFromWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		wordIDStr := c.Param("wordId")
		wordID, err := strconv.Atoi(wordIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		groupIDStr := c.Param("groupId")
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		if err := models.AddWordToGroup(db, wordID, groupID); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to add word to group")
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func RemoveWordFromGroupHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		wordIDStr := c.Param("wordId")
		wordID, err := strconv.Atoi(wordIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		groupIDStr := c.Param("groupId")
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}

		if err := models.RemoveWordFromGroup(db, wordID, groupID); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to remove word from group")
			return
		}

		c.Status(http.StatusNoContent)
	}
}
