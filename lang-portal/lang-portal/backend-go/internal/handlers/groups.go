package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/service"
)

// GroupHandler handles group-related requests
type GroupHandler struct {
	groupService *service.GroupService
}

// NewGroupHandler creates a new GroupHandler
func NewGroupHandler(gs *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: gs,
	}
}

// GetGroups returns a paginated list of groups
func (h *GroupHandler) GetGroups(c *gin.Context) {
	pagination := getPaginationParams(c)
	search := c.Query("search")

	groups, total, err := h.groupService.ListGroups(pagination.Page, &search)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get groups: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"pagination": PaginationResponse{
			CurrentPage:  pagination.Page,
			TotalPages:  (total + pagination.PageSize - 1) / pagination.PageSize,
			TotalItems:  total,
			ItemsPerPage: pagination.PageSize,
		},
	})
}

// GetGroup returns details of a specific group
func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.groupService.GetGroup(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	group, err := h.groupService.CreateGroup(req.Name)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create group: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateGroup updates an existing group
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	group, err := h.groupService.UpdateGroup(id, req.Name)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update group: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup deletes a group
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	if err := h.groupService.DeleteGroup(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete group: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetGroupWords returns words associated with a group
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	words, err := h.groupService.GetGroupWords(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group words: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, words)
}

// AddWordToGroup adds a word to a group
func (h *GroupHandler) AddWordToGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		WordID int64 `json:"word_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.groupService.AddWordToGroup(groupID, req.WordID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to add word to group: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveWordFromGroup removes a word from a group
func (h *GroupHandler) RemoveWordFromGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	if err := h.groupService.RemoveWordFromGroup(groupID, wordID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to remove word from group: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetGroupStudySessions returns study sessions for a group
func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Get study sessions for the group
	sessions, err := h.groupService.GetStudySessionsByGroupID(id)
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
			activity, err := h.groupService.GetStudyActivity(*session.StudyActivityID)
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