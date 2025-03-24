const { defineConfig } = require('cypress');

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:5173',
    setupNodeEvents(on, config) {
      // Add accessibility testing support
      on('task', {
        log(message) {
          console.log(message);
          return null;
        },
        table(message) {
          console.table(message);
          return null;
        }
      });
      return config;
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
    defaultCommandTimeout: 15000,
    // Increase timeout for page loads (important for SPA routing)
    pageLoadTimeout: 30000,
    // Increase timeout for requests (important for API mocking)
    requestTimeout: 10000,
    // Automatically handle uncaught exceptions
    experimentalRunAllSpecs: true,
    // Preserve cookies between tests for session testing
    testIsolation: false
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
