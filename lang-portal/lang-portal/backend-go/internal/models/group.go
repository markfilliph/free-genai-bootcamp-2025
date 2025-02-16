package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Group represents a collection of words
type Group struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

// GetGroups retrieves a paginated list of groups with optional filtering
func GetGroups(page, pageSize int, search string) ([]*Group, int, error) {
	offset := (page - 1) * pageSize

	// Base query with pagination
	query := `
		SELECT SQL_CALC_FOUND_ROWS 
			id, name, description, created_at, updated_at
		FROM word_groups
		WHERE 1=1
	`
	args := []interface{}{}

	// Add search condition if provided
	if search != "" {
		query += " AND (name LIKE ? OR description LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Add ordering and pagination
	query += `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, offset)

	// Execute the main query
	db := GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, &g)
	}

	// Get total count
	var total int
	err = db.QueryRow("SELECT FOUND_ROWS()").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return groups, total, nil
}

// GetGroup retrieves a single group by ID
func GetGroup(groupID int64) (*Group, error) {
	db := GetDB()

	var g Group
	err := db.QueryRow(`
		SELECT id, name, description, created_at, updated_at 
		FROM word_groups 
		WHERE id = ?`, groupID).Scan(&g.ID, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt)
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
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	result, err := GetDB().Exec(`
		INSERT INTO word_groups (name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?)`,
		g.Name, g.Description, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	g.ID = id
	return nil
}

// UpdateGroup updates an existing group
func UpdateGroup(g *Group) error {
	g.UpdatedAt = time.Now()

	result, err := GetDB().Exec(`
		UPDATE word_groups 
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?`,
		g.Name, g.Description, g.UpdatedAt, g.ID)
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

// AddWordsToGroup adds multiple words to a group in a single transaction
func AddWordsToGroup(groupID int64, wordIDs []int64) error {
	// Start transaction
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement
	stmt, err := tx.Prepare(`
		INSERT IGNORE INTO words_groups (word_id, group_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Execute for each word ID
	for _, wordID := range wordIDs {
		if _, err := stmt.Exec(groupID, wordID); err != nil {
			return fmt.Errorf("error adding word %d to group: %v", wordID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// RemoveWordsFromGroup removes multiple words from a group in a single transaction
func RemoveWordsFromGroup(groupID int64, wordIDs []int64) error {
	// Start transaction
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the delete statement
	stmt, err := tx.Prepare(`
		DELETE FROM words_groups 
		WHERE group_id = ? AND word_id = ?
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Execute for each word ID
	for _, wordID := range wordIDs {
		if _, err := stmt.Exec(groupID, wordID); err != nil {
			return fmt.Errorf("error removing word %d from group: %v", wordID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetGroupWords retrieves a paginated list of words for a specific group
func GetGroupWords(groupID int64, page, pageSize int) ([]Word, int, error) {
	offset := (page - 1) * pageSize

	// Base query with pagination and joins
	query := `
		SELECT SQL_CALC_FOUND_ROWS 
			w.id, w.japanese, w.romaji, w.english, w.parts,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = true), 
				0
			) as correct_count,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = false), 
				0
			) as incorrect_count
		FROM words w
		INNER JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.japanese DESC
		LIMIT ? OFFSET ?
	`

	// Execute the main query
	db := GetDB()
	rows, err := db.Query(query, groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying group words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(
			&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts,
			&w.CorrectCount, &w.IncorrectCount,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, w)
	}

	// Get total count
	var total int
	err = db.QueryRow("SELECT FOUND_ROWS()").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return words, total, nil
}

// GetGroupStats retrieves statistics for a group
func GetGroupStats(groupID int64) (*GroupStats, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	var stats GroupStats

	// Get total words and studied words
	err := GetDB().QueryRow(`
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
	rows, err := GetDB().Query(`
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