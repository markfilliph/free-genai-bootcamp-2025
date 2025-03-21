// Mock API implementation
export const API_BASE = 'http://localhost:8000';

// Mock decks data
export const mockDecks = [
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
];

// Mock flashcards data
export const mockFlashcards = [
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
];

// Mock API fetch function
export const apiFetch = jest.fn().mockImplementation((endpoint, options = {}) => {
  // Return mock data based on endpoint
  if (endpoint === '/decks') {
    return Promise.resolve(mockDecks);
  } else if (endpoint.startsWith('/decks/') && endpoint.endsWith('/flashcards')) {
    return Promise.resolve(mockFlashcards);
  } else if (endpoint === '/auth/login') {
    if (options.method === 'POST') {
      const body = JSON.parse(options.body);
      if (body.username === 'testuser' && body.password === 'password') {
        return Promise.resolve({ token: 'mock-token', user_id: '1' });
      } else {
        return Promise.reject(new Error('Invalid credentials'));
      }
    }
  }
  
  // Default response
  return Promise.resolve({ message: 'Mock API response' });
});
