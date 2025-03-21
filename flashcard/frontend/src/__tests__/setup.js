// This file will be automatically loaded by Jest

// Add custom Jest matchers for DOM elements
import '@testing-library/jest-dom/extend-expect';
import '@testing-library/jest-dom';
import { mockComponents } from './mocks/component-mocks';
import { mockAPI } from './mocks/api-simple';

beforeAll(() => {
  mockComponents();
  mockAPI();
});

// Mock the browser environment
global.MutationObserver = class {
  constructor(callback) {}
  disconnect() {}
  observe(element, initObject) {}
};

// Mock the browser's fetch
global.fetch = jest.fn(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ message: 'Mock API response' }),
    text: () => Promise.resolve('Mock text response')
  })
);

// Mock the browser's localStorage
const createLocalStorageMock = () => {
  const store = {};
  return {
    getItem: jest.fn(key => store[key] || null),
    setItem: jest.fn((key, value) => { store[key] = String(value); }),
    removeItem: jest.fn(key => { delete store[key]; }),
    clear: jest.fn(() => { Object.keys(store).forEach(key => delete store[key]); }),
    _store: store // For test inspection
  };
};

const localStorageMock = createLocalStorageMock();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true
});

// Mock the import.meta.env for Vite
global.import = {};
global.import.meta = {};
global.import.meta.env = {
  VITE_API_URL: 'http://localhost:8000'
};

// Mock the API_BASE for tests
jest.mock('../lib/api.js', () => {
  const originalModule = jest.requireActual('../lib/api.js');
  return {
    ...originalModule,
    API_BASE: 'http://localhost:8000',
    apiFetch: jest.fn(async (path, options = {}) => {
      const response = await global.fetch(`http://localhost:8000${path}`, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers
        },
        credentials: 'include',
        ...options
      });
      
      if (!response.ok) {
        const text = await response.text();
        throw new Error(`API Error (${response.status}): ${text}`);
      }
      
      return response.json();
    })
  };
});

// Mock svelte-routing
jest.mock('svelte-routing', () => {
  return {
    Link: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        const { to, class: className } = $$props;
        return `<a href="${to}" class="${className}">${$$slots.default ? $$slots.default({}) : ''}</a>`;
      }
    },
    navigate: jest.fn(),
    Router: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        return $$slots.default ? $$slots.default({}) : '';
      }
    },
    Route: {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        return $$slots.default ? $$slots.default({}) : '';
      }
    }
  };
});

// Create a helper for mocking Svelte components
global.mockSvelteComponent = (html) => ({
  $$render: () => html,
  $on: jest.fn(),
  dispatchEvent: jest.fn()
});

// Mock DeckList component for DeckManagement tests
jest.mock('../../components/DeckList.svelte', () => {
  return {
    default: {
      render: (props) => {
        const decks = props?.decks || [];
        let html = '<div class="deck-list">';
        
        if (decks.length === 0) {
          html += '<p>No decks found. Create your first deck!</p>';
        } else {
          decks.forEach(deck => {
            html += `<div class="deck-item" data-deck-id="${deck.id}">`;
            html += `<h3>${deck.name}</h3>`;
            html += `<p>Created: ${new Date(deck.created_at).toLocaleDateString()}</p>`;
            html += '</div>';
          });
        }
        
        html += '</div>';
        
        return {
          html,
          props,
          $on: jest.fn(),
          dispatchEvent: jest.fn()
        };
      }
    }
  };
});

// Create a global mock constructor for DeckList
global.mockDeckList = function(options) {
  const { props } = options || {};
  const decks = props?.decks || [];
  
  return {
    $$render: () => {
      let html = '<div class="deck-list">';
      
      if (decks.length === 0) {
        html += '<p>No decks found. Create your first deck!</p>';
      } else {
        decks.forEach(deck => {
          html += `<div class="deck-item" data-deck-id="${deck.id}">`;
          html += `<h3>${deck.name}</h3>`;
          html += `<p>Created: ${new Date(deck.created_at).toLocaleDateString()}</p>`;
          html += '</div>';
        });
      }
      
      html += '</div>';
      return html;
    }
  };
};

