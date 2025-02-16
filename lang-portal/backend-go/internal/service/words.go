package service

import (
	"database/sql"
	"errors"
	"strings"
	"lang-portal/internal/models"
)

// WordService handles business logic for word operations
type WordService struct {
	db *sql.DB
}

// NewWordService creates a new word service
func NewWordService(db *sql.DB) *WordService {
	return &WordService{
		db: db,
	}
}

// GetWords returns all words with their study statistics
func (s *WordService) GetWords() ([]models.WordResponse, error) {
	words, err := models.GetWords()
	if err != nil {
		return nil, err
	}

	var responses []models.WordResponse
	for _, word := range words {
		stats, err := models.GetWordReviewStats(word.ID)
		if err != nil {
			// Use default stats if error
			stats = map[string]int{
				"total_reviews": 0,
				"success_rate":  0,
			}
		}

		responses = append(responses, models.WordResponse{
			ID:           word.ID,
			Japanese:     word.Japanese,
			Romaji:       word.Romaji,
			English:      word.English,
			Parts:        word.Parts,
			TotalReviews: stats["total_reviews"],
			SuccessRate:  float64(stats["success_rate"]),
		})
	}

	return responses, nil
}

// GetWord returns details of a specific word with its study statistics
func (s *WordService) GetWord(id int64) (*models.WordResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid word ID")
	}

	word, err := models.GetWord(id)
	if err != nil {
		return nil, err
	}
	if word == nil {
		return nil, errors.New("word not found")
	}

	stats, err := models.GetWordReviewStats(word.ID)
	if err != nil {
		// Use default stats if error
		stats = map[string]int{
			"total_reviews": 0,
			"success_rate":  0,
		}
	}

	return &models.WordResponse{
		ID:           word.ID,
		Japanese:     word.Japanese,
		Romaji:       word.Romaji,
		English:      word.English,
		Parts:        word.Parts,
		TotalReviews: stats["total_reviews"],
		SuccessRate:  float64(stats["success_rate"]),
	}, nil
}

// CreateWord creates a new word with validation
func (s *WordService) CreateWord(japanese, romaji, english, parts string) (*models.Word, error) {
	// Validate input
	japanese = strings.TrimSpace(japanese)
	romaji = strings.TrimSpace(romaji)
	english = strings.TrimSpace(english)

	if japanese == "" {
		return nil, errors.New("japanese text cannot be empty")
	}
	if romaji == "" {
		return nil, errors.New("romaji cannot be empty")
	}
	if english == "" {
		return nil, errors.New("english translation cannot be empty")
	}

	// Check for duplicate Japanese text
	words, err := models.GetWords()
	if err != nil {
		return nil, err
	}
	for _, w := range words {
		if w.Japanese == japanese {
			return nil, errors.New("word with this Japanese text already exists")
		}
	}

	return models.CreateWord(japanese, romaji, english, parts)
}

// GetWordReviews returns review history for a specific word
func (s *WordService) GetWordReviews(wordID int64) ([]models.WordReviewResponse, error) {
	if wordID <= 0 {
		return nil, errors.New("invalid word ID")
	}

	// First verify word exists
	word, err := models.GetWord(wordID)
	if err != nil {
		return nil, err
	}
	if word == nil {
		return nil, errors.New("word not found")
	}

	reviews, err := models.GetWordReviews(wordID)
	if err != nil {
		return nil, err
	}

	var responses []models.WordReviewResponse
	for _, review := range reviews {
		session, err := models.GetStudySession(review.StudySessionID)
		if err != nil {
			continue
		}

		group, err := models.GetGroup(session.GroupID)
		if err != nil {
			continue
		}

		responses = append(responses, models.WordReviewResponse{
			ID:              review.ID,
			WordID:          review.WordID,
			StudySessionID:  review.StudySessionID,
			GroupName:       group.Name,
			Correct:         review.Correct,
			CreatedAt:       review.CreatedAt,
		})
	}

	return responses, nil
}