package handlers

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"net/http"
	"strconv"
	"time"
)

// GetStudyActivity returns details of a specific study activity
func GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	activity, err := models.GetStudyActivity(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study activity")
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetStudyActivitySessions returns paginated study sessions for an activity
func GetStudyActivitySessions(c *gin.Context) {
	pagination := getPaginationParams(c)

	sessions, err := models.GetStudySessions()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Filter and paginate sessions
	var filteredSessions []gin.H
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(sessions) {
		end = len(sessions)
	}

	for _, s := range sessions[start:end] {
		group, err := models.GetGroup(s.GroupID)
		if err != nil {
			continue
		}

		activity, err := models.GetStudyActivity(s.StudyActivityID)
		if err != nil {
			continue
		}

		filteredSessions = append(filteredSessions, gin.H{
			"id":                s.ID,
			"activity_name":     "Vocabulary Quiz", // TODO: Add activity name to model
			"group_name":        group.Name,
			"start_time":        s.CreatedAt,
			"end_time":          activity.CreatedAt,
			"review_items_count": len(sessions),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      filteredSessions,
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(sessions)),
	})
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(c *gin.Context) {
	var request struct {
		GroupID         int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	activity, err := models.CreateStudyActivity(request.GroupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create study activity")
		return
	}

	session, err := models.CreateStudySession(request.GroupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create study session")
		return
	}

	err = models.LinkStudyActivityToSession(activity.ID, session.ID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to link activity to session")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                activity.ID,
		"group_id":          activity.GroupID,
		"created_at":        activity.CreatedAt,
	})
}

// GetStudySessions returns a paginated list of study sessions
func GetStudySessions(c *gin.Context) {
	pagination := getPaginationParams(c)

	sessions, err := models.GetStudySessions()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Filter and paginate sessions
	var filteredSessions []gin.H
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(sessions) {
		end = len(sessions)
	}

	for _, s := range sessions[start:end] {
		group, err := models.GetGroup(s.GroupID)
		if err != nil {
			continue
		}

		stats, err := models.GetStudySessionStats(s.ID)
		if err != nil {
			continue
		}

		filteredSessions = append(filteredSessions, gin.H{
			"id":           s.ID,
			"group_name":   group.Name,
			"created_at":   s.CreatedAt,
			"total_words": stats["total"],
			"correct":     stats["correct"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      filteredSessions,
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(sessions)),
	})
}

// GetStudySession returns details of a specific study session
func GetStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := models.GetStudySession(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study session")
		return
	}

	group, err := models.GetGroup(session.GroupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group info")
		return
	}

	stats, err := models.GetStudySessionStats(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get session stats")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           session.ID,
		"group_name":   group.Name,
		"created_at":   session.CreatedAt,
		"total_words": stats["total"],
		"correct":     stats["correct"],
	})
}

// GetStudySessionWords returns words associated with a study session
func GetStudySessionWords(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	words, err := models.GetStudySessionWords(sessionID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get session words")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"words":      words,
	})
}

// ReviewWord records a word review in a study session
func ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var request struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	err = models.ReviewWord(wordID, sessionID, request.Correct)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to record word review")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"word_id":    wordID,
		"correct":    request.Correct,
		"created_at": time.Now(),
	})
}