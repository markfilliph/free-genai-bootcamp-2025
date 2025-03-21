# Svelte Component Testing Guide

This guide explains our approach to testing Svelte components in the Language Learning Flashcard Generator application.

## Challenges with Svelte Component Testing

Testing Svelte components presents several challenges:

1. **Component Dependencies**: Svelte components often depend on other components, making isolation difficult.
2. **External Libraries**: Components that use external libraries like svelte-routing can be hard to test.
3. **Lifecycle Methods**: Svelte's component lifecycle methods can be difficult to mock.
4. **Environment Dependencies**: Components that rely on browser APIs or environment variables need special handling.

## Our Testing Approach

We've adopted a pragmatic approach to testing Svelte components:

### 1. Direct Component Mocking

Instead of trying to render actual Svelte components with all their dependencies, we create simplified mock versions that return the expected HTML structure:

```javascript
jest.mock('../../components/ComponentName.svelte', () => ({
  default: {
    render: (props) => {
      return {
        html: `<div>Mocked Component HTML</div>`,
        instance: {
          // Mock component methods
          someMethod: jest.fn()
        }
      };
    }
  }
}));
```

This approach:
- Isolates component tests from their dependencies
- Avoids issues with external libraries
- Allows precise control over component output

### 2. Custom Render Functions

We use the `$$render` method to create custom renderers for Svelte components:

```javascript
jest.mock('svelte-routing', () => ({
  Link: {
    $$render: ($$result, $$props, $$bindings, $$slots) => {
      const { to, class: className } = $$props;
      return `<a href="${to}" class="${className}">${$$slots.default ? $$slots.default({}) : ''}</a>`;
    }
  }
}));
```

This allows us to control exactly what HTML is generated during tests.

### 3. Global Mock Helpers

We've added global helpers to create consistent mocks:

```javascript
global.mockSvelteComponent = (html) => ({
  $$render: () => html
});
```

### 4. Environment Mocking

For components that depend on environment variables or browser APIs:

```javascript
// Mock the import.meta.env for Vite
global.import = {};
global.import.meta = {};
global.import.meta.env = {
  VITE_API_URL: 'http://localhost:8000'
};

// Mock the browser's localStorage
const localStorageMock = { /* ... */ };
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});
```

## Example: Testing a FlashcardReview Component

Here's how we test the FlashcardReview component:

1. **Mock the Component**:
```javascript
jest.mock('../../components/FlashcardReview.svelte', () => {
  let showAnswer = false;
  
  return {
    default: {
      render: (props) => {
        // Generate HTML based on current state
        let html = '<div class="flashcard-review">';
        // ... generate HTML based on props and state
        html += '</div>';
        
        return {
          html,
          instance: {
            toggleAnswer: () => { showAnswer = !showAnswer; }
          }
        };
      }
    }
  };
});
```

2. **Write Tests**:
```javascript
test('shows answer when button is clicked', async () => {
  const { getByText, queryByText } = render(FlashcardReview, { 
    props: { flashcard: mockFlashcard } 
  });
  
  // Initially, answer should not be visible
  expect(queryByText('Hello, how are you?')).not.toBeInTheDocument();
  
  // Click "Show Answer" button
  await fireEvent.click(getByText('Show Answer'));
  
  // Now answer should be visible
  expect(getByText('Hello, how are you?')).toBeInTheDocument();
});
```

## Benefits of This Approach

1. **Isolation**: Each component can be tested in isolation without dependencies.
2. **Simplicity**: Tests are easier to write and maintain.
3. **Speed**: Tests run faster since they don't need to render complex component trees.
4. **Reliability**: Tests are less likely to break due to changes in dependencies.

## Limitations

1. **Not True Integration Tests**: These tests don't verify how components work together.
2. **Manual HTML Generation**: Requires manually generating HTML that matches component output.
3. **Maintenance**: Mock implementations need to be updated when components change.

## When to Use Different Approaches

- **Unit Testing**: Use the mocking approach described here.
- **Integration Testing**: Consider using Svelte's testing utilities directly.
- **End-to-End Testing**: Use tools like Cypress for full application testing.

## Conclusion

This approach provides a pragmatic balance between test coverage and maintainability. It allows us to test component behavior without getting bogged down in the complexities of Svelte's component system and external dependencies.
