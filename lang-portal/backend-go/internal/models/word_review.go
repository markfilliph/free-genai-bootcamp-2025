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

// CreateWordReview creates a new word review in the database
func CreateWordReview(wordID, sessionID int64, correct bool) (*WordReview, error) {
	var review WordReview
	err := DB.QueryRow(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, word_id, study_session_id, correct, created_at
	`, wordID, sessionID, correct).Scan(
		&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
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

// GetWordReviewsBySession retrieves all word reviews for a study session
func GetWordReviewsBySession(sessionID int64) ([]*WordReview, error) {
	rows, err := DB.Query(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at DESC
	`, sessionID)
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
		"total_reviews":   totalReviews,
		"correct_reviews": correctReviews,
		"success_rate":    successRate,
	}, nil
}

// GetStudyProgress returns study progress statistics
func GetStudyProgress() (map[string]int, error) {
	var totalWords, studiedWords int
	err := DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalWords)
	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&studiedWords)
	if err != nil {
		return nil, err
	}

	var completionPercentage int
	if totalWords > 0 {
		completionPercentage = (studiedWords * 100) / totalWords
	}

	return map[string]int{
		"total_words":          totalWords,
		"total_words_studied":  studiedWords,
		"remaining_words":      totalWords - studiedWords,
		"completion_percentage": completionPercentage,
	}, nil
}

// GetStudySessionWords returns all words reviewed in a study session
func GetStudySessionWords(sessionID int64) ([]Word, error) {
	rows, err := DB.Query(`
		SELECT DISTINCT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		INNER JOIN word_review_items wr ON w.id = wr.word_id
		WHERE wr.study_session_id = ?
		ORDER BY w.japanese
	`, sessionID)
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