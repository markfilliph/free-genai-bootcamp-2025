// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************

// Import API mocking utilities
import { 
  mockAuthEndpoints, 
  mockDeckEndpoints, 
  mockFlashcardEndpoints,
  mockStudySessionEndpoints,
  setupAllApiMocks,
  mockApiError
} from './apiMocks';

// Import mock data
import { mockUser, mockToken } from '../fixtures/mockData';

// -- Authentication Commands --

// Login with email and password
Cypress.Commands.add('login', (email, password) => {
  // Setup auth mocks
  mockAuthEndpoints();
  
  // Visit login page
  cy.visit('/login', { failOnStatusCode: false });
  
  // Fill login form
  cy.get('input[type="email"]', { timeout: 10000 }).should('be.visible').type(email || 'test@example.com');
  cy.get('input[type="password"]').type(password || 'password123');
  cy.get('form').submit();
  
  // Wait for login request
  cy.wait('@loginRequest', { timeout: 10000 });
  
  // Store token in localStorage to maintain session
  cy.window().then((win) => {
    win.localStorage.setItem('token', mockToken);
    win.localStorage.setItem('user', JSON.stringify(mockUser));
  });
});

// Register a new user
Cypress.Commands.add('register', (name, email, password) => {
  // Setup auth mocks
  mockAuthEndpoints();
  
  // Visit register page
  cy.visit('/register', { failOnStatusCode: false });
  
  // Fill registration form
  cy.get('input[name="name"]', { timeout: 10000 }).should('be.visible').type(name || 'Test User');
  cy.get('input[type="email"]').type(email || 'newuser@example.com');
  cy.get('input[type="password"]').type(password || 'password123');
  cy.get('input[type="password"][name="confirmPassword"]').type(password || 'password123');
  cy.get('form').submit();
  
  // Wait for register request
  cy.wait('@registerRequest', { timeout: 10000 });
});

// Request password reset
Cypress.Commands.add('requestPasswordReset', (email) => {
  // Setup auth mocks
  mockAuthEndpoints();
  
  // Visit forgot password page
  cy.visit('/forgot-password', { failOnStatusCode: false });
  
  // Fill forgot password form
  cy.get('input[type="email"]', { timeout: 10000 }).should('be.visible').type(email || 'test@example.com');
  cy.get('form').submit();
  
  // Wait for forgot password request
  cy.wait('@forgotPasswordRequest', { timeout: 10000 });
});

// Verify session persistence
Cypress.Commands.add('verifySessionPersistence', () => {
  // Setup auth mocks
  mockAuthEndpoints();
  
  // Set token in localStorage
  cy.window().then((win) => {
    win.localStorage.setItem('token', mockToken);
    win.localStorage.setItem('user', JSON.stringify(mockUser));
  });
  
  // Visit the home page and check if session is maintained
  cy.visit('/', { failOnStatusCode: false });
  cy.wait('@validateTokenRequest', { timeout: 10000 });
  
  // Verify user is logged in
  cy.contains(mockUser.name).should('be.visible');
});

// -- Deck Management Commands --

// Load decks list
Cypress.Commands.add('loadDecks', (customDecks = []) => {
  // Setup deck mocks with custom decks if provided
  if (customDecks.length > 0) {
    cy.intercept('GET', '**/decks', {
      statusCode: 200,
      body: customDecks
    }).as('getDecksRequest');
  } else {
    // Use default mocks
    mockDeckEndpoints();
  }
  
  // Ensure user is logged in
  cy.window().then((win) => {
    if (!win.localStorage.getItem('token')) {
      win.localStorage.setItem('token', mockToken);
      win.localStorage.setItem('user', JSON.stringify(mockUser));
    }
  });
});

// Create a new deck
Cypress.Commands.add('createDeck', (name, description, language) => {
  // Setup deck mocks if not already set up
  mockDeckEndpoints();
  
  // Click the create deck button
  cy.contains('Create Deck', { timeout: 10000 }).should('be.visible').click();
  
  // Fill out the form
  cy.get('input[name="name"]', { timeout: 10000 }).should('be.visible').type(name || 'New Test Deck');
  cy.get('textarea[name="description"]').type(description || 'Test deck description');
  cy.get('select[name="language"]').select(language || 'English');
  cy.get('form').submit();
  
  // Wait for the create request
  cy.wait('@createDeckRequest', { timeout: 10000 });
});

// Edit an existing deck
Cypress.Commands.add('editDeck', (deckId, newName, newDescription, newLanguage) => {
  // Setup deck mocks if not already set up
  mockDeckEndpoints();
  
  // Click the edit button
  cy.contains('Edit', { timeout: 10000 }).should('be.visible').click();
  
  // Update the form
  cy.get('input[name="name"]', { timeout: 10000 }).should('be.visible').clear().type(newName || 'Updated Deck Name');
  cy.get('textarea[name="description"]').clear().type(newDescription || 'Updated deck description');
  cy.get('select[name="language"]').select(newLanguage || 'Spanish');
  cy.get('form').submit();
  
  // Wait for the update request
  cy.wait('@updateDeckRequest', { timeout: 10000 });
});

// Delete a deck
Cypress.Commands.add('deleteDeck', (deckId) => {
  // Setup deck mocks if not already set up
  mockDeckEndpoints();
  
  // Click the delete button
  cy.contains('Delete', { timeout: 10000 }).should('be.visible').click();
  
  // Confirm deletion
  cy.get('[data-testid="confirm-delete"]', { timeout: 10000 }).should('be.visible').click();
  
  // Wait for the delete request
  cy.wait('@deleteDeckRequest', { timeout: 10000 });
});

// -- Flashcard Operations Commands --

// Load a deck with flashcards
Cypress.Commands.add('loadDeckWithFlashcards', (deckId = '1', customFlashcards = []) => {
  // Setup flashcard mocks with custom flashcards if provided
  if (customFlashcards.length > 0) {
    cy.intercept('GET', `**/decks/${deckId}/flashcards`, {
      statusCode: 200,
      body: customFlashcards
    }).as('getFlashcardsRequest');
  } else {
    // Use default mocks
    mockFlashcardEndpoints();
  }
  
  // Ensure user is logged in
  cy.window().then((win) => {
    if (!win.localStorage.getItem('token')) {
      win.localStorage.setItem('token', mockToken);
      win.localStorage.setItem('user', JSON.stringify(mockUser));
    }
  });
});

// Create a new flashcard
Cypress.Commands.add('createFlashcard', (deckId, word, translation, examples, notes, wordType) => {
  // Setup flashcard mocks if not already set up
  mockFlashcardEndpoints();
  
  // Click the add flashcard button
  cy.contains('Add Flashcard', { timeout: 10000 }).should('be.visible').click();
  
  // Fill out the form
  cy.get('input[name="word"]', { timeout: 10000 }).should('be.visible').type(word || 'nueva palabra');
  cy.get('input[name="translation"]').type(translation || 'new word');
  cy.get('input[name="examples"]').type(examples?.[0] || 'Esta es una nueva palabra');
  cy.get('textarea[name="notes"]').type(notes || 'A new Spanish word');
  cy.get('select[name="wordType"]').select(wordType || 'noun');
  cy.get('form').submit();
  
  // Wait for the create request
  cy.wait('@createFlashcardRequest', { timeout: 10000 });
});

// Edit an existing flashcard
Cypress.Commands.add('editFlashcard', (deckId, flashcardId, newWord, newTranslation) => {
  // Setup flashcard mocks if not already set up
  mockFlashcardEndpoints();
  
  // Click the edit button
  cy.contains('Edit', { timeout: 10000 }).should('be.visible').click();
  
  // Update the form
  cy.get('input[name="word"]', { timeout: 10000 }).should('be.visible').clear().type(newWord || 'updated word');
  cy.get('input[name="translation"]').clear().type(newTranslation || 'updated translation');
  cy.get('form').submit();
  
  // Wait for the update request
  cy.wait('@updateFlashcardRequest', { timeout: 10000 });
});

// Delete a flashcard
Cypress.Commands.add('deleteFlashcard', (flashcardId) => {
  cy.intercept('DELETE', `**/flashcards/${flashcardId}`, {
    statusCode: 200,
    body: {
      message: 'Flashcard deleted successfully'
    }
  }).as('deleteFlashcardRequest');
  
  cy.contains('Delete').click();
  cy.get('[data-testid="confirm-delete"]').click();
  cy.wait('@deleteFlashcardRequest');
});

// -- Study Session Commands --

// Complete a flashcard review session
Cypress.Commands.add('completeFlashcardSession', (ratings = []) => {
  // This will find all flashcards and review them until completion
  cy.get('body').then(($body) => {
    // Keep reviewing cards until we see the completion screen
    function reviewNextCard(index = 0) {
      if ($body.text().includes('Session Complete')) {
        return;
      }
      
      // Click show answer if it exists
      if ($body.text().includes('Show Answer')) {
        cy.contains('Show Answer').click();
      }
      
      // Use the provided rating if available, otherwise use 'Good'
      const rating = ratings[index] || 'Good';
      
      // Click the specified rating button
      if ($body.text().includes(rating)) {
        cy.contains(rating).click();
        
        // Check again after a short delay
        cy.wait(500).then(() => {
          cy.get('body').then(($newBody) => {
            $body = $newBody;
            if (!$body.text().includes('Session Complete')) {
              reviewNextCard(index + 1);
            }
          });
        });
      }
    }
    
    reviewNextCard();
  });
});

// Verify study session statistics
Cypress.Commands.add('verifySessionStats', () => {
  cy.contains('Session Statistics').should('be.visible');
  cy.contains('Cards Reviewed').should('be.visible');
  cy.contains('Performance').should('be.visible');
});

// -- Accessibility Testing Commands --

// Test keyboard navigation
Cypress.Commands.add('testKeyboardNavigation', () => {
  // Press Tab to navigate through elements
  cy.get('body').tab();
  cy.focused().should('be.visible');
  
  // Continue tabbing through all focusable elements
  for (let i = 0; i < 10; i++) {
    cy.focused().then($el => {
      const tagName = $el.prop('tagName').toLowerCase();
      const type = $el.attr('type');
      
      // If it's a button or link, press Enter to activate it
      if (tagName === 'button' || tagName === 'a' || (tagName === 'input' && type === 'submit')) {
        cy.focused().type('{enter}');
        cy.wait(300); // Wait for any actions to complete
      }
      
      // Move to next element
      cy.get('body').tab();
    });
  }
});

// Check for color contrast issues
Cypress.Commands.add('checkAccessibility', () => {
  cy.injectAxe();
  cy.checkA11y();
});
