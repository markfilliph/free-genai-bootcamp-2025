package models

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupStats struct {
	TotalWords      int `json:"total_words"`
	StudiedWords    int `json:"studied_words"`
	MasteredWords   int `json:"mastered_words"`
	StudySessions   int `json:"study_sessions"`
	LastStudyDays   int `json:"last_study_days"`
	StudyStreak     int `json:"study_streak"`
	CorrectAnswers  int `json:"correct_answers"`
	TotalAnswers    int `json:"total_answers"`
	SuccessRate     int `json:"success_rate"`
	StudyTimeMinute int `json:"study_time_minute"`
}

// GetGroups retrieves all groups with optional pagination
func GetGroups(offset, limit int) ([]Group, error) {
	query := "SELECT id, name, created_at FROM groups"
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		rows, err := DB.Query(query, limit, offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanGroups(rows)
	}

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGroups(rows)
}

// scanGroups scans rows into Group structs
func scanGroups(rows *sql.Rows) ([]Group, error) {
	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetGroup retrieves a single group by ID
func GetGroup(id int64) (*Group, error) {
	var g Group
	err := DB.QueryRow("SELECT id, name, created_at FROM groups WHERE id = ?", id).Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// CreateGroup creates a new group and returns it
func CreateGroup(name string) (*Group, error) {
	result, err := DB.Exec("INSERT INTO groups (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetGroup(id)
}

// GetGroupWords retrieves all words in a group
func GetGroupWords(groupID int64) ([]Word, error) {
	rows, err := DB.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN group_items gi ON w.id = gi.word_id
		WHERE gi.group_id = ?
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, nil
}

// AddWordToGroup adds a word to a group
func AddWordToGroup(wordID, groupID int64) error {
	_, err := DB.Exec("INSERT INTO group_items (word_id, group_id) VALUES (?, ?)", wordID, groupID)
	return err
}

// RemoveWordFromGroup removes a word from a group
func RemoveWordFromGroup(wordID, groupID int64) error {
	_, err := DB.Exec("DELETE FROM group_items WHERE word_id = ? AND group_id = ?", wordID, groupID)
	return err
}

// GetGroupStats retrieves statistics for a group
func GetGroupStats(groupID int64) (*GroupStats, error) {
	var stats GroupStats

	// Get total words in group
	err := DB.QueryRow(`
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN group_items gi ON w.id = gi.word_id
		WHERE gi.group_id = ?
	`, groupID).Scan(&stats.TotalWords)
	if err != nil {
		return nil, err
	}

	// Get study session stats
	err = DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT s.id) as study_sessions,
			COUNT(DISTINCT r.word_id) as studied_words,
			SUM(CASE WHEN r.correct = 1 THEN 1 ELSE 0 END) as correct_answers,
			COUNT(r.id) as total_answers
		FROM study_sessions s
		LEFT JOIN review_items r ON s.id = r.session_id
		WHERE s.group_id = ?
	`, groupID).Scan(&stats.StudySessions, &stats.StudiedWords, &stats.CorrectAnswers, &stats.TotalAnswers)
	if err != nil {
		return nil, err
	}

	// Calculate success rate
	if stats.TotalAnswers > 0 {
		stats.SuccessRate = int(float64(stats.CorrectAnswers) / float64(stats.TotalAnswers) * 100)
	}

	return &stats, nil
}