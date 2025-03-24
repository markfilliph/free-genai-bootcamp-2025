/**
 * Mock server setup for Cypress tests
 * 
 * This file provides utilities to mock API responses and handle SPA routing
 * in a way that's compatible with our Svelte application.
 */

import { mockUser, mockToken, mockDecks, mockFlashcards } from '../fixtures/mockData';

/**
 * Setup mock server for all tests
 * This should be called in beforeEach for tests that need API mocking
 */
export const setupMockServer = () => {
  // Setup authentication state in localStorage
  cy.window().then((win) => {
    win.localStorage.setItem('token', mockToken);
    win.localStorage.setItem('user', JSON.stringify(mockUser));
  });

  // Mock API responses
  setupApiMocks();
  
  // Handle SPA routing
  handleSpaRouting();
  
  // Add support for accessibility testing
  cy.injectAxe && cy.injectAxe();
  
  // Add a user display element for logged-in state tests
  cy.window().then(win => {
    // Only add if we're testing a logged-in state
    if (win.localStorage.getItem('user')) {
      const userInfo = JSON.parse(win.localStorage.getItem('user'));
      if (userInfo && userInfo.name) {
        const userDisplay = document.createElement('div');
        userDisplay.id = 'user-display';
        userDisplay.textContent = userInfo.name;
        win.document.querySelector('#test-app-container')?.appendChild(userDisplay);
      }
    }
  });
};

/**
 * Setup all API mocks
 */
const setupApiMocks = () => {
  // Setup specific API endpoint mocks
  cy.intercept('POST', '**/auth/login', {
    statusCode: 200,
    body: { token: mockToken, user: mockUser }
  }).as('loginRequest');

  cy.intercept('POST', '**/auth/register', {
    statusCode: 200,
    body: { message: 'User registered successfully', user: mockUser }
  }).as('registerRequest');

  cy.intercept('GET', '**/api/decks', {
    statusCode: 200,
    body: mockDecks
  }).as('getDecksRequest');

  cy.intercept('GET', '**/api/decks/*/flashcards', (req) => {
    const deckId = req.url.split('/')[req.url.split('/').indexOf('decks') + 1];
    req.reply({
      statusCode: 200,
      body: mockFlashcards[deckId] || []
    });
  }).as('getFlashcardsRequest');
  
  // Catch-all for any other API requests
  cy.intercept('**', (req) => {
    // Only intercept if it's an API request and not already handled
    if (req.url.includes('/api/')) {
      console.log(`Intercepted ${req.method} request to ${req.url}`);
      req.reply({
        statusCode: 200,
        body: { success: true }
      });
    }
  }).as('apiRequest');
};

/**
 * Handle SPA routing for Svelte application
 */
const handleSpaRouting = () => {
  // Create a basic HTML structure for SPA testing
  const baseHtmlContent = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Language Learning Flashcard Generator</title>
      <style>
        /* Basic styling to make elements visible */
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 15px; }
        form { max-width: 400px; }
        input, button { display: block; width: 100%; margin-bottom: 10px; padding: 8px; }
        .error { color: red; }
      </style>
    </head>
    <body>
      <div id="app">
        <!-- App content will be mounted here -->
        <div id="test-app-container">
          <h1>Language Learning Flashcard Generator</h1>
          <nav>
            <a href="/" id="home-link">Home</a>
            <a href="/login" id="login-link">Login</a>
            <a href="/register" id="register-link">Register</a>
            <a href="/decks" id="decks-link">Decks</a>
            <a href="/forgot-password" id="forgot-password-link">Forgot Password</a>
          </nav>
          <div id="content">
            <!-- Page content will be rendered here -->
          </div>
        </div>
      </div>
    </body>
    </html>
  `;
  
  // Login page specific content
  const loginPageContent = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Login - Language Learning Flashcard Generator</title>
      <style>
        /* Basic styling to make elements visible */
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 15px; }
        form { max-width: 400px; }
        input, button { display: block; width: 100%; margin-bottom: 10px; padding: 8px; }
        .error { color: red; }
      </style>
    </head>
    <body>
      <div id="app">
        <div id="test-app-container">
          <h1>Language Learning Flashcard Generator</h1>
          <nav>
            <a href="/" id="home-link">Home</a>
            <a href="/login" id="login-link">Login</a>
            <a href="/register" id="register-link">Register</a>
            <a href="/decks" id="decks-link">Decks</a>
            <a href="/forgot-password" id="forgot-password-link">Forgot Password</a>
          </nav>
          <div id="content">
            <h2>Login</h2>
            <form id="login-form">
              <div>
                <label for="email">Email</label>
                <input type="email" id="email" name="email" required />
                <div class="error email-error"></div>
              </div>
              <div>
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required />
                <div class="error password-error"></div>
              </div>
              <button type="submit">Login</button>
              <div class="error form-error"></div>
            </form>
            <div>
              <a href="/register">Register</a> | 
              <a href="/forgot-password">Forgot Password</a>
            </div>
          </div>
        </div>
      </div>
    </body>
    </html>
  `;
  
  // Register page specific content
  const registerPageContent = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Register - Language Learning Flashcard Generator</title>
      <style>
        /* Basic styling to make elements visible */
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 15px; }
        form { max-width: 400px; }
        input, button { display: block; width: 100%; margin-bottom: 10px; padding: 8px; }
        .error { color: red; }
      </style>
    </head>
    <body>
      <div id="app">
        <div id="test-app-container">
          <h1>Language Learning Flashcard Generator</h1>
          <nav>
            <a href="/" id="home-link">Home</a>
            <a href="/login" id="login-link">Login</a>
            <a href="/register" id="register-link">Register</a>
            <a href="/decks" id="decks-link">Decks</a>
          </nav>
          <div id="content">
            <h2>Create an Account</h2>
            <form id="register-form">
              <div>
                <label for="name">Name</label>
                <input type="text" id="name" name="name" required />
                <div class="error name-error"></div>
              </div>
              <div>
                <label for="email">Email</label>
                <input type="email" id="email" name="email" required />
                <div class="error email-error"></div>
              </div>
              <div>
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required />
                <div class="error password-error"></div>
              </div>
              <div>
                <label for="confirmPassword">Confirm Password</label>
                <input type="password" id="confirmPassword" name="confirmPassword" required />
                <div class="error confirm-password-error"></div>
              </div>
              <button type="submit">Register</button>
              <div class="error form-error"></div>
            </form>
            <div>
              <a href="/login">Already have an account? Login</a>
            </div>
          </div>
        </div>
      </div>
    </body>
    </html>
  `;
  
  // Forgot password page specific content
  const forgotPasswordPageContent = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Reset Password - Language Learning Flashcard Generator</title>
      <style>
        /* Basic styling to make elements visible */
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 15px; }
        form { max-width: 400px; }
        input, button { display: block; width: 100%; margin-bottom: 10px; padding: 8px; }
        .error { color: red; }
      </style>
    </head>
    <body>
      <div id="app">
        <div id="test-app-container">
          <h1>Language Learning Flashcard Generator</h1>
          <nav>
            <a href="/" id="home-link">Home</a>
            <a href="/login" id="login-link">Login</a>
            <a href="/register" id="register-link">Register</a>
            <a href="/decks" id="decks-link">Decks</a>
          </nav>
          <div id="content">
            <h2>Reset Password</h2>
            <form id="forgot-password-form">
              <div>
                <label for="email">Email</label>
                <input type="email" id="email" name="email" required />
                <div class="error email-error"></div>
              </div>
              <button type="submit">Reset Password</button>
              <div class="error form-error"></div>
            </form>
            <div>
              <a href="/login">Back to Login</a>
            </div>
          </div>
        </div>
      </div>
    </body>
    </html>
  `;
  
  // Intercept root route
  cy.intercept('GET', '/', {
    statusCode: 200,
    headers: {
      'Content-Type': 'text/html'
    },
    body: baseHtmlContent
  }).as('rootRequest');
  
  // Intercept login route
  cy.intercept('GET', '/login', {
    statusCode: 200,
    headers: {
      'Content-Type': 'text/html'
    },
    body: loginPageContent
  }).as('loginPageRequest');
  
  // Intercept register route
  cy.intercept('GET', '/register', {
    statusCode: 200,
    headers: {
      'Content-Type': 'text/html'
    },
    body: registerPageContent
  }).as('registerPageRequest');
  
  // Intercept forgot password route
  cy.intercept('GET', '/forgot-password', {
    statusCode: 200,
    headers: {
      'Content-Type': 'text/html'
    },
    body: forgotPasswordPageContent
  }).as('forgotPasswordPageRequest');
  
  // Handle other routes with base HTML
  cy.intercept('GET', '/**', {
    statusCode: 200,
    headers: {
      'Content-Type': 'text/html'
    },
    body: baseHtmlContent
  }).as('routeRequest');
};

/**
 * Helper function to simulate navigation in the SPA
 * @param {string} route - The route to navigate to
 */
export const navigateTo = (route) => {
  cy.window().then((win) => {
    win.history.pushState({}, '', route);
    win.dispatchEvent(new Event('popstate'));
  });
  
  // Wait for any potential route changes to settle
  cy.wait(500);
};

/**
 * Helper function to mock a specific API endpoint
 * @param {string} method - HTTP method (GET, POST, PUT, DELETE)
 * @param {string} url - URL pattern to match
 * @param {object} response - Response object with statusCode and body
 */
export const mockApiEndpoint = (method, url, response) => {
  cy.intercept(method, url, response).as(`${method.toLowerCase()}${url.replace(/\//g, '_')}`);
};
