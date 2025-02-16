package models

import (
	"fmt"
	"time"
)

// WordReview represents a review of a word
type WordReview struct {
	ID             *int       `json:"id"`
	WordID         *int       `json:"word_id"`
	StudySessionID *int       `json:"study_session_id"`
	Correct        bool       `json:"correct"`
	CreatedAt      *time.Time `json:"created_at"`
}

// WordReviewItem represents a word review item from a study session
type WordReviewItem struct {
	WordID  *int  `json:"word_id"`
	Correct bool  `json:"correct"`
}

// CreateWordReview creates a new word review
func CreateWordReview(review *WordReview) error {
	result, err := GetDB().Exec(`
		INSERT INTO word_reviews (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)`,
		review.WordID, review.StudySessionID,
		review.Correct, review.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating word review: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	iid := int(id)
	review.ID = &iid
	return nil
}

// GetWordReviewsBySession retrieves all word reviews for a study session
func GetWordReviewsBySession(sessionID int64) ([]*WordReview, error) {
	rows, err := GetDB().Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_reviews
		WHERE study_session_id = ?
		ORDER BY created_at`,
		sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying word reviews: %v", err)
	}
	defer rows.Close()

	var reviews []*WordReview
	for rows.Next() {
		var r WordReview
		if err := rows.Scan(&r.ID, &r.WordID, &r.StudySessionID, &r.Correct, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning word review: %v", err)
		}
		reviews = append(reviews, &r)
	}

	return reviews, nil
}

// GetWordReviewsByWord retrieves all reviews for a specific word
func GetWordReviewsByWord(wordID int64) ([]*WordReview, error) {
	rows, err := GetDB().Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_reviews
		WHERE word_id = ?
		ORDER BY created_at DESC`,
		wordID)
	if err != nil {
		return nil, fmt.Errorf("error querying word reviews: %v", err)
	}
	defer rows.Close()

	var reviews []*WordReview
	for rows.Next() {
		var r WordReview
		if err := rows.Scan(&r.ID, &r.WordID, &r.StudySessionID, &r.Correct, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning word review: %v", err)
		}
		reviews = append(reviews, &r)
	}

	return reviews, nil
}

// DeleteWordReview deletes a word review
func DeleteWordReview(id int64) error {
	result, err := GetDB().Exec(`DELETE FROM word_reviews WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting word review: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("word review not found: %d", id)
	}

	return nil
}

// CreateWordReviewsByActivity creates word reviews for a study session
func CreateWordReviewsByActivity(sessionID int64, items []*WordReviewItem) error {
	// Start transaction
	tx, err := GetDB().Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert each review
	for _, item := range items {
		sid := int(sessionID)
		now := time.Now()
		review := &WordReview{
			WordID:         item.WordID,
			StudySessionID: &sid,
			Correct:        item.Correct,
			CreatedAt:      &now,
		}

		result, err := tx.Exec(`
			INSERT INTO word_reviews (word_id, study_session_id, correct, created_at)
			VALUES (?, ?, ?, ?)`,
			review.WordID, review.StudySessionID,
			review.Correct, review.CreatedAt)
		if err != nil {
			return fmt.Errorf("error creating word review: %v", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert id: %v", err)
		}

		iid := int(id)
		review.ID = &iid
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
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