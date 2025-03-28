/**
 * Browser API Mocks
 * 
 * This file contains mocks for browser APIs that are used in the application
 * but not available in the Jest testing environment.
 */

// Mock MutationObserver
global.MutationObserver = class {
  constructor(callback) {}
  disconnect() {}
  observe(element, initObject) {}
};

// Mock document methods that might be used in components
document.createElement = document.createElement || jest.fn(() => ({
  setAttribute: jest.fn(),
  appendChild: jest.fn(),
  classList: {
    add: jest.fn(),
    remove: jest.fn(),
    contains: jest.fn(() => true)
  },
  style: {},
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  querySelector: jest.fn(),
  querySelectorAll: jest.fn(() => [])
}));

// Add custom Jest DOM matchers
expect.extend({
  toBeInTheDocument: (received) => {
    const pass = !!received;
    return {
      pass,
      message: () => pass
        ? `Expected element not to be in the document`
        : `Expected element to be in the document`
    };
  },
  toHaveAttribute: (received, name, value) => {
    const hasAttribute = received && 
      received.hasAttribute && 
      received.hasAttribute(name);
    
    const attributeMatches = value === undefined || 
      received.getAttribute(name) === value;
    
    const pass = hasAttribute && attributeMatches;
    
    return {
      pass,
      message: () => pass
        ? `Expected element not to have attribute "${name}"`
        : `Expected element to have attribute "${name}"`
    };
  },
  toHaveClass: (received, className) => {
    const pass = received && 
      received.classList && 
      received.classList.contains(className);
    
    return {
      pass,
      message: () => pass
        ? `Expected element not to have class "${className}"`
        : `Expected element to have class "${className}"`
    };
  },
  toHaveValue: (received, value) => {
    const pass = received && 
      received.value === value;
    
    return {
      pass,
      message: () => pass
        ? `Expected element not to have value "${value}"`
        : `Expected element to have value "${value}"`
    };
  }
});
