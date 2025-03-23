// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************

// -- This is a parent command --
Cypress.Commands.add('login', (email, password) => {
  cy.intercept('POST', '**/auth/login', {
    statusCode: 200,
    body: {
      token: 'fake-jwt-token',
      user: {
        id: '1',
        email: email || 'test@example.com',
        name: 'Test User'
      }
    }
  }).as('loginRequest');
  
  cy.visit('/login');
  cy.get('input[type="email"]').type(email || 'test@example.com');
  cy.get('input[type="password"]').type(password || 'password123');
  cy.get('form').submit();
  cy.wait('@loginRequest');
});

// Command to load a deck with flashcards
Cypress.Commands.add('loadDeckWithFlashcards', (deckId = '1', flashcards = []) => {
  // Default flashcards if none provided
  const defaultFlashcards = [
    {
      id: '1',
      word: 'hola',
      translation: 'hello',
      examples: ['¡Hola! ¿Cómo estás?'],
      notes: 'Common greeting in Spanish',
      wordType: 'noun',
      deckId: deckId
    },
    {
      id: '2',
      word: 'adiós',
      translation: 'goodbye',
      examples: ['Adiós, hasta mañana'],
      notes: 'Used when parting ways',
      wordType: 'noun',
      deckId: deckId
    }
  ];
  
  // Use provided flashcards or default ones
  const cardsToUse = flashcards.length > 0 ? flashcards : defaultFlashcards;
  
  cy.intercept('GET', `**/decks/${deckId}/flashcards`, {
    statusCode: 200,
    body: cardsToUse
  }).as('getFlashcardsRequest');
});

// Command to complete a flashcard review session
Cypress.Commands.add('completeFlashcardSession', () => {
  // This will find all flashcards and review them until completion
  cy.get('body').then(($body) => {
    // Keep reviewing cards until we see the completion screen
    function reviewNextCard() {
      if ($body.text().includes('Session Complete')) {
        return;
      }
      
      // Click show answer if it exists
      if ($body.text().includes('Show Answer')) {
        cy.contains('Show Answer').click();
      }
      
      // Click a rating button if any exists
      if ($body.text().includes('Good')) {
        cy.contains('Good').click();
        
        // Check again after a short delay
        cy.wait(500).then(() => {
          cy.get('body').then(($newBody) => {
            $body = $newBody;
            if (!$body.text().includes('Session Complete')) {
              reviewNextCard();
            }
          });
        });
      }
    }
    
    reviewNextCard();
  });
});
