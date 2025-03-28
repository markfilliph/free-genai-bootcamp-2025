# Custom Svelte Testing Approach

This document explains our custom approach to testing Svelte components in the Language Learning Flashcard Generator application.

## The Challenge

Testing Svelte components presents several unique challenges:

1. **Compilation Process**: Svelte components are compiled to JavaScript at build time, making traditional component mocking difficult.
2. **Component Dependencies**: Components often have complex dependencies on other components, making isolation challenging.
3. **Svelte-Specific Features**: Features like reactive statements, lifecycle methods, and bindings need special handling in tests.
4. **Testing Library Limitations**: Standard testing libraries may not fully support all Svelte features.

## Our Solution

We've developed a custom testing approach with four key strategies:

### 1. Direct Component Mocking

Instead of trying to mock the compiled output of Svelte components, we create simplified mock versions that return expected HTML:

```javascript
// Mock a component to return simple HTML
jest.mock('../../components/DeckList.svelte', () => ({
  default: {
    $$render: () => `
      <div class="deck-list">
        <ul>
          <li>Mock Deck 1</li>
          <li>Mock Deck 2</li>
        </ul>
      </div>
    `
  }
}));
```

This approach allows us to:
- Control exactly what HTML a component renders
- Avoid complex component dependencies
- Focus on testing the component under test, not its dependencies

### 2. Custom Render Functions

We use Svelte's `$$render` function to control component output based on props:

```javascript
jest.mock('../../components/FlashcardReview.svelte', () => ({
  default: {
    $$render: (result, props) => {
      const { card, onAnswer } = props || {};
      return `
        <div class="flashcard-review">
          <div class="card-front">${card ? card.frontText : ''}</div>
          <div class="card-back">${card ? card.backText : ''}</div>
          <button class="answer-button" data-answer="easy">Easy</button>
        </div>
      `;
    }
  }
}));
```

This approach allows us to:
- Create dynamic mock components that respond to props
- Test component interactions without complex setup
- Simulate different component states easily

### 3. Global Mock Helpers

We provide consistent mock implementations for common components and utilities:

```javascript
// In setup/component-mocks.js
export function mockComponent(name, html = '') {
  return {
    $$render: () => html || `<div data-testid="mock-${name}">${name} Component</div>`,
    render: (props) => html || `<div data-testid="mock-${name}">${name} Component</div>`
  };
}

// In a test file
import { mockComponent } from '../setup/component-mocks';

jest.mock('../../components/Navbar.svelte', () => ({
  default: mockComponent('Navbar')
}));
```

This approach allows us to:
- Maintain consistency across tests
- Reduce duplication of mock implementations
- Easily identify mock components in test output

### 4. Environment Mocking

We properly handle Vite environment variables and browser APIs:

```javascript
// In setup/environment-mocks.js
global.import = global.import || {};
global.import.meta = global.import.meta || {};
global.import.meta.env = {
  VITE_API_URL: 'http://localhost:8000',
  MODE: 'test',
  DEV: true
};
```

This approach allows us to:
- Test code that relies on environment variables
- Simulate different environment configurations
- Avoid errors related to missing browser APIs

## Real-World Examples

### Example 1: Testing StudySession Component

The StudySession component was particularly challenging to test because it:
- Renders a FlashcardReview component
- Manages state for the current card and study session
- Handles user interactions for answering cards

Our solution:

```javascript
// Mock the FlashcardReview component
jest.mock('../../components/FlashcardReview.svelte', () => ({
  default: {
    $$render: (result, props) => {
      const { card } = props || {};
      return `
        <div class="flashcard-review" data-testid="flashcard-review">
          ${card ? `<div>${card.frontText}</div>` : ''}
          <button data-testid="answer-button" data-answer="easy">Easy</button>
        </div>
      `;
    }
  }
}));

describe('StudySession', () => {
  const mockCards = [
    { id: '1', frontText: 'Card 1', backText: 'Translation 1' },
    { id: '2', frontText: 'Card 2', backText: 'Translation 2' }
  ];

  test('renders the FlashcardReview component with the current card', () => {
    const { getByTestId, getByText } = render(StudySession, {
      props: { cards: mockCards }
    });
    
    expect(getByTestId('flashcard-review')).toBeInTheDocument();
    expect(getByText('Card 1')).toBeInTheDocument();
  });
});
```

### Example 2: Testing Deck Component with Actions

The Deck component includes action buttons that needed to be tested:

```javascript
describe('Deck', () => {
  test('shows action buttons when showActions is true', () => {
    const { getByText } = render(Deck, {
      props: {
        deck: { id: '1', name: 'Test Deck', cardCount: 5 },
        showActions: true
      }
    });
    
    expect(getByText('Study')).toBeInTheDocument();
    expect(getByText('Edit')).toBeInTheDocument();
    expect(getByText('Delete')).toBeInTheDocument();
  });
  
  test('calls onDelete when delete button is clicked', async () => {
    const handleDelete = jest.fn();
    
    const { getByText } = render(Deck, {
      props: {
        deck: { id: '1', name: 'Test Deck', cardCount: 5 },
        showActions: true,
        onDelete: handleDelete
      }
    });
    
    await fireEvent.click(getByText('Delete'));
    expect(handleDelete).toHaveBeenCalledWith('1');
  });
});
```

### Example 3: Testing Login Component with API Calls

The Login component makes API calls that needed to be mocked:

```javascript
// Mock the API module
jest.mock('../../lib/api', () => ({
  login: jest.fn().mockResolvedValue({ token: 'mock-token', user: { id: 1, username: 'testuser' } })
}));

describe('Login', () => {
  test('calls login API and navigates on successful login', async () => {
    const navigateMock = jest.fn();
    jest.mock('svelte-routing', () => ({
      ...jest.requireActual('svelte-routing'),
      navigate: navigateMock
    }));
    
    const { getByLabelText, getByText } = render(Login);
    
    await fireEvent.input(getByLabelText('Username'), { target: { value: 'testuser' } });
    await fireEvent.input(getByLabelText('Password'), { target: { value: 'password' } });
    await fireEvent.click(getByText('Login'));
    
    expect(api.login).toHaveBeenCalledWith('testuser', 'password');
    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith('/decks');
    });
  });
});
```

## Benefits of Our Approach

1. **Simplicity**: Tests are easier to write and understand
2. **Isolation**: Components can be tested in isolation without complex dependencies
3. **Reliability**: Tests are less brittle and less likely to break with changes
4. **Coverage**: We can test components that would be difficult to test with standard approaches

## Limitations and Considerations

1. **Integration Testing**: This approach focuses on unit testing; integration tests require a different approach
2. **Svelte Features**: Some Svelte-specific features may not be fully tested with this approach
3. **Maintenance**: Mock implementations need to be kept in sync with actual components

## Conclusion

Our custom testing approach for Svelte components has allowed us to achieve high test coverage and reliable tests despite the challenges of testing Svelte components. By using direct component mocking, custom render functions, global mock helpers, and environment mocking, we've created a testing strategy that works well for our application's needs.
