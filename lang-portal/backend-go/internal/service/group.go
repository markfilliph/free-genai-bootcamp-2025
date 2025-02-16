package service

import (
	"fmt"

	"lang-portal/internal/models"
)

// GroupService handles business logic for word groups
type GroupService struct{}

// NewGroupService creates a new GroupService
func NewGroupService() *GroupService {
	return &GroupService{}
}

// CreateGroup creates a new word group
func (s *GroupService) CreateGroup(name string) (*models.Group, error) {
	group := &models.Group{
		Name: name,
	}

	if err := models.CreateGroup(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(groupID int64) (*models.Group, error) {
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	return group, nil
}

// ListGroups retrieves all groups with pagination
func (s *GroupService) ListGroups(page int, search *string) ([]*models.Group, int, error) {
	if page < 1 {
		page = 1
	}

	groups, total, err := models.GetGroups(page)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list groups: %v", err)
	}

	return groups, total, nil
}

// UpdateGroup updates a group's information
func (s *GroupService) UpdateGroup(groupID int64, name string) (*models.Group, error) {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	// Update group
	group.Name = name
	if err := models.UpdateGroup(group); err != nil {
		return nil, fmt.Errorf("failed to update group: %v", err)
	}

	return group, nil
}

// DeleteGroup deletes a group and its word associations
func (s *GroupService) DeleteGroup(groupID int64) error {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return fmt.Errorf("group not found: %d", groupID)
	}

	// Delete group
	if err := models.DeleteGroup(groupID); err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	return nil
}

// AddWordToGroup adds a word to a group
func (s *GroupService) AddWordToGroup(groupID, wordID int64) error {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return fmt.Errorf("group not found: %d", groupID)
	}

	// Verify word exists
	word, err := models.GetWord(wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return fmt.Errorf("word not found: %d", wordID)
	}

	// Add word to group
	if err := models.AddWordToGroup(groupID, wordID); err != nil {
		return fmt.Errorf("failed to add word to group: %v", err)
	}

	return nil
}

// RemoveWordFromGroup removes a word from a group
func (s *GroupService) RemoveWordFromGroup(groupID, wordID int64) error {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return fmt.Errorf("group not found: %d", groupID)
	}

	// Verify word exists
	word, err := models.GetWord(wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %v", err)
	}
	if word == nil {
		return fmt.Errorf("word not found: %d", wordID)
	}

	// Remove word from group
	if err := models.RemoveWordFromGroup(groupID, wordID); err != nil {
		return fmt.Errorf("failed to remove word from group: %v", err)
	}

	return nil
}

// GetGroupWords retrieves all words in a group
func (s *GroupService) GetGroupWords(groupID int64) ([]*models.Word, error) {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	// Get group words
	words, err := models.GetGroupWords(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group words: %v", err)
	}

	return words, nil
}

// SearchGroups searches for groups based on query
func (s *GroupService) SearchGroups(query string, page int) ([]models.Group, int, error) {
	if page < 1 {
		page = 1
	}

	groups, total, err := models.SearchGroups(query, page)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search groups: %v", err)
	}

	return groups, total, nil
}

// GetGroup retrieves a group by ID with optional word list
func (s *GroupService) GetGroup(id int64, withWords bool) (*models.Group, []models.Word, error) {
	group, err := models.GetGroup(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, nil, fmt.Errorf("group not found: %d", id)
	}

	if !withWords {
		return group, nil, nil
	}

	words, err := models.GetGroupWords(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get group words: %v", err)
	}

	return group, words, nil
}

// AddWordsToGroup adds multiple words to a group
func (s *GroupService) AddWordsToGroup(groupID int64, wordIDs []int64) error {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return fmt.Errorf("group not found: %d", groupID)
	}

	// Verify all words exist
	for _, wordID := range wordIDs {
		word, err := models.GetWord(wordID)
		if err != nil {
			return fmt.Errorf("failed to get word %d: %v", wordID, err)
		}
		if word == nil {
			return fmt.Errorf("word not found: %d", wordID)
		}
	}

	if err := models.AddWordsToGroup(groupID, wordIDs); err != nil {
		return fmt.Errorf("failed to add words to group: %v", err)
	}

	return nil
}

// RemoveWordsFromGroup removes multiple words from a group
func (s *GroupService) RemoveWordsFromGroup(groupID int64, wordIDs []int64) error {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return fmt.Errorf("group not found: %d", groupID)
	}

	if err := models.RemoveWordsFromGroup(groupID, wordIDs); err != nil {
		return fmt.Errorf("failed to remove words from group: %v", err)
	}

	return nil
}

// GetGroupProgress retrieves study progress for a group
func (s *GroupService) GetGroupProgress(groupID int64) (*models.GroupProgress, error) {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	progress, err := models.GetGroupProgress(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group progress: %v", err)
	}

	return progress, nil
}
