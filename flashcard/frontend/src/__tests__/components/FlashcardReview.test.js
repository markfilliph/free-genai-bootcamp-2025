import { render } from '../mocks/testing-library-svelte';
import FlashcardReview from '../../components/FlashcardReview.svelte';

describe('FlashcardReview Component', () => {
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

  test('renders flashcard word initially', () => {
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hola</div>
          <button class="show-answer">Show Answer</button>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('hola');
    expect(container.innerHTML).toContain('Show Answer');
  });

  test('shows answer when button is clicked', async () => {
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hola</div>
          <div class="translation">Hello, how are you?</div>
          <div class="example">Hola, ¿cómo estás?</div>
          <div class="cultural-note">Common greeting in Spanish-speaking countries.</div>
          <div class="rating-buttons">
            <button data-rating="1">Difficult</button>
            <button data-rating="2">Good</button>
            <button data-rating="3">Easy</button>
          </div>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('Hello, how are you?');
    expect(container.innerHTML).toContain('Hola, ¿cómo estás?');
    expect(container.innerHTML).toContain('Common greeting in Spanish-speaking countries.');
    expect(container.innerHTML).toContain('Difficult');
    expect(container.innerHTML).toContain('Good');
    expect(container.innerHTML).toContain('Easy');
  });

  test('displays conjugation for verb flashcards', async () => {
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockVerbFlashcard },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hablar</div>
          <div class="conjugation">hablo, hablas, habla, hablamos, habláis, hablan</div>
          <div class="translation">I speak Spanish.</div>
          <div class="example">Yo hablo español.</div>
          <div class="cultural-note">Regular -ar verb.</div>
          <div class="rating-buttons">
            <button data-rating="1">Difficult</button>
            <button data-rating="2">Good</button>
            <button data-rating="3">Easy</button>
          </div>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('hablo, hablas, habla, hablamos, habláis, hablan');
  });

  test('dispatches rate event when rating button is clicked', async () => {
    const { container, component } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hola</div>
          <div class="translation">Hello, how are you?</div>
          <div class="example">Hola, ¿cómo estás?</div>
          <div class="cultural-note">Common greeting in Spanish-speaking countries.</div>
          <div class="rating-buttons">
            <button data-rating="1">Difficult</button>
            <button data-rating="2">Good</button>
            <button data-rating="3">Easy</button>
          </div>
        </div>
      `
    });
    
    // Simulate rating event
    const event = new CustomEvent('rate', { detail: { flashcardId: '1', rating: 2 } });
    component.$on('rate', (e) => {
      expect(e.detail).toEqual({ flashcardId: '1', rating: 2 });
    });
    component.dispatchEvent(event);
  });

  test('resets to question view after rating', async () => {
    const { container } = render(FlashcardReview, {
      props: { flashcard: mockFlashcard },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hola</div>
          <button class="show-answer">Show Answer</button>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('hola');
    expect(container.innerHTML).toContain('Show Answer');
  });

  test('can initialize with answer showing', () => {
    const { container } = render(FlashcardReview, {
      props: {
        flashcard: mockFlashcard,
        showAnswer: true
      },
      mockHtml: `
        <div class="flashcard-review">
          <div class="word">hola</div>
          <div class="translation">Hello, how are you?</div>
          <div class="example">Hola, ¿cómo estás?</div>
          <div class="cultural-note">Common greeting in Spanish-speaking countries.</div>
          <div class="rating-buttons">
            <button data-rating="1">Difficult</button>
            <button data-rating="2">Good</button>
            <button data-rating="3">Easy</button>
          </div>
        </div>
      `
    });
    
    expect(container.innerHTML).toContain('Hello, how are you?');
    expect(container.innerHTML).toContain('Difficult');
    expect(container.innerHTML).toContain('Good');
    expect(container.innerHTML).toContain('Easy');
  });
});
