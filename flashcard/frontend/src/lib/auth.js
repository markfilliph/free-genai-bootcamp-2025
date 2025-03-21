import { writable } from 'svelte/store';

// Create a store for authentication state
export const authStore = writable({
  isAuthenticated: false,
  user: null,
  token: null
});

// Initialize auth state from localStorage
export function initAuth() {
  const token = localStorage.getItem('auth_token');
  const user = localStorage.getItem('auth_user');
  
  if (token && user) {
    authStore.set({
      isAuthenticated: true,
      token,
      user: JSON.parse(user)
    });
    return true;
  }
  
  return false;
}

// Login function
export function login(token, user) {
  // Save to localStorage
  localStorage.setItem('auth_token', token);
  localStorage.setItem('auth_user', JSON.stringify(user));
  
  // Update the store
  authStore.set({
    isAuthenticated: true,
    token,
    user
  });
}

// Logout function
export function logout() {
  // Clear localStorage
  localStorage.removeItem('auth_token');
  localStorage.removeItem('auth_user');
  
  // Update the store
  authStore.set({
    isAuthenticated: false,
    token: null,
    user: null
  });
}

// Get auth token
export function getToken() {
  let token;
  const unsubscribe = authStore.subscribe(state => {
    token = state.token;
  });
  unsubscribe();
  
  return token;
}
