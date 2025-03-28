<script>
    import { onMount } from 'svelte';
    import { navigate } from 'svelte-routing';
    import { decks, flashcards } from '../lib/stores.js';
    
    // Study session state
    let deckId = '';
    let deckName = '';
    let deckCards = [];
    let currentIndex = 0;
    let showAnswer = false;
    let loading = true;
    let error = null;
    let studyComplete = false;
    let correctCount = 0;
    let incorrectCount = 0;
    
    // Card status tracking
    let cardStatuses = [];
    
    // Subscribe to stores
    let decksList = [];
    const unsubscribeDecks = decks.subscribe(value => {
        decksList = value;
    });
    
    onMount(() => {
        // Get deck ID from URL
        const urlParams = new URLSearchParams(window.location.search);
        const urlDeckId = urlParams.get('deckId');
        
        if (urlDeckId) {
            deckId = urlDeckId;
            
            // Find the deck in the store
            const deck = decksList.find(d => d.id === deckId);
            if (deck) {
                deckName = deck.name;
                
                // Get flashcards for this deck from the store
                deckCards = flashcards.getFlashcardsByDeck(deckId);
                
                if (deckCards.length > 0) {
                    // Initialize card statuses
                    cardStatuses = deckCards.map(() => ({ 
                        seen: false, 
                        correct: false,
                        incorrect: false 
                    }));
                    
                    // Mark first card as seen
                    cardStatuses[0].seen = true;
                    loading = false;
                } else {
                    error = 'No flashcards found for this deck.';
                    loading = false;
                }
            } else {
                error = 'Deck not found.';
                loading = false;
            }
        } else {
            error = 'No deck specified.';
            loading = false;
        }
        
        // Clean up subscriptions on component destruction
        return () => {
            unsubscribeDecks();
        };
    });
    
    function toggleFlip(e) {
        // Prevent default if it's a keyboard event
        if (e && e.preventDefault) {
            e.preventDefault();
        }
        showAnswer = !showAnswer;
        console.log('Card flipped, showAnswer:', showAnswer);
    }

    function showCardAnswer() {
        showAnswer = true;
        console.log('Showing answer, showAnswer:', showAnswer);
    }
    
    function hideCardAnswer() {
        showAnswer = false;
        console.log('Hiding answer, showAnswer:', showAnswer);
    }
    
    function markCard(correct) {
        console.log('markCard called with:', correct);
        // Update card status
        cardStatuses[currentIndex].seen = true;
        
        if (correct) {
            cardStatuses[currentIndex].correct = true;
            correctCount++;
        } else {
            cardStatuses[currentIndex].incorrect = true;
            incorrectCount++;
        }
        
        // Add a small delay before moving to the next card
        setTimeout(() => {
            // Move to next card
            showAnswer = false;
            
            if (currentIndex < deckCards.length - 1) {
                currentIndex++;
                cardStatuses[currentIndex].seen = true;
            } else {
                // Study session complete
                studyComplete = true;
                
                // Update the last studied date for the deck
                const today = new Date().toISOString().split('T')[0]; // Format: YYYY-MM-DD
                decks.updateDeck(deckId, { lastReviewed: today });
                console.log(`Updated last studied date for deck ${deckId} to ${today}`);
            }
        }, 300);
    }
    
    function restartStudy() {
        currentIndex = 0;
        showAnswer = false;
        studyComplete = false;
        correctCount = 0;
        incorrectCount = 0;
        
        // Reset card statuses
        cardStatuses = deckCards.map(() => ({ 
            seen: false, 
            correct: false,
            incorrect: false 
        }));
        
        // Mark first card as seen
        if (deckCards.length > 0) {
            cardStatuses[0].seen = true;
        }
    }
    
    function goToDeckManagement() {
        navigate('/decks');
    }
</script>

<div class="study-deck">
    <div class="header">
        <h1>Study: {deckName}</h1>
        <button class="secondary-button" on:click={goToDeckManagement}>
            Back to Decks
        </button>
    </div>
    
    {#if loading}
        <div class="loading">
            <p>Loading flashcards...</p>
            <div class="spinner"></div>
        </div>
    {:else if error}
        <div class="error-message">
            <p>{error}</p>
            <button class="primary-button" on:click={goToDeckManagement}>Return to Decks</button>
        </div>
    {:else if studyComplete}
        <div class="study-complete">
            <h2>Study Session Complete!</h2>
            
            <div class="study-stats">
                <div class="stat-item correct">
                    <div class="stat-value">{correctCount}</div>
                    <div class="stat-label">Correct</div>
                </div>
                <div class="stat-item incorrect">
                    <div class="stat-value">{incorrectCount}</div>
                    <div class="stat-label">Incorrect</div>
                </div>
                <div class="stat-item total">
                    <div class="stat-value">{deckCards.length}</div>
                    <div class="stat-label">Total Cards</div>
                </div>
            </div>
            
            <div class="completion-actions">
                <button class="primary-button" on:click={restartStudy}>Study Again</button>
                <button class="secondary-button" on:click={goToDeckManagement}>Return to Decks</button>
            </div>
        </div>
    {:else}
        <div class="study-progress">
            <div class="progress-bar">
                <div class="progress-fill" style="width: {(currentIndex / deckCards.length) * 100}%"></div>
            </div>
            <div class="progress-text">
                Card {currentIndex + 1} of {deckCards.length}
            </div>
        </div>
        
        <!-- Flashcard View (Always visible) -->
        <div class="flashcard-container">
            <div 
              class="flashcard" 
              class:flipped={showAnswer}
              on:click={toggleFlip}
              on:keydown={e => e.key === 'Enter' && toggleFlip(e)}
              tabindex="0"
              role="button"
              aria-label="Flashcard, press Enter to flip"
            >
              <!-- Front (Question) -->
              <div class="flashcard-front">
                <p class="card-text">{deckCards[currentIndex].frontText}</p>
                <button class="show-answer-button" on:click|stopPropagation={showCardAnswer}>
                  Show Answer
                </button>
              </div>
              
              <!-- Back (Answer) -->
              <div class="flashcard-back">
                <p class="card-text">{deckCards[currentIndex].backText}</p>
                {#if deckCards[currentIndex].examples}
                  <div class="card-examples">
                    <h4>Examples:</h4>
                    <p>{deckCards[currentIndex].examples}</p>
                  </div>
                {/if}
                {#if deckCards[currentIndex].notes}
                  <div class="card-examples">
                    <h4>Notes:</h4>
                    <p>{deckCards[currentIndex].notes}</p>
                  </div>
                {/if}
              </div>
            </div>
        </div>
        
        <!-- Rating Section (Only visible when answer is shown) -->
        {#if showAnswer}
            <div class="rating-section">
                <h3>How did you do?</h3>
                <div class="rating-buttons">
                    <button class="incorrect-button" on:click={() => markCard(false)}>
                        I Got It Wrong
                    </button>
                    <button class="correct-button" on:click={() => markCard(true)}>
                        I Got It Right
                    </button>
                </div>
            </div>
        {/if}
        <div class="card-navigation">
            {#each cardStatuses as status, i}
                <div 
                    class="card-indicator" 
                    class:current={i === currentIndex}
                    class:seen={status.seen}
                    class:correct={status.correct}
                    class:incorrect={status.incorrect}
                ></div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .study-deck {
        max-width: 900px;
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
    
    .study-progress {
        margin-bottom: 2rem;
    }
    
    .progress-bar {
        height: 8px;
        background-color: #e9ecef;
        border-radius: 4px;
        overflow: hidden;
        margin-bottom: 0.5rem;
    }
    
    .progress-fill {
        height: 100%;
        background-color: #007bff;
        transition: width 0.3s ease;
    }
    
    .progress-text {
        text-align: center;
        font-size: 0.9rem;
        color: #6c757d;
    }
    
    .flashcard {
        width: 100%;
        height: 350px;
        position: relative;
        transform-style: preserve-3d;
        transition: transform 0.6s;
        cursor: pointer;
        perspective: 1000px;
    }

    .flashcard.flipped {
        transform: rotateY(180deg);
    }
    
    .flashcard-front, .flashcard-back {
        position: absolute;
        width: 100%;
        height: 100%;
        backface-visibility: hidden;
        -webkit-backface-visibility: hidden; /* Safari support */
        border-radius: 8px;
        padding: 1.5rem;
        box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        background: white;
        overflow-y: auto;
    }
    
    .flashcard-back {
        transform: rotateY(180deg);
        border-left: 4px solid #28a745;
    }
    
    @keyframes flipIn {
        0% { transform: rotateY(90deg); opacity: 0; }
        100% { transform: rotateY(0deg); opacity: 1; }
    }
    
    /* Ensure proper 3D rendering in different browsers */
    .flashcard-container {
        perspective: 1000px;
        width: 100%;
        max-width: 600px;
        margin: 0 auto;
    }
    

    
    .show-answer-button {
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        padding: 0.75rem 1.5rem;
        font-size: 1rem;
        cursor: pointer;
        transition: background-color 0.3s ease;
    }
    
    .show-answer-button:hover {
        background-color: #0069d9;
    }
    

    
    .card-text {
        font-size: 2rem;
        margin-bottom: 0.75rem;
        color: #333;
        text-align: center;
        width: 100%;
    }
    

    
    .card-examples {
        margin-top: 0.75rem;
        text-align: center;
        font-size: 1rem;
        width: 100%;
        line-height: 1.3;
    }
    
    .card-examples p {
        margin: 0;
        white-space: pre-line;
    }
    
    .card-examples h4 {
        margin: 0 0 0.3rem 0;
        color: #6c757d;
        font-size: 1rem;
        font-weight: 600;
        text-align: center;
    }
    

    
    .rating-section {
        background-color: #f8f9fa;
        border-radius: 8px;
        padding: 2rem;
        text-align: center;
        margin-top: 2rem;
        position: relative;
        z-index: 20;
        animation: fadeIn 0.5s ease-in-out;
    }
    
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }
    
    .rating-section h3 {
        margin-top: 0;
        margin-bottom: 1.5rem;
        color: #333;
    }
    
    .rating-buttons {
        display: flex;
        justify-content: center;
        gap: 2rem;
    }
    
    .correct-button, .incorrect-button {
        padding: 1rem 2rem;
        border-radius: 4px;
        font-weight: 600;
        font-size: 1.1rem;
        cursor: pointer;
        transition: all 0.3s ease;
        border: none;
        position: relative;
        z-index: 30;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    
    .correct-button:not([disabled]), .incorrect-button:not([disabled]) {
        cursor: pointer;
    }
    

    
    .correct-button {
        background: #28a745;
        color: white;
    }
    
    .correct-button:hover {
        background: #218838;
    }
    
    .incorrect-button {
        background: #dc3545;
        color: white;
    }
    
    .incorrect-button:hover {
        background: #c82333;
    }
    
    .card-navigation {
        display: flex;
        justify-content: center;
        gap: 0.5rem;
        margin-top: 2rem;
    }
    
    .card-indicator {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        background-color: #e9ecef;
    }
    
    .card-indicator.seen {
        background-color: #6c757d;
    }
    
    .card-indicator.current {
        transform: scale(1.3);
    }
    
    .card-indicator.correct {
        background-color: #28a745;
    }
    
    .card-indicator.incorrect {
        background-color: #dc3545;
    }
    
    .study-complete {
        text-align: center;
        background: #f8f9fa;
        border-radius: 8px;
        padding: 2rem;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .study-complete h2 {
        margin-top: 0;
        color: #333;
    }
    
    .study-stats {
        display: flex;
        justify-content: center;
        gap: 2rem;
        margin: 2rem 0;
    }
    
    .stat-item {
        text-align: center;
        padding: 1rem;
        border-radius: 8px;
        min-width: 100px;
    }
    
    .stat-value {
        font-size: 2.5rem;
        font-weight: bold;
        margin-bottom: 0.5rem;
    }
    
    .stat-label {
        font-size: 1rem;
        color: #6c757d;
    }
    
    .stat-item.correct {
        background-color: rgba(40, 167, 69, 0.1);
    }
    
    .stat-item.correct .stat-value {
        color: #28a745;
    }
    
    .stat-item.incorrect {
        background-color: rgba(220, 53, 69, 0.1);
    }
    
    .stat-item.incorrect .stat-value {
        color: #dc3545;
    }
    
    .stat-item.total {
        background-color: rgba(0, 123, 255, 0.1);
    }
    
    .stat-item.total .stat-value {
        color: #007bff;
    }
    
    .completion-actions {
        display: flex;
        justify-content: center;
        gap: 1rem;
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
    
    .secondary-button {
        background: #f8f9fa;
        color: #333;
        border: 1px solid #ddd;
    }
    
    .secondary-button:hover {
        background: #e2e6ea;
    }
    
    .loading, .error-message {
        text-align: center;
        padding: 2rem;
    }
    
    .spinner {
        margin: 1rem auto;
        width: 40px;
        height: 40px;
        border: 4px solid rgba(0, 123, 255, 0.1);
        border-radius: 50%;
        border-top-color: #007bff;
        animation: spin 1s ease-in-out infinite;
    }
    
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
    
    @media (max-width: 768px) {
        .flashcard {
            min-height: 200px;
            height: auto;
        }
        
        .card-text {
            font-size: 1.5rem;
        }
        
        .study-stats {
            flex-direction: column;
            gap: 1rem;
        }
        
        .rating-buttons {
            flex-direction: column;
            gap: 1rem;
        }
    }
</style>
