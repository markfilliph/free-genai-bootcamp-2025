package service

import (
	"fmt"
	"time"

	"lang-portal/internal/cache"
	"lang-portal/internal/models"
)

// GroupService handles business logic for word groups
type GroupService struct{}

// NewGroupService creates a new GroupService
func NewGroupService() *GroupService {
	return &GroupService{}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(name string, description string) (*models.Group, error) {
	group := &models.Group{
		Name:        name,
		Description: description,
	}

	if err := models.CreateGroup(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(groupID int64) (*models.Group, error) {
	cacheKey := fmt.Sprintf("group:%d", groupID)
	var group models.Group

	err := cache.GetOrSet(cacheKey, &group, 5*time.Minute, func() (interface{}, error) {
		return models.GetGroup(groupID)
	})

	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}

	return &group, nil
}

// ListGroups retrieves a paginated list of groups
func (s *GroupService) ListGroups(page int, search *string) ([]*models.Group, int, error) {
	if page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("groups:page:%d:search:%s", page, *search)
	var result struct {
		Groups []*models.Group
		Total  int
	}

	err := cache.GetOrSet(cacheKey, &result, 5*time.Minute, func() (interface{}, error) {
		groups, total, err := models.GetGroups(page, 10, *search)
		if err != nil {
			return nil, err
		}
		return struct {
			Groups []*models.Group
			Total  int
		}{
			Groups: groups,
			Total:  total,
		}, nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("error listing groups: %v", err)
	}

	return result.Groups, result.Total, nil
}

// UpdateGroup updates a group's information
func (s *GroupService) UpdateGroup(groupID int64, name, description string) (*models.Group, error) {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	group.Name = name
	group.Description = description

	if err := models.UpdateGroup(group); err != nil {
		return nil, fmt.Errorf("failed to update group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return group, nil
}

// DeleteGroup deletes a group
func (s *GroupService) DeleteGroup(groupID int64) error {
	if _, err := s.GetGroup(groupID); err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	if err := models.DeleteGroup(groupID); err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return nil
}

// AddWordToGroup adds a word to a group
func (s *GroupService) AddWordToGroup(groupID, wordID int64) error {
	if _, err := s.GetGroup(groupID); err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	if err := models.AddWordToGroup(groupID, wordID); err != nil {
		return fmt.Errorf("failed to add word to group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return nil
}

// RemoveWordFromGroup removes a word from a group
func (s *GroupService) RemoveWordFromGroup(groupID, wordID int64) error {
	if _, err := s.GetGroup(groupID); err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	if err := models.RemoveWordFromGroup(groupID, wordID); err != nil {
		return fmt.Errorf("failed to remove word from group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return nil
}

// GetGroupWords retrieves all words in a group
func (s *GroupService) GetGroupWords(groupID int64, page, pageSize int) ([]models.Word, int, error) {
	cacheKey := fmt.Sprintf("group_words:%d:page:%d:size:%d", groupID, page, pageSize)
	var result struct {
		Words []models.Word
		Total int
	}

	err := cache.GetOrSet(cacheKey, &result, 5*time.Minute, func() (interface{}, error) {
		words, total, err := models.GetGroupWords(groupID, page, pageSize)
		if err != nil {
			return nil, err
		}
		return struct {
			Words []models.Word
			Total int
		}{
			Words: words,
			Total: total,
		}, nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("error getting group words: %v", err)
	}

	return result.Words, result.Total, nil
}

// GetGroupProgress retrieves the progress for a group
func (s *GroupService) GetGroupProgress(groupID int64) (*models.GroupProgress, error) {
	cacheKey := fmt.Sprintf("group_progress:%d", groupID)
	var progress models.GroupProgress

	err := cache.GetOrSet(cacheKey, &progress, 5*time.Minute, func() (interface{}, error) {
		return models.GetGroupProgress(groupID)
	})

	if err != nil {
		return nil, fmt.Errorf("error getting group progress: %v", err)
	}

	return &progress, nil
}

// AddWordsToGroup adds multiple words to a group
func (s *GroupService) AddWordsToGroup(groupID int64, wordIDs []int64) error {
	if _, err := s.GetGroup(groupID); err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	if err := models.AddWordsToGroup(groupID, wordIDs); err != nil {
		return fmt.Errorf("failed to add words to group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return nil
}

// RemoveWordsFromGroup removes multiple words from a group
func (s *GroupService) RemoveWordsFromGroup(groupID int64, wordIDs []int64) error {
	if _, err := s.GetGroup(groupID); err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	if err := models.RemoveWordsFromGroup(groupID, wordIDs); err != nil {
		return fmt.Errorf("failed to remove words from group: %v", err)
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf("group:%d", groupID))

	return nil
}
