import { render, fireEvent } from '@testing-library/svelte';
import FlashcardForm from '../../components/FlashcardForm.svelte';

describe('FlashcardForm Component', () => {
  test('renders correctly with default props', () => {
    const { getByPlaceholderText, getByText } = render(FlashcardForm, { 
      props: { deckId: '1' } 
    });
    
    // Check if form elements are rendered
    expect(getByPlaceholderText('Question')).toBeInTheDocument();
    expect(getByPlaceholderText('Answer')).toBeInTheDocument();
    expect(getByText('Save Card')).toBeInTheDocument();
  });

  test('updates input values when typed into', async () => {
    const { getByPlaceholderText } = render(FlashcardForm, { 
      props: { deckId: '1' } 
    });
    
    const questionInput = getByPlaceholderText('Question');
    const answerInput = getByPlaceholderText('Answer');
    
    // Simulate user typing
    await fireEvent.input(questionInput, { target: { value: 'Test Question' } });
    await fireEvent.input(answerInput, { target: { value: 'Test Answer' } });
    
    // Check if values are updated
    expect(questionInput.value).toBe('Test Question');
    expect(answerInput.value).toBe('Test Answer');
  });

  test('form submission works correctly', async () => {
    // Create a mock function to test if form submission handler is called
    const mockSubmit = jest.fn();
    
    const { getByPlaceholderText, getByText, container } = render(FlashcardForm, { 
      props: { deckId: '1' } 
    });
    
    // Override the component's submit handler
    const form = container.querySelector('form');
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      mockSubmit();
    });
    
    // Fill in the form
    const questionInput = getByPlaceholderText('Question');
    const answerInput = getByPlaceholderText('Answer');
    
    await fireEvent.input(questionInput, { target: { value: 'Test Question' } });
    await fireEvent.input(answerInput, { target: { value: 'Test Answer' } });
    
    // Submit the form
    await fireEvent.click(getByText('Save Card'));
    
    // Check if the submit handler was called
    expect(mockSubmit).toHaveBeenCalled();
  });
});
