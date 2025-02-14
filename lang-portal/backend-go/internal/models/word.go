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