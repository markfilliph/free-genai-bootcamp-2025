const { defineConfig } = require('cypress');

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:5173',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    // Retry failed tests to reduce flakiness
    retries: {
      runMode: 2,
      openMode: 0
    },
    // Don't fail on status codes like 404
    // This is useful for tests that expect 404s or other non-200 responses
    failOnStatusCode: false,
    // Increase timeout for slow operations
    defaultCommandTimeout: 10000,
    // Automatically handle uncaught exceptions
    experimentalRunAllSpecs: true
  },
  component: {
    devServer: {
      framework: 'svelte',
      bundler: 'vite',
    },
  },
  // Enable video recording for CI environments only
  video: process.env.CI ? true : false,
  // Configure screenshots to only capture failures
  screenshotOnRunFailure: true,
  trashAssetsBeforeRuns: true
});
