package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"lang-portal/internal/models"
	"lang-portal/internal/service"
)

var (
	router           *gin.Engine
	dashboardService *service.DashboardService
	groupService     *service.GroupServiceImpl
	studyService     *service.StudyService
	wordService      *service.WordService
)

func TestMain(m *testing.M) {
	// Set up test environment
	setupTestEnv()

	// Run tests
	code := m.Run()

	// Clean up
	cleanupTestEnv()

	os.Exit(code)
}

func setupTestEnv() {
	// Set test database configuration
	testDBName := "lang_portal_test"
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "Mfs1985+")
	os.Setenv("DB_NAME", testDBName)

	// Create test database
	rootDB, err := sql.Open("mysql", fmt.Sprintf("root:Mfs1985+@tcp(localhost:3306)/"))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}
	defer rootDB.Close()

	// Drop test database if it exists
	_, err = rootDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to drop test database: %v", err))
	}

	// Create test database
	_, err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to create test database: %v", err))
	}

	// Initialize database with schema
	if err := models.InitDB(""); err != nil {
		panic(fmt.Sprintf("Failed to initialize test database: %v", err))
	}

	// Initialize services
	dashboardService = service.NewDashboardService()
	groupService = service.NewGroupServiceImpl()
	studyService = service.NewStudyService()
	wordService = service.NewWordService()

	// Set up router
	gin.SetMode(gin.TestMode)
	router = setupRouter()
}

func cleanupTestEnv() {
	models.CloseDB()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	// Dashboard routes
	api.GET("/dashboard/last-session", getLastStudySession)
	api.GET("/dashboard/progress", getStudyProgress)
	api.GET("/dashboard/stats", getQuickStats)

	// Study activity routes
	api.GET("/study/activities", getStudyActivities)
	api.GET("/study/activities/:id", getStudyActivity)
	api.GET("/study/activities/:id/sessions", getStudyActivitySessions)
	api.POST("/study/activities", createStudyActivity)

	// Word routes
	api.GET("/words", getWords)
	api.GET("/words/:id", getWord)

	// Group routes
	api.GET("/groups", getGroups)
	api.POST("/groups", createGroup)
	api.GET("/groups/:id", getGroup)
	api.GET("/groups/:id/words", getGroupWords)
	api.GET("/groups/:id/sessions", getGroupStudySessions)

	return r
}

// Test helpers
func performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(w, req)
	return w
}

// API Tests
func TestCreateAndGetGroup(t *testing.T) {
	// Create a group
	createBody := map[string]string{"name": "Test Group"}
	w := performRequest("POST", "/api/groups", createBody)
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Group
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Test Group", response.Name)

	// Get the created group
	w = performRequest("GET", fmt.Sprintf("/api/groups/%d", response.ID), nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse models.Group
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.Nil(t, err)
	assert.Equal(t, "Test Group", getResponse.Name)
}

func TestCreateAndGetWord(t *testing.T) {
	// Create a word
	createBody := map[string]interface{}{
		"japanese": "こんにちは",
		"romaji":   "konnichiwa",
		"english":  "hello",
		"parts":    map[string]string{"type": "greeting", "formality": "neutral"},
	}
	w := performRequest("POST", "/api/words", createBody)
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Word
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", response.Japanese)

	// Get the created word
	w = performRequest("GET", fmt.Sprintf("/api/words/%d", response.ID), nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse models.Word
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", getResponse.Japanese)
}

func TestStudyFlow(t *testing.T) {
	// Create a group
	group := map[string]string{"name": "Study Group"}
	w := performRequest("POST", "/api/groups", group)
	assert.Equal(t, http.StatusOK, w.Code)
	var groupResponse models.Group
	json.Unmarshal(w.Body.Bytes(), &groupResponse)

	// Create a study activity
	activity := map[string]interface{}{
		"name": "Basic Study",
		"type": "flashcard",
	}
	w = performRequest("POST", "/api/study/activities", activity)
	assert.Equal(t, http.StatusOK, w.Code)
	var activityResponse models.StudyActivity
	json.Unmarshal(w.Body.Bytes(), &activityResponse)

	// Create a study session
	session := map[string]interface{}{
		"group_id":    groupResponse.ID,
		"activity_id": activityResponse.ID,
	}
	w = performRequest("POST", "/api/study/sessions", session)
	assert.Equal(t, http.StatusOK, w.Code)
	var sessionResponse models.StudySession
	json.Unmarshal(w.Body.Bytes(), &sessionResponse)

	// Review a word
	review := map[string]interface{}{
		"word_id": 1,
		"correct": true,
	}
	w = performRequest("POST", fmt.Sprintf("/api/study/sessions/%d/reviews", sessionResponse.ID), review)
	assert.Equal(t, http.StatusOK, w.Code)
}
