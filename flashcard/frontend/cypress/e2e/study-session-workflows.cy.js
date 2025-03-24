/**
 * End-to-end tests for study session workflows
 * 
 * This test file covers all study session related operations:
 * - Starting a study session
 * - Rating cards (easy, good, difficult)
 * - Completing a session
 * - Viewing session statistics
 */

describe('Study Session Workflows', () => {
  beforeEach(() => {
    // Log in, load decks, and navigate to a specific deck
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    
    // Load a deck with more flashcards for thorough testing
    const testFlashcards = [
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
      },
      {
        id: '3',
        word: 'gracias',
        translation: 'thank you',
        examples: ['Muchas gracias por tu ayuda'],
        notes: 'Expression of gratitude',
        wordType: 'noun',
        deckId: '1'
      },
      {
        id: '4',
        word: 'por favor',
        translation: 'please',
        examples: ['Por favor, ayúdame'],
        notes: 'Used when making requests',
        wordType: 'phrase',
        deckId: '1'
      },
      {
        id: '5',
        word: 'buenos días',
        translation: 'good morning',
        examples: ['Buenos días, ¿cómo estás?'],
        notes: 'Morning greeting',
        wordType: 'phrase',
        deckId: '1'
      }
    ];
    
    cy.loadDeckWithFlashcards('1', testFlashcards);
  });

  it('should start a study session when selecting a deck', () => {
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Verify we're on the study session page
    cy.contains('Flashcards').should('be.visible');
    cy.contains('1 / 5').should('be.visible');
    
    // Verify the first flashcard is displayed
    cy.contains('hola').should('be.visible');
  });

  it('should show answer when clicking the Show Answer button', () => {
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Initially, the translation should not be visible
    cy.contains('hello').should('not.be.visible');
    
    // Click the Show Answer button
    cy.contains('Show Answer').click();
    
    // Now the translation should be visible
    cy.contains('hello').should('be.visible');
    
    // Rating buttons should be visible
    cy.contains('Easy').should('be.visible');
    cy.contains('Good').should('be.visible');
    cy.contains('Difficult').should('be.visible');
  });

  it('should progress to the next card when rating the current card', () => {
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Verify first card
    cy.contains('hola').should('be.visible');
    
    // Show answer and rate the card
    cy.contains('Show Answer').click();
    cy.contains('Good').click();
    
    // Verify second card is shown
    cy.contains('adiós').should('be.visible');
    cy.contains('2 / 5').should('be.visible');
  });

  it('should allow rating cards with different difficulty levels', () => {
    // Mock the API response for updating card ratings
    cy.intercept('POST', '**/flashcards/*/rating', {
      statusCode: 200,
      body: {
        message: 'Rating updated successfully'
      }
    }).as('updateRatingRequest');
    
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Rate first card as Easy
    cy.contains('Show Answer').click();
    cy.contains('Easy').click();
    cy.wait('@updateRatingRequest');
    
    // Rate second card as Good
    cy.contains('Show Answer').click();
    cy.contains('Good').click();
    cy.wait('@updateRatingRequest');
    
    // Rate third card as Difficult
    cy.contains('Show Answer').click();
    cy.contains('Difficult').click();
    cy.wait('@updateRatingRequest');
  });

  it('should complete a session after reviewing all cards', () => {
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Complete the session using custom command with specific ratings
    cy.completeFlashcardSession(['Easy', 'Good', 'Difficult', 'Good', 'Easy']);
    
    // Verify session completion screen
    cy.contains('Session Complete').should('be.visible');
    cy.contains('Restart Session').should('be.visible');
  });

  it('should display session statistics after completion', () => {
    // Mock the API response for session statistics
    cy.intercept('GET', '**/decks/1/statistics', {
      statusCode: 200,
      body: {
        totalCards: 5,
        cardsReviewed: 5,
        easyRatings: 2,
        goodRatings: 2,
        difficultRatings: 1,
        averageTimePerCard: 3.5
      }
    }).as('getStatisticsRequest');
    
    // Navigate to the deck and complete the session
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    cy.completeFlashcardSession();
    
    // View statistics
    cy.contains('View Statistics').click();
    cy.wait('@getStatisticsRequest');
    
    // Verify statistics are displayed
    cy.verifySessionStats();
    cy.contains('5 / 5').should('be.visible');
    cy.contains('Easy: 2').should('be.visible');
    cy.contains('Good: 2').should('be.visible');
    cy.contains('Difficult: 1').should('be.visible');
  });

  it('should allow restarting a session after completion', () => {
    // Navigate to the deck and complete the session
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    cy.completeFlashcardSession();
    
    // Restart the session
    cy.contains('Restart Session').click();
    
    // Verify we're back at the first card
    cy.contains('1 / 5').should('be.visible');
    cy.contains('hola').should('be.visible');
  });

  it('should handle empty deck scenario', () => {
    // Mock an empty flashcards response
    cy.intercept('GET', '**/decks/2/flashcards', {
      statusCode: 200,
      body: []
    }).as('getEmptyFlashcardsRequest');
    
    // Navigate to the empty deck
    cy.contains('French Phrases').click();
    cy.wait('@getEmptyFlashcardsRequest');
    
    // Verify empty state message
    cy.contains('0 / 0').should('be.visible');
    cy.contains('No flashcards found').should('be.visible');
    cy.contains('Add your first flashcard').should('be.visible');
  });

  it('should allow pausing and resuming a study session', () => {
    // Navigate to the deck
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Review a couple of cards
    cy.contains('Show Answer').click();
    cy.contains('Good').click();
    cy.contains('Show Answer').click();
    cy.contains('Good').click();
    
    // Pause the session
    cy.contains('Pause Session').click();
    
    // Verify we're on the pause screen
    cy.contains('Session Paused').should('be.visible');
    cy.contains('Resume').should('be.visible');
    cy.contains('Quit').should('be.visible');
    
    // Resume the session
    cy.contains('Resume').click();
    
    // Verify we're back at the correct card
    cy.contains('3 / 5').should('be.visible');
  });
});
