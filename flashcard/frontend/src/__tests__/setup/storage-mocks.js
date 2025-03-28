/**
 * Storage API Mocks
 * 
 * This file contains mocks for browser storage APIs (localStorage, sessionStorage)
 * that are used in the application but need to be mocked in the Jest environment.
 */

// Create a functional localStorage mock that maintains state during tests
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

// Create the localStorage mock instance
const localStorageMock = createLocalStorageMock();

// Attach localStorage mock to window
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true
});

// Also mock sessionStorage for completeness
Object.defineProperty(window, 'sessionStorage', {
  value: createLocalStorageMock(),
  writable: true
});
