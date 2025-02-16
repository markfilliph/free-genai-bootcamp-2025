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
			"id":               group.ID,
			"name":            group.Name,
			"total_words":     stats["total_words"],
			"studied_words":   stats["studied_words"],
			"mastered_words":  stats["mastered_words"],
			"study_sessions": stats["study_sessions"],
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
		"id":               group.ID,
		"name":            group.Name,
		"total_words":     stats["total_words"],
		"studied_words":   stats["studied_words"],
		"mastered_words":  stats["mastered_words"],
		"study_sessions": stats["study_sessions"],
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

// GetGroupStudySessions returns study sessions for a group with detailed statistics
func (s *GroupServiceImpl) GetGroupStudySessions(groupID int64) ([]map[string]interface{}, error) {
	// Get all sessions with a large limit since we'll filter by group
	sessions, err := models.GetStudySessions(0, 1000)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, session := range sessions {
		if session.GroupID != groupID {
			continue
		}

		stats, err := models.GetStudySessionStats(session.ID)
		if err != nil {
			continue
		}

		var activityName string
		if session.StudyActivityID != nil && *session.StudyActivityID != 0 {
			activity, err := models.GetStudyActivity(*session.StudyActivityID)
			if err == nil && activity != nil {
				activityName = "Vocabulary Quiz" // Default activity type for now
			}
		}

		result = append(result, map[string]interface{}{
			"id":             session.ID,
			"activity_name":  activityName,
			"created_at":     session.CreatedAt,
			"total_words":    stats["total"],
			"correct_words":  stats["correct"],
		})
	}

	return result, nil
}

// CreateGroup creates a new group
func (s *GroupServiceImpl) CreateGroup(name string) (map[string]interface{}, error) {
	group, err := models.CreateGroup(name)
	if err != nil {
		return nil, err
	}

	// Get group statistics
	stats, err := models.GetGroupStats(group.ID)
	if err != nil {
		return nil, err
	}

	// Combine group data with statistics
	return map[string]interface{}{
		"id":            group.ID,
		"name":          group.Name,
		"created_at":    group.CreatedAt,
		"total_words":   stats["total_words"],
		"studied_words": stats["studied_words"],
	}, nil
}