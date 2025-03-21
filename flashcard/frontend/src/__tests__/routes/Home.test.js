import { render } from '@testing-library/svelte';
import Home from '../../routes/Home.svelte';

describe('Home Component', () => {
  test('renders welcome message correctly', () => {
    const { getByText } = render(Home);
    
    // Check if welcome message is rendered
    expect(getByText('Welcome to Flashcard Master')).toBeInTheDocument();
    expect(getByText('Start creating and reviewing your flashcards!')).toBeInTheDocument();
  });

  test('renders call-to-action button with correct link', () => {
    const { getByText } = render(Home);
    
    // Check if CTA button is rendered with correct text and href
    const ctaButton = getByText('Get Started');
    expect(ctaButton).toBeInTheDocument();
    // Our updated mock now returns the correct href
    expect(ctaButton.getAttribute('href')).toBe('/decks');
    expect(ctaButton.classList.contains('cta-button')).toBe(true);
  });
});
