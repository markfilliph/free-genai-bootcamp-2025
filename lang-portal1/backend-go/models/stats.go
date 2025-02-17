package models

import (
	"database/sql"
)

type WordStudyStats struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords  int `json:"total_available_words"`
}

// GetWordStudyStats retrieves word study statistics
func GetWordStudyStats(db *sql.DB) (*WordStudyStats, error) {
	var progress WordStudyStats

	err := db.QueryRow(`
		SELECT 
			(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as studied,
			(SELECT COUNT(*) FROM words) as total`).Scan(
		&progress.TotalWordsStudied,
		&progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}
