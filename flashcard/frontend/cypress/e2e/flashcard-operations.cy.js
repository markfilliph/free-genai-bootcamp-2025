/**
 * End-to-end tests for flashcard operations
 * 
 * This test file covers all flashcard-related operations:
 * - Creating flashcards with various content types
 * - Editing flashcards
 * - Deleting flashcards
 * - Bulk operations on flashcards
 */

describe('Flashcard Operations', () => {
  beforeEach(() => {
    // Log in, load decks, and navigate to a specific deck
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    cy.loadDeckWithFlashcards('1');
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
  });

  it('should display a list of flashcards for a deck', () => {
    // Verify flashcards are displayed
    cy.contains('hola').should('be.visible');
    cy.contains('adiós').should('be.visible');
    
    // Verify flashcard details
    cy.contains('hello').should('be.visible');
    cy.contains('goodbye').should('be.visible');
  });

  it('should allow creating a new flashcard', () => {
    // Use the custom createFlashcard command
    cy.createFlashcard('1', 'gracias', 'thank you', ['Muchas gracias por tu ayuda'], 'Expression of gratitude', 'noun');
    
    // Verify the new flashcard appears in the list
    cy.contains('gracias').should('be.visible');
    cy.contains('thank you').should('be.visible');
  });

  it('should validate flashcard creation form', () => {
    // Mock a validation error response
    cy.intercept('POST', '**/decks/1/flashcards', {
      statusCode: 400,
      body: {
        error: 'Word and translation are required'
      }
    }).as('invalidFlashcardRequest');
    
    cy.contains('Add Flashcard').click();
    // Submit without filling required fields
    cy.get('form').submit();
    cy.wait('@invalidFlashcardRequest');
    
    // Verify error message is displayed
    cy.contains('Word and translation are required').should('be.visible');
  });

  it('should allow creating flashcards with different word types', () => {
    // Create a verb flashcard
    cy.createFlashcard('1', 'hablar', 'to speak', ['Yo hablo español'], 'Common verb', 'verb');
    
    // Create an adjective flashcard
    cy.createFlashcard('1', 'grande', 'big', ['Una casa grande'], 'Size adjective', 'adjective');
    
    // Verify both flashcards appear in the list
    cy.contains('hablar').should('be.visible');
    cy.contains('grande').should('be.visible');
    
    // Verify word types are displayed correctly
    cy.contains('verb').should('be.visible');
    cy.contains('adjective').should('be.visible');
  });

  it('should allow editing an existing flashcard', () => {
    // Use the custom editFlashcard command
    cy.editFlashcard('1', '1', 'hola (updated)', 'hello (updated)');
    
    // Verify the updated flashcard appears in the list
    cy.contains('hola (updated)').should('be.visible');
    cy.contains('hello (updated)').should('be.visible');
  });

  it('should allow deleting a flashcard', () => {
    // Use the custom deleteFlashcard command
    cy.deleteFlashcard('2');
    
    // Verify the flashcard is removed from the list
    cy.contains('adiós').should('not.exist');
  });

  it('should confirm before deleting a flashcard', () => {
    // Click delete but cancel the confirmation
    cy.contains('Delete').click();
    cy.contains('Cancel').click();
    
    // Verify the flashcard still exists
    cy.contains('adiós').should('be.visible');
  });

  it('should allow bulk selection of flashcards', () => {
    // Enable bulk selection mode
    cy.contains('Select Multiple').click();
    
    // Select multiple flashcards
    cy.get('[data-testid="flashcard-checkbox-1"]').check();
    cy.get('[data-testid="flashcard-checkbox-2"]').check();
    
    // Verify selection count
    cy.contains('2 selected').should('be.visible');
  });

  it('should allow bulk deletion of flashcards', () => {
    // Mock the bulk delete API response
    cy.intercept('DELETE', '**/flashcards/bulk', {
      statusCode: 200,
      body: {
        message: 'Flashcards deleted successfully'
      }
    }).as('bulkDeleteRequest');
    
    // Enable bulk selection mode
    cy.contains('Select Multiple').click();
    
    // Select multiple flashcards
    cy.get('[data-testid="flashcard-checkbox-1"]').check();
    cy.get('[data-testid="flashcard-checkbox-2"]').check();
    
    // Delete selected flashcards
    cy.contains('Delete Selected').click();
    cy.get('[data-testid="confirm-bulk-delete"]').click();
    cy.wait('@bulkDeleteRequest');
    
    // Verify flashcards are removed
    cy.contains('hola').should('not.exist');
    cy.contains('adiós').should('not.exist');
    
    // Verify success message
    cy.contains('Flashcards deleted successfully').should('be.visible');
  });

  it('should allow filtering flashcards by word type', () => {
    // Mock the filtered flashcards response
    cy.intercept('GET', '**/decks/1/flashcards?wordType=noun', {
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
        }
      ]
    }).as('filteredFlashcardsRequest');
    
    // Select noun from the word type filter
    cy.get('select[name="wordTypeFilter"]').select('noun');
    cy.wait('@filteredFlashcardsRequest');
    
    // Verify only noun flashcards are displayed
    cy.contains('hola').should('be.visible');
    cy.contains('adiós').should('not.exist');
  });

  it('should allow searching flashcards by word or translation', () => {
    // Mock the searched flashcards response
    cy.intercept('GET', '**/decks/1/flashcards?search=hello', {
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
        }
      ]
    }).as('searchedFlashcardsRequest');
    
    // Enter search term
    cy.get('input[name="flashcardSearchTerm"]').type('hello');
    cy.get('form[role="search"]').submit();
    cy.wait('@searchedFlashcardsRequest');
    
    // Verify only matching flashcards are displayed
    cy.contains('hola').should('be.visible');
    cy.contains('adiós').should('not.exist');
  });
});
