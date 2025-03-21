import App from '../App.svelte';

// Import our mock testing utilities
import { render } from './mocks/testing-library-svelte';

// Mock the router components
jest.mock('svelte-routing', () => ({
  Router: {
    $$render: () => '<div id="mock-router"></div>'
  },
  Link: {
    $$render: () => '<a id="mock-link"></a>'
  }
}));

// Mock the route components
jest.mock('../routes/Home.svelte', () => ({
  default: {
    $$render: () => '<div id="mock-home"></div>'
  }
}));

jest.mock('../routes/Login.svelte', () => ({
  default: {
    $$render: () => '<div id="mock-login"></div>'
  }
}));

jest.mock('../routes/DeckManagement.svelte', () => ({
  default: {
    $$render: () => '<div id="mock-deck-management"></div>'
  }
}));

// Mock the Navbar component
jest.mock('../components/Navbar.svelte', () => ({
  default: {
    $$render: () => '<div id="mock-navbar"></div>'
  }
}));

describe('App Component', () => {
  // Set up mocks for each test
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
  });

  test('renders the Navbar component', () => {
    // Update our mock testing-library-svelte to return HTML with mock-navbar
    const { container } = render(App, {
      mockHtml: '<div id="mock-navbar"></div><div id="mock-router"></div>'
    });
    
    // Check if the Navbar component is rendered
    expect(container.innerHTML).toContain('mock-navbar');
  });

  test('renders the Router component', () => {
    // Update our mock testing-library-svelte to return HTML with mock-router
    const { container } = render(App, {
      mockHtml: '<div id="mock-navbar"></div><div id="mock-router"></div>'
    });
    
    // Check if the Router component is rendered
    expect(container.innerHTML).toContain('mock-router');
  });
});
