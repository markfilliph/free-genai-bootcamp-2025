package service

import (
	"fmt"
	"time"

	"lang-portal/internal/models"
)

// StudySessionService handles business logic for study sessions
type StudySessionService struct {
	// Add any dependencies here
}

// NewStudySessionService creates a new StudySessionService
func NewStudySessionService() *StudySessionService {
	return &StudySessionService{}
}

// CreateStudySession creates a new study session
func (s *StudySessionService) CreateStudySession(groupID, activityID int64) (*models.StudySession, error) {
	now := time.Now()
	session := &models.StudySession{
		GroupID:         &groupID,
		StudyActivityID: &activityID,
		StartTime:       &now,
		CreatedAt:       &now,
	}

	if err := models.CreateStudySession(session); err != nil {
		return nil, fmt.Errorf("error creating study session: %v", err)
	}

	return session, nil
}

// GetStudySession retrieves a study session by ID
func (s *StudySessionService) GetStudySession(id int64) (*models.StudySession, error) {
	session, err := models.GetStudySession(id)
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found: %d", id)
	}
	return session, nil
}

// GetStudySessionsByGroup retrieves study sessions for a group with pagination
func (s *StudySessionService) GetStudySessionsByGroup(groupID int64, offset, limit int) ([]*models.StudySession, error) {
	sessions, err := models.GetStudySessionsByGroup(groupID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting study sessions: %v", err)
	}
	return sessions, nil
}

// GetStudySessionsByActivity retrieves study sessions for an activity with pagination
func (s *StudySessionService) GetStudySessionsByActivity(activityID int64, offset, limit int) ([]*models.StudySession, error) {
	sessions, err := models.GetStudySessionsByActivity(activityID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting study sessions: %v", err)
	}
	return sessions, nil
}

// CreateWordReview creates a word review for a study session
func (s *StudySessionService) CreateWordReview(sessionID, wordID int64, correct bool) error {
	now := time.Now()
	review := &models.WordReview{
		WordID:         &wordID,
		StudySessionID: &sessionID,
		Correct:        correct,
		CreatedAt:      &now,
	}

	if err := models.CreateWordReview(review); err != nil {
		return fmt.Errorf("error creating word review: %v", err)
	}

	return nil
}

// CreateWordReviews creates multiple word reviews for a study session
func (s *StudySessionService) CreateWordReviews(sessionID int64, items []*models.WordReviewItem) error {
	// Convert items to reviews
	var reviews []*models.WordReview
	for _, item := range items {
		now := time.Now()
		review := &models.WordReview{
			WordID:         item.WordID,
			StudySessionID: &sessionID,
			Correct:        item.Correct,
			CreatedAt:      &now,
		}
		reviews = append(reviews, review)
	}

	// Create reviews in a transaction
	if err := models.CreateWordReviewsByActivity(sessionID, items); err != nil {
		return fmt.Errorf("error creating word reviews: %v", err)
	}
	return nil
}

// GetWordReviewsBySession retrieves all word reviews for a study session
func (s *StudySessionService) GetWordReviewsBySession(sessionID int64) ([]*models.WordReview, error) {
	reviews, err := models.GetWordReviewsBySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting word reviews: %v", err)
	}
	return reviews, nil
}

// GetWordReviewsByWord retrieves all reviews for a word
func (s *StudySessionService) GetWordReviewsByWord(wordID int64) ([]*models.WordReview, error) {
	reviews, err := models.GetWordReviewsByWord(wordID)
	if err != nil {
		return nil, fmt.Errorf("error getting word reviews: %v", err)
	}
	return reviews, nil
}

// GetSession retrieves a study session by ID with stats
func (s *StudySessionService) GetSession(sessionID int64) (*models.StudySession, error) {
	session, err := models.GetStudySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found: %d", sessionID)
	}

	return session, nil
}

// GetLastGroupSession retrieves the most recent session for a group
func (s *StudySessionService) GetLastGroupSession(groupID int64) (*models.StudySession, error) {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	sessions, err := models.GetStudySessionsByGroup(groupID, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get last session: %v", err)
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return sessions[0], nil
}

// ListActivitySessions retrieves study sessions for an activity with pagination
func (s *StudySessionService) ListActivitySessions(activityID int64, page int) ([]*models.StudySession, int, error) {
	if page < 1 {
		page = 1
	}

	// Verify activity exists
	activity, err := models.GetStudyActivity(activityID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, 0, fmt.Errorf("study activity not found: %d", activityID)
	}

	limit := 10
	offset := (page - 1) * limit

	sessions, err := models.GetStudySessionsByActivity(activityID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sessions: %v", err)
	}

	total, err := models.GetTotalStudySessionsByActivity(activityID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total sessions: %v", err)
	}

	return sessions, total, nil
}

// RecordWordReview records a word review result in a session
func (s *StudySessionService) RecordWordReview(sessionID, wordID int64, correct bool) error {
	// Verify session exists
	session, err := models.GetStudySession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get study session: %v", err)
	}
	if session == nil {
		return fmt.Errorf("study session not found: %d", sessionID)
	}

	// Verify word exists
	word, err := models.GetWord(wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return fmt.Errorf("word not found: %d", wordID)
	}

	// Create word review
	review := &models.WordReviewItem{
		WordID:         wordID,
		StudySessionID: sessionID,
		Correct:        correct,
		CreatedAt:      time.Now(),
	}

	if err := models.CreateWordReview(review); err != nil {
		return fmt.Errorf("failed to record word review: %v", err)
	}

	return nil
}

// GetSessionReviews retrieves all word reviews in a session
func (s *StudySessionService) GetSessionReviews(sessionID int64) ([]models.WordReviewItem, error) {
	// Verify session exists
	session, err := models.GetStudySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found: %d", sessionID)
	}

	reviews, err := models.GetSessionReviews(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session reviews: %v", err)
	}

	return reviews, nil
}

// GetSessionStats retrieves detailed statistics for a session
func (s *StudySessionService) GetSessionStats(sessionID int64) (*models.StudySessionResponse, error) {
	// Get session
	session, err := models.GetStudySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found: %d", sessionID)
	}

	// Get group name
	group, err := models.GetGroup(session.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", session.GroupID)
	}

	// Get reviews
	reviews, err := models.GetSessionReviews(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session reviews: %v", err)
	}

	// Calculate statistics
	var correctWords int
	for _, review := range reviews {
		if review.Correct {
			correctWords++
		}
	}

	response := &models.StudySessionResponse{
		ID:              session.ID,
		GroupID:         session.GroupID,
		GroupName:       group.Name,
		StudyActivityID: session.StudyActivityID,
		TotalWords:      len(reviews),
		CorrectWords:    correctWords,
		CreatedAt:       session.CreatedAt,
	}

	return response, nil
}

// GetUserProgress retrieves overall study progress across all sessions
func (s *StudySessionService) GetUserProgress(since time.Time) (*models.StudyProgressResponse, error) {
	totalWords, err := models.GetTotalWordsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total words: %v", err)
	}

	studiedWords, err := models.GetTotalStudiedWordsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get studied words: %v", err)
	}

	var completionPercentage float64
	if totalWords > 0 {
		completionPercentage = float64(studiedWords) / float64(totalWords) * 100
	}

	progress := &models.StudyProgressResponse{
		TotalWords:          totalWords,
		TotalWordsStudied:   studiedWords,
		RemainingWords:      totalWords - studiedWords,
		CompletionPercentage: completionPercentage,
	}

	return progress, nil
}
