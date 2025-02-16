package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Group represents a collection of words
type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// GroupStats represents statistics for a word group
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

// GetGroups retrieves a paginated list of groups
func GetGroups(page int) ([]Group, int, error) {
	db := GetDB()
	
	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM word_groups").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting groups: %v", err)
	}

	// Calculate offset
	offset := (page - 1) * 100
	if offset < 0 {
		offset = 0
	}

	// Get paginated groups
	rows, err := db.Query(`
		SELECT id, name, created_at 
		FROM word_groups 
		ORDER BY name 
		LIMIT 100 OFFSET ?`, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, g)
	}

	return groups, total, nil
}

// GetGroup retrieves a single group by ID
func GetGroup(groupID int64) (*Group, error) {
	db := GetDB()

	var g Group
	err := db.QueryRow(`
		SELECT id, name, created_at 
		FROM word_groups 
		WHERE id = ?`, groupID).Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying group: %v", err)
	}

	return &g, nil
}

// CreateGroup creates a new group
func CreateGroup(g *Group) error {
	db := GetDB()

	result, err := db.Exec(`
		INSERT INTO word_groups (name, created_at)
		VALUES (?, CURRENT_TIMESTAMP)`,
		g.Name)
	if err != nil {
		return fmt.Errorf("error inserting group: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	g.ID = id
	return nil
}

// UpdateGroup updates an existing group
func UpdateGroup(g *Group) error {
	db := GetDB()

	result, err := db.Exec(`
		UPDATE word_groups 
		SET name = ?
		WHERE id = ?`,
		g.Name, g.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found: %d", g.ID)
	}

	return nil
}

// DeleteGroup deletes a group and its word associations
func DeleteGroup(groupID int64) error {
	db := GetDB()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Delete word associations first
	_, err = tx.Exec("DELETE FROM words_groups WHERE group_id = ?", groupID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting word associations: %v", err)
	}

	// Delete the group
	result, err := tx.Exec("DELETE FROM word_groups WHERE id = ?", groupID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("group not found: %d", groupID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// AddWordsToGroup adds words to a group
func AddWordsToGroup(groupID int64, wordIDs []int64) error {
	if len(wordIDs) == 0 {
		return nil
	}

	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	for _, wordID := range wordIDs {
		_, err := tx.Exec(`
			INSERT INTO words_groups (word_id, group_id)
			VALUES (?, ?)
			ON CONFLICT (word_id, group_id) DO NOTHING`,
			wordID, groupID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding word %d to group: %v", wordID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// RemoveWordsFromGroup removes words from a group
func RemoveWordsFromGroup(groupID int64, wordIDs []int64) error {
	if len(wordIDs) == 0 {
		return nil
	}

	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	for _, wordID := range wordIDs {
		_, err := tx.Exec(`
			DELETE FROM words_groups 
			WHERE word_id = ? AND group_id = ?`,
			wordID, groupID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error removing word %d from group: %v", wordID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetGroupWords retrieves all words in a group
func GetGroupWords(groupID int64) ([]Word, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	rows, err := DB.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.japanese
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// GetGroupStats retrieves statistics for a group
func GetGroupStats(groupID int64) (*GroupStats, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	var stats GroupStats

	// Get total words and studied words
	err := DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT wg.word_id) as total_words,
			COUNT(DISTINCT CASE WHEN wri.id IS NOT NULL THEN wg.word_id END) as studied_words,
			COUNT(DISTINCT CASE WHEN wri.correct = 1 THEN wg.word_id END) as mastered_words,
			COUNT(DISTINCT ss.id) as study_sessions,
			COALESCE(DATEDIFF(CURRENT_TIMESTAMP, MAX(ss.created_at)), 0) as last_study_days,
			COUNT(wri.id) as total_answers,
			SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_answers
		FROM words_groups wg
		LEFT JOIN word_review_items wri ON wg.word_id = wri.word_id
		LEFT JOIN study_sessions ss ON wri.study_session_id = ss.id
		WHERE wg.group_id = ?
	`, groupID).Scan(
		&stats.TotalWords,
		&stats.StudiedWords,
		&stats.MasteredWords,
		&stats.StudySessions,
		&stats.LastStudyDays,
		&stats.TotalAnswers,
		&stats.CorrectAnswers,
	)
	if err != nil {
		return nil, err
	}

	// Calculate success rate
	if stats.TotalAnswers > 0 {
		stats.SuccessRate = (stats.CorrectAnswers * 100) / stats.TotalAnswers
	}

	// Calculate study streak
	streak, err := calculateGroupStudyStreak(groupID)
	if err != nil {
		return nil, err
	}
	stats.StudyStreak = streak

	return &stats, nil
}

// calculateGroupStudyStreak calculates the current study streak for a group
func calculateGroupStudyStreak(groupID int64) (int, error) {
	rows, err := DB.Query(`
		SELECT DATE(created_at) as study_date
		FROM study_sessions
		WHERE group_id = ?
		GROUP BY DATE(created_at)
		ORDER BY study_date DESC
	`, groupID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return 0, err
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, nil
	}

	streak := 1
	currentDate := dates[0]

	for i := 1; i < len(dates); i++ {
		expectedDate := currentDate.AddDate(0, 0, -1)
		if dates[i].Equal(expectedDate) {
			streak++
			currentDate = dates[i]
		} else {
			break
		}
	}

	return streak, nil
}