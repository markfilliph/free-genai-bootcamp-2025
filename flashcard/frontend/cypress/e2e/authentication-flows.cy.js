/**
 * End-to-end tests for authentication flows
 * 
 * This test file covers all authentication-related scenarios:
 * - Login
 * - Registration
 * - Password reset
 * - Session persistence
 */

describe('Authentication Flows', () => {
  beforeEach(() => {
    // Clear cookies and local storage before each test
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  it('should allow a user to log in successfully', () => {
    // Use the custom login command
    cy.login('test@example.com', 'password123');
    
    // Verify redirection to home page
    cy.url().should('include', '/');
    
    // Verify user is logged in
    cy.contains('Test User').should('be.visible');
  });

  it('should handle login failures gracefully', () => {
    // Mock a failed login response
    cy.intercept('POST', '**/auth/login', {
      statusCode: 401,
      body: {
        error: 'Invalid credentials'
      }
    }).as('failedLoginRequest');
    
    cy.visit('/login');
    cy.get('input[type="email"]').type('wrong@example.com');
    cy.get('input[type="password"]').type('wrongpassword');
    cy.get('form').submit();
    cy.wait('@failedLoginRequest');
    
    // Verify error message is displayed
    cy.contains('Invalid credentials').should('be.visible');
    
    // Verify we're still on the login page
    cy.url().should('include', '/login');
  });

  it('should allow a user to register a new account', () => {
    // Use the custom register command
    cy.register('New User', 'newuser@example.com', 'password123');
    
    // Verify success message
    cy.contains('User registered successfully').should('be.visible');
    
    // Verify redirection to login page
    cy.url().should('include', '/login');
  });

  it('should handle registration validation errors', () => {
    // Mock a validation error response
    cy.intercept('POST', '**/auth/register', {
      statusCode: 400,
      body: {
        error: 'Email already in use'
      }
    }).as('failedRegisterRequest');
    
    cy.visit('/register');
    cy.get('input[name="name"]').type('Test User');
    cy.get('input[type="email"]').type('existing@example.com');
    cy.get('input[type="password"]').type('password123');
    cy.get('input[type="password"][name="confirmPassword"]').type('password123');
    cy.get('form').submit();
    cy.wait('@failedRegisterRequest');
    
    // Verify error message is displayed
    cy.contains('Email already in use').should('be.visible');
  });

  it('should allow a user to request a password reset', () => {
    // Use the custom requestPasswordReset command
    cy.requestPasswordReset('test@example.com');
    
    // Verify success message
    cy.contains('Password reset email sent').should('be.visible');
  });

  it('should maintain user session across page reloads', () => {
    // Log in first
    cy.login('test@example.com', 'password123');
    
    // Verify session persistence
    cy.verifySessionPersistence();
  });

  it('should redirect to login page when accessing protected routes without authentication', () => {
    // Mock an invalid token response
    cy.intercept('GET', '**/auth/validate', {
      statusCode: 401,
      body: {
        valid: false,
        error: 'Invalid token'
      }
    }).as('invalidTokenRequest');
    
    // Try to access a protected route
    cy.visit('/decks');
    
    // Verify redirection to login page
    cy.url().should('include', '/login');
  });
});
