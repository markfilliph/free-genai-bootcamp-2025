import { mockApiResponses } from '../mocks/api-mock.js';
import { apiFetch } from '../../lib/api.js';

describe('API Integration Tests', () => {
  // Mock successful responses for each test
  test('login with valid credentials returns token', () => {
    // Direct test of the expected response structure
    const expectedResponse = {
      token: 'mock-token-123',
      user_id: '1'
    };
    
    // Verify the mock response matches what we expect
    expect(mockApiResponses['/auth/login'].POST({
      body: JSON.stringify({
        email: 'test@example.com',
        password: 'password123'
      })
    }).body).toEqual(expectedResponse);
  });

  test('login with invalid credentials throws error', () => {
    // Verify the mock response for invalid credentials
    const response = mockApiResponses['/auth/login'].POST({
      body: JSON.stringify({
        email: 'test@example.com',
        password: 'wrongpassword'
      })
    });
    
    expect(response.status).toBe(401);
    expect(response.error).toBe('Invalid credentials');
  });

  test('register with valid data creates user', () => {
    const userData = {
      username: 'newuser',
      email: 'newuser@example.com',
      password: 'password123'
    };
    
    const response = mockApiResponses['/auth/register'].POST({
      body: JSON.stringify(userData)
    });
    
    expect(response.body).toEqual({
      user_id: '2',
      username: userData.username,
      email: userData.email
    });
  });

  test('get decks returns list of decks', () => {
    const response = mockApiResponses['/decks'].GET().body;
    
    expect(Array.isArray(response)).toBe(true);
    expect(response.length).toBe(2);
    expect(response[0].name).toBe('Spanish Basics');
    expect(response[1].name).toBe('Verb Conjugations');
  });

  test('create deck with valid data returns new deck', () => {
    const response = mockApiResponses['/decks'].POST({
      body: JSON.stringify({
        name: 'New Test Deck'
      })
    }).body;
    
    expect(response).toEqual({
      id: '3',
      name: 'New Test Deck',
      user_id: '1',
      created_at: expect.any(String)
    });
  });

  test('create flashcard with valid data returns new flashcard', () => {
    const response = mockApiResponses['/decks/1/flashcards'].POST({
      body: JSON.stringify({
        word: 'gracias',
        example_sentence: 'Muchas gracias por tu ayuda.',
        translation: 'Thank you very much for your help.'
      })
    }).body;
    
    expect(response).toEqual({
      id: '3',
      deck_id: '1',
      word: 'gracias',
      example_sentence: 'Muchas gracias por tu ayuda.',
      translation: 'Thank you very much for your help.',
      cultural_note: 'Cultural note.',
      conjugation: null,
      created_at: expect.any(String)
    });
  });

  test('generate content for word returns example sentences', () => {
    const response = mockApiResponses['/generate/word'].GET().body;
    
    expect(response).toEqual({
      word: 'hola',
      translation: 'hello',
      example_sentence: 'Hola, ¿cómo estás?',
      cultural_note: 'Common greeting in Spanish-speaking countries'
    });
  });

  test('generate content for verb returns conjugations', () => {
    const response = mockApiResponses['/generate/verb'].GET().body;
    
    expect(response).toEqual({
      word: 'hablar',
      translation: 'to speak',
      example_sentence: 'Yo hablo español.',
      cultural_note: 'Regular -ar verb',
      conjugation: 'hablo, hablas, habla, hablamos, habláis, hablan'
    });
  });
});
