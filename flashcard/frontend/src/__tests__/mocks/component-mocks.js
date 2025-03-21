// Component mocks
export function mockComponents() {
  jest.mock('svelte-routing', () => ({
    Link: jest.fn().mockImplementation(() => ({
      $$render: () => '<a data-testid="mock-link"></a>'
    }))
  }));

  // Add DeckList mock
  jest.mock('../../components/DeckList.svelte', () => ({
    default: {
      render: (props) => `<div data-testid="mock-decklist"></div>`
    }
  }));

  // Add comprehensive mocks
jest.mock('svelte-routing', () => ({
  Link: {
    render: () => '<a data-testid="mock-link"></a>'
  }
}));

jest.mock('../../components/DeckList.svelte', () => ({
  default: {
    render: () => '<div data-testid="mock-decklist"></div>'
  }
}));

}
