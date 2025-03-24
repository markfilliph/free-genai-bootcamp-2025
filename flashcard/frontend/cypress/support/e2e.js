// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

// Import cypress-axe for accessibility testing
import 'cypress-axe'

// Import Testing Library Cypress commands
import '@testing-library/cypress/add-commands'

// Configure Cypress for SPA testing
Cypress.on('uncaught:exception', (err, runnable) => {
  // Returning false here prevents Cypress from failing the test when
  // uncaught exceptions occur in the application code
  // This is often necessary for SPAs using client-side routing
  return false
})

// Add support for tab navigation in tests
import 'cypress-plugin-tab'

// Set up default viewport size for consistent testing
beforeEach(() => {
  // Use a standard desktop viewport
  cy.viewport(1280, 720)
})

// Add tab command for keyboard navigation testing
Cypress.Commands.add('tab', { prevSubject: 'optional' }, (subject) => {
  const tabKey = { key: 'Tab', code: 'Tab', which: 9 }
  if (subject) {
    cy.wrap(subject).trigger('keydown', tabKey)
  } else {
    cy.focused().trigger('keydown', tabKey)
  }
  return cy.focused()
})

// Alternatively you can use CommonJS syntax:
// require('./commands')

// Prevent uncaught exceptions from failing tests
Cypress.on('uncaught:exception', (err, runnable) => {
  // returning false here prevents Cypress from
  // failing the test on uncaught exceptions
  return false
})
