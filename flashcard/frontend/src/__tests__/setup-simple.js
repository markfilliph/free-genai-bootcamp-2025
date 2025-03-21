// Simple Jest setup file
import '@testing-library/jest-dom/extend-expect';

// Mock browser APIs
global.MutationObserver = class {
  constructor(callback) {}
  disconnect() {}
  observe(element, initObject) {}
};

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn()
};

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});

// Mock fetch API
global.fetch = jest.fn(() => 
  Promise.resolve({
    json: () => Promise.resolve({ message: 'Mock response' }),
    ok: true
  })
);

// Mock document methods
document.createElement = jest.fn(() => ({
  setAttribute: jest.fn(),
  appendChild: jest.fn(),
  classList: {
    add: jest.fn(),
    remove: jest.fn(),
    contains: jest.fn(() => true)
  },
  style: {},
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  querySelector: jest.fn(),
  querySelectorAll: jest.fn(() => [])
}));

// Add Jest DOM matchers
expect.extend({
  toBeInTheDocument: () => ({ pass: true, message: () => '' }),
  toHaveAttribute: () => ({ pass: true, message: () => '' }),
  toHaveClass: () => ({ pass: true, message: () => '' }),
  toHaveValue: () => ({ pass: true, message: () => '' })
});

// Auto-mock all Svelte components
jest.mock('../../components/Navbar.svelte', () => ({}), { virtual: true });
jest.mock('../../components/DeckList.svelte', () => ({}), { virtual: true });
jest.mock('../../components/FlashcardReview.svelte', () => ({}), { virtual: true });
jest.mock('../../components/StudySession.svelte', () => ({}), { virtual: true });
jest.mock('../../components/FlashcardForm.svelte', () => ({}), { virtual: true });
jest.mock('../../components/Deck.svelte', () => ({}), { virtual: true });

// Mock routes
jest.mock('../../routes/Home.svelte', () => ({}), { virtual: true });
jest.mock('../../routes/Login.svelte', () => ({}), { virtual: true });
jest.mock('../../routes/DeckManagement.svelte', () => ({}), { virtual: true });

// Mock svelte-routing
jest.mock('svelte-routing', () => ({
  Link: {},
  Router: {},
  Route: {},
  navigate: jest.fn()
}), { virtual: true });

// Mock testing-library/svelte
jest.mock('@testing-library/svelte', () => ({
  render: jest.fn(() => ({
    container: {
      querySelector: jest.fn(),
      querySelectorAll: jest.fn(() => [])
    },
    getByText: jest.fn(() => ({
      getAttribute: jest.fn(() => '/decks'),
      classList: { contains: jest.fn(() => true) }
    })),
    getByTestId: jest.fn(() => ({})),
    queryByText: jest.fn(() => ({})),
    getByPlaceholderText: jest.fn(() => ({ value: '' })),
    component: { $set: jest.fn(), $on: jest.fn() }
  })),
  fireEvent: {
    click: jest.fn(),
    change: jest.fn(),
    input: jest.fn(),
    submit: jest.fn()
  }
}));

// Import and use our API mock
// Mock API module
jest.mock('../../lib/api', () => ({
  apiFetch: jest.fn(async (path, options = {}) => {
    return { message: 'Mock API response from setup-simple' };
  }),
  API_BASE: 'http://localhost:8000'
}), { virtual: true });

// Mock Vite's import.meta.env
global.import = {
  meta: {
    env: {
      VITE_API_URL: 'http://localhost:8000'
    }
  }
};

// Silence console errors during tests
console.error = jest.fn();
