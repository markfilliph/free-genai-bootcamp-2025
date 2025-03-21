import { mockFetch } from '../mocks/api-mock.js';
import { apiFetch } from '../../lib/api.js';

describe('API Integration Tests', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
    // Replace the global fetch with our mock implementation
    global.fetch = jest.fn().mockImplementation(mockFetch);
  });

  test('login with valid credentials returns token', async () => {
    const response = await apiFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        email: 'test@example.com',
        password: 'password123'
      })
    });
    
    expect(response).toEqual({
      token: 'mock-token-123',
      user_id: '1'
    });
  });

  test('login with invalid credentials throws error', async () => {
    await expect(apiFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        email: 'test@example.com',
        password: 'wrong-password'
      })
    })).rejects.toThrow('API Error (401): Invalid credentials');
  });

  test('register with valid data creates user', async () => {
    const response = await apiFetch('/auth/register', {
      method: 'POST',
      body: JSON.stringify({
        username: 'newuser',
        email: 'newuser@example.com',
        password: 'password123'
      })
    });
    
    expect(response).toEqual({
      user_id: '2',
      username: 'newuser',
      email: 'newuser@example.com'
    });
  });

  test('get decks returns list of decks', async () => {
    const response = await apiFetch('/decks');
    expect(Array.isArray(response)).toBe(true);
    expect(response.length).toBe(2);
    expect(response[0].name).toBe('Spanish Basics');
    expect(response[1].name).toBe('Verb Conjugations');
  });

  test('create deck with valid data returns new deck', async () => {
    const response = await apiFetch('/decks', {
      method: 'POST',
      body: JSON.stringify({
        name: 'New Test Deck'
      })
    });
    
    expect(response).toEqual({
      id: '3',
      name: 'New Test Deck',
      user_id: '1',
      created_at: expect.any(String)
    });
  });

  test('create flashcard with valid data returns new flashcard', async () => {
    const response = await apiFetch('/flashcards', {
      method: 'POST',
      body: JSON.stringify({
        word: 'gracias',
        example_sentence: 'Muchas gracias por tu ayuda.',
        translation: 'Thank you very much for your help.',
        deck_id: '1'
      })
    });
    
    expect(response).toEqual({
      id: '4',
      word: 'gracias',
      example_sentence: 'Muchas gracias por tu ayuda.',
      translation: 'Thank you very much for your help.',
      deck_id: '1',
      created_at: expect.any(String)
    });
  });

  test('generate content for word returns example sentences', async () => {
    const response = await apiFetch('/generate', {
      method: 'POST',
      body: JSON.stringify({
        word: 'casa',
        is_verb: false
      })
    });
    
    expect(response).toEqual({
      example_sentences: [
        'This is an example sentence with the word "casa".',
        'Here\'s another example using "casa" in context.'
      ],
      conjugations: null,
      cultural_note: 'Cultural note about "casa".',
    });
  });

  test('generate content for verb returns conjugations', async () => {
    const response = await apiFetch('/generate', {
      method: 'POST',
      body: JSON.stringify({
        word: 'comer',
        is_verb: true
      })
    });
    
    expect(response).toEqual({
      example_sentences: [
        'This is an example sentence with the word "comer".',
        'Here\'s another example using "comer" in context.'
      ],
      conjugations: 'Mock conjugations for verb',
      cultural_note: 'Cultural note about "comer".',
    });
  });
});
