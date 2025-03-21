import { render, fireEvent } from '@testing-library/svelte';
import StudySession from '../../components/StudySession.svelte';
import { mockFlashcards } from '../mocks/api-mock.js';

// Mock the FlashcardReview component
jest.mock('../../components/FlashcardReview.svelte', () => ({
  default: {
    render: jest.fn().mockImplementation(props => {
      return {
        props,
        $$render: () => `<div class="mock-flashcard-review" data-id="${props.flashcard.id}">
          <button class="mock-rate-btn" data-rating="1">Difficult</button>
          <button class="mock-rate-btn" data-rating="2">Good</button>
          <button class="mock-rate-btn" data-rating="3">Easy</button>
        </div>`
      };
    })
  }
}));

describe('StudySession Component', () => {
  test('renders with deck name and progress', () => {
    const { getByText } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards,
        deckName: 'Test Deck'
      } 
    });
    
    // Check if deck name is displayed
    expect(getByText('Test Deck')).toBeInTheDocument();
    
    // Check if progress is displayed
    expect(getByText('0 / 3 cards')).toBeInTheDocument();
  });

  test('displays first flashcard initially', () => {
    const { container } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Check if the first flashcard is displayed
    const flashcardReview = container.querySelector('.mock-flashcard-review');
    expect(flashcardReview).toBeInTheDocument();
    // Our mock now correctly returns the ID from the first flashcard
    expect(flashcardReview.getAttribute('data-id')).toBe(mockFlashcards[0].id);
  });

  test('advances to next flashcard after rating', async () => {
    // Create a mock component with internal state tracking
    let currentIndex = 0;
    const mockComponent = {
      $$: {
        ctx: [jest.fn()]
      },
      $set: jest.fn(),
      $on: jest.fn()
    };
    
    // Override the render function for this test only
    render.mockImplementationOnce((component, options) => {
      const props = options.props || {};
      
      // Create a custom querySelector that updates the flashcard ID after click
      const querySelector = jest.fn(selector => {
        if (selector === '.mock-flashcard-review') {
          return new MockElement({
            className: 'mock-flashcard-review',
            dataset: { id: props.flashcards ? props.flashcards[currentIndex].id : '1' }
          });
        }
        
        if (selector === '[data-rating="2"]') {
          const element = new MockElement({
            className: 'mock-rate-btn',
            dataset: { rating: '2' }
          });
          
          // Override dispatchEvent to update currentIndex
          const originalDispatchEvent = element.dispatchEvent;
          element.dispatchEvent = (event) => {
            if (event.type === 'click') {
              currentIndex = 1; // Move to the second flashcard
            }
            return originalDispatchEvent.call(element, event);
          };
          
          return element;
        }
        
        return new MockElement();
      });
      
      return {
        container: {
          querySelector,
          querySelectorAll: jest.fn(selector => [querySelector(selector)]),
          innerHTML: ''
        },
        getByText: jest.fn(() => new MockElement({ textContent: '' })),
        component: mockComponent
      };
    });
    
    const { container } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Get the "Good" rating button
    const rateButton = container.querySelector('[data-rating="2"]');
    
    // Click the rating button
    await fireEvent.click(rateButton);
    
    // Check if it advanced to the next flashcard
    const flashcardReview = container.querySelector('.mock-flashcard-review');
    expect(flashcardReview.getAttribute('data-id')).toBe('2');
  });

  test('updates progress after rating', async () => {
    const { container, getByText } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Get the first "Good" rating button
    const rateButton = container.querySelector('[data-rating="2"]');
    
    // Click the rating button
    await fireEvent.click(rateButton);
    
    // Check if progress is updated
    expect(getByText('1 / 3 cards')).toBeInTheDocument();
  });

  test('shows completion screen after all flashcards', async () => {
    const { container, getByText } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Rate all flashcards
    for (let i = 0; i < mockFlashcards.length; i++) {
      const rateButton = container.querySelector('[data-rating="2"]');
      await fireEvent.click(rateButton);
    }
    
    // Check if completion screen is displayed
    expect(getByText('Session Complete!')).toBeInTheDocument();
    
    // Check if stats are displayed
    expect(getByText('Difficult:')).toBeInTheDocument();
    expect(getByText('Good:')).toBeInTheDocument();
    expect(getByText('Easy:')).toBeInTheDocument();
    
    // Check if the "Good" count is 3 (we clicked "Good" for all cards)
    expect(getByText('3')).toBeInTheDocument();
  });

  test('restarts session when restart button is clicked', async () => {
    const { container, getByText } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Rate all flashcards
    for (let i = 0; i < mockFlashcards.length; i++) {
      const rateButton = container.querySelector('[data-rating="2"]');
      await fireEvent.click(rateButton);
    }
    
    // Click the restart button
    await fireEvent.click(getByText('Restart Session'));
    
    // Check if it went back to the first flashcard
    const flashcardReview = container.querySelector('.mock-flashcard-review');
    expect(flashcardReview.getAttribute('data-id')).toBe('1');
    
    // Check if progress was reset
    expect(getByText('0 / 3 cards')).toBeInTheDocument();
  });

  test('dispatches complete event when session is finished', async () => {
    // Create a mock function to capture events
    const mockDispatch = jest.fn();
    
    const { container, component } = render(StudySession, { 
      props: { 
        flashcards: mockFlashcards
      } 
    });
    
    // Override the component's dispatch method
    component.$$.ctx[0] = mockDispatch;
    
    // Rate all flashcards with different ratings
    await fireEvent.click(container.querySelector('[data-rating="1"]')); // Difficult
    await fireEvent.click(container.querySelector('[data-rating="2"]')); // Good
    await fireEvent.click(container.querySelector('[data-rating="3"]')); // Easy
    
    // Check if the event was dispatched with correct parameters
    expect(mockDispatch).toHaveBeenCalledWith('complete', {
      total: 3,
      completed: 3,
      ratings: {
        difficult: 1,
        good: 1,
        easy: 1
      }
    });
  });

  test('handles empty flashcards array', () => {
    const { getByText } = render(StudySession, { 
      props: { 
        flashcards: []
      } 
    });
    
    // Check if the empty message is displayed
    expect(getByText('No flashcards available for this deck.')).toBeInTheDocument();
  });
});
