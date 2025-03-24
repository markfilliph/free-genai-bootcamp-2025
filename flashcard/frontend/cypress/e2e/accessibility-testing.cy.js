/**
 * End-to-end tests for accessibility
 * 
 * This test file covers accessibility testing:
 * - Keyboard navigation throughout the application
 * - Screen reader compatibility
 * - Focus management
 * - Color contrast and visual accessibility
 */

describe('Accessibility Testing', () => {
  beforeEach(() => {
    // Load axe before each test
    cy.injectAxe();
  });

  it('should have no accessibility violations on the login page', () => {
    cy.visit('/login');
    cy.wait(500); // Wait for page to fully load
    cy.checkA11y();
  });

  it('should have no accessibility violations on the home page', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    cy.wait(500); // Wait for page to fully load
    cy.checkA11y();
  });

  it('should have no accessibility violations on the flashcard review page', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    cy.loadDeckWithFlashcards('1');
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    cy.wait(500); // Wait for page to fully load
    cy.checkA11y();
  });

  it('should be navigable using only the keyboard on the login page', () => {
    cy.visit('/login');
    
    // Test keyboard navigation
    cy.focused().should('have.attr', 'name', 'email');
    cy.tab().should('have.attr', 'name', 'password');
    cy.tab().should('have.attr', 'type', 'submit');
    
    // Test form submission with keyboard
    cy.get('input[name="email"]').type('test@example.com');
    cy.get('input[name="password"]').type('password123{enter}');
    cy.wait('@loginRequest');
    cy.url().should('not.include', '/login');
  });

  it('should be navigable using only the keyboard on the deck list', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    
    // Test keyboard navigation through deck list
    cy.contains('Create Deck').focus();
    cy.tab().should('contain', 'Spanish Basics');
    cy.focused().type('{enter}');
    cy.wait('@getFlashcardsRequest');
    cy.url().should('include', '/decks/1');
  });

  it('should have proper focus management when opening and closing modals', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    
    // Open create deck modal
    cy.contains('Create Deck').click();
    
    // Focus should be trapped in the modal
    cy.focused().should('exist');
    cy.tab().should('exist');
    cy.tab().should('exist');
    cy.tab().should('exist');
    
    // Close modal with keyboard
    cy.get('button[aria-label="Close"]').focus();
    cy.focused().type('{enter}');
    
    // Focus should return to the element that opened the modal
    cy.focused().should('contain', 'Create Deck');
  });

  it('should have appropriate ARIA attributes on interactive elements', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    
    // Check ARIA attributes on navigation
    cy.get('nav').should('have.attr', 'role', 'navigation');
    
    // Check ARIA attributes on search
    cy.get('form[role="search"]').should('exist');
    
    // Check ARIA attributes on buttons
    cy.contains('Create Deck').should('have.attr', 'aria-haspopup', 'dialog');
  });

  it('should announce dynamic content changes to screen readers', () => {
    cy.login('test@example.com', 'password123');
    cy.loadDecks();
    cy.wait('@getDecksRequest');
    cy.loadDeckWithFlashcards('1');
    cy.contains('Spanish Basics').click();
    cy.wait('@getFlashcardsRequest');
    
    // Check for live regions that announce changes
    cy.get('[aria-live]').should('exist');
    
    // Show answer and check that the content is announced
    cy.contains('Show Answer').click();
    cy.get('[aria-live="polite"]').should('contain', 'hello');
  });
});
