package service

import (
	"lang-portal/internal/models"
)

// GroupService handles business logic for group operations
type GroupService struct{}

// NewGroupService creates a new group service
func NewGroupService() *GroupService {
	return &GroupService{}
}

// GetGroups returns all groups with additional statistics
func (s *GroupService) GetGroups() ([]map[string]interface{}, error) {
	groups, err := models.GetGroups()
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
func (s *GroupService) GetGroup(id int64) (map[string]interface{}, error) {
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
func (s *GroupService) GetGroupWords(groupID int64) ([]map[string]interface{}, error) {
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
func (s *GroupService) GetGroupStudySessions(groupID int64) ([]map[string]interface{}, error) {
	sessions, err := models.GetStudySessions()
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
		if session.StudyActivityID != 0 {
			activity, err := models.GetStudyActivity(session.StudyActivityID)
			if err == nil {
				activityName = "Vocabulary Quiz" // TODO: Add activity types
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