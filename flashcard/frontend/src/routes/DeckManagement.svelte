<script>
    import { onMount } from 'svelte';
    import { navigate, Link } from 'svelte-routing';
    import { apiFetch } from '../lib/api.js';
    import DeckList from '../components/DeckList.svelte';
    
    // Mock decks data
    let mockDecks = [
        {
            id: '1',
            name: 'Spanish Verbs',
            description: 'Common Spanish verbs with conjugations',
            cardCount: 25,
            lastReviewed: '2025-03-20'
        },
        {
            id: '2',
            name: 'Food Vocabulary',
            description: 'Words related to food and dining',
            cardCount: 18,
            lastReviewed: '2025-03-22'
        },
        {
            id: '3',
            name: 'Travel Phrases',
            description: 'Useful phrases for traveling',
            cardCount: 15,
            lastReviewed: '2025-03-18'
        }
    ];
    
    let decks = [];
    let newDeckName = '';
    let newDeckDescription = '';
    let loading = true;
    let error = null;
    let showAddDeckForm = false;
    
    // Load data on mount
    onMount(() => {
        // Load mock decks after a short delay to simulate API call
        setTimeout(() => {
            decks = [...mockDecks];
            loading = false;
        }, 800);
    });
    
    function toggleAddDeckForm() {
        showAddDeckForm = !showAddDeckForm;
        if (!showAddDeckForm) {
            // Reset form when hiding
            newDeckName = '';
            newDeckDescription = '';
        }
    }
    
    async function createDeck() {
        if (!newDeckName.trim()) return;
        
        try {
            // Create a new mock deck
            const newDeck = {
                id: Date.now().toString(),
                name: newDeckName,
                description: newDeckDescription || 'No description provided',
                cardCount: 0,
                lastReviewed: 'Never'
            };
            
            // Add to decks list
            decks = [...decks, newDeck];
            
            // Reset form
            newDeckName = '';
            newDeckDescription = '';
            showAddDeckForm = false;
        } catch (err) {
            error = err.message;
        }
    }
    
    function deleteDeck(deckId) {
        decks = decks.filter(deck => deck.id !== deckId);
    }
    
    function addCardsToDeck(deckId) {
        // Navigate to the create flashcards page with the deck ID
        navigate(`/create?deckId=${deckId}`);
    }
    
    function studyDeck(deckId) {
        // Navigate to the study deck page with the deck ID
        navigate(`/study?deckId=${deckId}`);
    }
</script>

<div class="deck-management">
    <div class="header">
        <h1>My Flashcard Decks</h1>
        <button class="add-deck-button" on:click={toggleAddDeckForm}>
            {showAddDeckForm ? 'Cancel' : '+ Add New Deck'}
        </button>
    </div>
    
    {#if showAddDeckForm}
        <div class="add-deck-form">
            <h2>Create New Deck</h2>
            <form on:submit|preventDefault={createDeck}>
                <div class="form-group">
                    <label for="deck-name">Deck Name</label>
                    <input 
                        id="deck-name"
                        bind:value={newDeckName} 
                        placeholder="Enter deck name" 
                        required
                    />
                </div>
                
                <div class="form-group">
                    <label for="deck-description">Description (optional)</label>
                    <textarea 
                        id="deck-description"
                        bind:value={newDeckDescription} 
                        placeholder="Enter a description for this deck"
                        rows="3"
                    ></textarea>
                </div>
                
                <div class="form-actions">
                    <button type="button" class="cancel-button" on:click={toggleAddDeckForm}>Cancel</button>
                    <button type="submit" class="create-button">Create Deck</button>
                </div>
            </form>
        </div>
    {/if}
    
    {#if loading}
        <div class="loading">
            <p>Loading your decks...</p>
            <div class="spinner"></div>
        </div>
    {:else if error}
        <div class="error-message">
            <p>Error loading decks: {error}</p>
            <button on:click={() => { loading = true; error = null; }}>Try Again</button>
        </div>
    {:else if decks.length === 0}
        <div class="empty-state">
            <h2>You don't have any decks yet</h2>
            <p>Create your first flashcard deck to start learning!</p>
            <button class="add-deck-button" on:click={toggleAddDeckForm}>+ Create Your First Deck</button>
        </div>
    {:else}
        <div class="decks-grid">
            {#each decks as deck (deck.id)}
                <div class="deck-card">
                    <div class="deck-info">
                        <h3>{deck.name}</h3>
                        <p class="description">{deck.description}</p>
                        <div class="deck-meta">
                            <span class="card-count">{deck.cardCount} cards</span>
                            <span class="last-reviewed">Last studied: {deck.lastReviewed}</span>
                        </div>
                    </div>
                    <div class="deck-actions">
                        <button class="study-button" on:click={() => studyDeck(deck.id)}>Study</button>
                        <button class="add-cards-button" on:click={() => addCardsToDeck(deck.id)}>Add Cards</button>
                        <button class="delete-button" on:click={() => deleteDeck(deck.id)}>Delete</button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .deck-management {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem;
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
    
    .add-deck-button {
        background: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        padding: 0.6rem 1.2rem;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.3s ease;
    }
    
    .add-deck-button:hover {
        background: #0069d9;
    }
    
    .add-deck-form {
        background: #f8f9fa;
        border-radius: 8px;
        padding: 1.5rem;
        margin-bottom: 2rem;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    
    .add-deck-form h2 {
        margin-top: 0;
        margin-bottom: 1.5rem;
        color: #333;
    }
    
    .form-group {
        margin-bottom: 1rem;
    }
    
    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    
    input, textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
    }
    
    input:focus, textarea:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
    }
    
    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: 1rem;
        margin-top: 1.5rem;
    }
    
    .cancel-button {
        background: #f8f9fa;
        color: #333;
        border: 1px solid #ddd;
    }
    
    .cancel-button:hover {
        background: #e2e6ea;
    }
    
    .create-button {
        background: #007bff;
        color: white;
    }
    
    .create-button:hover {
        background: #0069d9;
    }
    
    button {
        padding: 0.6rem 1.2rem;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        border: none;
    }
    
    .loading {
        text-align: center;
        padding: 3rem 0;
    }
    
    .spinner {
        display: inline-block;
        width: 40px;
        height: 40px;
        border: 4px solid rgba(0, 123, 255, 0.3);
        border-radius: 50%;
        border-top-color: #007bff;
        animation: spin 1s ease-in-out infinite;
        margin-top: 1rem;
    }
    
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
    
    .error-message {
        text-align: center;
        padding: 2rem;
        background: #f8d7da;
        border-radius: 8px;
        color: #721c24;
    }
    
    .empty-state {
        text-align: center;
        padding: 3rem 1rem;
        background: #f8f9fa;
        border-radius: 8px;
        margin-top: 2rem;
    }
    
    .empty-state h2 {
        margin-bottom: 1rem;
        color: #333;
    }
    
    .empty-state p {
        margin-bottom: 2rem;
        color: #6c757d;
    }
    
    .decks-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
    }
    
    .deck-card {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        overflow: hidden;
        transition: transform 0.3s ease, box-shadow 0.3s ease;
    }
    
    .deck-card:hover {
        transform: translateY(-5px);
        box-shadow: 0 5px 15px rgba(0,0,0,0.1);
    }
    
    .deck-info {
        padding: 1.5rem;
    }
    
    .deck-info h3 {
        margin-top: 0;
        margin-bottom: 0.5rem;
        color: #333;
    }
    
    .description {
        color: #6c757d;
        margin-bottom: 1rem;
        font-size: 0.9rem;
    }
    
    .deck-meta {
        display: flex;
        justify-content: space-between;
        font-size: 0.8rem;
        color: #6c757d;
    }
    
    .deck-actions {
        display: flex;
        border-top: 1px solid #eee;
    }
    
    .deck-actions button {
        flex: 1;
        padding: 0.75rem;
        border: none;
        background: #f8f9fa;
        cursor: pointer;
        transition: background 0.3s ease;
    }
    
    .study-button {
        color: #28a745;
    }
    
    .study-button:hover {
        background: #e2f0e7;
    }
    
    .add-cards-button {
        color: #007bff;
        border-left: 1px solid #eee;
        border-right: 1px solid #eee;
    }
    
    .add-cards-button:hover {
        background: #e6f2ff;
    }
    
    .delete-button {
        color: #dc3545;
    }
    
    .delete-button:hover {
        background: #f8d7da;
    }
    
    @media (max-width: 768px) {
        .header {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }
        
        .decks-grid {
            grid-template-columns: 1fr;
        }
    }

</style>
