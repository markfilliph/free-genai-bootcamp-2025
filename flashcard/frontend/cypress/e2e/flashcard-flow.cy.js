/**
 * End-to-end tests for the flashcard application flow
 * 
 * This test file simulates a complete user journey through the flashcard application,
 * from login to deck selection, card review, and session completion.
 */

describe('Flashcard Application Flow', () => {
  beforeEach(() => {
    // Mock the API responses for authentication
    cy.intercept('POST', '**/auth/login', {
      statusCode: 200,
      body: {
        token: 'fake-jwt-token',
        user: {
          id: '1',
          email: 'test@example.com',
          name: 'Test User'
        }
      }
    }).as('loginRequest');

    // Mock the API responses for decks
    cy.intercept('GET', '**/decks', {
      statusCode: 200,
      body: [
        {
          id: '1',
          name: 'Spanish Basics',
          description: 'Essential Spanish vocabulary',
          language: 'Spanish',
          createdAt: '2023-01-01T00:00:00.000Z'
        }
      ]
    }).as('getDecksRequest');

    // Mock the API responses for flashcards
    cy.intercept('GET', '**/decks/1/flashcards', {
      statusCode: 200,
      body: [
        {
          id: '1',
          word: 'hola',
          translation: 'hello',
          examples: ['¡Hola! ¿Cómo estás?'],
          notes: 'Common greeting in Spanish',
          wordType: 'noun',
          deckId: '1'
        },
        {
          id: '2',
          word: 'adiós',
          translation: 'goodbye',
          examples: ['Adiós, hasta mañana'],
          notes: 'Used when parting ways',
          wordType: 'noun',
          deckId: '1'
        }
      ]
    }).as('getFlashcardsRequest');
  });

  it('should allow a user to log in, select a deck, and review flashcards', () => {
    // Visit the login page
    cy.visit('/login');
    
    // Fill in login credentials
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the login form
    cy.get('form').submit();
    
    // Wait for login request to complete
    cy.wait('@loginRequest');
    
    // Verify redirection to home page
    cy.url().should('include', '/');
    
    // Wait for decks to load
    cy.wait('@getDecksRequest');
    
    // Click on the Spanish Basics deck
    cy.contains('Spanish Basics').click();
    
    // Wait for flashcards to load
    cy.wait('@getFlashcardsRequest');
    
    // Verify we're on the study session page
    cy.contains('Flashcards').should('be.visible');
    
    // Verify the first flashcard is displayed
    cy.contains('hola').should('be.visible');
    
    // Click the Show Answer button
    cy.contains('Show Answer').click();
    
    // Verify the translation is shown
    cy.contains('hello').should('be.visible');
    
    // Rate the card
    cy.contains('Good').click();
    
    // Verify the second flashcard is displayed
    cy.contains('adiós').should('be.visible');
    
    // Click the Show Answer button again
    cy.contains('Show Answer').click();
    
    // Rate the second card
    cy.contains('Easy').click();
    
    // Verify the session completion screen
    cy.contains('Session Complete').should('be.visible');
    
    // Verify the restart button is present
    cy.contains('Restart Session').should('be.visible');
  });

  it('should handle empty deck scenario', () => {
    // Mock an empty flashcards response
    cy.intercept('GET', '**/decks/1/flashcards', {
      statusCode: 200,
      body: []
    }).as('getEmptyFlashcardsRequest');
    
    // Login and navigate to deck
    cy.visit('/login');
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    cy.get('form').submit();
    cy.wait('@loginRequest');
    cy.wait('@getDecksRequest');
    cy.contains('Spanish Basics').click();
    
    // Wait for empty flashcards to load
    cy.wait('@getEmptyFlashcardsRequest');
    
    // Verify empty state message
    cy.contains('0 / 0').should('be.visible');
    cy.contains('Session Complete').should('be.visible');
  });

  it('should handle API errors gracefully', () => {
    // Mock a failed API response
    cy.intercept('GET', '**/decks/1/flashcards', {
      statusCode: 500,
      body: { error: 'Server error' }
    }).as('getFailedFlashcardsRequest');
    
    // Login and navigate to deck
    cy.visit('/login');
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    cy.get('form').submit();
    cy.wait('@loginRequest');
    cy.wait('@getDecksRequest');
    cy.contains('Spanish Basics').click();
    
    // Wait for failed request
    cy.wait('@getFailedFlashcardsRequest');
    
    // Verify error message is displayed
    // This will depend on how your app handles errors
    cy.contains(/error|failed|unable/i).should('be.visible');
  });
});
