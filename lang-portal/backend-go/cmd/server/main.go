package main

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/database"
	"lang-portal/internal/service"
	"log"
	"net/http"
	"strconv"
)

var (
	dashboardService *service.DashboardService
	groupService     *service.GroupServiceImpl
	studyService     *service.StudyService
	wordService      *service.WordService
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Initialize services
	dashboardService = service.NewDashboardService()
	groupService = service.NewGroupServiceImpl()
	studyService = service.NewStudyService()
	wordService = service.NewWordService()

	// Initialize Gin router
	r := gin.Default()

	// Initialize routes
	initializeRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeRoutes(r *gin.Engine) {
	// API group
	api := r.Group("/api")
	
	// Dashboard routes
	api.GET("/dashboard/last-study-session", getLastStudySession)
	api.GET("/dashboard/study-progress", getStudyProgress)
	api.GET("/dashboard/quick-stats", getQuickStats)

	// Study activities routes
	api.GET("/study-activities/:id", getStudyActivity)
	api.GET("/study-activities/:id/sessions", getStudyActivitySessions)
	api.POST("/study-activities", createStudyActivity)

	// Words routes
	api.GET("/words", getWords)
	api.GET("/words/:id", getWord)

	// Groups routes
	api.GET("/groups", getGroups)
	api.POST("/groups", createGroup)
	api.GET("/groups/:id", getGroup)
	api.GET("/groups/:id/words", getGroupWords)
	api.GET("/groups/:id/study-sessions", getGroupStudySessions)

	// Study sessions routes
	api.GET("/study-sessions", getStudySessions)
	api.GET("/study-sessions/:id", getStudySession)
	api.GET("/study-sessions/:id/words", getStudySessionWords)
	api.POST("/study-sessions/:id/words/:word_id/review", reviewWordInSession)

	// Reset routes
	api.POST("/reset/history", resetHistory)
	api.POST("/reset/full", fullReset)
}

// Dashboard handlers
func getLastStudySession(c *gin.Context) {
	session, err := dashboardService.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func getStudyProgress(c *gin.Context) {
	progress, err := dashboardService.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func getQuickStats(c *gin.Context) {
	stats, err := dashboardService.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Study activity handlers
func getStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := studyService.GetStudyActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func getStudyActivitySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	sessions, err := studyService.GetStudyActivitySessions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func createStudyActivity(c *gin.Context) {
	var request struct {
		GroupID int64 `json:"group_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	session, err := studyService.CreateStudySession(request.GroupID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

// Word handlers
func getWords(c *gin.Context) {
	words, err := wordService.GetWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, words)
}

func getWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := wordService.GetWord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

// Group handlers
func getGroups(c *gin.Context) {
	offset := 0
	limit := 10 // Default limit

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	groups, err := groupService.GetGroups(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func createGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := groupService.CreateGroup(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func getGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := groupService.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

func getGroupWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	words, err := groupService.GetGroupWords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, words)
}

func getGroupStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	sessions, err := groupService.GetGroupStudySessions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// Study session handlers
func getStudySessions(c *gin.Context) {
	activities, err := studyService.GetStudyActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func getStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	activity, err := studyService.GetStudyActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func getStudySessionWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	words, err := wordService.GetWordReviews(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, words)
}

func reviewWordInSession(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var request struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	review, err := studyService.ReviewWord(sessionID, wordID, request.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, review)
}

// Reset handlers
func resetHistory(c *gin.Context) {
	// Reset all study history but keep words and groups
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec(`
		DELETE FROM word_reviews;
		DELETE FROM study_sessions;
		DELETE FROM study_activities;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Study history has been reset"})
}

func fullReset(c *gin.Context) {
	// Reset everything in the database
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec(`
		DELETE FROM word_reviews;
		DELETE FROM study_sessions;
		DELETE FROM study_activities;
		DELETE FROM words;
		DELETE FROM groups;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database has been fully reset"})
}