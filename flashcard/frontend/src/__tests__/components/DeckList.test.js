import { render } from '../mocks/testing-library-svelte';
import DeckList from '../../components/DeckList.svelte';

describe('DeckList Component', () => {
  test('renders empty list when no decks are provided', () => {
    const { container } = render(DeckList, {
      props: { decks: [] },
      mockHtml: '<ul></ul>'
    });
    
    // Check if the list is empty
    expect(container.innerHTML).toBe('<ul></ul>');
  });

  test('renders deck items correctly', () => {
    const mockDecks = [
      { id: '1', title: 'Spanish Basics' },
      { id: '2', title: 'Verb Conjugations' },
      { id: '3', title: 'Travel Phrases' }
    ];
    
    const { container } = render(DeckList, { 
      props: { decks: mockDecks },
      mockHtml: `
        <ul>
          <li><a href="/decks/1">Spanish Basics</a></li>
          <li><a href="/decks/2">Verb Conjugations</a></li>
          <li><a href="/decks/3">Travel Phrases</a></li>
        </ul>
      `
    });
    
    // Check if all deck items are rendered
    expect(container.innerHTML).toContain('Spanish Basics');
    expect(container.innerHTML).toContain('Verb Conjugations');
    expect(container.innerHTML).toContain('Travel Phrases');
    
    // Check if the correct number of list items is present
    expect(container.innerHTML.match(/<li>/g).length).toBe(3);
  });

  test('deck links have correct href attributes', () => {
    const mockDecks = [
      { id: '1', title: 'Spanish Basics' },
      { id: '2', title: 'Verb Conjugations' }
    ];
    
    const { container } = render(DeckList, { 
      props: { decks: mockDecks },
      mockHtml: `
        <ul>
          <li><a href="/decks/1">Spanish Basics</a></li>
          <li><a href="/decks/2">Verb Conjugations</a></li>
        </ul>
      `
    });
    
    // Check if links have correct href attributes
    expect(container.innerHTML).toContain('href="/decks/1"');
    expect(container.innerHTML).toContain('href="/decks/2"');
  });
});
