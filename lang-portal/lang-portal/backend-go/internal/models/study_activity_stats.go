package models

import "time"

// StudyActivityStats represents statistics for a study activity
type StudyActivityStats struct {
	TotalSessions      int       `json:"total_sessions"`
	TotalWords         int       `json:"total_words"`
	CorrectWords       int       `json:"correct_words"`
	IncorrectWords     int       `json:"incorrect_words"`
	LastStudiedAt      time.Time `json:"last_studied_at"`
	AverageAccuracy    float64   `json:"average_accuracy"`
	AverageTimePerWord float64   `json:"average_time_per_word"`
}

// GetStudyActivityStats retrieves statistics for a study activity
func GetStudyActivityStats(activityID int64) (*StudyActivityStats, error) {
	var stats StudyActivityStats

	// Get session counts and last studied time
	err := GetDB().QueryRow(`
		SELECT 
			COUNT(DISTINCT ss.id) as total_sessions,
			COUNT(DISTINCT wr.word_id) as total_words,
			COUNT(DISTINCT CASE WHEN wr.correct = true THEN wr.word_id END) as correct_words,
			COUNT(DISTINCT CASE WHEN wr.correct = false THEN wr.word_id END) as incorrect_words,
			MAX(ss.created_at) as last_studied_at,
			AVG(CASE WHEN wr.correct = true THEN 1.0 ELSE 0.0 END) * 100 as average_accuracy,
			AVG(TIMESTAMPDIFF(SECOND, ss.created_at, wr.created_at)) as average_time_per_word
		FROM study_sessions ss
		LEFT JOIN word_reviews wr ON ss.id = wr.study_session_id
		WHERE ss.study_activity_id = ?`,
		activityID).Scan(
		&stats.TotalSessions,
		&stats.TotalWords,
		&stats.CorrectWords,
		&stats.IncorrectWords,
		&stats.LastStudiedAt,
		&stats.AverageAccuracy,
		&stats.AverageTimePerWord)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// StudyActivityProgress represents progress in a study activity
type StudyActivityProgress struct {
	TotalWords      int     `json:"total_words"`
	StudiedWords    int     `json:"studied_words"`
	CorrectWords    int     `json:"correct_words"`
	IncorrectWords  int     `json:"incorrect_words"`
	CompletionRate  float64 `json:"completion_rate"`
	AccuracyRate    float64 `json:"accuracy_rate"`
}

// GetStudyActivityProgress retrieves progress for a study activity
func GetStudyActivityProgress(activityID int64) (*StudyActivityProgress, error) {
	var progress StudyActivityProgress

	// Get total words in activity
	err := GetDB().QueryRow(`
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN study_activities sa ON sa.group_id = w.group_id
		WHERE sa.id = ?`,
		activityID).Scan(&progress.TotalWords)
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
		JOIN study_activities sa ON sa.group_id = w.group_id
		LEFT JOIN word_reviews wr ON w.id = wr.word_id
		WHERE sa.id = ?`,
		activityID).Scan(&progress.StudiedWords, &progress.CorrectWords, &progress.IncorrectWords)
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
