/**
 * Component Mocks
 * 
 * This file contains mocks for Svelte components and routing libraries
 * that are used in the application but need to be mocked in the Jest environment.
 */

// Helper function to create a mock Svelte component
export function mockComponent(name, html = '') {
  return {
    $$render: () => html || `<div data-testid="mock-${name}">${name} Component</div>`,
    render: (props) => html || `<div data-testid="mock-${name}">${name} Component</div>`
  };
}

// Helper function to create HTML content for component mocks
export function mockHtml(content) {
  return content;
}

// Mock svelte-routing
jest.mock('svelte-routing', () => {
  return {
    Link: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        const { to, class: className } = $$props || {};
        return `<a href="${to || '/'}" class="${className || ''}">${$$slots.default ? $$slots.default() : 'Link'}</a>`;
      }
    },
    Router: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        return $$slots.default ? $$slots.default() : '';
      }
    },
    Route: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        return $$slots.default ? $$slots.default() : '';
      }
    },
    navigate: jest.fn()
  };
});

// Mock common components
const mockComponents = {
  // Navigation components
  '../../components/Navbar.svelte': () => mockComponent('Navbar'),
  
  // Deck components
  '../../components/DeckList.svelte': () => mockComponent('DeckList', `
    <div data-testid="mock-deck-list">
      <ul>
        <li>Mock Deck 1</li>
        <li>Mock Deck 2</li>
      </ul>
    </div>
  `),
  '../../components/Deck.svelte': () => mockComponent('Deck'),
  
  // Flashcard components
  '../../components/FlashcardForm.svelte': () => mockComponent('FlashcardForm'),
  '../../components/FlashcardReview.svelte': () => mockComponent('FlashcardReview'),
  '../../components/StudySession.svelte': () => mockComponent('StudySession')
};

// Apply all component mocks
Object.entries(mockComponents).forEach(([path, mockFn]) => {
  jest.mock(path, () => ({ default: mockFn() }), { virtual: true });
});
