// Simple mock API implementation
export const API_BASE = 'http://localhost:8000';

// Mock response data
const mockResponses = {
  '/auth/login': {
    success: { token: 'mock-token-123', user_id: '1' },
    error: new Error('API Error (401): Unauthorized')
  },
  '/auth/register': {
    success: { user_id: '2', username: 'newuser', email: 'newuser@example.com' },
    error: new Error('API Error (400): Email already exists')
  },
  '/decks': [
    {
      id: '1',
      name: 'Spanish Basics',
      user_id: '1',
      created_at: '2025-03-19T10:00:00Z'
    },
    {
      id: '2',
      name: 'Verb Conjugations',
      user_id: '1',
      created_at: '2025-03-18T09:30:00Z'
    }
  ],
  '/decks/1/flashcards': [
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
  ],
  '/generate/word': {
    word: 'hola',
    translation: 'hello',
    example_sentence: 'Hola, ¿cómo estás?',
    cultural_note: 'Common greeting in Spanish-speaking countries'
  },
  '/generate/verb': {
    word: 'hablar',
    translation: 'to speak',
    example_sentence: 'Yo hablo español.',
    cultural_note: 'Regular -ar verb',
    conjugation: 'hablo, hablas, habla, hablamos, habláis, hablan'
  }
};

// Mock API fetch function
export const apiFetch = jest.fn().mockImplementation((endpoint, options = {}) => {
  // Extract the base endpoint without query parameters
  const baseEndpoint = endpoint.split('?')[0];
  
  // Handle login endpoint with special logic
  if (baseEndpoint === '/auth/login') {
    const body = options.body ? JSON.parse(options.body) : {};
    
    // Check if this is a test for invalid credentials
    if (body.password === 'wrongpassword') {
      return Promise.reject(mockResponses['/auth/login'].error);
    }
    
    return Promise.resolve(mockResponses['/auth/login'].success);
  }
  
  // Handle register endpoint with special logic
  if (baseEndpoint === '/auth/register') {
    const body = options.body ? JSON.parse(options.body) : {};
    
    // Check if this is a test for existing email
    if (body.email === 'existing@example.com') {
      return Promise.reject(mockResponses['/auth/register'].error);
    }
    
    return Promise.resolve(mockResponses['/auth/register'].success);
  }
  
  // Handle network error test
  if (baseEndpoint === '/test-network-error') {
    return Promise.reject(new Error('Network error'));
  }
  
  // For other endpoints, return the mock response if available
  if (mockResponses[baseEndpoint]) {
    return Promise.resolve(mockResponses[baseEndpoint]);
  }
  
  // For deck creation
  if (baseEndpoint === '/decks' && options.method === 'POST') {
    const body = options.body ? JSON.parse(options.body) : {};
    return Promise.resolve({
      id: '3',
      name: body.name || 'New Deck',
      user_id: '1',
      created_at: new Date().toISOString()
    });
  }
  
  // For flashcard creation
  if (baseEndpoint.match(/\/decks\/\d+\/flashcards/) && options.method === 'POST') {
    const body = options.body ? JSON.parse(options.body) : {};
    return Promise.resolve({
      id: '3',
      deck_id: baseEndpoint.split('/')[2],
      word: body.word || 'new word',
      translation: body.translation || 'new translation',
      example_sentence: body.example_sentence || 'Example sentence.',
      cultural_note: body.cultural_note || 'Cultural note.',
      conjugation: body.conjugation || null,
      created_at: new Date().toISOString()
    });
  }
  
  // Default response
  return Promise.resolve({ message: 'Mock API response' });
});
