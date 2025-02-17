package models

import (
	"database/sql"
	"time"
)

type QuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

type StudyProgress struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords  int `json:"total_available_words"`
}

type LastStudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
}

// GetQuickStats retrieves dashboard quick statistics
func GetQuickStats(db *sql.DB) (*QuickStats, error) {
	var stats QuickStats

	// Get success rate and total study sessions
	err := db.QueryRow(`
		SELECT 
			COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END) * 100, 0) as success_rate,
			COUNT(DISTINCT study_session_id) as total_sessions
		FROM word_review_items`).Scan(&stats.SuccessRate, &stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	// Get total active groups (groups with at least one study session)
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT group_id) 
		FROM study_sessions`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	// Calculate study streak
	err = db.QueryRow(`
		WITH RECURSIVE dates(date) AS (
			SELECT date(MAX(created_at)) FROM study_sessions
			UNION ALL
			SELECT date(date, '-1 day')
			FROM dates
			WHERE date > date((
				SELECT MIN(created_at) FROM study_sessions
			))
		),
		daily_sessions AS (
			SELECT date(created_at) as session_date
			FROM study_sessions
			GROUP BY date(created_at)
		)
		SELECT COUNT(*)
		FROM (
			SELECT dates.date
			FROM dates
			LEFT JOIN daily_sessions ON dates.date = daily_sessions.session_date
			WHERE daily_sessions.session_date IS NOT NULL
			ORDER BY dates.date DESC
		) streak
		WHERE rowid <= (
			SELECT MIN(rowid) - 1
			FROM (
				SELECT dates.date, rowid
				FROM dates
				LEFT JOIN daily_sessions ON dates.date = daily_sessions.session_date
				WHERE daily_sessions.session_date IS NULL
			)
		)`).Scan(&stats.StudyStreakDays)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetStudyProgress retrieves study progress statistics
func GetStudyProgress(db *sql.DB) (*StudyProgress, error) {
	var progress StudyProgress

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

// GetLastStudySession retrieves the most recent study session
func GetLastStudySession(db *sql.DB) (*LastStudySession, error) {
	var session LastStudySession

	err := db.QueryRow(`
		SELECT 
			ss.id,
			ss.group_id,
			ss.created_at,
			ss.study_activity_id,
			g.name as group_name
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		ORDER BY ss.created_at DESC
		LIMIT 1`).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
