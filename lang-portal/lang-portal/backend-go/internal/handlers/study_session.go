package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/service"
)

// StudySessionHandler handles study session-related requests
type StudySessionHandler struct {
	sessionService *service.StudySessionService
}

// NewStudySessionHandler creates a new StudySessionHandler
func NewStudySessionHandler(ss *service.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{
		sessionService: ss,
	}
}

// StartStudySession starts a new study session
func (h *StudySessionHandler) StartStudySession(c *gin.Context) {
	var req struct {
		GroupID int64 `json:"group_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.sessionService.StartSession(req.GroupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to start study session: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, session)
}

// EndStudySession ends an active study session
func (h *StudySessionHandler) EndStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := h.sessionService.EndSession(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to end study session: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudySession returns details of a specific study session
func (h *StudySessionHandler) GetStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := h.sessionService.GetSession(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study session: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudySessions returns a list of study sessions
func (h *StudySessionHandler) GetStudySessions(c *gin.Context) {
	pagination := getPaginationParams(c)
	groupIDStr := c.Query("group_id")

	var groupID int64
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}
	}

	sessions, total, err := h.sessionService.ListSessions(pagination.Page, groupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sessions,
		"pagination": PaginationResponse{
			CurrentPage:  pagination.Page,
			TotalPages:  (total + pagination.PageSize - 1) / pagination.PageSize,
			TotalItems:  total,
			ItemsPerPage: pagination.PageSize,
		},
	})
}

// GetStudySessionStats returns statistics for a study session
func (h *StudySessionHandler) GetStudySessionStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	stats, err := h.sessionService.GetSessionStats(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get session stats: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetStudyProgress returns study progress over time
func (h *StudySessionHandler) GetStudyProgress(c *gin.Context) {
	// Default to last 30 days if not specified
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}

	since := time.Now().AddDate(0, 0, -days)
	progress, err := h.sessionService.GetUserProgress(since)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study progress: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, progress)
}
