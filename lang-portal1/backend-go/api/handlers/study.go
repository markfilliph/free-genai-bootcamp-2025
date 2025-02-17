package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/models"
)

func GetStudyActivity(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid activity ID")
			return
		}

		activity, err := models.GetStudyActivity(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Study activity not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get study activity")
			return
		}

		c.JSON(http.StatusOK, activity)
	}
}

func GetStudyActivitySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid activity ID")
			return
		}

		page, perPage := getPaginationParams(c)

		sessions, total, err := models.GetStudySessions(db, id, page, perPage)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(sessions, page, total, perPage))
	}
}

type CreateStudySessionRequest struct {
	GroupID         int `json:"group_id" binding:"required"`
	StudyActivityID int `json:"study_activity_id" binding:"required"`
}

func CreateStudySession(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateStudySessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		session, err := models.CreateStudySession(db, req.GroupID, req.StudyActivityID)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create study session")
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

type AddWordReviewRequest struct {
	Correct bool `json:"correct" binding:"required"`
}

func GetStudySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, perPage := getPaginationParams(c)

		sessions, total, err := models.GetStudySessions(db, 0, page, perPage) // 0 means all activities
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(sessions, page, total, perPage))
	}
}

func GetStudySession(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid session ID")
			return
		}

		session, err := models.GetStudySession(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Study session not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func AddWordReview(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionIDStr := c.Param("session_id")
		sessionID, err := strconv.Atoi(sessionIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid session ID")
			return
		}

		wordIDStr := c.Param("word_id")
		wordID, err := strconv.Atoi(wordIDStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		var req AddWordReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		err = models.AddWordReview(db, sessionID, wordID, req.Correct)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to add word review")
			return
		}

		c.Status(http.StatusCreated)
	}
}
