package models

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name string `json:"name" binding:"required,name"`
}

// UpdateGroupRequest represents the request to update a group
type UpdateGroupRequest struct {
	Name string `json:"name" binding:"required,name"`
}

// CreateWordRequest represents the request to create a new word
type CreateWordRequest struct {
	Original    string `json:"original" binding:"required,word"`
	Translation string `json:"translation" binding:"required,word"`
}

// UpdateWordRequest represents the request to update a word
type UpdateWordRequest struct {
	Original    string `json:"original" binding:"required,word"`
	Translation string `json:"translation" binding:"required,word"`
}

// StartStudySessionRequest represents the request to start a new study session
type StartStudySessionRequest struct {
	GroupID int64 `json:"group_id" binding:"required,min=1"`
}

// ReviewWordRequest represents the request to review a word in a study session
type ReviewWordRequest struct {
	Correct bool `json:"correct" binding:"required"`
}

// AddWordToGroupRequest represents the request to add a word to a group
type AddWordToGroupRequest struct {
	WordID int64 `json:"word_id" binding:"required,min=1"`
}

// BatchAddWordsRequest represents the request to add multiple words to a group
type BatchAddWordsRequest struct {
	WordIDs []int64 `json:"word_ids" binding:"required,min=1,dive,min=1"`
}

// BatchRemoveWordsRequest represents the request to remove multiple words from a group
type BatchRemoveWordsRequest struct {
	WordIDs []int64 `json:"word_ids" binding:"required,min=1,dive,min=1"`
}
