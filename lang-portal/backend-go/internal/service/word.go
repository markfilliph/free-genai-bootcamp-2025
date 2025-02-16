package service

import (
	"fmt"

	"lang-portal/internal/models"
)

// WordService handles business logic for words
type WordService struct{}

// NewWordService creates a new WordService
func NewWordService() *WordService {
	return &WordService{}
}

// CreateWord creates a new word with validation
func (s *WordService) CreateWord(original, translation string) (*models.Word, error) {
	if original == "" || translation == "" {
		return nil, fmt.Errorf("original and translation are required")
	}

	word := &models.Word{
		Original:    original,
		Translation: translation,
	}

	if err := models.CreateWord(word); err != nil {
		return nil, fmt.Errorf("failed to create word: %v", err)
	}

	return word, nil
}

// GetWord retrieves a word by ID with optional stats
func (s *WordService) GetWord(wordID int64, withStats bool) (*models.Word, *models.WordReviewStats, error) {
	word, err := models.GetWord(wordID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return nil, nil, fmt.Errorf("word not found: %d", wordID)
	}

	if !withStats {
		return word, nil, nil
	}

	stats, err := models.GetWordReviewStats(wordID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get word stats: %v", err)
	}

	return word, stats, nil
}

// ListWords retrieves a paginated list of words with optional filters
func (s *WordService) ListWords(page int, search *string) ([]*models.Word, int, error) {
	if page < 1 {
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	words, total, err := models.ListWords(page, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list words: %v", err)
	}

	return words, total, nil
}

// UpdateWord updates an existing word
func (s *WordService) UpdateWord(wordID int64, original, translation string) (*models.Word, error) {
	if original == "" || translation == "" {
		return nil, fmt.Errorf("original and translation are required")
	}

	word, err := models.GetWord(wordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return nil, fmt.Errorf("word not found: %d", wordID)
	}

	word.Original = original
	word.Translation = translation

	if err := models.UpdateWord(word); err != nil {
		return nil, fmt.Errorf("failed to update word: %v", err)
	}

	return word, nil
}

// DeleteWord deletes a word and its associations
func (s *WordService) DeleteWord(wordID int64) error {
	word, err := models.GetWord(wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return fmt.Errorf("word not found: %d", wordID)
	}

	if err := models.DeleteWord(wordID); err != nil {
		return fmt.Errorf("failed to delete word: %v", err)
	}

	return nil
}

// SearchWords searches for words based on query
func (s *WordService) SearchWords(query string, page int) ([]*models.Word, int, error) {
	if page < 1 {
		page = 1
	}

	words, total, err := models.SearchWords(query, page)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search words: %v", err)
	}

	return words, total, nil
}

// GetWordGroups retrieves all groups containing a word
func (s *WordService) GetWordGroups(wordID int64) ([]*models.Group, error) {
	word, err := models.GetWord(wordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return nil, fmt.Errorf("word not found: %d", wordID)
	}

	groups, err := models.GetWordGroups(wordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get word groups: %v", err)
	}

	return groups, nil
}
