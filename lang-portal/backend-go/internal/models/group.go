package models

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetGroups returns all groups
func GetGroups() ([]Group, error) {
	rows, err := DB.Query("SELECT id, name FROM groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetGroup returns a single group by ID
func GetGroup(id int64) (*Group, error) {
	var g Group
	err := DB.QueryRow("SELECT id, name FROM groups WHERE id = ?", id).Scan(&g.ID, &g.Name)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// GetGroupWords returns all words in a group
func GetGroupWords(groupID int64) ([]Word, error) {
	query := `
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	rows, err := DB.Query(query, groupID)
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
	result, err := DB.Exec("INSERT INTO groups (name) VALUES (?)", name)
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
	_, err := DB.Exec("UPDATE groups SET name = ? WHERE id = ?", name, id)
	return err
}

// DeleteGroup deletes a group and its relationships
func DeleteGroup(id int64) error {
	tx, err := DB.Begin()
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
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// AddWordToGroup adds a word to a group
func AddWordToGroup(wordID, groupID int64) error {
	_, err := DB.Exec(`
		INSERT INTO words_groups (word_id, group_id)
		VALUES (?, ?)
	`, wordID, groupID)
	return err
}

// RemoveWordFromGroup removes a word from a group
func RemoveWordFromGroup(wordID, groupID int64) error {
	_, err := DB.Exec(`
		DELETE FROM words_groups 
		WHERE word_id = ? AND group_id = ?
	`, wordID, groupID)
	return err
}

// GetGroupStats returns statistics about a group
func GetGroupStats(groupID int64) (map[string]int, error) {
	var totalWords, studiedWords int
	
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM words_groups 
		WHERE group_id = ?
	`, groupID).Scan(&totalWords)
	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(`
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