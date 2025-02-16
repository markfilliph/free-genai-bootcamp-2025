package handlers

import (
	"lang-portal/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	pagination := getPaginationParams(c)
	offset := (pagination.Page - 1) * pagination.PageSize

	// Get total count for pagination
	total, err := models.GetTotalStudySessionsByActivity(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get total sessions count")
		return
	}

	// Get paginated sessions
	sessions, err := models.GetStudySessionsByActivity(id, offset, pagination.PageSize)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	var result []gin.H
	for _, s := range sessions {
		group, err := models.GetGroup(s.GroupID)
		if err != nil {
			continue
		}

		stats, err := models.GetStudySessionStats(s.ID)
		if err != nil {
			continue
		}

		result = append(result, gin.H{
			"id":               s.ID,
			"group_id":         s.GroupID,
			"group_name":       group.Name,
			"created_at":       s.CreatedAt,
			"total_words":      stats["total"],
			"correct_words":    stats["correct"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      result,
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, total),
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
		"id":         activity.ID,
		"group_id":   activity.GroupID,
		"created_at": activity.CreatedAt,
	})
}

// GetStudySessions returns a paginated list of study sessions
func GetStudySessions(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20")) // Default to 20 sessions per page

	sessions, err := models.GetStudySessions(offset, limit)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	var filteredSessions []gin.H
	for _, s := range sessions {
		group, err := models.GetGroup(s.GroupID)
		if err != nil {
			continue
		}

		stats, err := models.GetStudySessionStats(s.ID)
		if err != nil {
			continue
		}

		filteredSessions = append(filteredSessions, gin.H{
			"id":          s.ID,
			"group_name":  group.Name,
			"created_at":  s.CreatedAt,
			"total_words": stats["total"],
			"correct":     stats["correct"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": filteredSessions,
		"pagination": gin.H{
			"offset": offset,
			"limit":  limit,
		},
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
		"id":          session.ID,
		"group_name":  group.Name,
		"created_at":  session.CreatedAt,
		"total_words": stats["total"],
		"correct":     stats["correct"],
	})
}

// GetStudySessionByActivity returns the study session associated with an activity
func GetStudySessionByActivity(c *gin.Context) {
	// Parse activity ID from URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	// Get all study sessions (we'll optimize this later with a direct query)
	const maxSessions = 100 // Limit to last 100 sessions for performance
	sessions, err := models.GetStudySessions(0, maxSessions)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Find the session with this activity ID
	for _, s := range sessions {
		if s.StudyActivityID != nil && *s.StudyActivityID == id {
			c.JSON(http.StatusOK, s)
			return
		}
	}

	respondWithError(c, http.StatusNotFound, "No session found for this activity")
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

// ReviewWordInSession records a word review in a study session
func ReviewWordInSession(c *gin.Context) {
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
		respondWithError(c, http.StatusInternalServerError, "Failed to review word")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Word review recorded successfully",
	})
}
