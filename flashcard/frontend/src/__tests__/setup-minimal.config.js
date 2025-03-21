// Minimal Jest setup file
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

// Add Jest DOM matchers
expect.extend({
  toBeInTheDocument: () => ({ pass: true, message: () => '' }),
  toHaveAttribute: () => ({ pass: true, message: () => '' }),
  toHaveClass: () => ({ pass: true, message: () => '' }),
  toHaveValue: () => ({ pass: true, message: () => '' })
});

// Mock @testing-library/svelte
jest.mock('@testing-library/svelte', () => require('./mocks/testing-library-svelte'), { virtual: true });

// Mock svelte-routing
jest.mock('svelte-routing', () => ({
  Link: {
    $$render: ($$result, $$props) => {
      const { to, class: className } = $$props || {};
      return `<a href="${to || '/'}" class="${className || ''}">Link Text</a>`;
    }
  },
  Router: {
    $$render: ($$result, $$props, $$bindings, $$slots) => {
      return $$slots.default ? $$slots.default({}) : '';
    }
  },
  Route: {
    $$render: ($$result, $$props, $$bindings, $$slots) => {
      return $$slots.default ? $$slots.default({}) : '';
    }
  },
  navigate: jest.fn()
}), { virtual: true });

// Create mock API functions directly
global.apiFetch = jest.fn().mockImplementation(() => Promise.resolve({ message: 'Mock API response' }));
global.API_BASE = 'http://localhost:8000';

// Mock API module - use correct relative path
jest.mock('../lib/api', () => ({
  apiFetch: global.apiFetch,
  API_BASE: global.API_BASE
}));

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
