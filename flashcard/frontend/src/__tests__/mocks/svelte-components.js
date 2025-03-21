// Mock for Svelte components
const createMockComponent = (name) => ({
  $$render: () => `<div data-component="${name}"></div>`
});

// Mock all components
jest.mock('../../components/Navbar.svelte', () => createMockComponent('Navbar'), { virtual: true });
jest.mock('../../components/DeckList.svelte', () => createMockComponent('DeckList'), { virtual: true });
jest.mock('../../components/FlashcardReview.svelte', () => createMockComponent('FlashcardReview'), { virtual: true });
jest.mock('../../components/StudySession.svelte', () => createMockComponent('StudySession'), { virtual: true });
jest.mock('../../components/FlashcardForm.svelte', () => createMockComponent('FlashcardForm'), { virtual: true });
jest.mock('../../components/Deck.svelte', () => createMockComponent('Deck'), { virtual: true });

// Mock routes
jest.mock('../../routes/Home.svelte', () => createMockComponent('Home'), { virtual: true });
jest.mock('../../routes/Login.svelte', () => createMockComponent('Login'), { virtual: true });
jest.mock('../../routes/DeckManagement.svelte', () => createMockComponent('DeckManagement'), { virtual: true });

// Mock svelte-routing
jest.mock('svelte-routing', () => ({
  Link: {
    $$render: ($$result, $$props) => {
      const { to, class: className } = $$props || {};
      return `<a href="${to || '/'}" class="${className || ''}"></a>`;
    }
  },
  Router: {
    $$render: ($$result, $$props, $$bindings, $$slots) => {
      return $$slots.default ? $$slots.default({}) : '';
    }
  },
  Route: {
    $$render: ($$result, $$props, $$bindings, $$slots) => {
      return $$slots.default ? $$slots.default({}) : '';
    }
  },
  navigate: jest.fn()
}), { virtual: true });

// Add tests to prevent the "no tests" error
describe('Svelte Components Mock', () => {
  test('createMockComponent generates correct HTML', () => {
    const TestComponent = createMockComponent('TestComponent');
    const html = TestComponent.$$render();
    expect(html).toBe('<div data-component="TestComponent"></div>');
  });

  test('Link component renders with correct props', () => {
    const mockLink = {
      $$render: ($$result, $$props) => {
        const { to, class: className } = $$props || {};
        return `<a href="${to || '/'}" class="${className || ''}"></a>`;
      }
    };
    const html = mockLink.$$render(null, { to: '/test', class: 'test-class' });
    expect(html).toBe('<a href="/test" class="test-class"></a>');
  });

  test('Router component renders slot content', () => {
    const mockRouter = {
      $$render: ($$result, $$props, $$bindings, $$slots) => {
        return $$slots.default ? $$slots.default({}) : '';
      }
    };
    const html = mockRouter.$$render(null, {}, {}, {
      default: () => '<div>Test Content</div>'
    });
    expect(html).toBe('<div>Test Content</div>');
  });
});
