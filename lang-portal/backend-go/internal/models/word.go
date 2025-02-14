package models

type Word struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts,omitempty"` // JSON field
}

// GetWords returns all words
func GetWords() ([]Word, error) {
	rows, err := DB.Query("SELECT id, japanese, romaji, english, parts FROM words")
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

// GetWord returns a single word by ID
func GetWord(id int64) (*Word, error) {
	var w Word
	err := DB.QueryRow("SELECT id, japanese, romaji, english, parts FROM words WHERE id = ?", id).Scan(
		&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// ReviewWord updates the word review status
func ReviewWord(wordID int64, studySessionID int64, correct bool) error {
	_, err := DB.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct)
		VALUES (?, ?, ?)
	`, wordID, studySessionID, correct)
	return err
}

// CreateWord creates a new word
func CreateWord(japanese, romaji, english, parts string) (*Word, error) {
	result, err := DB.Exec(`
		INSERT INTO words (japanese, romaji, english, parts)
		VALUES (?, ?, ?, ?)
	`, japanese, romaji, english, parts)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Word{
		ID:       id,
		Japanese: japanese,
		Romaji:   romaji,
		English:  english,
		Parts:    parts,
	}, nil
}

// UpdateWord updates an existing word
func UpdateWord(id int64, japanese, romaji, english, parts string) error {
	_, err := DB.Exec(`
		UPDATE words 
		SET japanese = ?, romaji = ?, english = ?, parts = ?
		WHERE id = ?
	`, japanese, romaji, english, parts, id)
	return err
}

// DeleteWord deletes a word and its relationships
func DeleteWord(id int64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete from words_groups
	_, err = tx.Exec("DELETE FROM words_groups WHERE word_id = ?", id)
	if err != nil {
		return err
	}

	// Delete from word_review_items
	_, err = tx.Exec("DELETE FROM word_review_items WHERE word_id = ?", id)
	if err != nil {
		return err
	}

	// Delete the word
	_, err = tx.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetWordGroups returns all groups that contain this word
func GetWordGroups(wordID int64) ([]Group, error) {
	query := `
		SELECT g.id, g.name
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
	`
	rows, err := DB.Query(query, wordID)
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

// GetWordReviewStats returns review statistics for a word
func GetWordReviewStats(wordID int64) (map[string]int, error) {
	var totalReviews, correctReviews int

	err := DB.QueryRow(`
		SELECT COUNT(*), SUM(CASE WHEN correct THEN 1 ELSE 0 END)
		FROM word_review_items
		WHERE word_id = ?
	`, wordID).Scan(&totalReviews, &correctReviews)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_reviews":   totalReviews,
		"correct_reviews": correctReviews,
	}, nil
}