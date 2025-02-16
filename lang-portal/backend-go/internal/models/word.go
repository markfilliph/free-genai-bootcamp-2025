package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Word represents a vocabulary word in the system
type Word struct {
	ID       int64           `json:"id"`
	Japanese string         `json:"japanese"`
	Romaji   string         `json:"romaji"`
	English  string         `json:"english"`
	Parts    map[string]any `json:"parts"`
}

// GetWords retrieves a paginated list of words
func GetWords(page int) ([]Word, int, error) {
	db := GetDB()
	
	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting words: %v", err)
	}

	// Calculate offset
	offset := (page - 1) * 100
	if offset < 0 {
		offset = 0
	}

	// Get paginated words
	rows, err := db.Query(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		ORDER BY id 
		LIMIT 100 OFFSET ?`, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsJSON string
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &partsJSON)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}

		// Parse JSON parts
		if err := json.Unmarshal([]byte(partsJSON), &w.Parts); err != nil {
			return nil, 0, fmt.Errorf("error parsing parts JSON: %v", err)
		}

		words = append(words, w)
	}

	return words, total, nil
}

// GetWord retrieves a single word by ID
func GetWord(wordID int64) (*Word, error) {
	db := GetDB()

	var w Word
	var partsJSON string
	err := db.QueryRow(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		WHERE id = ?`, wordID).Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &partsJSON)
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
		INSERT INTO words (japanese, romaji, english, parts)
		VALUES (?, ?, ?, ?)`,
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
		SET japanese = ?, romaji = ?, english = ?, parts = ?
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
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
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
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &partsJSON)
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