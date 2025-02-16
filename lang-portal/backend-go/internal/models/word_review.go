package models

import (
	"database/sql"
	"time"
)

// WordReview represents a word review in a study session
type WordReview struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateWordReview creates a new word review
func CreateWordReview(tx *sql.Tx, wordID, studySessionID int64, correct bool) error {
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	
	var err error
	if tx != nil {
		_, err = tx.Exec(query, wordID, studySessionID, correct)
	} else {
		_, err = DB.Exec(query, wordID, studySessionID, correct)
	}
	
	return err
}

// GetWordReview retrieves a word review by ID
func GetWordReview(id int64) (*WordReview, error) {
	var review WordReview
	err := DB.QueryRow(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE id = ?
	`, id).Scan(&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetWordReviewsBySession returns all reviews for a specific study session
func GetWordReviewsBySession(studySessionID int64) ([]*WordReview, error) {
	rows, err := DB.Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at DESC
	`, studySessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*WordReview
	for rows.Next() {
		var review WordReview
		err := rows.Scan(
			&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

// GetWordReviews returns all reviews for a specific word
func GetWordReviews(wordID int64) ([]*WordReview, error) {
	rows, err := DB.Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE word_id = ?
		ORDER BY created_at DESC
	`, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*WordReview
	for rows.Next() {
		var review WordReview
		err := rows.Scan(&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

// GetWordReviewStats gets statistics for a word's review history
func GetWordReviewStats(wordID int64) (map[string]int, error) {
	var totalReviews, correctReviews int
	err := DB.QueryRow(`
		SELECT 
			COUNT(*) as total_reviews,
			SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_reviews
		FROM word_review_items
		WHERE word_id = ?
	`, wordID).Scan(&totalReviews, &correctReviews)

	if err != nil {
		return nil, err
	}

	var successRate int
	if totalReviews > 0 {
		successRate = (correctReviews * 100) / totalReviews
	}

	return map[string]int{
		"total_reviews": totalReviews,
		"success_rate":  successRate,
	}, nil
}

// GetRecentReviewsSuccessRate gets the success rate from recent reviews
func GetRecentReviewsSuccessRate() (float64, error) {
	var totalReviews, correctReviews int
	err := DB.QueryRow(`
		SELECT 
			COUNT(*) as total_reviews,
			SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_reviews
		FROM word_review_items
		WHERE created_at >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 30 DAY)
	`).Scan(&totalReviews, &correctReviews)

	if err != nil {
		return 0, err
	}

	if totalReviews == 0 {
		return 0, nil
	}

	return float64(correctReviews) * 100 / float64(totalReviews), nil
}

// GetStudySessionWords returns all words reviewed in a study session
func GetStudySessionWords(studySessionID int64) ([]Word, error) {
	rows, err := DB.Query(`
		SELECT DISTINCT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		ORDER BY w.japanese
	`, studySessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts)
		if err != nil {
			return nil, err
		}
		words = append(words, w)
	}

	return words, nil
}