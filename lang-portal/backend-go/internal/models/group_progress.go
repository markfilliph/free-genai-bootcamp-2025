package models

// GroupProgress represents the learning progress for a group
type GroupProgress struct {
	TotalWords      int     `json:"total_words"`
	StudiedWords    int     `json:"studied_words"`
	CorrectWords    int     `json:"correct_words"`
	IncorrectWords  int     `json:"incorrect_words"`
	CompletionRate  float64 `json:"completion_rate"`
	AccuracyRate    float64 `json:"accuracy_rate"`
}

// GetGroupProgress calculates the study progress for a group
func GetGroupProgress(groupID int64) (*GroupProgress, error) {
	var progress GroupProgress

	// Get total words in group
	err := GetDB().QueryRow(`
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?`,
		groupID).Scan(&progress.TotalWords)
	if err != nil {
		return nil, err
	}

	// Get studied words statistics
	err = GetDB().QueryRow(`
		SELECT 
			COUNT(DISTINCT wr.word_id) as studied_words,
			COUNT(DISTINCT CASE WHEN wr.correct = true THEN wr.word_id END) as correct_words,
			COUNT(DISTINCT CASE WHEN wr.correct = false THEN wr.word_id END) as incorrect_words
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_reviews wr ON w.id = wr.word_id
		WHERE wg.group_id = ?`,
		groupID).Scan(&progress.StudiedWords, &progress.CorrectWords, &progress.IncorrectWords)
	if err != nil {
		return nil, err
	}

	// Calculate rates
	if progress.TotalWords > 0 {
		progress.CompletionRate = float64(progress.StudiedWords) / float64(progress.TotalWords) * 100
	}
	if progress.StudiedWords > 0 {
		progress.AccuracyRate = float64(progress.CorrectWords) / float64(progress.StudiedWords) * 100
	}

	return &progress, nil
}
