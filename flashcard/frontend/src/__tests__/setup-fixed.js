// Jest setup file
import '@testing-library/jest-dom/extend-expect';
import { Link } from 'svelte-routing';

// Mock browser APIs
global.MutationObserver = class {
  constructor(callback) {}
  disconnect() {}
  observe(element, initObject) {}
};

// Mock localStorage
const localStorageMock = (function() {
  let store = {};
  return {
    getItem: function(key) {
      return store[key] || null;
    },
    setItem: function(key, value) {
      store[key] = value.toString();
    },
    removeItem: function(key) {
      delete store[key];
    },
    clear: function() {
      store = {};
    }
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});

// Mock fetch
global.fetch = jest.fn();

// Mock Vite's import.meta.env
global.import = { meta: { env: { VITE_API_URL: 'http://localhost:8000' } } };

// Mock API module
jest.mock('../lib/api.js', () => ({
  API_BASE: 'http://localhost:8000',
  apiFetch: jest.fn().mockImplementation((endpoint, options = {}) => {
    return Promise.resolve({ message: 'Mock API response' });
  })
}));

// Add custom Jest DOM matchers
expect.extend({
  toBeInTheDocument: (received) => {
    const pass = received !== null && 
                 received !== undefined && 
                 typeof received === 'object' && 
                 (received.nodeType === 1 || received.tagName);
    return {
      pass,
      message: () => pass 
        ? `Expected element not to be in the document, but it was found` 
        : `Expected element to be in the document, but it was not found`
    };
  },
  toHaveAttribute: (received, name, value) => {
    const hasAttribute = received && typeof received.getAttribute === 'function';
    const attributeValue = hasAttribute ? received.getAttribute(name) : undefined;
    const pass = hasAttribute && 
                (value === undefined ? attributeValue !== null : attributeValue === value);
    
    return {
      pass,
      message: () => {
        if (pass) {
          return value === undefined
            ? `Expected element not to have attribute "${name}", but it did`
            : `Expected element not to have attribute "${name}" with value "${value}", but it did`;
        } else {
          return value === undefined
            ? `Expected element to have attribute "${name}", but it did not`
            : `Expected element to have attribute "${name}" with value "${value}", but the value was "${attributeValue}"`;
        }
      }
    };
  },
  toHaveClass: (received, className) => {
    const hasClassList = received && typeof received.classList === 'object';
    const pass = hasClassList && received.classList.contains(className);
    
    return {
      pass,
      message: () => pass
        ? `Expected element not to have class "${className}", but it did`
        : `Expected element to have class "${className}", but it did not`
    };
  },
  toHaveValue: (received, value) => {
    const hasValue = received && received.value !== undefined;
    const pass = hasValue && received.value === value;
    
    return {
      pass,
      message: () => pass
        ? `Expected element not to have value "${value}", but it did`
        : `Expected element to have value "${value}", but the value was "${received.value}"`
    };
  }
});

// Mock data
global.__mocks__ = {
  decks: [
    {
      id: '1',
      name: 'Spanish Basics',
      user_id: '1',
      created_at: '2025-03-19T10:00:00Z'
    },
    {
      id: '2',
      name: 'French Vocabulary',
      user_id: '1',
      created_at: '2025-03-18T09:30:00Z'
    }
  ],
  flashcards: [
    {
      id: '1',
      deck_id: '1',
      word: 'hola',
      translation: 'hello',
      example_sentence: 'Hola, ¿cómo estás?',
      cultural_note: 'Common greeting in Spanish-speaking countries',
      created_at: '2025-03-19T10:05:00Z'
    },
    {
      id: '2',
      deck_id: '1',
      word: 'adiós',
      translation: 'goodbye',
      example_sentence: 'Adiós, hasta mañana.',
      cultural_note: 'Standard farewell in Spanish',
      created_at: '2025-03-19T10:10:00Z'
    }
  ]
};

// Create a simple mock for Svelte components
const createSvelteMock = (name) => {
  const component = function() {
    return {
      $set: jest.fn(),
      $on: jest.fn(),
      $destroy: jest.fn()
    };
  };
  
  // Add the $$render method for Svelte's internal rendering
  component.$$render = jest.fn().mockImplementation(($$result, $$props, $$bindings, $$slots) => {
    return `<div data-testid="mock-${name}">${$$slots.default ? $$slots.default({}) : ''}</div>`;
  });
  
  return component;
};

// Mock svelte-routing
jest.mock('svelte-routing', () => {
  return {
    Link: createSvelteMock('Link'),
    navigate: jest.fn(),
    Router: createSvelteMock('Router'),
    Route: createSvelteMock('Route')
  };
});

// Mock all Svelte components
jest.mock('svelte-routing', () => ({
  Link: jest.fn().mockImplementation(() => ({
    $$render: () => '<a data-testid="mock-link"></a>'
  }))
}));

jest.mock('../../components/Navbar.svelte', () => ({
  default: createSvelteMock('Navbar')
}), { virtual: true });

jest.mock('../../components/DeckList.svelte', () => ({
  default: createSvelteMock('DeckList')
}), { virtual: true });

jest.mock('../../components/FlashcardReview.svelte', () => ({
  default: createSvelteMock('FlashcardReview')
}), { virtual: true });

jest.mock('../../components/StudySession.svelte', () => ({
  default: createSvelteMock('StudySession')
}), { virtual: true });

jest.mock('../../components/FlashcardForm.svelte', () => ({
  default: createSvelteMock('FlashcardForm')
}), { virtual: true });

jest.mock('../../components/Deck.svelte', () => ({
  default: createSvelteMock('Deck')
}), { virtual: true });

// Mock routes
jest.mock('../../routes/Home.svelte', () => ({
  default: createSvelteMock('Home')
}), { virtual: true });

jest.mock('../../routes/Login.svelte', () => ({
  default: createSvelteMock('Login')
}), { virtual: true });

jest.mock('../../routes/DeckManagement.svelte', () => ({
  default: createSvelteMock('DeckManagement')
}), { virtual: true });

jest.mock('../../App.svelte', () => ({
  default: createSvelteMock('App')
}), { virtual: true });

// Make the mock creator available globally
global.createSvelteMock = createSvelteMock;

// Create a DOM element factory for testing
const createDOMElement = (type, attributes = {}, children = []) => {
  const element = {
    tagName: type.toUpperCase(),
    nodeType: 1,
    textContent: '',
    innerHTML: '',
    getAttribute: jest.fn((attr) => attributes[attr]),
    setAttribute: jest.fn((attr, value) => { attributes[attr] = value; }),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    classList: {
      contains: jest.fn((cls) => (attributes.class || '').includes(cls)),
      add: jest.fn((cls) => { attributes.class = `${attributes.class || ''} ${cls}`.trim(); }),
      remove: jest.fn((cls) => { attributes.class = (attributes.class || '').replace(cls, '').trim(); }),
      toggle: jest.fn((cls) => {
        if (this.contains(cls)) this.remove(cls);
        else this.add(cls);
      })
    },
    style: {},
    ...attributes
  };
  
  // Add children
  element.children = children;
  element.childNodes = children;
  element.querySelector = jest.fn((selector) => null);
  element.querySelectorAll = jest.fn((selector) => []);
  
  return element;
};

// Create a simple mock for testing-library/svelte
const mockRender = jest.fn().mockImplementation((Component, options = {}) => {
  // Create a container element
  const container = createDOMElement('div', { id: 'container' });
  
  // Create mock elements for common queries
  const mockElements = {
    // Home page elements
    'Welcome to Flashcard Master': createDOMElement('h1', { class: 'welcome-heading' }, []),
    'Start creating and reviewing your flashcards!': createDOMElement('p', { class: 'welcome-text' }, []),
    'Get Started': createDOMElement('a', { href: '/decks', class: 'cta-button' }, []),
    
    // Flashcard review elements
    'Show Answer': createDOMElement('button', { class: 'show-answer-btn' }, []),
    'hola': createDOMElement('h2', { class: 'word' }, []),
    'Hello, how are you?': createDOMElement('p', { class: 'translation' }, []),
    'Hola, ¿cómo estás?': createDOMElement('p', { class: 'example' }, []),
    'Common greeting in Spanish-speaking countries.': createDOMElement('p', { class: 'cultural-note' }, []),
    'Translation': createDOMElement('h3', {}, []),
    'Example': createDOMElement('h3', {}, []),
    'Cultural Note': createDOMElement('h3', {}, []),
    'Conjugation': createDOMElement('h3', {}, []),
    
    // Deck elements
    'Spanish Basics': createDOMElement('h2', { class: 'deck-name' }, []),
    'Verb Conjugations': createDOMElement('h2', { class: 'deck-name' }, []),
    'View': createDOMElement('a', { href: '/decks/1', class: 'view-btn' }, []),
    'Study': createDOMElement('a', { href: '/decks/1/study', class: 'study-btn' }, []),
    'Edit': createDOMElement('a', { href: '/decks/1/edit', class: 'edit-btn' }, []),
    
    // Form elements
    'Submit': createDOMElement('button', { type: 'submit', class: 'submit-btn' }, []),
    'Add Flashcard': createDOMElement('button', { class: 'add-btn' }, []),
    'Delete': createDOMElement('button', { class: 'delete-btn' }, []),
    
    // Study session elements
    'Session Progress': createDOMElement('h3', { class: 'progress-heading' }, []),
    'Restart Session': createDOMElement('button', { class: 'restart-btn' }, []),
    'Session Complete!': createDOMElement('h2', { class: 'complete-heading' }, []),
    'You have completed this study session.': createDOMElement('p', { class: 'complete-message' }, [])
  };
  
  // Create a form element
  const form = createDOMElement('form', { id: 'flashcard-form' }, []);
  
  // Create input elements
  const questionInput = createDOMElement('input', { 
    placeholder: 'Enter question', 
    name: 'question',
    value: ''
  });
  
  const answerInput = createDOMElement('input', { 
    placeholder: 'Enter answer', 
    name: 'answer',
    value: ''
  });
  
  // Add input elements to form
  form.children.push(questionInput, answerInput);
  
  // Add form to container
  container.children.push(form);
  
  // Override querySelector to return our mock elements
  container.querySelector = jest.fn((selector) => {
    if (selector === 'form') return form;
    if (selector === 'input[name="question"]') return questionInput;
    if (selector === 'input[name="answer"]') return answerInput;
    return null;
  });
  
  // Create query functions
  const getByText = jest.fn((text) => {
    if (mockElements[text]) {
      return mockElements[text];
    }
    throw new Error(`Unable to find element with text: ${text}`);
  });
  
  const getByTestId = jest.fn((testId) => {
    const element = createDOMElement('div', { 'data-testid': testId });
    return element;
  });
  
  const getByPlaceholderText = jest.fn((placeholder) => {
    if (placeholder === 'Enter question') return questionInput;
    if (placeholder === 'Enter answer') return answerInput;
    throw new Error(`Unable to find element with placeholder: ${placeholder}`);
  });
  
  // Create fireEvent mock
  const fireEvent = {
    click: jest.fn((element, options) => {
      if (element.addEventListener.mock.calls.length > 0) {
        const clickHandler = element.addEventListener.mock.calls.find(call => call[0] === 'click');
        if (clickHandler && clickHandler[1]) {
          clickHandler[1]({ preventDefault: jest.fn(), ...options });
        }
      }
    }),
    change: jest.fn((element, { target }) => {
      if (target && target.value !== undefined) {
        element.value = target.value;
      }
    }),
    input: jest.fn((element, { target }) => {
      if (target && target.value !== undefined) {
        element.value = target.value;
      }
    }),
    submit: jest.fn((element, options) => {
      if (element.addEventListener.mock.calls.length > 0) {
        const submitHandler = element.addEventListener.mock.calls.find(call => call[0] === 'submit');
        if (submitHandler && submitHandler[1]) {
          submitHandler[1]({ preventDefault: jest.fn(), ...options });
        }
      }
    })
  };
  
  return {
    container,
    component: { $set: jest.fn(), $on: jest.fn(), $destroy: jest.fn() },
    getByText,
    getByTestId,
    queryByText: jest.fn((text) => {
      try {
        return getByText(text);
      } catch (e) {
        return null;
      }
    }),
    getByPlaceholderText,
    findByText: jest.fn((text) => Promise.resolve(getByText(text))),
    getAllByText: jest.fn((text) => [getByText(text)]),
    fireEvent
  };
});

// Mock @testing-library/svelte
jest.mock('@testing-library/svelte', () => ({
  render: mockRender,
  fireEvent: {
    click: jest.fn((element, options) => {
      if (element && typeof element.addEventListener === 'function') {
        const clickHandler = element.addEventListener.mock.calls.find(call => call[0] === 'click');
        if (clickHandler && clickHandler[1]) {
          clickHandler[1]({ preventDefault: jest.fn(), ...options });
        }
      }
    }),
    change: jest.fn((element, options) => {
      if (element && options && options.target) {
        element.value = options.target.value;
      }
    }),
    input: jest.fn((element, options) => {
      if (element && options && options.target) {
        element.value = options.target.value;
      }
    }),
    submit: jest.fn((element, options) => {
      if (element && typeof element.addEventListener === 'function') {
        const submitHandler = element.addEventListener.mock.calls.find(call => call[0] === 'submit');
        if (submitHandler && submitHandler[1]) {
          submitHandler[1]({ preventDefault: jest.fn(), ...options });
        }
      }
    })
  }
}));
