package models

import (
	"database/sql"
	"fmt"
	"time"
)

// WordReviewItem represents a single word review attempt
type WordReviewItem struct {
	ID              int64     `json:"id"`
	WordID          int64     `json:"word_id"`
	StudySessionID  int64     `json:"study_session_id"`
	Correct         bool      `json:"correct"`
	CreatedAt       time.Time `json:"created_at"`
}

// WordReviewStats contains statistics for a word's review history
type WordReviewStats struct {
	WordID          int64   `json:"word_id"`
	TotalReviews    int     `json:"total_reviews"`
	CorrectReviews  int     `json:"correct_reviews"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	LastReviewedAt  string  `json:"last_reviewed_at"`
}

// CreateWordReview records a word review attempt
func CreateWordReview(wordID, studySessionID int64, correct bool) error {
	db := GetDB()

	// Verify word exists
	exists, err := wordExists(wordID)
	if err != nil {
		return fmt.Errorf("error checking word existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("word not found: %d", wordID)
	}

	// Verify study session exists
	exists, err = studySessionExists(studySessionID)
	if err != nil {
		return fmt.Errorf("error checking study session existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("study session not found: %d", studySessionID)
	}

	// Create review
	_, err = db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		wordID, studySessionID, correct)
	if err != nil {
		return fmt.Errorf("error creating word review: %v", err)
	}

	return nil
}

// CreateBatchWordReviews records multiple word reviews in a single transaction
func CreateBatchWordReviews(reviews []WordReviewItem) error {
	if len(reviews) == 0 {
		return nil
	}

	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, review := range reviews {
		// Verify word exists
		exists, err := wordExists(review.WordID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error checking word existence: %v", err)
		}
		if !exists {
			tx.Rollback()
			return fmt.Errorf("word not found: %d", review.WordID)
		}

		// Verify study session exists
		exists, err = studySessionExists(review.StudySessionID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error checking study session existence: %v", err)
		}
		if !exists {
			tx.Rollback()
			return fmt.Errorf("study session not found: %d", review.StudySessionID)
		}

		_, err = stmt.Exec(review.WordID, review.StudySessionID, review.Correct)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating word review: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetWordReviewStats retrieves review statistics for a word
func GetWordReviewStats(wordID int64) (*WordReviewStats, error) {
	db := GetDB()

	var stats WordReviewStats
	stats.WordID = wordID

	err := db.QueryRow(`
		SELECT 
			COUNT(*) as total_reviews,
			SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_reviews,
			COALESCE(MAX(created_at), '') as last_reviewed_at
		FROM word_review_items
		WHERE word_id = ?`, wordID).
		Scan(&stats.TotalReviews, &stats.CorrectReviews, &stats.LastReviewedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting word review stats: %v", err)
	}

	if stats.TotalReviews > 0 {
		stats.AccuracyRate = float64(stats.CorrectReviews) / float64(stats.TotalReviews)
	}

	return &stats, nil
}

// GetWordReviewHistory retrieves review history for a word with pagination
func GetWordReviewHistory(wordID int64, page int) ([]WordReviewItem, int, error) {
	db := GetDB()

	// Get total count
	var total int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM word_review_items 
		WHERE word_id = ?`, wordID).
		Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting reviews: %v", err)
	}

	// Calculate offset
	offset := (page - 1) * 100
	if offset < 0 {
		offset = 0
	}

	// Get reviews with pagination
	rows, err := db.Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE word_id = ?
		ORDER BY created_at DESC
		LIMIT 100 OFFSET ?`, wordID, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying reviews: %v", err)
	}
	defer rows.Close()

	var reviews []WordReviewItem
	for rows.Next() {
		var r WordReviewItem
		if err := rows.Scan(&r.ID, &r.WordID, &r.StudySessionID, &r.Correct, &r.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning review: %v", err)
		}
		reviews = append(reviews, r)
	}

	return reviews, total, nil
}

// GetSessionReviews retrieves all reviews for a study session
func GetSessionReviews(sessionID int64) ([]WordReviewItem, error) {
	db := GetDB()

	rows, err := db.Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying session reviews: %v", err)
	}
	defer rows.Close()

	var reviews []WordReviewItem
	for rows.Next() {
		var r WordReviewItem
		if err := rows.Scan(&r.ID, &r.WordID, &r.StudySessionID, &r.Correct, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning review: %v", err)
		}
		reviews = append(reviews, r)
	}

	return reviews, nil
}

// Helper functions

func wordExists(wordID int64) (bool, error) {
	db := GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)", wordID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking word existence: %v", err)
	}

	return exists, nil
}

func studySessionExists(sessionID int64) (bool, error) {
	db := GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM study_sessions WHERE id = ?)", sessionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking study session existence: %v", err)
	}

	return exists, nil
}