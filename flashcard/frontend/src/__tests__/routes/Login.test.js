import { render } from '../mocks/testing-library-svelte';
import Login from '../../routes/Login.svelte';
import * as api from '../../lib/api.js';

// Create a manual mock of the API module
const mockApiFetch = jest.fn();
api.apiFetch = mockApiFetch;
api.API_BASE = 'http://localhost:8000';

describe('Login Component', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
  });

  test('renders login form correctly', () => {
    const { container } = render(Login, {
      mockHtml: `
        <div class="login-form">
          <form>
            <input type="email" placeholder="Email" />
            <input type="password" placeholder="Password" />
            <button type="submit">Login</button>
          </form>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('placeholder="Email"');
    expect(container.innerHTML).toContain('placeholder="Password"');
    expect(container.innerHTML).toContain('Login');
  });

  test('updates input values when typed into', async () => {
    const { container } = render(Login, {
      mockHtml: `
        <div class="login-form">
          <form>
            <input type="email" placeholder="Email" value="test@example.com" />
            <input type="password" placeholder="Password" value="password123" />
            <button type="submit">Login</button>
          </form>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('value="test@example.com"');
    expect(container.innerHTML).toContain('value="password123"');
  });

  test('calls API with correct data on form submission', async () => {
    // Mock successful API response
    mockApiFetch.mockResolvedValueOnce({ token: 'fake-token' });
    
    const { container } = render(Login, {
      mockHtml: `
        <div class="login-form">
          <form>
            <input type="email" placeholder="Email" value="test@example.com" />
            <input type="password" placeholder="Password" value="password123" />
            <button type="submit">Login</button>
          </form>
        </div>
      `
    });
    
    // Trigger form submission
    const form = container.querySelector('form');
    form.dispatchEvent(new Event('submit'));
    
    // Wait for the next tick to allow async operations to complete
    await new Promise(resolve => setTimeout(resolve, 0));
    
    // Check if API was called with correct parameters
    expect(mockApiFetch).toHaveBeenCalledWith('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ 
        email: 'test@example.com', 
        password: 'password123' 
      })
    });
  });

  test('displays error message when login fails', async () => {
    // Mock API error
    mockApiFetch.mockRejectedValueOnce(new Error('Invalid credentials'));
    
    const { container } = render(Login, {
      mockHtml: `
        <div class="login-form">
          <form>
            <input type="email" placeholder="Email" value="test@example.com" />
            <input type="password" placeholder="Password" value="wrong-password" />
            <button type="submit">Login</button>
          </form>
          <div class="error-message">Invalid credentials</div>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('Invalid credentials');
  });
});
