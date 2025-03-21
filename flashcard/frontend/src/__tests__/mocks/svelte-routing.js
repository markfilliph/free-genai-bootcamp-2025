// Mock for svelte-routing

// Create a mock Link component
export const Link = {
  $$render: ($$result, $$props, $$bindings, $$slots) => {
    const { to, class: className } = $$props;
    return `<a href="${to || '/'}" class="${className || ''}">${$$slots.default ? $$slots.default({}) : ''}</a>`;
  }
};

// Create a mock Router component
export const Router = {
  $$render: ($$result, $$props, $$bindings, $$slots) => {
    return $$slots.default ? $$slots.default({}) : '';
  }
};

// Create a mock Route component
export const Route = {
  $$render: ($$result, $$props, $$bindings, $$slots) => {
    return $$slots.default ? $$slots.default({}) : '';
  },
  props: {
    path: '',
    component: null
  }
};

// Mock navigate function
export const navigate = jest.fn();

// Mock useLocation hook
export const useLocation = jest.fn(() => ({
  pathname: '/',
  search: '',
  hash: '',
  state: {}
}));

// Mock useNavigate hook
export const useNavigate = jest.fn(() => navigate);

// Mock useParams hook
export const useParams = jest.fn(() => ({}));

// Add tests to prevent the "no tests" error
describe('Svelte Routing Mock', () => {
  test('Link component renders correctly', () => {
    const html = Link.$$render(null, { to: '/test', class: 'test-class' }, {}, {
      default: () => 'Link Text'
    });
    expect(html).toBe('<a href="/test" class="test-class">Link Text</a>');
  });

  test('Router component renders slot content', () => {
    const html = Router.$$render(null, {}, {}, {
      default: () => '<div>Router Content</div>'
    });
    expect(html).toBe('<div>Router Content</div>');
  });

  test('Route component renders slot content', () => {
    const html = Route.$$render(null, {}, {}, {
      default: () => '<div>Route Content</div>'
    });
    expect(html).toBe('<div>Route Content</div>');
  });

  test('useLocation returns expected mock location', () => {
    const location = useLocation();
    expect(location).toEqual({
      pathname: '/',
      search: '',
      hash: '',
      state: {}
    });
  });

  test('useNavigate returns navigate function', () => {
    const navigateFunc = useNavigate();
    expect(navigateFunc).toBe(navigate);
  });
});
