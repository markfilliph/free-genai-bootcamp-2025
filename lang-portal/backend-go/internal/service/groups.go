package service

import (
	"lang-portal/internal/models"
)

// GroupServiceImpl handles business logic for group operations
type GroupServiceImpl struct{}

// NewGroupServiceImpl creates a new group service
func NewGroupServiceImpl() *GroupServiceImpl {
	return &GroupServiceImpl{}
}

// GetGroups returns all groups with additional statistics
func (s *GroupServiceImpl) GetGroups(offset, limit int) ([]map[string]interface{}, error) {
	groups, err := models.GetGroups(offset, limit)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, group := range groups {
		stats, err := models.GetGroupStats(group.ID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":              group.ID,
			"name":            group.Name,
			"total_words":     stats.TotalWords,
			"studied_words":   stats.StudiedWords,
			"mastered_words":  stats.MasteredWords,
			"study_sessions": stats.StudySessions,
		})
	}

	return result, nil
}

// GetGroup returns details of a specific group
func (s *GroupServiceImpl) GetGroup(id int64) (map[string]interface{}, error) {
	group, err := models.GetGroup(id)
	if err != nil {
		return nil, err
	}

	stats, err := models.GetGroupStats(group.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":              group.ID,
		"name":            group.Name,
		"total_words":     stats.TotalWords,
		"studied_words":   stats.StudiedWords,
		"mastered_words":  stats.MasteredWords,
		"study_sessions": stats.StudySessions,
		"success_rate":    stats.SuccessRate,
	}, nil
}

// GetGroupWords returns all words in a group with their study statistics
func (s *GroupServiceImpl) GetGroupWords(groupID int64) ([]map[string]interface{}, error) {
	words, err := models.GetGroupWords(groupID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, word := range words {
		result = append(result, map[string]interface{}{
			"id":       word.ID,
			"japanese": word.Japanese,
			"romaji":   word.Romaji,
			"english":  word.English,
			"parts":    word.Parts,
		})
	}

	return result, nil
}

// GetGroupStudySessions returns study sessions for a group with detailed statistics
func (s *GroupServiceImpl) GetGroupStudySessions(groupID int64) ([]map[string]interface{}, error) {
	// Get the group first to ensure it exists
	_, err := models.GetGroup(groupID)
	if err != nil {
		return nil, err
	}

	// Get study sessions for the group
	sessions, err := models.GetStudySessionsByGroup(groupID, 0, 100) // Default pagination: first 100 sessions
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, session := range sessions {
		// Get activity details
		if session.StudyActivityID != nil {
			_, err := models.GetStudyActivity(*session.StudyActivityID)
			if err != nil {
				continue
			}

			// Get session statistics
			stats, err := models.GetStudySessionStats(session.ID)
			if err != nil {
				continue
			}

			result = append(result, map[string]interface{}{
				"id":           session.ID,
				"activity":     "Vocabulary Quiz", // Default activity type for now
				"created_at":   session.CreatedAt,
				"total_words":  stats["total_reviews"],
				"correct":      stats["correct_reviews"],
				"success_rate": stats["success_rate"],
			})
		}
	}

	return result, nil
}

// CreateGroup creates a new group
func (s *GroupServiceImpl) CreateGroup(name string) (map[string]interface{}, error) {
	group, err := models.CreateGroup(name)
	if err != nil {
		return nil, err
	}

	stats, err := models.GetGroupStats(group.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":              group.ID,
		"name":            group.Name,
		"created_at":      group.CreatedAt,
		"total_words":     stats.TotalWords,
		"studied_words":   stats.StudiedWords,
		"mastered_words":  stats.MasteredWords,
		"study_sessions": stats.StudySessions,
		"success_rate":    stats.SuccessRate,
	}, nil
}