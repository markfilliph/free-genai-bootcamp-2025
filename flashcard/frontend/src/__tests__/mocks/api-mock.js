// Mock API responses for testing

// Mock user data
export const mockUsers = [
  {
    id: '1',
    username: 'testuser',
    email: 'test@example.com'
  }
];

// Mock deck data
export const mockDecks = [
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
    created_at: '2025-03-19T11:00:00Z'
  }
];

// Mock flashcard data
export const mockFlashcards = [
  {
    id: '1',
    word: 'hola',
    example_sentence: 'Hola, ¿cómo estás?',
    translation: 'Hello, how are you?',
    conjugation: null,
    cultural_note: 'Common greeting in Spanish-speaking countries.',
    deck_id: '1',
    created_at: '2025-03-19T10:30:00Z'
  },
  {
    id: '2',
    word: 'adiós',
    example_sentence: 'Adiós, hasta mañana.',
    translation: 'Goodbye, see you tomorrow.',
    conjugation: null,
    cultural_note: 'Common farewell in Spanish-speaking countries.',
    deck_id: '1',
    created_at: '2025-03-19T10:35:00Z'
  },
  {
    id: '3',
    word: 'hablar',
    example_sentence: 'Yo hablo español.',
    translation: 'I speak Spanish.',
    conjugation: 'hablo, hablas, habla, hablamos, habláis, hablan',
    cultural_note: 'Regular -ar verb.',
    deck_id: '2',
    created_at: '2025-03-19T11:30:00Z'
  }
];

// Mock API responses
export const mockApiResponses = {
  // Auth endpoints
  '/auth/login': {
    POST: (data) => {
      const { email, password } = JSON.parse(data.body);
      if (email === 'test@example.com' && password === 'password123') {
        return {
          status: 200,
          body: {
            token: 'mock-token-123',
            user_id: '1'
          }
        };
      }
      return {
        status: 401,
        error: 'Invalid credentials'
      };
    }
  },
  
  '/auth/register': {
    POST: (data) => {
      const { username, email, password } = JSON.parse(data.body);
      if (username && email && password) {
        return {
          status: 200,
          body: {
            user_id: '2',
            username,
            email
          }
        };
      }
      return {
        status: 400,
        error: 'Missing required fields'
      };
    }
  },
  
  // Decks endpoints
  '/decks': {
    GET: () => ({
      status: 200,
      body: mockDecks
    }),
    POST: (data) => {
      const { name } = JSON.parse(data.body);
      if (name) {
        return {
          status: 200,
          body: {
            id: '3',
            name,
            user_id: '1',
            created_at: new Date().toISOString()
          }
        };
      }
      return {
        status: 400,
        error: 'Missing required fields'
      };
    }
  },
  
  // Flashcards endpoints
  '/flashcards': {
    POST: (data) => {
      const { word, example_sentence, translation, deck_id } = JSON.parse(data.body);
      if (word && example_sentence && translation && deck_id) {
        return {
          status: 200,
          body: {
            id: '4',
            word,
            example_sentence,
            translation,
            deck_id,
            created_at: new Date().toISOString()
          }
        };
      }
      return {
        status: 400,
        error: 'Missing required fields'
      };
    }
  },
  
  // Deck flashcards endpoint
  '/decks/1/flashcards': {
    POST: (data) => {
      const { word, example_sentence, translation } = JSON.parse(data.body);
      if (word && example_sentence && translation) {
        return {
          status: 200,
          body: {
            id: '3',
            deck_id: '1',
            word,
            example_sentence,
            translation,
            cultural_note: 'Cultural note.',
            conjugation: null,
            created_at: new Date().toISOString()
          }
        };
      }
      return {
        status: 400,
        error: 'Missing required fields'
      };
    }
  },
  
  // Generate endpoints
  '/generate/word': {
    GET: () => ({
      status: 200,
      body: {
        word: 'hola',
        translation: 'hello',
        example_sentence: 'Hola, ¿cómo estás?',
        cultural_note: 'Common greeting in Spanish-speaking countries'
      }
    })
  },
  
  '/generate/verb': {
    GET: () => ({
      status: 200,
      body: {
        word: 'hablar',
        translation: 'to speak',
        example_sentence: 'Yo hablo español.',
        cultural_note: 'Regular -ar verb',
        conjugation: 'hablo, hablas, habla, hablamos, habláis, hablan'
      }
    })
  },
  
  // Generate endpoint (generic)
  '/generate': {
    POST: (data) => {
      const { word, is_verb } = JSON.parse(data.body);
      if (word) {
        return {
          status: 200,
          body: {
            example_sentences: [
              `This is an example sentence with the word "${word}".`,
              `Here's another example using "${word}" in context.`
            ],
            conjugations: is_verb ? 'Mock conjugations for verb' : null,
            cultural_note: `Cultural note about "${word}".`
          }
        };
      }
      return {
        status: 400,
        error: 'Missing required fields'
      };
    }
  }
};

// Mock fetch implementation
export const mockFetch = (url, options = {}) => {
  return new Promise((resolve, reject) => {
    // Extract the path from the URL
    const path = url.replace(/^.*\/\/[^/]+/, '');
    
    // Get the mock response for this path and method
    const mockResponse = mockApiResponses[path];
    if (!mockResponse) {
      // Use a default response if no specific mock is defined
      resolve({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ message: 'Default mock response' }),
        text: () => Promise.resolve(JSON.stringify({ message: 'Default mock response' }))
      });
      return;
    }
    
    const method = options.method || 'GET';
    const handler = mockResponse[method];
    if (!handler) {
      reject(new Error(`No mock response for ${method} ${path}`));
      return;
    }
    
    // Get the response data
    const responseData = handler(options);
    
    // If there's an error, return an error response
    if (responseData.error) {
      resolve({
        ok: false,
        status: responseData.status,
        text: () => Promise.resolve(responseData.error)
      });
      return;
    }
    
    // Return a successful response
    resolve({
      ok: true,
      status: responseData.status || 200,
      json: () => Promise.resolve(responseData.body),
      text: () => Promise.resolve(JSON.stringify(responseData.body))
    });
  });
};

// Add a simple test to prevent the "no tests" error
describe('API Mock', () => {
  test('mockUsers contains expected test data', () => {
    expect(mockUsers).toBeDefined();
    expect(Array.isArray(mockUsers)).toBe(true);
    expect(mockUsers.length).toBeGreaterThan(0);
    expect(mockUsers[0]).toHaveProperty('username');
  });
  
  test('mockDecks contains expected test data', () => {
    expect(mockDecks).toBeDefined();
    expect(Array.isArray(mockDecks)).toBe(true);
    expect(mockDecks.length).toBeGreaterThan(0);
    expect(mockDecks[0]).toHaveProperty('name');
  });
  
  test('mockFlashcards contains expected test data', () => {
    expect(mockFlashcards).toBeDefined();
    expect(Array.isArray(mockFlashcards)).toBe(true);
    expect(mockFlashcards.length).toBeGreaterThan(0);
    expect(mockFlashcards[0]).toHaveProperty('word');
  });
});
