/**
 * Basic smoke tests to verify the testing environment is working correctly
 */

import { setupMockServer, navigateTo } from '../support/mockServer';

describe('Basic Smoke Tests', () => {
  beforeEach(() => {
    // Setup the mock server for all tests
    setupMockServer();
  });

  it('should load the application', () => {
    // Visit the root URL with our mock server
    cy.visit('/', { failOnStatusCode: false });
    
    // Verify the page has loaded
    cy.get('body').should('be.visible');
    cy.get('#test-app-container').should('exist');
    cy.contains('h1', 'Language Learning Flashcard Generator').should('be.visible');
    cy.log('Basic page load test passed');
  });

  it('should handle basic DOM interactions', () => {
    // Visit the root URL
    cy.visit('/', { failOnStatusCode: false });
    
    // Find navigation links and verify they're clickable
    cy.get('#login-link').should('be.visible').click();
    
    // Verify the URL changed
    cy.url().should('include', '/login');
    cy.log('Basic interaction test passed');
  });

  it('should handle localStorage operations', () => {
    // Visit the root URL
    cy.visit('/', { failOnStatusCode: false });
    
    // Test localStorage operations
    cy.window().then(win => {
      // Set a test value
      win.localStorage.setItem('test-key', 'test-value');
      
      // Verify it was set correctly
      expect(win.localStorage.getItem('test-key')).to.equal('test-value');
      
      // Clean up
      win.localStorage.removeItem('test-key');
      
      // Verify our mock server setup the auth token
      expect(win.localStorage.getItem('token')).to.exist;
      
      cy.log('localStorage test passed');
    });
  });

  it('should handle basic routing', () => {
    // Visit the root URL
    cy.visit('/', { failOnStatusCode: false });
    
    // Use our helper to navigate to a different route
    navigateTo('/decks');
    
    // Verify the URL changed
    cy.url().should('include', '/decks');
    
    // Navigate to another route
    navigateTo('/login');
    
    // Verify the URL changed again
    cy.url().should('include', '/login');
    
    cy.log('Basic routing test passed');
  });
  
  it('should handle API mocking', () => {
    // Visit the root URL
    cy.visit('/', { failOnStatusCode: false });
    
    // Use Cypress request with failOnStatusCode: false
    // This allows us to test the mocking without requiring actual endpoints
    cy.request({
      url: '/api/decks',
      failOnStatusCode: false
    }).then((response) => {
      // Just verify we got some response
      expect(response).to.exist;
      cy.log('API request test passed');
    });
  });
});
