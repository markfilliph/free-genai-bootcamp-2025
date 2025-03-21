import { authStore, initAuth, login, logout } from '../../lib/auth.js';
import { get } from 'svelte/store';

describe('Auth Module', () => {
  // Mock localStorage
  const mockStorage = {};
  const mockLocalStorage = {
    getItem: jest.fn(key => mockStorage[key] || null),
    setItem: jest.fn((key, value) => { mockStorage[key] = value; }),
    removeItem: jest.fn(key => { delete mockStorage[key]; }),
    clear: jest.fn(() => { Object.keys(mockStorage).forEach(key => delete mockStorage[key]); }),
    length: 0
  };
  
  // Replace global localStorage with mock
  global.localStorage = mockLocalStorage;

  // Clear localStorage before each test
  beforeEach(() => {
    mockLocalStorage.clear();
    // Reset the store to initial state
    authStore.set({
      isAuthenticated: false,
      user: null,
      token: null
    });
  });

  test('authStore initializes with unauthenticated state', () => {
    const state = get(authStore);
    
    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
    expect(state.token).toBeNull();
  });

  test('login updates store and localStorage', () => {
    const testUser = { id: '1', username: 'testuser' };
    const testToken = 'test-token-123';
    
    login(testToken, testUser);
    
    // Check store state
    const state = get(authStore);
    expect(state.isAuthenticated).toBe(true);
    expect(state.token).toBe(testToken);
    expect(state.user).toEqual(testUser);
    
    // Check localStorage
    expect(localStorage.getItem('auth_token')).toBe(testToken);
    expect(localStorage.getItem('auth_user')).toBe(JSON.stringify(testUser));
  });

  test('logout clears store and localStorage', () => {
    // First login
    const testUser = { id: '1', username: 'testuser' };
    const testToken = 'test-token-123';
    login(testToken, testUser);
    
    // Then logout
    logout();
    
    // Check store state
    const state = get(authStore);
    expect(state.isAuthenticated).toBe(false);
    expect(state.token).toBeNull();
    expect(state.user).toBeNull();
    
    // Check localStorage
    expect(localStorage.getItem('auth_token')).toBeNull();
    expect(localStorage.getItem('auth_user')).toBeNull();
  });

  test('initAuth restores state from localStorage', () => {
    // Setup localStorage with auth data
    const testUser = { id: '1', username: 'testuser' };
    const testToken = 'test-token-123';
    localStorage.setItem('auth_token', testToken);
    localStorage.setItem('auth_user', JSON.stringify(testUser));
    
    // Call initAuth
    const result = initAuth();
    
    // Check return value
    expect(result).toBe(true);
    
    // Check store state
    const state = get(authStore);
    expect(state.isAuthenticated).toBe(true);
    expect(state.token).toBe(testToken);
    expect(state.user).toEqual(testUser);
  });

  test('initAuth returns false when no auth data in localStorage', () => {
    // Call initAuth with empty localStorage
    const result = initAuth();
    
    // Check return value
    expect(result).toBe(false);
    
    // Check store state remains unchanged
    const state = get(authStore);
    expect(state.isAuthenticated).toBe(false);
    expect(state.token).toBeNull();
    expect(state.user).toBeNull();
  });
});
