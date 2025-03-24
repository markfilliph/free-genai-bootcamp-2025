/**
 * E2E tests for the StudySession component
 * 
 * These tests focus on the study session functionality, which was identified
 * as a critical priority with only 1.93% coverage.
 */

import { setupAllApiMocks, mockDeckEndpoints, mockFlashcardEndpoints, mockStudySessionEndpoints, mockApiError } from '../support/apiMocks';

describe('Study Session Tests', () => {
  beforeEach(() => {
    // Setup all API mocks before each test
    setupAllApiMocks();
    
    // Set up localStorage with user session
    cy.window().then((win) => {
      win.localStorage.setItem('token', 'fake-jwt-token');
      win.localStorage.setItem('user', JSON.stringify({
        id: '1',
        name: 'Test User',
        email: 'test@example.com'
      }));
    });
  });
  
  it('should start a study session successfully', () => {
    // Start a study session for deck 1
    cy.startStudySession('1');
    
    // Verify the flashcard review component is displayed
    cy.contains('Flashcard Review').should('be.visible');
    cy.contains('Show Answer').should('be.visible');
  });
  
  it('should display flashcard content correctly', () => {
    // Start a study session
    cy.startStudySession('1');
    
    // Check that the flashcard displays the word
    cy.contains('hola').should('be.visible');
    
    // Click show answer
    cy.contains('Show Answer').click();
    
    // Verify translation and examples are shown
    cy.contains('hello').should('be.visible');
    cy.contains('¡Hola! ¿Cómo estás?').should('be.visible');
  });
  
  it('should process user ratings correctly', () => {
    // Start a study session
    cy.startStudySession('1');
    
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Rate the card as "Easy"
    cy.contains('Easy').click();
    
    // Verify the rating was processed
    cy.wait('@rateFlashcardRequest').then((interception) => {
      expect(interception.request.body.rating).to.equal('easy');
    });
  });
  
  it('should complete a full study session', () => {
    // Start a study session
    cy.startStudySession('1');
    
    // Complete the session with specific ratings
    cy.completeFlashcardSession(['Easy', 'Hard', 'Medium']);
    
    // Verify we reached the completion screen
    cy.contains('Session Complete').should('be.visible');
    cy.contains('Session Statistics').should('be.visible');
  });
  
  it('should display session statistics after completion', () => {
    // Start and complete a study session
    cy.startStudySession('1');
    cy.completeFlashcardSession();
    
    // Verify statistics are displayed
    cy.verifySessionStats();
    
    // Check specific statistics
    cy.contains('Cards Reviewed').should('be.visible');
    cy.contains('Performance').should('be.visible');
    
    // Verify the "Study Again" button is available
    cy.contains('Study Again').should('be.visible');
  });
  
  it('should handle empty decks gracefully', () => {
    // Mock an empty deck of flashcards
    cy.intercept('GET', '**/decks/*/flashcards', {
      statusCode: 200,
      body: []
    }).as('getEmptyFlashcardsRequest');
    
    // Try to start a study session
    cy.visit('/decks/1', { failOnStatusCode: false });
    cy.contains('Study').click();
    
    // Should show a message about empty deck
    cy.contains('No flashcards available').should('be.visible');
  });
  
  it('should handle API errors during study session', () => {
    // Start a session
    cy.startStudySession('1');
    
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Mock an error for the rating endpoint
    cy.intercept('PUT', '**/study-sessions/*/flashcards/*', {
      statusCode: 500,
      body: {
        error: 'Internal Server Error',
        message: 'Failed to save rating'
      }
    }).as('ratingErrorRequest');
    
    // Try to rate the card
    cy.contains('Easy').click();
    
    // Should show an error message
    cy.contains('Failed to save rating').should('be.visible');
  });
  
  it('should be accessible via keyboard navigation', () => {
    // Start a study session
    cy.startStudySession('1');
    
    // Test keyboard navigation
    cy.focused().type('{enter}'); // Show answer
    cy.contains('hello').should('be.visible'); // Answer should be visible
    
    // Navigate to a rating button and press it
    cy.focused().tab().tab(); // Navigate to first rating button
    cy.focused().type('{enter}'); // Select rating
    
    // Should move to next card or completion screen
    cy.wait('@rateFlashcardRequest');
  });
  
  it('should pass accessibility checks', () => {
    // Start a study session
    cy.startStudySession('1');
    
    // Check accessibility
    cy.checkAccessibility();
    
    // Show answer and check again
    cy.contains('Show Answer').click();
    cy.checkAccessibility();
    
    // Complete session and check completion screen
    cy.completeFlashcardSession();
    cy.checkAccessibility();
  });
});
