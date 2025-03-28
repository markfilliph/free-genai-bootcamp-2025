# Frontend Testing Guide

This guide explains the testing approach used in the Language Learning Flashcard Generator frontend application.

## Testing Structure

The frontend testing is organized into a modular structure to make tests more maintainable and easier to understand.

### Directory Structure

```
frontend/
├── src/
│   ├── __tests__/           # All test files
│   │   ├── setup/           # Consolidated test setup files
│   │   │   ├── index.js     # Main entry point for test setup
│   │   │   ├── browser-mocks.js     # Browser API mocks
│   │   │   ├── component-mocks.js   # Svelte component mocks
│   │   │   ├── environment-mocks.js # Vite environment mocks
│   │   │   ├── fetch-mock.js        # Fetch API and API module mocks
│   │   │   ├── storage-mocks.js     # localStorage mocks
│   │   │   └── test-utils.js        # Testing utility functions
│   │   ├── components/      # Component tests
│   │   ├── routes/          # Route component tests
│   │   ├── lib/             # Utility function tests
│   │   ├── integration/     # Integration tests
│   │   └── mocks/           # Mock implementations (legacy)
```

## Testing Approach

Our testing approach for Svelte components follows these principles:

1. **Direct Component Mocking**: We create simplified mock versions of components that return expected HTML
2. **Custom Render Functions**: We use `$$render` to control component output
3. **Global Mock Helpers**: We provide consistent mocks for common components
4. **Environment Mocking**: We properly handle Vite environment variables and browser APIs

### Why This Approach?

Testing Svelte components presents unique challenges:

- Svelte's compilation process makes traditional component mocking difficult
- Components often have complex dependencies on other components
- Svelte-specific features like reactive statements and lifecycle methods need special handling
- Testing libraries may not fully support all Svelte features

Our custom approach addresses these challenges by providing a simplified way to test components in isolation without complex dependencies.

## Test Setup Modules

The test setup has been consolidated into modular files, each with a specific purpose:

### 1. `index.js`

The main entry point for Jest setup. This file imports all other setup modules and exports utility functions for use in tests.

### 2. `browser-mocks.js`

Mocks for browser APIs that are used in the application but not available in the Jest environment:
- MutationObserver
- DOM manipulation methods
- Custom Jest DOM matchers

### 3. `storage-mocks.js`

Mocks for browser storage APIs:
- localStorage
- sessionStorage

### 4. `fetch-mock.js`

Mocks for network-related functionality:
- Fetch API
- Application API module and its methods

### 5. `component-mocks.js`

Mocks for Svelte components and routing:
- svelte-routing (Link, Router, Route)
- Common application components (Navbar, DeckList, etc.)

### 6. `environment-mocks.js`

Mocks for environment-specific features:
- Vite's import.meta.env
- Environment variables

### 7. `test-utils.js`

Utility functions for testing:
- mockComponent: Creates mock Svelte components
- mockHtml: Creates HTML content for component mocks
- mockEvent: Creates mock event objects
- mockFormSubmit: Creates mock form submission events

## Writing Tests

When writing tests, you can use the provided utilities to create mocks and test components. Here are examples of our custom testing approaches:

### 1. Direct Component Mocking

```javascript
// Example: Mocking the DeckList component
import { render } from '@testing-library/svelte';
import { mockComponent } from '../setup/test-utils';
import DeckManagement from '../../routes/DeckManagement.svelte';

// Mock the DeckList component
jest.mock('../../components/DeckList.svelte', () => ({
  default: mockComponent('DeckList', `
    <div data-testid="mock-deck-list">
      <ul>
        <li>Mock Deck 1</li>
        <li>Mock Deck 2</li>
      </ul>
    </div>
  `)
}));

describe('DeckManagement with mocked DeckList', () => {
  test('renders the DeckList component', () => {
    const { getByTestId } = render(DeckManagement);
    expect(getByTestId('mock-deck-list')).toBeInTheDocument();
  });
});
```

### 2. Custom Render Functions

```javascript
// Example: Using $$render to control component output
import { mockHtml } from '../setup/test-utils';

// Mock a Svelte component with custom render function
jest.mock('../../components/FlashcardReview.svelte', () => ({
  default: {
    $$render: (result, props) => {
      const { card, onAnswer } = props;
      return mockHtml(`
        <div class="flashcard-review" data-testid="flashcard-review">
          <div class="card-front">${card.frontText}</div>
          <div class="card-back">${card.backText}</div>
          <button class="answer-button" data-answer="easy">Easy</button>
        </div>
      `);
    }
  }
}));
```

### 3. Testing Component Events

```javascript
// Example: Testing component events
import { render, fireEvent } from '@testing-library/svelte';
import { mockEvent } from '../setup/test-utils';
import FlashcardForm from '../../components/FlashcardForm.svelte';

describe('FlashcardForm', () => {
  test('submits the form with correct data', async () => {
    // Create a mock submit handler
    const handleSubmit = jest.fn();
    
    // Render the component with the mock handler
    const { getByLabelText, getByText } = render(FlashcardForm, {
      props: { onSubmit: handleSubmit }
    });
    
    // Fill in the form
    await fireEvent.input(getByLabelText('Front Text'), { target: { value: 'Test Word' } });
    await fireEvent.input(getByLabelText('Back Text'), { target: { value: 'Test Translation' } });
    
    // Submit the form
    await fireEvent.click(getByText('Create Flashcard'));
    
    // Check that the handler was called with the correct data
    expect(handleSubmit).toHaveBeenCalledWith(expect.objectContaining({
      frontText: 'Test Word',
      backText: 'Test Translation'
    }));
  });
});
```

### 4. Testing with API Mocks

```javascript
// Example: Testing components that use API calls
import { render, waitFor } from '@testing-library/svelte';
import * as api from '../../lib/api';
import CreateFlashcards from '../../routes/CreateFlashcards.svelte';

// Mock the API module
jest.mock('../../lib/api');

describe('CreateFlashcards', () => {
  beforeEach(() => {
    // Setup API mock responses
    api.getDecks.mockResolvedValue([
      { id: '1', name: 'Test Deck', cardCount: 0 }
    ]);
    api.generateContent.mockResolvedValue({
      example_sentences: ['Example sentence'],
      conjugations: 'Verb conjugations',
      cultural_note: 'Cultural note'
    });
  });
  
  test('loads decks on mount', async () => {
    const { getByText } = render(CreateFlashcards);
    
    await waitFor(() => {
      expect(api.getDecks).toHaveBeenCalled();
      expect(getByText('Test Deck')).toBeInTheDocument();
    });
  });
});
```

## Running Tests

To run the tests, use the following commands:

```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

## Coverage Thresholds

The project has the following coverage thresholds:

- Global: 50% statements, 40% branches, 50% functions, 50% lines
- Critical components have lower thresholds but are prioritized for improvement:
  - StudySession.svelte: 30% statements, 20% branches, 30% functions, 30% lines
  - FlashcardReview.svelte: 30% statements, 20% branches, 30% functions, 30% lines
  - Login.svelte: 30% statements, 20% branches, 30% functions, 30% lines

## Legacy Files

The following files have been consolidated and are kept for reference but are no longer used:
- `setup.js`
- `setup.new.js`
- `setup-fixed.config.js`
- `setup-minimal.config.js`
- `setup-simple.js`
