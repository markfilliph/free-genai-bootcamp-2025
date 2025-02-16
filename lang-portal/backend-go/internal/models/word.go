package models

type Word struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts,omitempty"`
}

// GetWords returns all words
func GetWords() ([]Word, error) {
	rows, err := DB.Query(`
		SELECT id, japanese, romaji, english, parts
		FROM words
		ORDER BY japanese
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts)
		if err != nil {
			return nil, err
		}
		words = append(words, w)
	}

	return words, nil
}

// GetWord returns a single word by ID
func GetWord(id int64) (*Word, error) {
	var w Word
	err := DB.QueryRow(`
		SELECT id, japanese, romaji, english, parts
		FROM words
		WHERE id = ?
	`, id).Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts)
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

	return GetWord(id)
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

	// Delete word reviews
	_, err = tx.Exec(`
		DELETE FROM word_review_items
		WHERE word_id = ?
	`, id)
	if err != nil {
		return err
	}

	// Delete the word
	_, err = tx.Exec(`
		DELETE FROM words
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}