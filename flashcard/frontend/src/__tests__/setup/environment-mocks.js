/**
 * Environment Mocks
 * 
 * This file contains mocks for environment variables and Vite-specific features
 * that are used in the application but need to be mocked in the Jest environment.
 */

// Mock Vite's import.meta.env
global.import = global.import || {};
global.import.meta = global.import.meta || {};
global.import.meta.env = {
  VITE_API_URL: 'http://localhost:8000',
  MODE: 'test',
  DEV: true,
  PROD: false,
  SSR: false
};
