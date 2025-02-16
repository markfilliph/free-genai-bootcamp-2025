package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Word represents a vocabulary word in the system
type Word struct {
	ID              int64     `json:"id"`
	Japanese        string    `json:"japanese"`
	Romaji          string    `json:"romaji"`
	English         string    `json:"english"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CorrectCount    int       `json:"correct_count"`
	IncorrectCount  int       `json:"incorrect_count"`
	Parts           map[string]any `json:"parts"`
}

// GetWords retrieves a paginated list of words with optional filtering
func GetWords(page, pageSize int, search string) ([]Word, int, error) {
	offset := (page - 1) * pageSize

	// Base query with pagination and review counts
	query := `
		SELECT SQL_CALC_FOUND_ROWS 
			w.id, w.japanese, w.romaji, w.english, w.created_at, w.updated_at,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = true), 
				0
			) as correct_count,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = false), 
				0
			) as incorrect_count,
			w.parts
		FROM words w
		WHERE 1=1
	`
	args := []interface{}{}

	// Add search condition if provided
	if search != "" {
		query += " AND (w.japanese LIKE ? OR w.romaji LIKE ? OR w.english LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Add ordering and pagination
	query += `
		ORDER BY w.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, offset)

	// Execute the main query
	db := GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsJSON string
		if err := rows.Scan(
			&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt,
			&w.CorrectCount, &w.IncorrectCount, &partsJSON,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}

		// Parse JSON parts
		if err := json.Unmarshal([]byte(partsJSON), &w.Parts); err != nil {
			return nil, 0, fmt.Errorf("error parsing parts JSON: %v", err)
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

// GetWordsByIDs retrieves multiple words by their IDs in a single query
func GetWordsByIDs(wordIDs []int64) ([]Word, error) {
	// Convert []int64 to string for IN clause
	query := `
		SELECT 
			w.id, w.japanese, w.romaji, w.english, w.created_at, w.updated_at,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = true), 
				0
			) as correct_count,
			COALESCE(
				(SELECT COUNT(*) FROM word_reviews wr 
				WHERE wr.word_id = w.id AND wr.correct = false), 
				0
			) as incorrect_count,
			w.parts
		FROM words w
		WHERE w.id IN (?)
		ORDER BY w.created_at DESC
	`

	// Use a placeholder for each ID
	placeholders := make([]string, len(wordIDs))
	args := make([]interface{}, len(wordIDs))
	for i, id := range wordIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// Replace the (?) with actual placeholders
	query = fmt.Sprintf(query, joinStrings(placeholders, ","))

	// Execute query
	db := GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying words by IDs: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsJSON string
		if err := rows.Scan(
			&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt,
			&w.CorrectCount, &w.IncorrectCount, &partsJSON,
		); err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}

		// Parse JSON parts
		if err := json.Unmarshal([]byte(partsJSON), &w.Parts); err != nil {
			return nil, fmt.Errorf("error parsing parts JSON: %v", err)
		}

		words = append(words, w)
	}

	return words, nil
}

// AddWordReview adds a review for a word in a study session
func AddWordReview(wordID, studySessionID int64, correct bool) error {
	// Start transaction
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Add the review
	_, err = tx.Exec(`
		INSERT INTO word_reviews (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, NOW())
	`, wordID, studySessionID, correct)
	if err != nil {
		return fmt.Errorf("error adding word review: %v", err)
	}

	// Update word statistics
	_, err = tx.Exec(`
		UPDATE words w
		SET w.updated_at = NOW()
		WHERE w.id = ?
	`, wordID)
	if err != nil {
		return fmt.Errorf("error updating word statistics: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Helper function to join strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// GetWord retrieves a single word by ID
func GetWord(wordID int64) (*Word, error) {
	db := GetDB()

	var w Word
	var partsJSON string
	err := db.QueryRow(`
		SELECT id, japanese, romaji, english, created_at, updated_at, parts
		FROM words 
		WHERE id = ?`, wordID).Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt, &partsJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying word: %v", err)
	}

	// Parse JSON parts
	if err := json.Unmarshal([]byte(partsJSON), &w.Parts); err != nil {
		return nil, fmt.Errorf("error parsing parts JSON: %v", err)
	}

	return &w, nil
}

// CreateWord creates a new word
func CreateWord(w *Word) error {
	db := GetDB()

	// Convert parts to JSON
	partsJSON, err := json.Marshal(w.Parts)
	if err != nil {
		return fmt.Errorf("error marshaling parts: %v", err)
	}

	result, err := db.Exec(`
		INSERT INTO words (japanese, romaji, english, created_at, updated_at, parts)
		VALUES (?, ?, ?, NOW(), NOW(), ?)`,
		w.Japanese, w.Romaji, w.English, partsJSON)
	if err != nil {
		return fmt.Errorf("error inserting word: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	w.ID = id
	return nil
}

// UpdateWord updates an existing word
func UpdateWord(w *Word) error {
	db := GetDB()

	// Convert parts to JSON
	partsJSON, err := json.Marshal(w.Parts)
	if err != nil {
		return fmt.Errorf("error marshaling parts: %v", err)
	}

	result, err := db.Exec(`
		UPDATE words 
		SET japanese = ?, romaji = ?, english = ?, updated_at = NOW(), parts = ?
		WHERE id = ?`,
		w.Japanese, w.Romaji, w.English, partsJSON, w.ID)
	if err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found: %d", w.ID)
	}

	return nil
}

// DeleteWord deletes a word by ID
func DeleteWord(wordID int64) error {
	db := GetDB()

	result, err := db.Exec("DELETE FROM words WHERE id = ?", wordID)
	if err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found: %d", wordID)
	}

	return nil
}

// GetWordsByGroupID retrieves all words in a group
func GetWordsByGroupID(groupID int64) ([]Word, error) {
	db := GetDB()

	rows, err := db.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.created_at, w.updated_at, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.id`, groupID)
	if err != nil {
		return nil, fmt.Errorf("error querying words by group: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsJSON string
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt, &partsJSON)
		if err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}

		// Parse JSON parts
		if err := json.Unmarshal([]byte(partsJSON), &w.Parts); err != nil {
			return nil, fmt.Errorf("error parsing parts JSON: %v", err)
		}

		words = append(words, w)
	}

	return words, nil
}