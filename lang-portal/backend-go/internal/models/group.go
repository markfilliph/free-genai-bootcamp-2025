package models

import (
	"database/sql"
	"time"

	"lang-portal/internal/database"
)

// Group represents a word group
type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// GetGroups returns all groups with pagination
func GetGroups(offset, limit int) ([]Group, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	query := "SELECT id, name, created_at FROM word_groups"
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
	}

	var rows *sql.Rows
	var queryErr error
	if limit > 0 {
		rows, queryErr = db.Query(query, limit, offset)
	} else {
		rows, queryErr = db.Query(query)
	}
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

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

// GetGroup returns a single group by ID
func GetGroup(id int64) (*Group, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	var g Group
	err = db.QueryRow("SELECT id, name, created_at FROM word_groups WHERE id = ?", id).Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// GetGroupWords returns all words in a group
func GetGroupWords(groupID int64) ([]Word, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	rows, err := db.Query(query, groupID)
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

// CreateGroup creates a new group
func CreateGroup(name string) (*Group, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec("INSERT INTO word_groups (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Group{ID: id, Name: name}, nil
}

// UpdateGroup updates an existing group
func UpdateGroup(id int64, name string) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE word_groups SET name = ? WHERE id = ?", name, id)
	return err
}

// DeleteGroup deletes a group and its relationships
func DeleteGroup(id int64) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete from words_groups
	_, err = tx.Exec("DELETE FROM words_groups WHERE group_id = ?", id)
	if err != nil {
		return err
	}

	// Delete from study_activities
	_, err = tx.Exec("DELETE FROM study_activities WHERE group_id = ?", id)
	if err != nil {
		return err
	}

	// Delete from study_sessions
	_, err = tx.Exec("DELETE FROM study_sessions WHERE group_id = ?", id)
	if err != nil {
		return err
	}

	// Delete the group
	_, err = tx.Exec("DELETE FROM word_groups WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// AddWordToGroup adds a word to a group
func AddWordToGroup(wordID, groupID int64) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO words_groups (word_id, group_id)
		VALUES (?, ?)
	`, wordID, groupID)
	return err
}

// RemoveWordFromGroup removes a word from a group
func RemoveWordFromGroup(wordID, groupID int64) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE FROM words_groups 
		WHERE word_id = ? AND group_id = ?
	`, wordID, groupID)
	return err
}

// GetGroupStats returns statistics about a group
func GetGroupStats(groupID int64) (map[string]int, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	var totalWords, studiedWords int
	
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM words_groups 
		WHERE group_id = ?
	`, groupID).Scan(&totalWords)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
		SELECT COUNT(DISTINCT wri.word_id)
		FROM word_review_items wri
		JOIN words_groups wg ON wri.word_id = wg.word_id
		WHERE wg.group_id = ?
	`, groupID).Scan(&studiedWords)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_words": totalWords,
		"studied_words": studiedWords,
	}, nil
}