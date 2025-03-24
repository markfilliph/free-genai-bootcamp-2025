/**
 * API mocking utilities for Cypress tests
 * This file contains functions to consistently mock API responses
 */

import { 
  mockUser, 
  mockToken, 
  mockDecks, 
  mockFlashcards, 
  mockStudySession,
  mockStudySessionResults,
  mockErrors
} from '../fixtures/mockData';

/**
 * Setup all API mocks for a complete test run
 * This is useful for tests that need to interact with multiple endpoints
 */
export const setupAllApiMocks = () => {
  // Authentication endpoints
  mockAuthEndpoints();
  
  // Deck endpoints
  mockDeckEndpoints();
  
  // Flashcard endpoints
  mockFlashcardEndpoints();
  
  // Study session endpoints
  mockStudySessionEndpoints();
};

/**
 * Mock all authentication-related endpoints
 */
export const mockAuthEndpoints = () => {
  // Login endpoint
  cy.intercept('POST', '**/auth/login', {
    statusCode: 200,
    body: {
      token: mockToken,
      user: mockUser
    }
  }).as('loginRequest');
  
  // Register endpoint
  cy.intercept('POST', '**/auth/register', {
    statusCode: 201,
    body: {
      message: 'User registered successfully',
      user: {
        ...mockUser,
        email: 'newuser@example.com'
      }
    }
  }).as('registerRequest');
  
  // Forgot password endpoint
  cy.intercept('POST', '**/auth/forgot-password', {
    statusCode: 200,
    body: {
      message: 'Password reset email sent'
    }
  }).as('forgotPasswordRequest');
  
  // Reset password endpoint
  cy.intercept('POST', '**/auth/reset-password', {
    statusCode: 200,
    body: {
      message: 'Password reset successfully'
    }
  }).as('resetPasswordRequest');
  
  // Validate token endpoint
  cy.intercept('GET', '**/auth/validate', {
    statusCode: 200,
    body: {
      valid: true,
      user: mockUser
    }
  }).as('validateTokenRequest');
  
  // Logout endpoint
  cy.intercept('POST', '**/auth/logout', {
    statusCode: 200,
    body: {
      message: 'Logged out successfully'
    }
  }).as('logoutRequest');
};

/**
 * Mock all deck-related endpoints
 */
export const mockDeckEndpoints = () => {
  // Get all decks
  cy.intercept('GET', '**/decks', {
    statusCode: 200,
    body: mockDecks
  }).as('getDecksRequest');
  
  // Get a specific deck
  cy.intercept('GET', '**/decks/*', (req) => {
    const deckId = req.url.split('/').pop();
    const deck = mockDecks.find(d => d.id === deckId);
    
    if (deck) {
      req.reply({
        statusCode: 200,
        body: deck
      });
    } else {
      req.reply(mockErrors.notFound);
    }
  }).as('getDeckRequest');
  
  // Create a new deck
  cy.intercept('POST', '**/decks', (req) => {
    const newDeck = {
      id: `deck-${Date.now()}`,
      ...req.body,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      userId: mockUser.id,
      flashcardCount: 0
    };
    
    req.reply({
      statusCode: 201,
      body: newDeck
    });
  }).as('createDeckRequest');
  
  // Update a deck
  cy.intercept('PUT', '**/decks/*', (req) => {
    const deckId = req.url.split('/').pop();
    const deck = mockDecks.find(d => d.id === deckId);
    
    if (deck) {
      const updatedDeck = {
        ...deck,
        ...req.body,
        updatedAt: new Date().toISOString()
      };
      
      req.reply({
        statusCode: 200,
        body: updatedDeck
      });
    } else {
      req.reply(mockErrors.notFound);
    }
  }).as('updateDeckRequest');
  
  // Delete a deck
  cy.intercept('DELETE', '**/decks/*', {
    statusCode: 200,
    body: {
      message: 'Deck deleted successfully'
    }
  }).as('deleteDeckRequest');
};

/**
 * Mock all flashcard-related endpoints
 */
export const mockFlashcardEndpoints = () => {
  // Get flashcards for a deck
  cy.intercept('GET', '**/decks/*/flashcards', (req) => {
    const deckId = req.url.split('/')[req.url.split('/').indexOf('decks') + 1];
    const cards = mockFlashcards[deckId] || [];
    
    req.reply({
      statusCode: 200,
      body: cards
    });
  }).as('getFlashcardsRequest');
  
  // Get a specific flashcard
  cy.intercept('GET', '**/flashcards/*', (req) => {
    const cardId = req.url.split('/').pop();
    let foundCard = null;
    
    // Search through all decks for the card
    Object.values(mockFlashcards).forEach(cards => {
      const card = cards.find(c => c.id === cardId);
      if (card) foundCard = card;
    });
    
    if (foundCard) {
      req.reply({
        statusCode: 200,
        body: foundCard
      });
    } else {
      req.reply(mockErrors.notFound);
    }
  }).as('getFlashcardRequest');
  
  // Create a new flashcard
  cy.intercept('POST', '**/flashcards', (req) => {
    const newCard = {
      id: `card-${Date.now()}`,
      ...req.body,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    
    req.reply({
      statusCode: 201,
      body: newCard
    });
  }).as('createFlashcardRequest');
  
  // Update a flashcard
  cy.intercept('PUT', '**/flashcards/*', (req) => {
    const cardId = req.url.split('/').pop();
    let foundCard = null;
    
    // Search through all decks for the card
    Object.values(mockFlashcards).forEach(cards => {
      const card = cards.find(c => c.id === cardId);
      if (card) foundCard = card;
    });
    
    if (foundCard) {
      const updatedCard = {
        ...foundCard,
        ...req.body,
        updatedAt: new Date().toISOString()
      };
      
      req.reply({
        statusCode: 200,
        body: updatedCard
      });
    } else {
      req.reply(mockErrors.notFound);
    }
  }).as('updateFlashcardRequest');
  
  // Delete a flashcard
  cy.intercept('DELETE', '**/flashcards/*', {
    statusCode: 200,
    body: {
      message: 'Flashcard deleted successfully'
    }
  }).as('deleteFlashcardRequest');
  
  // Bulk operations on flashcards
  cy.intercept('POST', '**/flashcards/bulk', {
    statusCode: 200,
    body: {
      message: 'Bulk operation completed successfully',
      affectedCount: 3
    }
  }).as('bulkFlashcardRequest');
};

/**
 * Mock all study session-related endpoints
 */
export const mockStudySessionEndpoints = () => {
  // Start a study session
  cy.intercept('POST', '**/study-sessions', {
    statusCode: 201,
    body: mockStudySession
  }).as('startStudySessionRequest');
  
  // Get a study session
  cy.intercept('GET', '**/study-sessions/*', {
    statusCode: 200,
    body: mockStudySession
  }).as('getStudySessionRequest');
  
  // Update a study session (rate a card)
  cy.intercept('PUT', '**/study-sessions/*/flashcards/*', (req) => {
    const sessionId = req.url.split('/')[req.url.split('/').indexOf('study-sessions') + 1];
    const cardId = req.url.split('/').pop();
    
    req.reply({
      statusCode: 200,
      body: {
        message: 'Flashcard rating saved',
        sessionId,
        flashcardId: cardId,
        rating: req.body.rating
      }
    });
  }).as('rateFlashcardRequest');
  
  // Complete a study session
  cy.intercept('PUT', '**/study-sessions/*/complete', {
    statusCode: 200,
    body: mockStudySessionResults
  }).as('completeStudySessionRequest');
  
  // Get study statistics
  cy.intercept('GET', '**/users/*/statistics', {
    statusCode: 200,
    body: {
      totalStudySessions: 10,
      totalCardsStudied: 150,
      averageRating: 1.8,
      sessionsPerDay: [
        { date: '2023-01-01', count: 2 },
        { date: '2023-01-02', count: 3 },
        { date: '2023-01-03', count: 1 },
        { date: '2023-01-04', count: 4 }
      ],
      ratingDistribution: {
        easy: 80,
        medium: 50,
        hard: 20
      }
    }
  }).as('getUserStatisticsRequest');
};

/**
 * Mock API error responses
 * @param {string} endpoint - The endpoint to mock
 * @param {string} errorType - The type of error to mock (unauthorized, notFound, badRequest, serverError)
 */
export const mockApiError = (endpoint, errorType = 'serverError') => {
  cy.intercept(endpoint, mockErrors[errorType]).as(`${errorType}ErrorRequest`);
};
