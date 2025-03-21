// Mock for the api.js module to avoid import.meta.env issues in tests
export const API_BASE = 'http://localhost:8000';

// Create a jest.fn() mock that can be spied on in tests
export const apiFetch = jest.fn(async (path, options = {}) => {
  try {
    const response = await fetch(`${API_BASE}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      credentials: 'include',
      ...options
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`API Error (${response.status}): ${error}`);
    }

    return await response.json();
  } catch (error) {
    if (error.message.startsWith('API Error')) {
      throw error;
    }
    throw new Error(`Network error: ${error.message}`);
  }
});
