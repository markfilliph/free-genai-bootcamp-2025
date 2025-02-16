package models

import "time"

// WordResponse represents a word with its study statistics
type WordResponse struct {
	ID           int64       `json:"id"`
	Japanese     string      `json:"japanese"`
	Romaji       string      `json:"romaji"`
	English      string      `json:"english"`
	Parts        interface{} `json:"parts,omitempty"`
	TotalReviews int         `json:"total_reviews"`
	SuccessRate  float64     `json:"success_rate"`
}

// WordReviewResponse represents a word review with session details
type WordReviewResponse struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	GroupName      string    `json:"group_name"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// StudyActivityResponse represents a study activity with its statistics
type StudyActivityResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	TotalSessions int       `json:"total_sessions"`
	TotalWords    int       `json:"total_words"`
	SuccessRate   float64   `json:"success_rate"`
	CreatedAt     time.Time `json:"created_at"`
	Session       *StudySession `json:"session,omitempty"`
}

// StudySessionResponse represents a study session with its statistics
type StudySessionResponse struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	GroupName       string    `json:"group_name"`
	StudyActivityID *int64    `json:"study_activity_id,omitempty"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	CreatedAt       time.Time `json:"created_at"`
}

// StudyProgressResponse represents overall study progress
type StudyProgressResponse struct {
	TotalWords          int     `json:"total_words"`
	TotalWordsStudied   int     `json:"total_words_studied"`
	RemainingWords      int     `json:"remaining_words"`
	CompletionPercentage float64 `json:"completion_percentage"`
}

// QuickStatsResponse represents quick overview statistics
type QuickStatsResponse struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}
