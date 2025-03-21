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
}

// Add simple test to prevent Jest from complaining
describe('Component Mocks', () => {
  test('mockComponents function exists', () => {
    expect(typeof mockComponents).toBe('function');
  });
});
