package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/service"
)

// WordHandler handles word-related requests
type WordHandler struct {
	wordService *service.WordService
}

// NewWordHandler creates a new WordHandler
func NewWordHandler(ws *service.WordService) *WordHandler {
	return &WordHandler{
		wordService: ws,
	}
}

// GetWords returns a paginated list of words
func (h *WordHandler) GetWords(c *gin.Context) {
	pagination := getPaginationParams(c)
	search := c.Query("search")

	words, total, err := h.wordService.ListWords(pagination.Page, &search)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get words: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": words,
		"pagination": PaginationResponse{
			CurrentPage:  pagination.Page,
			TotalPages:  (total + pagination.PageSize - 1) / pagination.PageSize,
			TotalItems:  total,
			ItemsPerPage: pagination.PageSize,
		},
	})
}

// GetWord returns details of a specific word
func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word, err := h.wordService.GetWord(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get word: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, word)
}

// CreateWord creates a new word
func (h *WordHandler) CreateWord(c *gin.Context) {
	var req struct {
		Original    string `json:"original" binding:"required"`
		Translation string `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	word, err := h.wordService.CreateWord(req.Original, req.Translation)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create word: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, word)
}

// UpdateWord updates an existing word
func (h *WordHandler) UpdateWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var req struct {
		Original    string `json:"original" binding:"required"`
		Translation string `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	word, err := h.wordService.UpdateWord(id, req.Original, req.Translation)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update word: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, word)
}

// DeleteWord deletes a word
func (h *WordHandler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	if err := h.wordService.DeleteWord(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete word: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetWordGroups returns groups containing a word
func (h *WordHandler) GetWordGroups(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	groups, err := h.wordService.GetWordGroups(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get word groups: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, groups)
}

// ReviewWord handles word review in a study session
func (h *WordHandler) ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var request struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	review, err := h.wordService.CreateWordReview(wordID, sessionID, request.Correct)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to record word review")
		return
	}

	c.JSON(http.StatusOK, review)
}