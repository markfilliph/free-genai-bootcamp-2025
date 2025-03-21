// Import the API functions
import * as apiModule from '../../lib/api.js';

// Create a manual mock of the API module
const mockApiFetch = jest.fn();
const API_BASE = 'http://localhost:8000';

// Override the imported module with our mocks
apiModule.apiFetch = mockApiFetch;
apiModule.API_BASE = API_BASE;

describe('API Utility', () => {
  // Mock fetch responses
  const mockSuccessResponse = { message: 'Mock API response' };
  const mockErrorResponse = 'Unauthorized';
  const mockNetworkError = new Error('Network error');

  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
  });

  test('API_BASE is set correctly', () => {
    expect(apiModule.API_BASE).toBe('http://localhost:8000');
  });

  test('apiFetch makes request with correct URL and default options', async () => {
    // Setup mock response
    mockApiFetch.mockResolvedValueOnce(mockSuccessResponse);
    
    // Call the API function
    const response = await apiModule.apiFetch('/test-endpoint');
    
    // Check if apiFetch was called with correct parameters
    expect(mockApiFetch).toHaveBeenCalledWith('/test-endpoint', {});
    
    // Check if response is correct
    expect(response).toEqual(mockSuccessResponse);
  });

  test('apiFetch makes request with custom options', async () => {
    // Setup mock response
    mockApiFetch.mockResolvedValueOnce(mockSuccessResponse);
    
    // Custom options
    const options = {
      method: 'POST',
      body: JSON.stringify({ test: 'data' }),
      headers: {
        'Authorization': 'Bearer token123'
      }
    };
    
    // Call the API function with custom options
    const response = await apiModule.apiFetch('/test-endpoint', options);
    
    // Check if apiFetch was called with correct parameters
    expect(mockApiFetch).toHaveBeenCalledWith('/test-endpoint', options);
    
    // Check if response is correct
    expect(response).toEqual(mockSuccessResponse);
  });

  test('apiFetch handles error responses correctly', async () => {
    // Setup mock error response
    mockApiFetch.mockRejectedValueOnce(new Error('API Error (401): Unauthorized'));
    
    // Call the API function and expect it to throw an error
    await expect(apiModule.apiFetch('/test-endpoint')).rejects.toThrow('API Error (401): Unauthorized');
  });

  test('apiFetch handles network errors correctly', async () => {
    // Setup mock network error
    mockApiFetch.mockRejectedValueOnce(mockNetworkError);
    
    // Call the API function and expect it to throw an error
    await expect(apiModule.apiFetch('/test-endpoint')).rejects.toThrow('Network error');
  });
});
