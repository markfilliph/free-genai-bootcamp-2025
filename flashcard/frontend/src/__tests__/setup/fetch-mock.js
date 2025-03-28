/**
 * Fetch API Mocks
 * 
 * This file contains mocks for the Fetch API that is used for network requests
 * in the application but needs to be mocked in the Jest environment.
 */

// Create a basic fetch mock that returns a successful response
global.fetch = jest.fn(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ message: 'Mock API response' }),
    text: () => Promise.resolve('Mock text response'),
    status: 200,
    headers: new Map()
  })
);

// Mock the API module
jest.mock('../../lib/api.js', () => ({
  API_BASE: 'http://localhost:8000',
  apiFetch: jest.fn().mockImplementation((endpoint, options = {}) => {
    // Default mock implementation
    return Promise.resolve({ message: 'Mock API response' });
  }),
  
  // Add specific API function mocks
  login: jest.fn().mockResolvedValue({ token: 'mock-token', user: { id: 1, username: 'testuser' } }),
  register: jest.fn().mockResolvedValue({ success: true }),
  getDecks: jest.fn().mockResolvedValue([]),
  createDeck: jest.fn().mockResolvedValue({ id: '1', name: 'New Deck' }),
  getFlashcards: jest.fn().mockResolvedValue([]),
  createFlashcard: jest.fn().mockResolvedValue({ id: '1', word: 'test' }),
  generateContent: jest.fn().mockResolvedValue({
    example_sentences: ['Example sentence'],
    conjugations: 'Verb conjugations',
    cultural_note: 'Cultural note'
  })
}));
