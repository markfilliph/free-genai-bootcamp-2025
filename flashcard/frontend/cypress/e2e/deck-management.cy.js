/**
 * End-to-end tests for deck management
 * 
 * This test file covers all deck management operations:
 * - Creating new decks
 * - Editing existing decks
 * - Deleting decks
 * - Filtering and searching decks
 */

describe('Deck Management', () => {
  beforeEach(() => {
    // Log in and load decks before each test
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
  });

  it('should display a list of decks on the home page', () => {
    // Verify decks are displayed
    cy.contains('Spanish Basics').should('be.visible');
    cy.contains('French Phrases').should('be.visible');
    
    // Verify deck details are displayed
    cy.contains('Essential Spanish vocabulary').should('be.visible');
    cy.contains('Common French expressions').should('be.visible');
  });

  it('should allow creating a new deck', () => {
    // Use the custom createDeck command
    cy.createDeck('German Vocabulary', 'Basic German words and phrases', 'German');
    
    // Verify the new deck appears in the list
    cy.contains('German Vocabulary').should('be.visible');
    cy.contains('Basic German words and phrases').should('be.visible');
  });

  it('should validate deck creation form', () => {
    // Mock a validation error response
    cy.intercept('POST', '**/decks', {
      statusCode: 400,
      body: {
        error: 'Name is required'
      }
    }).as('invalidDeckRequest');
    
    cy.contains('Create Deck').click();
    // Submit without filling required fields
    cy.get('form').submit();
    cy.wait('@invalidDeckRequest');
    
    // Verify error message is displayed
    cy.contains('Name is required').should('be.visible');
  });

  it('should allow editing an existing deck', () => {
    // Use the custom editDeck command
    cy.editDeck('1', 'Updated Spanish Basics', 'Revised Spanish vocabulary', 'Spanish');
    
    // Verify the updated deck appears in the list
    cy.contains('Updated Spanish Basics').should('be.visible');
    cy.contains('Revised Spanish vocabulary').should('be.visible');
  });

  it('should allow deleting a deck', () => {
    // Use the custom deleteDeck command
    cy.deleteDeck('2');
    
    // Verify the deck is removed from the list
    cy.contains('French Phrases').should('not.exist');
  });

  it('should confirm before deleting a deck', () => {
    // Click delete but cancel the confirmation
    cy.contains('Delete').click();
    cy.contains('Cancel').click();
    
    // Verify the deck still exists
    cy.contains('French Phrases').should('be.visible');
  });

  it('should allow filtering decks by language', () => {
    // Mock the API response for filtered decks
    cy.intercept('GET', '**/decks?language=Spanish', {
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
    }).as('filteredDecksRequest');
    
    // Select Spanish from the language filter
    cy.get('select[name="languageFilter"]').select('Spanish');
    cy.wait('@filteredDecksRequest');
    
    // Verify only Spanish decks are displayed
    cy.contains('Spanish Basics').should('be.visible');
    cy.contains('French Phrases').should('not.exist');
  });

  it('should allow searching decks by name', () => {
    // Mock the API response for searched decks
    cy.intercept('GET', '**/decks?search=French', {
      statusCode: 200,
      body: [
        {
          id: '2',
          name: 'French Phrases',
          description: 'Common French expressions',
          language: 'French',
          createdAt: '2023-01-02T00:00:00.000Z'
        }
      ]
    }).as('searchedDecksRequest');
    
    // Enter search term
    cy.get('input[name="searchTerm"]').type('French');
    cy.get('form[role="search"]').submit();
    cy.wait('@searchedDecksRequest');
    
    // Verify only matching decks are displayed
    cy.contains('French Phrases').should('be.visible');
    cy.contains('Spanish Basics').should('not.exist');
  });

  it('should handle empty search results gracefully', () => {
    // Mock an empty search response
    cy.intercept('GET', '**/decks?search=NonExistent', {
      statusCode: 200,
      body: []
    }).as('emptySearchRequest');
    
    // Enter search term
    cy.get('input[name="searchTerm"]').type('NonExistent');
    cy.get('form[role="search"]').submit();
    cy.wait('@emptySearchRequest');
    
    // Verify empty state message
    cy.contains('No decks found').should('be.visible');
  });
});
