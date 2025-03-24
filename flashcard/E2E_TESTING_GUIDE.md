# E2E Testing Guide for Language Learning Flashcard Generator

This guide documents the end-to-end (E2E) testing approach for the Language Learning Flashcard Generator application. It covers the testing framework, test utilities, test scenarios, and best practices for writing and maintaining E2E tests.

## Table of Contents

1. [Testing Framework](#testing-framework)
2. [Test Structure](#test-structure)
3. [Custom Commands](#custom-commands)
4. [Test Scenarios](#test-scenarios)
5. [Accessibility Testing](#accessibility-testing)
6. [Continuous Integration](#continuous-integration)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

## Testing Framework

We use [Cypress](https://www.cypress.io/) as our E2E testing framework. Cypress provides a complete end-to-end testing experience with:

- Real-time reloading
- Time travel debugging
- Automatic waiting
- Network traffic control
- Consistent results

### Setup

The Cypress configuration is defined in `frontend/cypress.config.js`. Key settings include:

```javascript
{
  e2e: {
    baseUrl: 'http://localhost:5173',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
  component: {
    devServer: {
      framework: 'svelte',
      bundler: 'vite',
    },
  },
}
```

### Running Tests

To run the E2E tests:

1. Start the development server:
   ```
   cd frontend
   npm run dev
   ```

2. Run Cypress tests:
   ```
   cd frontend
   npm run cypress:open  # For interactive mode
   npm run cypress:run   # For headless mode
   ```

## Test Structure

Our E2E tests are organized by feature area:

- `authentication-flows.cy.js`: Tests for login, registration, password reset, and session management
- `deck-management.cy.js`: Tests for creating, editing, deleting, and searching decks
- `flashcard-operations.cy.js`: Tests for creating, editing, deleting, and managing flashcards
- `flashcard-flow.cy.js`: Tests for the complete user journey
- `study-session-workflows.cy.js`: Tests for study sessions, card ratings, and statistics
- `accessibility-testing.cy.js`: Tests for keyboard navigation, screen reader compatibility, and ARIA compliance

Each test file follows this structure:

```javascript
describe('Feature Area', () => {
  beforeEach(() => {
    // Setup code that runs before each test
  });

  it('should perform a specific action', () => {
    // Test steps
    // Assertions
  });
});
```

## Custom Commands

We've created custom Cypress commands to simplify common operations. These are defined in `frontend/cypress/support/commands.js`.

### Authentication Commands

```javascript
// Login with email and password
cy.login(email, password);

// Register a new user
cy.register(name, email, password);

// Request password reset
cy.requestPasswordReset(email);

// Verify session persistence
cy.verifySessionPersistence();
```

### Deck Management Commands

```javascript
// Load decks list
cy.loadDecks(decks);

// Create a new deck
cy.createDeck(name, description, language);

// Edit an existing deck
cy.editDeck(deckId, newName, newDescription, newLanguage);

// Delete a deck
cy.deleteDeck(deckId);
```

### Flashcard Operations Commands

```javascript
// Load a deck with flashcards
cy.loadDeckWithFlashcards(deckId, flashcards);

// Create a new flashcard
cy.createFlashcard(deckId, word, translation, examples, notes, wordType);

// Edit an existing flashcard
cy.editFlashcard(deckId, flashcardId, newWord, newTranslation);

// Delete a flashcard
cy.deleteFlashcard(flashcardId);
```

### Study Session Commands

```javascript
// Complete a flashcard review session
cy.completeFlashcardSession(ratings);

// Verify study session statistics
cy.verifySessionStats();
```

### Accessibility Testing Commands

```javascript
// Test keyboard navigation
cy.testKeyboardNavigation();

// Check for accessibility issues
cy.checkAccessibility();
```

## Test Scenarios

### Authentication Flows

- Login with valid credentials
- Handle login failures
- Register a new account
- Handle registration validation errors
- Request password reset
- Maintain user session across page reloads
- Redirect to login page for protected routes

### Deck Management

- Display list of decks
- Create a new deck
- Validate deck creation form
- Edit an existing deck
- Delete a deck
- Confirm before deleting
- Filter decks by language
- Search decks by name
- Handle empty search results

### Flashcard Operations

- Display list of flashcards
- Create a new flashcard
- Validate flashcard creation form
- Create flashcards with different word types
- Edit an existing flashcard
- Delete a flashcard
- Confirm before deleting
- Bulk select flashcards
- Bulk delete flashcards
- Filter flashcards by word type
- Search flashcards by word or translation

### Study Session Workflows

- Start a study session
- Show answer when clicking button
- Progress to next card after rating
- Rate cards with different difficulty levels
- Complete a session
- Display session statistics
- Restart a session
- Handle empty deck scenario
- Pause and resume a session

### Accessibility Testing

- Check for accessibility violations on key pages
- Test keyboard navigation
- Verify focus management in modals
- Check ARIA attributes
- Verify screen reader announcements

## Accessibility Testing

We use [cypress-axe](https://github.com/component-driven/cypress-axe) to test for accessibility issues. This integrates the [axe-core](https://github.com/dequelabs/axe-core) accessibility testing engine with Cypress.

### Setup

The accessibility testing setup is defined in `frontend/cypress/support/e2e.js`:

```javascript
// Import cypress-axe for accessibility testing
import 'cypress-axe'

// Add tab command for keyboard navigation testing
Cypress.Commands.add('tab', { prevSubject: 'optional' }, (subject) => {
  const tabKey = { key: 'Tab', code: 'Tab', which: 9 }
  if (subject) {
    cy.wrap(subject).trigger('keydown', tabKey)
  } else {
    cy.focused().trigger('keydown', tabKey)
  }
  return cy.focused()
})
```

### Usage

To test for accessibility violations:

```javascript
it('should have no accessibility violations', () => {
  cy.visit('/login');
  cy.injectAxe();
  cy.checkA11y();
});
```

To test keyboard navigation:

```javascript
it('should be navigable using only the keyboard', () => {
  cy.visit('/login');
  cy.focused().should('have.attr', 'name', 'email');
  cy.tab().should('have.attr', 'name', 'password');
  cy.tab().should('have.attr', 'type', 'submit');
});
```

## Continuous Integration

We use GitHub Actions to run our E2E tests automatically on pull requests and pushes to the main branch. The workflow is defined in `.github/workflows/cypress-tests.yml`.

Key features of our CI setup:

- Runs on Ubuntu latest
- Uses Node.js 18
- Installs dependencies with npm ci
- Starts the frontend server
- Runs Cypress tests in Chrome
- Uploads screenshots and videos as artifacts on failure

## Best Practices

### Writing Effective E2E Tests

1. **Focus on user flows**: Test complete user journeys rather than isolated UI elements.
2. **Use custom commands**: Create reusable commands for common operations.
3. **Mock API responses**: Use `cy.intercept()` to mock API responses for consistent testing.
4. **Be specific with selectors**: Use data-testid attributes for stable test selectors.
5. **Test error states**: Verify that the application handles errors gracefully.
6. **Keep tests independent**: Each test should be able to run independently.
7. **Minimize test duplication**: Use beforeEach hooks and custom commands to reduce duplication.

### Maintaining E2E Tests

1. **Update tests when the UI changes**: Keep tests in sync with UI changes.
2. **Review test failures promptly**: Address test failures as soon as they occur.
3. **Refactor tests regularly**: Keep test code clean and maintainable.
4. **Document test assumptions**: Comment on any assumptions or special conditions.
5. **Use descriptive test names**: Make it clear what each test is verifying.

## Troubleshooting

### Common Issues

1. **Tests timing out**: Use `cy.wait()` to wait for specific events rather than arbitrary timeouts.
2. **Selector not found**: Use more specific selectors or add data-testid attributes.
3. **API mocking issues**: Verify that the intercept patterns match the actual API calls.
4. **Flaky tests**: Identify and fix sources of non-determinism in tests.

### Debugging Tips

1. **Use Cypress time travel**: Click on commands in the Cypress runner to see the state at that point.
2. **Add .debug()**: Insert `cy.debug()` to pause execution and inspect the state.
3. **Check screenshots and videos**: Review artifacts from failed CI runs.
4. **Use console logs**: Add `cy.log()` statements to track test progress.
5. **Isolate the issue**: Run specific tests with `it.only()` to focus on problematic tests.

---

This guide is a living document. Please update it as our testing approach evolves.
