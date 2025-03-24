# E2E Testing Guide for Language Learning Flashcard Generator

This guide documents the end-to-end (E2E) testing approach for the Language Learning Flashcard Generator application using Cypress.

## Table of Contents

1. [Testing Architecture](#testing-architecture)
2. [Mock Data Strategy](#mock-data-strategy)
3. [Custom Commands](#custom-commands)
4. [Test Organization](#test-organization)
5. [Accessibility Testing](#accessibility-testing)
6. [Error Handling Tests](#error-handling-tests)
7. [Running Tests](#running-tests)
8. [Troubleshooting](#troubleshooting)

## Testing Architecture

Our E2E testing architecture consists of several key components:

1. **Cypress Configuration**: Located in `cypress.config.js`, this sets up the testing environment with appropriate timeouts, retry settings, and other configurations to ensure reliable tests.

2. **API Mocking Utilities**: Located in `cypress/support/apiMocks.js`, these utilities provide consistent mocking of API responses across all tests.

3. **Custom Commands**: Located in `cypress/support/commands.js`, these provide reusable test actions for common workflows.

4. **Mock Data**: Located in `cypress/fixtures/mockData.js`, this centralizes all mock data used in tests.

5. **Test Specs**: Located in `cypress/e2e/`, these contain the actual test cases organized by feature.

## Mock Data Strategy

We use a comprehensive mocking strategy to ensure tests can run independently of the actual backend API:

1. **Centralized Mock Data**: All mock data is defined in `cypress/fixtures/mockData.js` to ensure consistency across tests.

2. **API Endpoint Mocking**: API responses are mocked using Cypress's `cy.intercept()` function, organized by endpoint type in `apiMocks.js`.

3. **Dynamic Mock Generation**: Helper functions can generate mock data dynamically when needed.

4. **Error Scenario Mocking**: We provide utilities to mock various error responses to test error handling.

Example of using mock data:

```javascript
// Import mock utilities
import { mockAuthEndpoints } from '../support/apiMocks';

// Setup mocks before test
beforeEach(() => {
  mockAuthEndpoints();
});

// Test with mocked API
it('should login successfully', () => {
  cy.get('input[type="email"]').type('test@example.com');
  cy.get('input[type="password"]').type('password123');
  cy.get('form').submit();
  cy.wait('@loginRequest');
});
```

## Custom Commands

We've created custom Cypress commands to simplify common test actions:

### Authentication Commands

- `cy.login(email, password)`: Logs in with the specified credentials
- `cy.register(name, email, password)`: Registers a new user
- `cy.requestPasswordReset(email)`: Requests a password reset
- `cy.verifySessionPersistence()`: Verifies that the user session is maintained

### Deck Management Commands

- `cy.loadDecks(customDecks)`: Loads a list of decks
- `cy.createDeck(name, description, language)`: Creates a new deck
- `cy.editDeck(deckId, newName, newDescription, newLanguage)`: Edits an existing deck
- `cy.deleteDeck(deckId)`: Deletes a deck

### Flashcard Operations Commands

- `cy.loadDeckWithFlashcards(deckId, customFlashcards)`: Loads a deck with flashcards
- `cy.createFlashcard(deckId, word, translation, examples, notes, wordType)`: Creates a new flashcard
- `cy.editFlashcard(deckId, flashcardId, newWord, newTranslation)`: Edits an existing flashcard
- `cy.deleteFlashcard(flashcardId)`: Deletes a flashcard

### Study Session Commands

- `cy.startStudySession(deckId)`: Starts a new study session
- `cy.completeFlashcardSession(ratings)`: Completes a flashcard review session
- `cy.verifySessionStats()`: Verifies the study session statistics

### Accessibility Testing Commands

- `cy.testKeyboardNavigation(maxElements)`: Tests keyboard navigation
- `cy.checkAccessibility(context, options)`: Checks for accessibility issues
- `cy.checkScreenReaderAccessibility()`: Tests screen reader accessibility

Example usage:

```javascript
describe('Deck Management', () => {
  beforeEach(() => {
    cy.login('test@example.com', 'password123');
  });

  it('should create a new deck', () => {
    cy.createDeck('Spanish Basics', 'Essential Spanish vocabulary', 'Spanish');
    cy.contains('Spanish Basics').should('be.visible');
  });
});
```

## Test Organization

Tests are organized by feature, with each test file focusing on a specific component or workflow:

1. **Authentication Tests**: Login, registration, password reset
2. **Deck Management Tests**: Creating, editing, and deleting decks
3. **Flashcard Tests**: Creating, editing, and deleting flashcards
4. **Study Session Tests**: Flashcard review and study sessions
5. **Error Handling Tests**: Testing API failures and edge cases

Each test file follows this structure:
- Setup (beforeEach)
- Happy path tests
- Validation tests
- Error handling tests
- Accessibility tests

## Accessibility Testing

We've implemented comprehensive accessibility testing:

1. **Automated Checks**: Using axe-core via Cypress-axe to check WCAG compliance
2. **Keyboard Navigation**: Testing that all functionality is accessible via keyboard
3. **Screen Reader Testing**: Verifying proper ARIA attributes and alt text

Example:

```javascript
it('should be accessible', () => {
  cy.checkAccessibility();
  cy.testKeyboardNavigation();
  cy.checkScreenReaderAccessibility();
});
```

## Error Handling Tests

We test various error scenarios to ensure the application handles them gracefully:

1. **API Failures**: Testing how the UI responds to different API error codes
2. **Validation Errors**: Testing form validation and error messages
3. **Authentication Errors**: Testing invalid credentials and session expiration
4. **Network Issues**: Testing offline behavior and reconnection

Example:

```javascript
it('should handle server error during login', () => {
  // Mock a server error
  cy.intercept('POST', '**/auth/login', {
    statusCode: 500,
    body: {
      error: 'Internal Server Error',
      message: 'Something went wrong'
    }
  }).as('loginServerErrorRequest');
  
  // Fill and submit the form
  cy.get('input[type="email"]').type('test@example.com');
  cy.get('input[type="password"]').type('password123');
  cy.get('form').submit();
  
  // Verify error message
  cy.contains('Something went wrong').should('be.visible');
});
```

## Running Tests

To run the E2E tests:

1. **Start the development server**:
   ```
   npm run dev
   ```

2. **Open Cypress**:
   ```
   npx cypress open
   ```

3. **Run all tests**:
   ```
   npx cypress run
   ```

4. **Run specific tests**:
   ```
   npx cypress run --spec "cypress/e2e/login.cy.js"
   ```

## Troubleshooting

Common issues and solutions:

1. **Tests failing due to timeouts**: Increase the `defaultCommandTimeout` in `cypress.config.js`

2. **Element not found errors**: Add `{ timeout: 10000 }` to your `cy.get()` or `cy.contains()` commands

3. **API mocking issues**: Verify the API endpoint patterns in `apiMocks.js` match the actual endpoints used by the application

4. **Missing dependencies**: Ensure all required packages are installed:
   ```
   npm install cypress cypress-axe @testing-library/cypress
   ```

5. **Linux environment issues**: Install required dependencies:
   ```
   apt-get update && apt-get install -y libgtk2.0-0 libgtk-3-0 libgbm-dev libnotify-dev libgconf-2-4 libnss3 libxss1 libasound2 libxtst6 xauth xvfb
   ```
