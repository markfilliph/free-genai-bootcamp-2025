/**
 * E2E tests for the FlashcardReview component
 * 
 * These tests focus on the flashcard review functionality, which was identified
 * as a critical priority with only 1.98% coverage.
 */

import { setupAllApiMocks, mockFlashcardEndpoints, mockStudySessionEndpoints, mockApiError } from '../support/apiMocks';

describe('Flashcard Review Tests', () => {
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
    
    // Start a study session to get to the flashcard review component
    cy.startStudySession('1');
  });
  
  it('should display flashcard front correctly', () => {
    // Verify the flashcard front is displayed correctly
    cy.get('[data-testid="flashcard-front"]').should('be.visible');
    cy.contains('hola').should('be.visible');
    cy.contains('Show Answer').should('be.visible');
  });
  
  it('should show flashcard back when "Show Answer" is clicked', () => {
    // Click the "Show Answer" button
    cy.contains('Show Answer').click();
    
    // Verify the flashcard back is displayed correctly
    cy.get('[data-testid="flashcard-back"]').should('be.visible');
    cy.contains('hello').should('be.visible');
    cy.contains('¡Hola! ¿Cómo estás?').should('be.visible');
    
    // Verify rating buttons are displayed
    cy.contains('Easy').should('be.visible');
    cy.contains('Medium').should('be.visible');
    cy.contains('Hard').should('be.visible');
  });
  
  it('should handle "Easy" rating correctly', () => {
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Click the "Easy" rating button
    cy.contains('Easy').click();
    
    // Verify the rating request was made correctly
    cy.wait('@rateFlashcardRequest').then((interception) => {
      expect(interception.request.body.rating).to.equal('easy');
    });
  });
  
  it('should handle "Medium" rating correctly', () => {
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Click the "Medium" rating button
    cy.contains('Medium').click();
    
    // Verify the rating request was made correctly
    cy.wait('@rateFlashcardRequest').then((interception) => {
      expect(interception.request.body.rating).to.equal('medium');
    });
  });
  
  it('should handle "Hard" rating correctly', () => {
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Click the "Hard" rating button
    cy.contains('Hard').click();
    
    // Verify the rating request was made correctly
    cy.wait('@rateFlashcardRequest').then((interception) => {
      expect(interception.request.body.rating).to.equal('hard');
    });
  });
  
  it('should display word type and notes if available', () => {
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Verify word type and notes are displayed
    cy.contains('Word Type:').should('be.visible');
    cy.contains('noun').should('be.visible');
    cy.contains('Notes:').should('be.visible');
    cy.contains('Common greeting in Spanish').should('be.visible');
  });
  
  it('should show progress indicator', () => {
    // Verify the progress indicator is displayed
    cy.get('[data-testid="progress-indicator"]').should('be.visible');
    cy.contains('Card 1 of').should('be.visible');
  });
  
  it('should handle API errors gracefully', () => {
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
    cy.contains('Error').should('be.visible');
    cy.contains('Failed to save rating').should('be.visible');
    
    // Should still allow continuing
    cy.contains('Continue').click();
    
    // Should move to the next card
    cy.contains('Card 2 of').should('be.visible');
  });
  
  it('should be keyboard accessible', () => {
    // Test keyboard navigation
    cy.focused().type('{enter}'); // Show answer
    
    // Verify the answer is shown
    cy.contains('hello').should('be.visible');
    
    // Navigate to rating buttons using tab
    cy.focused().tab().tab(); // Navigate to first rating button
    cy.focused().should('contain', 'Easy');
    
    // Select the rating
    cy.focused().type('{enter}');
    
    // Should move to the next card
    cy.wait('@rateFlashcardRequest');
  });
  
  it('should handle cards with images correctly', () => {
    // Mock a flashcard with an image
    cy.intercept('GET', '**/decks/*/flashcards', {
      statusCode: 200,
      body: [{
        id: 'img-card-1',
        word: 'gato',
        translation: 'cat',
        examples: ['El gato es negro.'],
        notes: 'A common animal',
        wordType: 'noun',
        imageUrl: 'https://example.com/cat.jpg'
      }]
    }).as('getFlashcardsWithImageRequest');
    
    // Reload the page to get the new flashcard
    cy.reload();
    
    // Wait for the flashcard to load
    cy.wait('@getFlashcardsWithImageRequest');
    
    // Verify the image is displayed
    cy.get('img').should('have.attr', 'src', 'https://example.com/cat.jpg');
    cy.get('img').should('have.attr', 'alt', 'gato');
    
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Verify the translation is displayed
    cy.contains('cat').should('be.visible');
  });
  
  it('should pass accessibility checks', () => {
    // Check accessibility on the front of the card
    cy.checkAccessibility();
    
    // Show the answer
    cy.contains('Show Answer').click();
    
    // Check accessibility on the back of the card
    cy.checkAccessibility();
    
    // Check screen reader accessibility
    cy.checkScreenReaderAccessibility();
  });
});
