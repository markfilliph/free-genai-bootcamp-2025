<script>
    import { navigate } from 'svelte-routing';
    import { onMount } from 'svelte';
    import { decks, flashcards } from '../lib/stores.js';
    import * as api from '../lib/api.js';
    
    // Flashcard data
    let frontText = '';
    let backText = '';
    let examples = '';
    let notes = '';
    let deckId = '';
    let isLoading = false;
    let error = null;
    let success = false;
    let selectedDeckName = '';
    let createdFlashcards = [];
    let decksList = [];
    
    // AI generation states
    let isGeneratingExamples = false;
    let isGeneratingNotes = false;
    let isVerb = false; // Toggle for verb conjugation
    let generationError = null;
    
    // Function to refresh decks list from localStorage
    function refreshDecksList() {
        decksList = decks.getAllDecks();
        console.log('Refreshed decks in CreateFlashcards:', decksList);
    }
    
    // Get the initial decks list
    refreshDecksList();
    
    // Subscribe to the decks store for updates
    const unsubscribeDecks = decks.subscribe(value => {
        decksList = [...value]; // Create a new array to ensure reactivity
        console.log('Decks updated in CreateFlashcards:', decksList);
    });
    
    onMount(() => {
        // Force refresh the decks list from localStorage
        refreshDecksList();
        
        // Check if there's a deckId in the URL
        const urlParams = new URLSearchParams(window.location.search);
        const urlDeckId = urlParams.get('deckId');
        
        if (urlDeckId) {
            deckId = urlDeckId;
            
            // Find the deck directly from localStorage
            const allDecks = decks.getAllDecks();
            const deck = allDecks.find(d => d.id === urlDeckId);
            
            if (deck) {
                selectedDeckName = deck.name;
                console.log('Found deck in localStorage:', deck);
            } else {
                console.log('Deck not found in localStorage:', urlDeckId);
            }
        }
        
        // Clean up subscriptions on component destruction
        return () => {
            unsubscribeDecks();
        };
    });
    
    // Form validation
    function validateForm() {
        if (!frontText.trim()) {
            error = 'Front text is required';
            return false;
        }
        if (!backText.trim()) {
            error = 'Back text is required';
            return false;
        }
        if (!deckId) {
            error = 'Please select a deck';
            return false;
        }
        return true;
    }
    
    // AI Generation functions
    async function generateExamples() {
        if (!frontText) {
            generationError = 'Please enter a Spanish word first';
            return;
        }
        
        generationError = null;
        isGeneratingExamples = true;
        
        try {
            const generatedExamples = await api.generateExampleSentences(frontText);
            if (generatedExamples && generatedExamples.length > 0) {
                examples = generatedExamples.join('\n');
            } else {
                generationError = 'No examples generated. Please try again.';
            }
        } catch (err) {
            console.error('Error generating examples:', err);
            generationError = `Failed to generate examples: ${err.message}`;
        } finally {
            isGeneratingExamples = false;
        }
    }
    
    async function generateCulturalNotes() {
        if (!frontText) {
            generationError = 'Please enter a Spanish word first';
            return;
        }
        
        generationError = null;
        isGeneratingNotes = true;
        
        try {
            const culturalNote = await api.generateCulturalNote(frontText);
            if (culturalNote) {
                notes = culturalNote;
            } else {
                generationError = 'No cultural notes generated. Please try again.';
            }
        } catch (err) {
            console.error('Error generating cultural notes:', err);
            generationError = `Failed to generate cultural notes: ${err.message}`;
        } finally {
            isGeneratingNotes = false;
        }
    }
    
    async function generateVerbConjugations() {
        if (!frontText) {
            generationError = 'Please enter a Spanish verb first';
            return;
        }
        
        generationError = null;
        isGeneratingNotes = true;
        
        try {
            const conjugations = await api.generateVerbConjugations(frontText);
            if (conjugations) {
                notes = `Verb Conjugations:\n${conjugations}`;
            } else {
                generationError = 'No conjugations generated. Please try again.';
            }
        } catch (err) {
            console.error('Error generating conjugations:', err);
            generationError = `Failed to generate conjugations: ${err.message}`;
        } finally {
            isGeneratingNotes = false;
        }
    }
    
    // Create a new flashcard
    async function createFlashcard() {
        error = null;
        success = false;
        
        if (!validateForm()) return;
        
        isLoading = true;
        
        try {
            console.log('Creating flashcard for deck:', deckId);
            
            // Create a new flashcard using the store
            const newCard = flashcards.addFlashcard(deckId, {
                frontText,
                backText,
                examples,
                notes
            });
            
            console.log('New flashcard created:', newCard);
            
            // Get the latest flashcards for this deck
            createdFlashcards = flashcards.getFlashcardsByDeck(deckId);
            console.log('Updated flashcards for deck:', createdFlashcards);
            
            // Show success message
            success = true;
            
            // Reset form fields but keep the deck selection
            setTimeout(() => {
                frontText = '';
                backText = '';
                examples = '';
                notes = '';
                isVerb = false;
                success = false;
            }, 2000);
        } catch (err) {
            error = 'Failed to create flashcard. Please try again.';
            console.error('Error creating flashcard:', err);
        } finally {
            isLoading = false;
        }
    }
    
    // Update selected deck name when deck ID changes
    $: {
        if (deckId) {
            const selectedDeck = decksList.find(d => d.id === deckId);
            if (selectedDeck) {
                selectedDeckName = selectedDeck.name;
            }
        } else {
            selectedDeckName = '';
        }
    }
    
    // Create a new deck
    function goToDeckManagement() {
        navigate('/decks');
    }
</script>

<div class="create-flashcards">
    <div class="header">
        <h1>
            {#if selectedDeckName}
                Add Cards to "{selectedDeckName}"
            {:else}
                Create New Flashcards
            {/if}
        </h1>
        <button class="secondary-button" on:click={goToDeckManagement}>
            Manage Decks
        </button>
    </div>
    
    <div class="flashcard-form-container">
        <form on:submit|preventDefault={createFlashcard}>
            <div class="form-group">
                <label for="deck-select">Select Deck</label>
                <select id="deck-select" bind:value={deckId} required>
                    <option value="">-- Select a Deck --</option>
                    {#each decksList as deck}
                        <option value={deck.id}>{deck.name}</option>
                    {/each}
                </select>
            </div>
            
            <div class="flashcard-preview">
                <div class="card-side front" class:has-content={frontText}>
                    <div class="card-content">
                        {#if frontText}
                            <p>{frontText}</p>
                        {:else}
                            <p class="placeholder">Front Text</p>
                        {/if}
                    </div>
                </div>
                <div class="card-side back" class:has-content={backText}>
                    <div class="card-content">
                        {#if backText}
                            <p>{backText}</p>
                            {#if examples}
                                <div class="examples">
                                    <h4>Examples:</h4>
                                    <p>{examples}</p>
                                </div>
                            {/if}
                            {#if notes}
                                <div class="notes">
                                    <h4>Notes:</h4>
                                    <p>{notes}</p>
                                </div>
                            {/if}
                        {:else}
                            <p class="placeholder">Back Text</p>
                        {/if}
                    </div>
                </div>
            </div>
            
            <div class="form-columns">
                <div class="form-column">
                    <div class="form-group">
                        <label for="front-text">Front Text (Spanish)</label>
                        <input 
                            id="front-text" 
                            bind:value={frontText} 
                            placeholder="e.g., hablar" 
                            required
                        />
                    </div>
                    
                    <div class="form-group">
                        <label for="back-text">Back Text (English)</label>
                        <input 
                            id="back-text" 
                            bind:value={backText} 
                            placeholder="e.g., to speak" 
                            required
                        />
                    </div>
                </div>
                
                <div class="form-column">
                    <div class="form-group">
                        <div class="label-with-button">
                            <label for="examples">Examples (Optional)</label>
                            <button 
                                type="button" 
                                class="ai-generate-button" 
                                on:click={generateExamples}
                                disabled={isGeneratingExamples || !frontText}
                            >
                                {#if isGeneratingExamples}
                                    Generating...
                                {:else}
                                    Generate with AI
                                {/if}
                            </button>
                        </div>
                        <textarea 
                            id="examples" 
                            bind:value={examples} 
                            placeholder="e.g., Yo hablo español. (I speak Spanish.)"
                            rows="3"
                        ></textarea>
                    </div>
                    
                    <div class="form-group">
                        <div class="label-with-button">
                            <label for="notes">Notes (Optional)</label>
                            <div class="ai-buttons-group">
                                <button 
                                    type="button" 
                                    class="ai-generate-button" 
                                    on:click={generateCulturalNotes}
                                    disabled={isGeneratingNotes || !frontText}
                                >
                                    {#if isGeneratingNotes}
                                        Generating...
                                    {:else}
                                        Cultural Notes
                                    {/if}
                                </button>
                                <button 
                                    type="button" 
                                    class="ai-generate-button" 
                                    on:click={generateVerbConjugations}
                                    disabled={isGeneratingNotes || !frontText}
                                >
                                    Verb Conjugations
                                </button>
                            </div>
                        </div>
                        <textarea 
                            id="notes" 
                            bind:value={notes} 
                            placeholder="e.g., Regular -ar verb"
                            rows="3"
                        ></textarea>
                    </div>
                </div>
            </div>
            
            {#if error}
                <div class="error-message" role="alert">
                    {error}
                </div>
            {/if}
            
            {#if generationError}
                <div class="error-message" role="alert">
                    <strong>AI Generation Error:</strong> {generationError}
                </div>
            {/if}
            
            {#if success}
                <div class="success-message" role="status">
                    Flashcard created successfully!
                </div>
            {/if}
            
            <div class="form-actions">
                <button type="submit" class="primary-button" disabled={isLoading}>
                    {isLoading ? 'Creating...' : 'Create Flashcard'}
                </button>
                <button type="button" class="secondary-button" on:click={() => {
                    frontText = '';
                    backText = '';
                    examples = '';
                    notes = '';
                    error = null;
                }}>
                    Clear Form
                </button>
            </div>
            
            {#if createdFlashcards.length > 0}
                <div class="created-flashcards">
                    <h3>Created Flashcards</h3>
                    <div class="flashcards-list">
                        {#each createdFlashcards as flashcard}
                            <div class="flashcard-item">
                                <div class="flashcard-content">
                                    <div class="flashcard-front"><p>{flashcard.frontText}</p></div>
                                    <div class="flashcard-back">
                                        <p>{flashcard.backText}</p>
                                        {#if flashcard.examples}
                                            <div class="examples">
                                                <h4>Examples:</h4>
                                                <p>{flashcard.examples}</p>
                                            </div>
                                        {/if}
                                        {#if flashcard.notes}
                                            <div class="notes">
                                                <h4>Notes:</h4>
                                                <p>{flashcard.notes}</p>
                                            </div>
                                        {/if}
                                    </div>
                                </div>
                                <div class="flashcard-meta">
                                    Created: {new Date(flashcard.createdAt).toLocaleString()}
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </form>
    </div>
</div>

<style>
    .create-flashcards {
        max-width: 1000px;
        margin: 0 auto;
        padding: 1rem;
    }
    
    .created-flashcards {
        margin-top: 2rem;
        background-color: #f8f9fa;
        border-radius: 8px;
        padding: 1.5rem;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .created-flashcards h3 {
        margin-top: 0;
        margin-bottom: 1rem;
        color: #333;
    }
    
    .flashcards-list {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
        gap: 1rem;
    }
    
    .flashcard-item {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        padding: 1rem;
        transition: all 0.3s ease;
        min-height: 100px;
        max-height: 300px;
        overflow-y: auto;
    }
    
    .flashcard-item p {
        margin: 0;
        line-height: 1.2;
    }
    
    .flashcard-item:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    
    .flashcard-content {
        display: flex;
        justify-content: space-between;
        margin-bottom: 0.5rem;
    }
    
    .flashcard-front, .flashcard-back {
        flex: 1;
        padding: 0.5rem;
    }
    
    .flashcard-front {
        font-weight: bold;
        color: #007bff;
        border-right: 1px solid #eee;
    }
    
    .flashcard-meta {
        font-size: 0.8rem;
        color: #6c757d;
        text-align: right;
        margin-top: 0.5rem;
    }
    
    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }
    
    h1 {
        margin: 0;
        color: #333;
    }
    
    .flashcard-form-container {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        padding: 2rem;
    }
    
    .form-group {
        margin-bottom: 1.5rem;
    }
    
    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #333;
    }
    
    input, textarea, select {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
    }
    
    textarea {
        max-height: 150px;
        overflow-y: auto;
        resize: vertical;
    }
    
    input:focus, textarea:focus, select:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
    }
    
    .form-columns {
        display: flex;
        gap: 2rem;
        margin-bottom: 1.5rem;
    }
    
    .form-column {
        flex: 1;
    }
    
    .flashcard-preview {
        display: flex;
        gap: 2rem;
        margin-bottom: 2rem;
    }
    
    .card-side {
        flex: 1;
        min-height: 200px;
        max-height: 300px;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        display: flex;
        align-items: flex-start;
        justify-content: center;
        padding: 1rem;
        text-align: center;
        transition: all 0.3s ease;
        overflow-y: auto;
    }
    
    .front {
        background-color: #f8f9fa;
        border: 2px dashed #ddd;
    }
    
    .back {
        background-color: #f8f9fa;
        border: 2px dashed #ddd;
    }
    
    .has-content {
        border: 2px solid #007bff;
        background-color: white;
    }
    
    .card-content {
        width: 100%;
    }
    
    .placeholder {
        color: #adb5bd;
        font-style: italic;
    }
    
    .examples, .notes {
        margin-top: 0.5rem;
        font-size: 0.9rem;
        text-align: left;
        max-height: 120px;
        overflow-y: auto;
        padding-right: 5px;
    }
    
    .examples p, .notes p {
        margin: 0;
        line-height: 1.2;
        white-space: pre-line;
    }
    
    .examples h4, .notes h4 {
        margin: 0 0 0.1rem 0;
        color: #6c757d;
        font-size: 0.8rem;
    }
    
    .form-actions {
        display: flex;
        gap: 1rem;
        margin-top: 1.5rem;
    }
    
    .primary-button, .secondary-button {
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        border: none;
    }
    
    .primary-button {
        background: #007bff;
        color: white;
    }
    
    .primary-button:hover {
        background: #0069d9;
    }
    
    .primary-button:disabled {
        background: #6c757d;
        cursor: not-allowed;
    }
    
    .secondary-button {
        background: #f8f9fa;
        color: #333;
        border: 1px solid #ddd;
    }
    
    .secondary-button:hover {
        background: #e2e6ea;
    }
    
    .error-message {
        background: #f8d7da;
        color: #721c24;
        padding: 0.75rem;
        border-radius: 4px;
        margin-bottom: 1rem;
    }
    
    .success-message {
        background: #d4edda;
        color: #155724;
        padding: 0.75rem;
        border-radius: 4px;
        margin-bottom: 1rem;
    }
    
    @media (max-width: 768px) {
        .form-columns, .flashcard-preview {
            flex-direction: column;
            gap: 1rem;
        }
        
        .card-side {
            height: 150px;
        }
    }
    
    /* AI Generation Buttons Styles */
    .label-with-button {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 5px;
    }
    
    .ai-generate-button {
        font-size: 0.8rem;
        padding: 4px 8px;
        background-color: #4a6fa5;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }
    
    .ai-generate-button:hover:not(:disabled) {
        background-color: #3a5a8a;
    }
    
    .ai-generate-button:disabled {
        background-color: #cccccc;
        cursor: not-allowed;
    }
    
    .ai-buttons-group {
        display: flex;
        gap: 5px;
    }
</style>
