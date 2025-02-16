package service

import (
	"lang-portal/internal/models"
)

// WordService handles business logic for word operations
type WordService struct{}

// NewWordService creates a new word service
func NewWordService() *WordService {
	return &WordService{}
}

// GetWords returns all words with their study statistics
func (s *WordService) GetWords() ([]map[string]interface{}, error) {
	words, err := models.GetWords()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, word := range words {
		stats, err := models.GetWordReviewStats(word.ID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":            word.ID,
			"japanese":      word.Japanese,
			"romaji":       word.Romaji,
			"english":      word.English,
			"parts":        word.Parts,
			"total_reviews": stats["total_reviews"],
			"success_rate":  stats["success_rate"],
		})
	}

	return result, nil
}

// GetWord returns details of a specific word with its study statistics
func (s *WordService) GetWord(id int64) (map[string]interface{}, error) {
	word, err := models.GetWord(id)
	if err != nil {
		return nil, err
	}

	stats, err := models.GetWordReviewStats(word.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":            word.ID,
		"japanese":      word.Japanese,
		"romaji":       word.Romaji,
		"english":      word.English,
		"parts":        word.Parts,
		"total_reviews": stats["total_reviews"],
		"success_rate":  stats["success_rate"],
	}, nil
}

// GetWordReviews returns review history for a specific word
func (s *WordService) GetWordReviews(wordID int64) ([]map[string]interface{}, error) {
	reviews, err := models.GetWordReviews(wordID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, review := range reviews {
		session, err := models.GetStudySession(review.StudySessionID)
		if err != nil {
			continue
		}

		group, err := models.GetGroup(session.GroupID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":               review.ID,
			"word_id":          review.WordID,
			"study_session_id": review.StudySessionID,
			"group_name":       group.Name,
			"correct":          review.Correct,
			"created_at":       review.CreatedAt,
		})
	}

	return result, nil
}