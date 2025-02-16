package handlers

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"net/http"
	"strconv"
)

// GetGroups returns a paginated list of groups
func GetGroups(c *gin.Context) {
	pagination := getPaginationParams(c)

	// Get all groups
	groups, err := models.GetGroups(0, 0) // Pass 0 for limit to get all groups
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	// Filter and paginate groups
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(groups) {
		end = len(groups)
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      groups[start:end],
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(groups)),
	})
}

// GetGroup returns details of a specific group
func GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := models.GetGroup(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group")
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords returns words associated with a group
func GetGroupWords(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	pagination := getPaginationParams(c)

	words, err := models.GetGroupWords(groupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group words")
		return
	}

	// Filter and paginate words
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(words) {
		end = len(words)
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":   groupID,
		"items":      words[start:end],
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(words)),
	})
}

// GetGroupStudySessions returns study sessions for a group
func GetGroupStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Get study sessions for the group
	sessions, err := models.GetStudySessionsByGroupID(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	type EnrichedSession struct {
		*models.StudySession
		Activity *models.StudyActivity `json:"activity,omitempty"`
	}

	enrichedSessions := make([]EnrichedSession, len(sessions))
	for i, session := range sessions {
		enrichedSession := EnrichedSession{
			StudySession: session,
		}

		// Only get activity if StudyActivityID is not nil
		if session.StudyActivityID != nil {
			activity, err := models.GetStudyActivity(*session.StudyActivityID)
			if err != nil {
				respondWithError(c, http.StatusInternalServerError, "Failed to get study activity")
				return
			}
			enrichedSession.Activity = activity
		}

		enrichedSessions[i] = enrichedSession
	}

	c.JSON(http.StatusOK, enrichedSessions)
}