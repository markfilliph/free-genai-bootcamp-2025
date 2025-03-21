// This file will be automatically loaded by Jest

// Add custom Jest matchers for DOM elements
import '@testing-library/jest-dom/extend-expect';

// Mock the browser environment
global.MutationObserver = class {
  constructor(callback) {}
  disconnect() {}
  observe(element, initObject) {}
};

// Mock the browser's localStorage
const localStorageMock = (function() {
  let store = {};
  return {
    getItem: function(key) {
      return store[key] || null;
    },
    setItem: function(key, value) {
      store[key] = value.toString();
    },
    removeItem: function(key) {
      delete store[key];
    },
    clear: function() {
      store = {};
    }
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});

// Mock fetch API
global.fetch = jest.fn();

// Mock the import.meta.env for Vite
global.import = {};
global.import.meta = {};
global.import.meta.env = {
  VITE_API_URL: 'http://localhost:8000'
};

// Mock the API module
jest.mock('../lib/api.js', () => ({
  API_BASE: 'http://localhost:8000',
  apiFetch: jest.fn().mockImplementation((endpoint, options = {}) => {
    // Mock implementation will be added in individual tests
    return Promise.resolve({ message: 'Mock API response' });
  })
}));

// Define mock data
global.__mocks__ = {
  decks: [
    {
      id: '1',
      name: 'Spanish Basics',
      user_id: '1',
      created_at: '2025-03-19T10:00:00Z'
    },
    {
      id: '2',
      name: 'French Vocabulary',
      user_id: '1',
      created_at: '2025-03-18T09:30:00Z'
    }
  ],
  flashcards: [
    {
      id: '1',
      deck_id: '1',
      word: 'hola',
      translation: 'hello',
      example_sentence: 'Hola, ¿cómo estás?',
      cultural_note: 'Common greeting in Spanish-speaking countries',
      conjugation: null,
      created_at: '2025-03-19T10:05:00Z'
    },
    {
      id: '2',
      deck_id: '1',
      word: 'adiós',
      translation: 'goodbye',
      example_sentence: 'Adiós, hasta mañana.',
      cultural_note: 'Standard farewell in Spanish',
      conjugation: null,
      created_at: '2025-03-19T10:10:00Z'
    }
  ]
};

// Simple mock for svelte-routing
jest.mock('svelte-routing', () => ({
  Link: jest.fn().mockImplementation(() => ({})),
  navigate: jest.fn(),
  Router: jest.fn().mockImplementation(() => ({})),
  Route: jest.fn().mockImplementation(() => ({}))
}));

// Create a helper for mocking Svelte components
global.mockComponent = (name) => {
  return jest.fn().mockImplementation(() => ({
    $set: jest.fn(),
    $on: jest.fn(),
    $destroy: jest.fn()
  }));
};

// Mock all components
const componentMocks = {
  // Components
  '../components/DeckList.svelte': () => mockComponent('DeckList'),
  '../components/FlashcardReview.svelte': () => mockComponent('FlashcardReview'),
  '../components/StudySession.svelte': () => mockComponent('StudySession'),
  '../components/FlashcardForm.svelte': () => mockComponent('FlashcardForm'),
  '../components/Navbar.svelte': () => mockComponent('Navbar'),
  '../components/Deck.svelte': () => mockComponent('Deck'),
  
  // Routes
  '../routes/Home.svelte': () => mockComponent('Home'),
  '../routes/Login.svelte': () => mockComponent('Login'),
  '../routes/DeckManagement.svelte': () => mockComponent('DeckManagement'),
  '../routes/App.svelte': () => mockComponent('App')
};

// Apply all mocks
Object.entries(componentMocks).forEach(([path, mockFn]) => {
  jest.mock(path, () => ({
    default: mockFn()
  }), { virtual: true });
});
