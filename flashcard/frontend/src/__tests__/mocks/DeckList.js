// Mock implementation for DeckList component
export default function DeckList(options) {
  const { props } = options || {};
  const decks = props?.decks || [];
  
  // No DOM manipulation in the mock
  
  return {
    $$: {
      fragment: document.createDocumentFragment()
    },
    $$render: () => {
      let html = '<div class="deck-list">';
      
      if (decks.length === 0) {
        html += '<p>No decks found. Create your first deck!</p>';
      } else {
        decks.forEach(deck => {
          html += `<div class="deck-item" data-deck-id="${deck.id}">`;
          html += `<h3>${deck.name}</h3>`;
          html += `<p>Created: ${new Date(deck.created_at).toLocaleDateString()}</p>`;
          html += '</div>';
        });
      }
      
      html += '</div>';
      return html;
    },
    // Add any component methods that might be called in tests
    update: (newProps) => {
      if (newProps.decks) {
        props.decks = newProps.decks;
      }
    }
  };
}
