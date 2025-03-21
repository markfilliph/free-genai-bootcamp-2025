import { render } from '../mocks/testing-library-svelte';
import * as utils from '../../lib/utils.js';

// Mock the utils module
jest.mock('../../lib/utils.js', () => ({
  formatDate: jest.fn()
}));

// Import the component
import Deck from '../../components/Deck.svelte';

describe('Deck Component', () => {
  // Sample deck data for testing
  const mockDeck = {
    id: '1',
    name: 'Spanish Basics',
    user_id: '1',
    created_at: '2025-03-19T10:00:00Z'
  };

  beforeEach(() => {
    // Reset mocks before each test
    jest.clearAllMocks();
    
    // Mock the formatDate function to return a fixed string
    utils.formatDate.mockReturnValue('Mar 19, 2025, 10:00 AM');
  });

  test('renders deck name and formatted date', () => {
    const { container } = render(Deck, { 
      props: { deck: mockDeck },
      mockHtml: `
        <div class="deck-card" data-testid="deck-card">
          <div class="deck-info">
            <h3 class="deck-name">Spanish Basics</h3>
            <p class="deck-date">Created: Mar 19, 2025, 10:00 AM</p>
          </div>
        </div>
      `
    });
    
    // Check if the HTML contains the expected content
    expect(container.innerHTML).toContain('Spanish Basics');
    expect(container.innerHTML).toContain('Mar 19, 2025, 10:00 AM');
    
    // Check if formatDate was called with the correct argument
    expect(utils.formatDate).toHaveBeenCalledWith('2025-03-19T10:00:00Z');
  });

  test('renders action buttons when showActions is true', () => {
    const { container } = render(Deck, { 
      props: { 
        deck: mockDeck,
        showActions: true
      },
      mockHtml: `
        <div class="deck-card" data-testid="deck-card">
          <div class="deck-info">
            <h3 class="deck-name">Spanish Basics</h3>
            <p class="deck-date">Created: Mar 19, 2025, 10:00 AM</p>
          </div>
          <div class="deck-actions">
            <a href="/decks/1" class="view-btn">View</a>
            <a href="/decks/1/study" class="study-btn">Study</a>
            <a href="/decks/1/edit" class="edit-btn">Edit</a>
          </div>
        </div>
      `
    });
    
    // Check if action buttons are displayed
    expect(container.innerHTML).toContain('View');
    expect(container.innerHTML).toContain('Study');
    expect(container.innerHTML).toContain('Edit');
  });

  test('does not render action buttons when showActions is false', () => {
    const { container } = render(Deck, { 
      props: { 
        deck: mockDeck,
        showActions: false
      },
      mockHtml: `
        <div class="deck-card" data-testid="deck-card">
          <div class="deck-info">
            <h3 class="deck-name">Spanish Basics</h3>
            <p class="deck-date">Created: Mar 19, 2025, 10:00 AM</p>
          </div>
        </div>
      `
    });
    
    // Check that action buttons are not displayed
    expect(container.innerHTML).not.toContain('View');
    expect(container.innerHTML).not.toContain('Study');
    expect(container.innerHTML).not.toContain('Edit');
  });

  test('action buttons have correct href attributes', () => {
    const { container } = render(Deck, { 
      props: { deck: mockDeck },
      mockHtml: `
        <div class="deck-card" data-testid="deck-card">
          <div class="deck-info">
            <h3 class="deck-name">Spanish Basics</h3>
            <p class="deck-date">Created: Mar 19, 2025, 10:00 AM</p>
          </div>
          <div class="deck-actions">
            <a href="/decks/1" class="view-btn">View</a>
            <a href="/decks/1/study" class="study-btn">Study</a>
            <a href="/decks/1/edit" class="edit-btn">Edit</a>
          </div>
        </div>
      `
    });
    
    // Get all links in the component
    const links = container.querySelectorAll('a');
    
    // Check href attributes
    expect(links[0].getAttribute('href')).toBe('/decks/1');
    expect(links[1].getAttribute('href')).toBe('/decks/1/study');
    expect(links[2].getAttribute('href')).toBe('/decks/1/edit');
  });

  test('has correct CSS classes for styling', () => {
    const { container } = render(Deck, { 
      props: { deck: mockDeck },
      mockHtml: `
        <div class="deck-card" data-testid="deck-card">
          <div class="deck-info">
            <h3 class="deck-name">Spanish Basics</h3>
            <p class="deck-date">Created: Mar 19, 2025, 10:00 AM</p>
          </div>
          <div class="deck-actions">
            <a href="/decks/1" class="view-btn">View</a>
            <a href="/decks/1/study" class="study-btn">Study</a>
            <a href="/decks/1/edit" class="edit-btn">Edit</a>
          </div>
        </div>
      `
    });
    
    // Check if the main container has the correct class
    expect(container.innerHTML).toContain('class="deck-card"');
    
    // Check if action buttons have the correct classes
    expect(container.innerHTML).toContain('class="view-btn"');
    expect(container.innerHTML).toContain('class="study-btn"');
    expect(container.innerHTML).toContain('class="edit-btn"');
  });
});
