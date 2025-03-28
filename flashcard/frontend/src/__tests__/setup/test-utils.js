/**
 * Test Utilities
 * 
 * This file contains helper functions and utilities for testing Svelte components
 * in a consistent way across the application.
 */

/**
 * Creates a mock Svelte component with customizable render output
 * 
 * @param {string} name - The name of the component to mock
 * @param {string} html - Optional HTML to return when the component is rendered
 * @returns {Object} A mock Svelte component object
 */
export function mockComponent(name, html = '') {
  return {
    $$render: () => html || `<div data-testid="mock-${name}">${name} Component</div>`,
    render: (props) => html || `<div data-testid="mock-${name}">${name} Component</div>`
  };
}

/**
 * Creates HTML content for component mocks
 * This is useful for creating more complex mock HTML structures
 * 
 * @param {string} content - The HTML content to return
 * @returns {string} The HTML content
 */
export function mockHtml(content) {
  return content;
}

/**
 * Creates a mock event object for testing event handlers
 * 
 * @param {string} type - The event type (e.g., 'click', 'input')
 * @param {Object} props - Additional properties to add to the event object
 * @returns {Object} A mock event object
 */
export function mockEvent(type = 'click', props = {}) {
  return {
    type,
    preventDefault: jest.fn(),
    stopPropagation: jest.fn(),
    target: {},
    currentTarget: {},
    ...props
  };
}

/**
 * Creates a mock form submission event
 * 
 * @param {Object} formData - Form data to include in the event
 * @returns {Object} A mock form submission event
 */
export function mockFormSubmit(formData = {}) {
  const formElements = {};
  
  // Create form elements for each form data field
  Object.entries(formData).forEach(([name, value]) => {
    formElements[name] = { value };
  });
  
  return mockEvent('submit', {
    preventDefault: jest.fn(),
    target: {
      elements: formElements
    }
  });
}
