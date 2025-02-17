package models

import (
	"database/sql"
	"encoding/json"
)

type Word struct {
	ID       int             `json:"id"`
	Japanese string          `json:"japanese"`
	Romaji   string          `json:"romaji"`
	English  string          `json:"english"`
	Parts    json.RawMessage `json:"parts"` // Store as JSON string
}

type WordWithGroups struct {
	Word
	Groups []Group `json:"groups"`
}

// GetWords retrieves a paginated list of words
func GetWords(db *sql.DB, page, perPage int) ([]Word, int, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated words
	rows, err := db.Query(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		LIMIT ? OFFSET ?`, 
		perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts); err != nil {
			return nil, 0, err
		}
		words = append(words, w)
	}

	return words, total, nil
}

// GetWord retrieves a single word by ID
func GetWord(db *sql.DB, id int) (*WordWithGroups, error) {
	var w WordWithGroups
	err := db.QueryRow(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		WHERE id = ?`,
		id).Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts)
	if err != nil {
		return nil, err
	}

	// Get associated groups
	rows, err := db.Query(`
		SELECT g.id, g.name 
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?`,
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		w.Groups = append(w.Groups, g)
	}

	return &w, nil
}

// CreateWord creates a new word
func CreateWord(db *sql.DB, word *Word) error {
	result, err := db.Exec(`
		INSERT INTO words (japanese, romaji, english, parts)
		VALUES (?, ?, ?, ?)`,
		word.Japanese, word.Romaji, word.English, word.Parts)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	word.ID = int(id)
	return nil
}

// UpdateWord updates an existing word
func UpdateWord(db *sql.DB, word *Word) error {
	_, err := db.Exec(`
		UPDATE words 
		SET japanese = ?, romaji = ?, english = ?, parts = ?
		WHERE id = ?`,
		word.Japanese, word.Romaji, word.English, word.Parts, word.ID)
	return err
}

// DeleteWord deletes a word and its group associations
func DeleteWord(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete from word_groups first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM words_groups WHERE word_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the word
	_, err = tx.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
