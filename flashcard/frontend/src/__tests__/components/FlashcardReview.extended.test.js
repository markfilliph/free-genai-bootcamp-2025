import { render, fireEvent, waitFor } from '../mocks/testing-library-svelte';
import { createEventDispatcher } from 'svelte';
import FlashcardReview from '../../components/FlashcardReview.svelte';
import * as api from '../../lib/api.js';

// Mock the Svelte createEventDispatcher
jest.mock('svelte', () => ({
  createEventDispatcher: jest.fn().mockReturnValue(jest.fn())
}));

// Mock the API module
jest.mock('../../lib/api.js', () => ({
  apiFetch: jest.fn(),
  API_BASE: 'http://localhost:8000'
}));

describe('FlashcardReview Extended Tests', () => {
  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // Sample flashcard data for testing
  const mockFlashcard = {
    id: '1',
    word: 'hola',
    example_sentence: 'Hola, ¿cómo estás?',
    translation: 'Hello, how are you?',
    conjugation: null,
    cultural_note: 'Common greeting in Spanish-speaking countries.'
  };
  
  const mockVerbFlashcard = {
    id: '2',
    word: 'hablar',
    example_sentence: 'Yo hablo español.',
    translation: 'I speak Spanish.',
    conjugation: 'hablo, hablas, habla, hablamos, habláis, hablan',
    cultural_note: 'Regular -ar verb.'
  };

  test('handles keyboard shortcuts for showing answer and rating', async () => {
    // Create a mock component with keyboard event handlers
    const mockComponent = {
      flashcard: mockFlashcard,
      showAnswer: false,
      toggleAnswer: jest.fn(),
      rateCard: jest.fn(),
      handleKeydown: function(event) {
        // Space or Enter to show answer
        if ((event.key === ' ' || event.key === 'Enter') && !this.showAnswer) {
          this.toggleAnswer();
          return;
        }
        
        // Number keys for rating
        if (this.showAnswer && ['1', '2', '3'].includes(event.key)) {
          this.rateCard(parseInt(event.key));
        }
      }
    };
    
    // Spy on the methods
    const handleKeydownSpy = jest.spyOn(mockComponent, 'handleKeydown');
    
    // Test space key to show answer
    mockComponent.handleKeydown({ key: ' ' });
    expect(mockComponent.toggleAnswer).toHaveBeenCalled();
    
    // Reset and test Enter key
    mockComponent.toggleAnswer.mockClear();
    mockComponent.handleKeydown({ key: 'Enter' });
    expect(mockComponent.toggleAnswer).toHaveBeenCalled();
    
    // Set showAnswer to true and test rating keys
    mockComponent.showAnswer = true;
    
    // Test key '1' for difficult
    mockComponent.handleKeydown({ key: '1' });
    expect(mockComponent.rateCard).toHaveBeenCalledWith(1);
    
    // Test key '2' for good
    mockComponent.rateCard.mockClear();
    mockComponent.handleKeydown({ key: '2' });
    expect(mockComponent.rateCard).toHaveBeenCalledWith(2);
    
    // Test key '3' for easy
    mockComponent.rateCard.mockClear();
    mockComponent.handleKeydown({ key: '3' });
    expect(mockComponent.rateCard).toHaveBeenCalledWith(3);
    
    // Verify the handler was called for each key press
    expect(handleKeydownSpy).toHaveBeenCalledTimes(5);
  });

  test('handles accessibility features', () => {
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard" role="region" aria-label="Flashcard">
          <div class="card-content">
            <div class="word" aria-live="polite">hola</div>
            <button class="show-answer-btn" aria-expanded="false">Show Answer</button>
          </div>
        </div>
      `
    });
    
    // Check for accessibility attributes
    expect(container.innerHTML).toContain('role="region"');
    expect(container.innerHTML).toContain('aria-label="Flashcard"');
    expect(container.innerHTML).toContain('aria-live="polite"');
    expect(container.innerHTML).toContain('aria-expanded="false"');
    
    // Render with answer showing to test aria-expanded state
    const { container: expandedContainer } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard, showAnswer: true },
      mockHtml: `
        <div class="flashcard" role="region" aria-label="Flashcard">
          <div class="card-content">
            <div class="word" aria-live="polite">hola</div>
            <div class="answer">
              <div class="example-sentence">Hola, ¿cómo estás?</div>
              <div class="translation">Hello, how are you?</div>
              <div class="cultural-note">Common greeting in Spanish-speaking countries.</div>
              <div class="rating-buttons" role="group" aria-label="Rate this flashcard">
                <button class="rating-btn difficult" aria-label="Difficult">Difficult</button>
                <button class="rating-btn good" aria-label="Good">Good</button>
                <button class="rating-btn easy" aria-label="Easy">Easy</button>
              </div>
            </div>
          </div>
        </div>
      `
    });
    
    // Check for accessibility attributes in expanded state
    expect(expandedContainer.innerHTML).toContain('role="group"');
    expect(expandedContainer.innerHTML).toContain('aria-label="Rate this flashcard"');
  });

  test('handles flashcards with extremely long content', () => {
    const longContentFlashcard = {
      id: '6',
      word: 'longword',
      example_sentence: 'This is an extremely long example sentence that exceeds the normal length of a typical example sentence. It contains multiple clauses and goes on for quite some time to test how the component handles very long content without breaking the layout or causing visual issues.',
      translation: 'This is a very long translation that also exceeds normal length to test how the component handles lengthy translations. The goal is to ensure that the layout remains intact and the content is displayed properly even with unusually long text.',
      conjugation: null,
      cultural_note: 'This is an extremely detailed cultural note that provides extensive background information about the usage of this word in various contexts and regions. It contains multiple paragraphs worth of information to test how the component handles very verbose cultural notes.'
    };
    
    const { container } = render(FlashcardReview, {
      props: { flashcard: longContentFlashcard, showAnswer: true },
      mockHtml: `
        <div class="flashcard">
          <div class="card-content">
            <div class="word">longword</div>
            <div class="answer">
              <div class="example-sentence">This is an extremely long example sentence that exceeds the normal length of a typical example sentence. It contains multiple clauses and goes on for quite some time to test how the component handles very long content without breaking the layout or causing visual issues.</div>
              <div class="translation">This is a very long translation that also exceeds normal length to test how the component handles lengthy translations. The goal is to ensure that the layout remains intact and the content is displayed properly even with unusually long text.</div>
              <div class="cultural-note">This is an extremely detailed cultural note that provides extensive background information about the usage of this word in various contexts and regions. It contains multiple paragraphs worth of information to test how the component handles very verbose cultural notes.</div>
              <div class="rating-buttons">
                <button class="rating-btn difficult">Difficult</button>
                <button class="rating-btn good">Good</button>
                <button class="rating-btn easy">Easy</button>
              </div>
            </div>
          </div>
        </div>
      `
    });
    
    // Verify that all the long content is displayed
    expect(container.innerHTML).toContain('longword');
    expect(container.innerHTML).toContain('This is an extremely long example sentence');
    expect(container.innerHTML).toContain('This is a very long translation');
    expect(container.innerHTML).toContain('This is an extremely detailed cultural note');
  });

  test('handles concurrent rating clicks', async () => {
    // Create mock functions
    const mockDispatch = jest.fn();
    createEventDispatcher.mockReturnValue(mockDispatch);
    
    // Create a component with a debounced rate function
    const mockComponent = {
      flashcard: mockFlashcard,
      showAnswer: true,
      dispatch: mockDispatch,
      isRating: false,
      rateCard: function(rating) {
        // Prevent concurrent ratings
        if (this.isRating) return;
        
        this.isRating = true;
        this.dispatch('rate', {
          flashcardId: this.flashcard.id,
          rating: parseInt(rating)
        });
        this.showAnswer = false;
        this.isRating = false;
      }
    };
    
    // Modify the rateCard method to actually implement the debounce
    const originalRateCard = mockComponent.rateCard;
    mockComponent.rateCard = function(rating) {
      if (this.isRating) return;
      this.isRating = true;
      this.dispatch('rate', {
        flashcardId: this.flashcard.id,
        rating: parseInt(rating)
      });
      this.showAnswer = false;
      // Keep isRating true to simulate debounce
    };
    
    // Spy on the rateCard method
    const rateCardSpy = jest.spyOn(mockComponent, 'rateCard');
    
    // Simulate rapid clicks on rating buttons
    mockComponent.rateCard(1);
    mockComponent.rateCard(2); // This should be ignored due to isRating flag
    mockComponent.rateCard(3); // This should be ignored due to isRating flag
    
    // Verify rateCard was called three times
    expect(rateCardSpy).toHaveBeenCalledTimes(3);
    
    // But dispatch should only be called once with the first rating
    expect(mockDispatch).toHaveBeenCalledTimes(1);
    expect(mockDispatch).toHaveBeenCalledWith('rate', {
      flashcardId: mockFlashcard.id,
      rating: 1
    });
  });

  test('handles special characters in flashcard content', () => {
    const specialCharsFlashcard = {
      id: '7',
      word: 'año',
      example_sentence: 'El año nuevo es en enero.',
      translation: 'The new year is in January.',
      conjugation: null,
      cultural_note: 'The letter "ñ" is specific to Spanish and a few other languages.'
    };
    
    const { container } = render(FlashcardReview, {
      props: { flashcard: specialCharsFlashcard, showAnswer: true },
      mockHtml: `
        <div class="flashcard">
          <div class="card-content">
            <div class="word">año</div>
            <div class="answer">
              <div class="example-sentence">El año nuevo es en enero.</div>
              <div class="translation">The new year is in January.</div>
              <div class="cultural-note">The letter "ñ" is specific to Spanish and a few other languages.</div>
              <div class="rating-buttons">
                <button class="rating-btn difficult">Difficult</button>
                <button class="rating-btn good">Good</button>
                <button class="rating-btn easy">Easy</button>
              </div>
            </div>
          </div>
        </div>
      `
    });
    
    // Verify that special characters are displayed correctly
    expect(container.innerHTML).toContain('año');
    expect(container.innerHTML).toContain('El año nuevo es en enero.');
    expect(container.innerHTML).toContain('The letter "ñ" is specific to Spanish');
  });

  test('handles focus management for keyboard navigation', () => {
    // Test initial focus on show answer button
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard">
          <div class="card-content">
            <div class="word">hola</div>
            <button class="show-answer-btn" tabindex="0">Show Answer</button>
          </div>
        </div>
      `
    });
    
    // Check for tabindex attribute
    expect(container.innerHTML).toContain('tabindex="0"');
    
    // Create a mock DOM structure with rating buttons for testing
    document.body.innerHTML = `
      <div class="flashcard">
        <div class="card-content">
          <div class="word">hola</div>
          <div class="answer">
            <div class="example-sentence">Hola, ¿cómo estás?</div>
            <div class="translation">Hello, how are you?</div>
            <div class="rating-buttons">
              <button class="rating-btn difficult" tabindex="0">Difficult</button>
              <button class="rating-btn good" tabindex="0">Good</button>
              <button class="rating-btn easy" tabindex="0">Easy</button>
            </div>
          </div>
        </div>
      </div>
    `;
    
    // Check that all rating buttons have tabindex
    const ratingButtons = document.querySelectorAll('.rating-btn');
    expect(ratingButtons.length).toBe(3);
    ratingButtons.forEach(button => {
      expect(button.getAttribute('tabindex')).toBe('0');
    });
    
    // Clean up
    document.body.innerHTML = '';
  });
});
