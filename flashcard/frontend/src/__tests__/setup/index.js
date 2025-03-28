/**
 * Main Jest setup file
 * 
 * This is the primary setup file for Jest tests in the Flashcard application.
 * It imports and configures all necessary mocks and utilities for testing.
 */

// Import testing libraries
import '@testing-library/jest-dom/extend-expect';

// Import specific setup modules
import './browser-mocks';
import './storage-mocks';
import './fetch-mock';
import './component-mocks';
import './environment-mocks';

// Import custom test utilities
import { mockComponent, mockHtml } from './test-utils';

// Export utilities for use in tests
export {
  mockComponent,
  mockHtml
};
