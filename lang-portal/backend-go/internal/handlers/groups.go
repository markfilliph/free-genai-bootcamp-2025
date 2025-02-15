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

	groups, err := models.GetGroups()
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
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	pagination := getPaginationParams(c)

	sessions, err := models.GetStudySessions()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Filter sessions by group ID
	var groupSessions []models.StudySession
	for _, s := range sessions {
		if s.GroupID == groupID {
			groupSessions = append(groupSessions, s)
		}
	}

	// Paginate filtered sessions
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(groupSessions) {
		end = len(groupSessions)
	}

	var filteredSessions []gin.H
	for _, s := range groupSessions[start:end] {
		stats, err := models.GetStudySessionStats(s.ID)
		if err != nil {
			continue
		}

		filteredSessions = append(filteredSessions, gin.H{
			"id":           s.ID,
			"created_at":   s.CreatedAt,
			"total_words": stats["total"],
			"correct":     stats["correct"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":   groupID,
		"items":      filteredSessions,
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(groupSessions)),
	})
}