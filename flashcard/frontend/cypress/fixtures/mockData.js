/**
 * Mock data for Cypress tests
 * This file contains all the mock data used in the E2E tests
 */

// Authentication mocks
export const mockUser = {
  id: '1',
  email: 'test@example.com',
  name: 'Test User'
};

export const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaWF0IjoxNjE2MTYyMjIyLCJleHAiOjE2MTYyNDg2MjJ9.fake-token';

// Deck mocks
export const mockDecks = [
  {
    id: '1',
    name: 'Spanish Basics',
    description: 'Essential Spanish vocabulary',
    language: 'Spanish',
    createdAt: '2023-01-01T00:00:00.000Z',
    updatedAt: '2023-01-01T00:00:00.000Z',
    userId: '1',
    flashcardCount: 10
  },
  {
    id: '2',
    name: 'French Phrases',
    description: 'Common French expressions',
    language: 'French',
    createdAt: '2023-01-02T00:00:00.000Z',
    updatedAt: '2023-01-02T00:00:00.000Z',
    userId: '1',
    flashcardCount: 5
  },
  {
    id: '3',
    name: 'German Vocabulary',
    description: 'Basic German words',
    language: 'German',
    createdAt: '2023-01-03T00:00:00.000Z',
    updatedAt: '2023-01-03T00:00:00.000Z',
    userId: '1',
    flashcardCount: 8
  }
];

// Flashcard mocks
export const mockFlashcards = {
  '1': [
    {
      id: '101',
      word: 'hola',
      translation: 'hello',
      examples: ['¡Hola! ¿Cómo estás?'],
      notes: 'Common greeting in Spanish',
      wordType: 'noun',
      deckId: '1',
      createdAt: '2023-01-01T00:00:00.000Z',
      updatedAt: '2023-01-01T00:00:00.000Z',
      difficulty: 'easy'
    },
    {
      id: '102',
      word: 'adiós',
      translation: 'goodbye',
      examples: ['Adiós, hasta mañana'],
      notes: 'Used when parting ways',
      wordType: 'noun',
      deckId: '1',
      createdAt: '2023-01-01T00:00:00.000Z',
      updatedAt: '2023-01-01T00:00:00.000Z',
      difficulty: 'medium'
    },
    {
      id: '103',
      word: 'gracias',
      translation: 'thank you',
      examples: ['Muchas gracias por tu ayuda'],
      notes: 'Expression of gratitude',
      wordType: 'noun',
      deckId: '1',
      createdAt: '2023-01-01T00:00:00.000Z',
      updatedAt: '2023-01-01T00:00:00.000Z',
      difficulty: 'easy'
    }
  ],
  '2': [
    {
      id: '201',
      word: 'bonjour',
      translation: 'hello',
      examples: ['Bonjour, comment ça va?'],
      notes: 'Common greeting in French',
      wordType: 'noun',
      deckId: '2',
      createdAt: '2023-01-02T00:00:00.000Z',
      updatedAt: '2023-01-02T00:00:00.000Z',
      difficulty: 'easy'
    },
    {
      id: '202',
      word: 'au revoir',
      translation: 'goodbye',
      examples: ['Au revoir, à bientôt!'],
      notes: 'Used when parting ways',
      wordType: 'noun',
      deckId: '2',
      createdAt: '2023-01-02T00:00:00.000Z',
      updatedAt: '2023-01-02T00:00:00.000Z',
      difficulty: 'medium'
    }
  ]
};

// Study session mocks
export const mockStudySession = {
  id: '1',
  deckId: '1',
  userId: '1',
  startedAt: '2023-01-10T00:00:00.000Z',
  completedAt: null,
  flashcards: mockFlashcards['1'].map(card => ({
    ...card,
    sessionRating: null
  }))
};

export const mockStudySessionResults = {
  id: '1',
  deckId: '1',
  userId: '1',
  startedAt: '2023-01-10T00:00:00.000Z',
  completedAt: '2023-01-10T00:10:00.000Z',
  stats: {
    totalCards: 3,
    easyRatings: 2,
    mediumRatings: 1,
    hardRatings: 0,
    timeSpent: 600 // seconds
  }
};

// Error responses
export const mockErrors = {
  unauthorized: {
    statusCode: 401,
    body: {
      error: 'Unauthorized',
      message: 'You must be logged in to access this resource'
    }
  },
  notFound: {
    statusCode: 404,
    body: {
      error: 'Not Found',
      message: 'The requested resource was not found'
    }
  },
  badRequest: {
    statusCode: 400,
    body: {
      error: 'Bad Request',
      message: 'Invalid input data'
    }
  },
  serverError: {
    statusCode: 500,
    body: {
      error: 'Internal Server Error',
      message: 'Something went wrong on the server'
    }
  }
};

// Helper function to generate a new deck
export const generateDeck = (overrides = {}) => {
  const id = overrides.id || `deck-${Date.now()}`;
  return {
    id,
    name: overrides.name || `Test Deck ${id}`,
    description: overrides.description || 'Auto-generated test deck',
    language: overrides.language || 'English',
    createdAt: overrides.createdAt || new Date().toISOString(),
    updatedAt: overrides.updatedAt || new Date().toISOString(),
    userId: overrides.userId || '1',
    flashcardCount: overrides.flashcardCount || 0
  };
};

// Helper function to generate a new flashcard
export const generateFlashcard = (deckId, overrides = {}) => {
  const id = overrides.id || `card-${Date.now()}`;
  return {
    id,
    word: overrides.word || `Word ${id}`,
    translation: overrides.translation || `Translation ${id}`,
    examples: overrides.examples || [`Example sentence for ${id}`],
    notes: overrides.notes || `Notes for ${id}`,
    wordType: overrides.wordType || 'noun',
    deckId: deckId,
    createdAt: overrides.createdAt || new Date().toISOString(),
    updatedAt: overrides.updatedAt || new Date().toISOString(),
    difficulty: overrides.difficulty || 'medium'
  };
};
