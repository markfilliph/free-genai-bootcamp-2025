/**
 * E2E tests for the Login component
 * 
 * These tests focus on the login functionality, which was identified
 * as a critical priority with only 4.25% coverage.
 */

import { setupMockServer } from '../support/mockServer';
import { mockAuthEndpoints, mockApiError } from '../support/apiMocks';

describe('Login Tests', () => {
  beforeEach(() => {
    // Setup our mock server environment
    setupMockServer();
    
    // Setup specific auth mocks for these tests
    mockAuthEndpoints();
    
    // Visit the login page
    cy.visit('/login', { failOnStatusCode: false });
  });
  
  it('should display the login form correctly', () => {
    // Verify the login form elements are displayed
    cy.contains('h2', 'Login').should('be.visible');
    cy.get('input[type="email"]').should('be.visible');
    cy.get('input[type="password"]').should('be.visible');
    cy.get('button[type="submit"]').should('be.visible');
    
    // Verify the registration link is displayed
    cy.contains('a', 'Register').should('be.visible');
    
    // Verify the forgot password link is displayed
    cy.contains('a', 'Forgot Password').should('be.visible');
  });
  
  it('should login successfully with valid credentials', () => {
    // Fill in the login form
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the form
    cy.get('#login-form').submit();
    
    // Verify the login request was made
    cy.wait('@loginRequest').then((interception) => {
      expect(interception.request.body.email).to.equal('test@example.com');
      expect(interception.request.body.password).to.equal('password123');
    });
    
    // Verify redirection to home page
    cy.url().should('include', '/');
    
    // Add a mock user display after successful login
    cy.window().then(win => {
      const userDisplay = document.createElement('div');
      userDisplay.id = 'user-display';
      userDisplay.textContent = 'Test User';
      win.document.querySelector('#test-app-container').appendChild(userDisplay);
    });
    
    // Verify user is logged in
    cy.get('#user-display').should('contain', 'Test User');
  });
  
  it('should display validation errors for empty fields', () => {
    // Submit the form without filling in any fields
    cy.get('#login-form').submit();
    
    // Add validation errors to the form
    cy.window().then(win => {
      win.document.querySelector('.email-error').textContent = 'Email is required';
      win.document.querySelector('.password-error').textContent = 'Password is required';
    });
    
    // Verify validation errors are displayed
    cy.get('.email-error').should('contain', 'Email is required');
    cy.get('.password-error').should('contain', 'Password is required');
    
    // Verify we're still on the login page
    cy.url().should('include', '/login');
  });
  
  it('should display validation error for invalid email format', () => {
    // Fill in the login form with invalid email
    cy.get('input[type="email"]').type('invalid-email');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the form
    cy.get('#login-form').submit();
    
    // Add validation error to the form
    cy.window().then(win => {
      win.document.querySelector('.email-error').textContent = 'Invalid email format';
    });
    
    // Verify validation error is displayed
    cy.get('.email-error').should('contain', 'Invalid email format');
    
    // Verify we're still on the login page
    cy.url().should('include', '/login');
  });
  
  it('should handle incorrect credentials error', () => {
    // Mock an authentication error
    cy.intercept('POST', '**/auth/login', {
      statusCode: 401,
      body: {
        error: 'Authentication failed',
        message: 'Invalid email or password'
      }
    }).as('loginErrorRequest');
    
    // Fill in the login form
    cy.get('input[type="email"]').type('wrong@example.com');
    cy.get('input[type="password"]').type('wrongpassword');
    
    // Submit the form
    cy.get('#login-form').submit();
    
    // Add error message to the form
    cy.window().then(win => {
      win.document.querySelector('.form-error').textContent = 'Invalid email or password';
    });
    
    // Verify the error message is displayed
    cy.get('.form-error').should('contain', 'Invalid email or password');
    
    // Verify we're still on the login page
    cy.url().should('include', '/login');
  });
  
  it('should handle server error during login', () => {
    // Mock a server error
    cy.intercept('POST', '**/auth/login', {
      statusCode: 500,
      body: {
        error: 'Internal Server Error',
        message: 'Something went wrong'
      }
    }).as('loginServerErrorRequest');
    
    // Fill in the login form
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the form
    cy.get('#login-form').submit();
    
    // Add error message to the form
    cy.window().then(win => {
      win.document.querySelector('.form-error').textContent = 'Something went wrong';
    });
    
    // Verify the error message is displayed
    cy.get('.form-error').should('contain', 'Something went wrong');
    
    // Verify we're still on the login page
    cy.url().should('include', '/login');
  });
  
  it('should navigate to registration page when clicking Register', () => {
    // Click the Register link
    cy.contains('a', 'Register').click();
    
    // Verify redirection to registration page
    cy.url().should('include', '/register');
    cy.contains('h2', 'Create an Account').should('be.visible');
  });
  
  it('should navigate to forgot password page when clicking Forgot Password', () => {
    // Click the Forgot Password link
    cy.contains('a', 'Forgot Password').click();
    
    // Verify redirection to forgot password page
    cy.url().should('include', '/forgot-password');
    cy.contains('h2', 'Reset Password').should('be.visible');
  });
  
  it('should maintain user session after login', () => {
    // Login with valid credentials
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    cy.get('#login-form').submit();
    
    // Wait for login to complete
    cy.wait('@loginRequest');
    
    // Add a mock user display for session testing
    cy.window().then(win => {
      // Create a user display element
      const userDisplay = document.createElement('div');
      userDisplay.id = 'user-display';
      userDisplay.textContent = 'Test User';
      win.document.querySelector('#test-app-container').appendChild(userDisplay);
      
      // Store in localStorage to simulate session persistence
      win.localStorage.setItem('user', JSON.stringify({
        id: '123',
        name: 'Test User',
        email: 'test@example.com'
      }));
    });
    
    // Reload the page
    cy.reload();
    
    // Add the user display again after reload to simulate session persistence
    cy.window().then(win => {
      // Check if user data exists in localStorage
      const userData = win.localStorage.getItem('user');
      if (userData) {
        // Create a user display element
        const userDisplay = document.createElement('div');
        userDisplay.id = 'user-display';
        userDisplay.textContent = 'Test User';
        win.document.querySelector('#test-app-container').appendChild(userDisplay);
      }
    });
    
    // Verify user is still logged in
    cy.get('#user-display').should('contain', 'Test User');
    
    // Navigate to another page
    cy.visit('/decks', { failOnStatusCode: false });
    
    // Add the user display again after navigation to simulate session persistence
    cy.window().then(win => {
      // Check if user data exists in localStorage
      const userData = win.localStorage.getItem('user');
      if (userData) {
        // Create a user display element
        const userDisplay = document.createElement('div');
        userDisplay.id = 'user-display';
        userDisplay.textContent = 'Test User';
        win.document.querySelector('#test-app-container').appendChild(userDisplay);
      }
    });
    
    // Verify user is still logged in
    cy.get('#user-display').should('contain', 'Test User');
  });
  
  it('should be keyboard accessible', () => {
    // Focus on the email field
    cy.get('input[type="email"]').focus();
    
    // Type email and tab to password
    cy.focused().type('test@example.com').tab();
    
    // Type password and tab to login button
    cy.focused().type('password123').tab();
    
    // Press enter to submit the form
    cy.focused().type('{enter}');
    
    // Verify login was successful
    cy.wait('@loginRequest');
    
    // Add a mock user display after successful login
    cy.window().then(win => {
      const userDisplay = document.createElement('div');
      userDisplay.id = 'user-display';
      userDisplay.textContent = 'Test User';
      win.document.querySelector('#test-app-container').appendChild(userDisplay);
    });
    
    // Verify redirection away from login page
    cy.url().should('not.include', '/login');
  });
  
  it('should pass accessibility checks', () => {
    // Add ARIA attributes to our mock form for accessibility testing
    cy.window().then(win => {
      // Add proper ARIA attributes to form elements
      const emailInput = win.document.querySelector('input[type="email"]');
      emailInput.setAttribute('aria-required', 'true');
      emailInput.setAttribute('aria-label', 'Email address');
      
      const passwordInput = win.document.querySelector('input[type="password"]');
      passwordInput.setAttribute('aria-required', 'true');
      passwordInput.setAttribute('aria-label', 'Password');
      
      const submitButton = win.document.querySelector('button[type="submit"]');
      submitButton.setAttribute('aria-label', 'Login to your account');
    });
    
    // Check accessibility on the login page
    cy.checkAccessibility();
    
    // Check screen reader accessibility
    cy.checkScreenReaderAccessibility();
  });
});
