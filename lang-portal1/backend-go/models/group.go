package models

import (
	"database/sql"
)

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GroupWithStats struct {
	Group
	WordCount        int     `json:"word_count"`
	StudySessionCount int     `json:"study_session_count"`
	SuccessRate      float64 `json:"success_rate"`
}

// GetGroups retrieves a paginated list of groups
func GetGroups(db *sql.DB, page, perPage int) ([]Group, int, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated groups
	rows, err := db.Query(`
		SELECT id, name 
		FROM groups 
		LIMIT ? OFFSET ?`,
		perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}

	return groups, total, nil
}

// GetGroup retrieves a single group with stats
func GetGroup(db *sql.DB, id int) (*GroupWithStats, error) {
	var g GroupWithStats
	err := db.QueryRow(`
		SELECT g.id, g.name,
			COUNT(DISTINCT w.id) as word_count,
			COUNT(DISTINCT ss.id) as study_session_count,
			COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END) * 100, 0) as success_rate
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		LEFT JOIN words w ON wg.word_id = w.id
		LEFT JOIN study_sessions ss ON g.id = ss.group_id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE g.id = ?
		GROUP BY g.id`,
		id).Scan(&g.ID, &g.Name, &g.WordCount, &g.StudySessionCount, &g.SuccessRate)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// GetGroupWords retrieves all words in a group
func GetGroupWords(db *sql.DB, groupID int) ([]Word, error) {
	rows, err := db.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?`,
		groupID)
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
func CreateGroup(db *sql.DB, group *Group) error {
	result, err := db.Exec(`
		INSERT INTO groups (name)
		VALUES (?)`,
		group.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	group.ID = int(id)
	return nil
}

// UpdateGroup updates an existing group
func UpdateGroup(db *sql.DB, group *Group) error {
	_, err := db.Exec(`
		UPDATE groups 
		SET name = ?
		WHERE id = ?`,
		group.Name, group.ID)
	return err
}

// DeleteGroup deletes a group and its word associations
func DeleteGroup(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete from word_groups first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM words_groups WHERE group_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the group
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
