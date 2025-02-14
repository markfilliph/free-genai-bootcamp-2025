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