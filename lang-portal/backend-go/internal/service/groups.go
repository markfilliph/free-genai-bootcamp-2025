package service

import (
	"database/sql"
	"errors"
	"strings"
	"lang-portal/internal/models"
)

// GroupService handles business logic for group operations
type GroupService struct {
	db *sql.DB
}

// NewGroupService creates a new group service
func NewGroupService(db *sql.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// GetGroups returns all groups with additional statistics
func (s *GroupService) GetGroups(offset, limit int) ([]models.Group, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	return models.GetGroups(offset, limit)
}

// GetGroup returns details of a specific group
func (s *GroupService) GetGroup(id int64) (*models.Group, error) {
	if id <= 0 {
		return nil, errors.New("invalid group ID")
	}
	return models.GetGroup(id)
}

// GetGroupWords returns all words in a group
func (s *GroupService) GetGroupWords(groupID int64) ([]models.Word, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	// First verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	return models.GetGroupWords(groupID)
}

// GetGroupStudySessions returns study sessions for a group
func (s *GroupService) GetGroupStudySessions(groupID int64) ([]models.StudySessionResponse, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	// First verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	sessions, err := models.GetStudySessionsByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	var responses []models.StudySessionResponse
	for _, session := range sessions {
		stats, err := models.GetStudySessionStats(session.ID)
		if err != nil {
			continue
		}

		responses = append(responses, models.StudySessionResponse{
			ID:              session.ID,
			GroupID:         session.GroupID,
			GroupName:       group.Name,
			StudyActivityID: session.StudyActivityID,
			TotalWords:      stats["total"],
			CorrectWords:    stats["correct"],
			CreatedAt:       session.CreatedAt,
		})
	}

	return responses, nil
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(name string) (*models.Group, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("group name cannot be empty")
	}
	if len(name) > 255 {
		return nil, errors.New("group name too long")
	}

	// Check for duplicate name
	groups, err := models.GetGroups(0, 0)
	if err != nil {
		return nil, err
	}
	for _, g := range groups {
		if strings.EqualFold(g.Name, name) {
			return nil, errors.New("group with this name already exists")
		}
	}

	return models.CreateGroup(name)
}