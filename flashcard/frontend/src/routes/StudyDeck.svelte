<script>
    import { onMount } from 'svelte';
    import { navigate } from 'svelte-routing';
    
    // Mock decks data
    const mockDecks = [
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
    
    // Mock flashcards for each deck
    const mockFlashcards = {
        '1': [
            { id: '101', frontText: 'hablar', backText: 'to speak', examples: 'Yo hablo español. (I speak Spanish.)', notes: 'Regular -ar verb' },
            { id: '102', frontText: 'comer', backText: 'to eat', examples: 'Ellos comen pizza. (They eat pizza.)', notes: 'Regular -er verb' },
            { id: '103', frontText: 'vivir', backText: 'to live', examples: 'Nosotros vivimos en Madrid. (We live in Madrid.)', notes: 'Regular -ir verb' },
            { id: '104', frontText: 'ser', backText: 'to be (permanent)', examples: 'Ella es doctora. (She is a doctor.)', notes: 'Irregular verb' },
            { id: '105', frontText: 'estar', backText: 'to be (temporary)', examples: 'Estoy cansado. (I am tired.)', notes: 'Irregular verb' }
        ],
        '2': [
            { id: '201', frontText: 'la manzana', backText: 'apple', examples: 'Me gusta comer manzanas. (I like to eat apples.)', notes: 'Feminine noun' },
            { id: '202', frontText: 'el pan', backText: 'bread', examples: 'Quiero comprar pan. (I want to buy bread.)', notes: 'Masculine noun' },
            { id: '203', frontText: 'la leche', backText: 'milk', examples: 'Bebo leche cada mañana. (I drink milk every morning.)', notes: 'Feminine noun' }
        ],
        '3': [
            { id: '301', frontText: '¿Dónde está...?', backText: 'Where is...?', examples: '¿Dónde está el baño? (Where is the bathroom?)', notes: 'Question phrase' },
            { id: '302', frontText: '¿Cuánto cuesta?', backText: 'How much does it cost?', examples: '¿Cuánto cuesta este libro? (How much does this book cost?)', notes: 'Question phrase' }
        ]
    };
    
    // Study session state
    let deckId = '';
    let deckName = '';
    let flashcards = [];
    let currentIndex = 0;
    let showAnswer = false;
    let loading = true;
    let error = null;
    let studyComplete = false;
    let correctCount = 0;
    let incorrectCount = 0;
    
    // Card status tracking
    let cardStatuses = [];
    
    onMount(() => {
        // Get deck ID from URL
        const urlParams = new URLSearchParams(window.location.search);
        const urlDeckId = urlParams.get('deckId');
        
        if (urlDeckId) {
            deckId = urlDeckId;
            
            // Find the deck
            const deck = mockDecks.find(d => d.id === deckId);
            if (deck) {
                deckName = deck.name;
                
                // Get flashcards for this deck
                if (mockFlashcards[deckId]) {
                    flashcards = [...mockFlashcards[deckId]];
                    
                    // Initialize card statuses
                    cardStatuses = flashcards.map(() => ({ 
                        seen: false, 
                        correct: false,
                        incorrect: false 
                    }));
                    
                    // Mark first card as seen
                    if (flashcards.length > 0) {
                        cardStatuses[0].seen = true;
                    }
                    
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
    });
    
    function flipCard() {
        showAnswer = !showAnswer;
    }
    
    function markCard(correct) {
        // Update card status
        cardStatuses[currentIndex].seen = true;
        
        if (correct) {
            cardStatuses[currentIndex].correct = true;
            correctCount++;
        } else {
            cardStatuses[currentIndex].incorrect = true;
            incorrectCount++;
        }
        
        // Move to next card
        showAnswer = false;
        
        if (currentIndex < flashcards.length - 1) {
            currentIndex++;
            cardStatuses[currentIndex].seen = true;
        } else {
            // Study session complete
            studyComplete = true;
        }
    }
    
    function restartStudy() {
        currentIndex = 0;
        showAnswer = false;
        studyComplete = false;
        correctCount = 0;
        incorrectCount = 0;
        
        // Reset card statuses
        cardStatuses = flashcards.map(() => ({ 
            seen: false, 
            correct: false,
            incorrect: false 
        }));
        
        // Mark first card as seen
        if (flashcards.length > 0) {
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
                    <div class="stat-value">{flashcards.length}</div>
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
                <div class="progress-fill" style="width: {(currentIndex / flashcards.length) * 100}%"></div>
            </div>
            <div class="progress-text">
                Card {currentIndex + 1} of {flashcards.length}
            </div>
        </div>
        
        <div class="flashcard-container">
            <div class="flashcard" class:flipped={showAnswer} on:click={flipCard} on:keydown={(e) => e.key === 'Enter' && flipCard()} tabindex="0" role="button" aria-label="Flip flashcard">
                <div class="card-side front">
                    <div class="card-content">
                        <p class="card-text">{flashcards[currentIndex].frontText}</p>
                        <p class="card-hint">Click to flip</p>
                    </div>
                </div>
                <div class="card-side back">
                    <div class="card-content">
                        <p class="card-text">{flashcards[currentIndex].backText}</p>
                        
                        {#if flashcards[currentIndex].examples}
                            <div class="card-examples">
                                <h4>Examples:</h4>
                                <p>{flashcards[currentIndex].examples}</p>
                            </div>
                        {/if}
                        
                        {#if flashcards[currentIndex].notes}
                            <div class="card-notes">
                                <h4>Notes:</h4>
                                <p>{flashcards[currentIndex].notes}</p>
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
        
        {#if showAnswer}
            <div class="rating-buttons">
                <button class="incorrect-button" on:click={() => markCard(false)}>
                    Incorrect
                </button>
                <button class="correct-button" on:click={() => markCard(true)}>
                    Correct
                </button>
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
    
    .flashcard-container {
        perspective: 1000px;
        margin-bottom: 2rem;
    }
    
    .flashcard {
        position: relative;
        width: 100%;
        height: 300px;
        cursor: pointer;
        transform-style: preserve-3d;
        transition: transform 0.6s;
    }
    
    .flashcard.flipped {
        transform: rotateY(180deg);
    }
    
    .card-side {
        position: absolute;
        width: 100%;
        height: 100%;
        backface-visibility: hidden;
        border-radius: 8px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 2rem;
    }
    
    .front {
        background-color: #ffffff;
        border: 2px solid #007bff;
    }
    
    .back {
        background-color: #ffffff;
        border: 2px solid #28a745;
        transform: rotateY(180deg);
    }
    
    .card-content {
        width: 100%;
        text-align: center;
    }
    
    .card-text {
        font-size: 2rem;
        margin-bottom: 1rem;
        color: #333;
    }
    
    .card-hint {
        font-size: 0.9rem;
        color: #6c757d;
        font-style: italic;
    }
    
    .card-examples, .card-notes {
        margin-top: 1.5rem;
        text-align: left;
        font-size: 1rem;
    }
    
    .card-examples h4, .card-notes h4 {
        margin: 0 0 0.5rem 0;
        color: #6c757d;
        font-size: 1rem;
    }
    
    .rating-buttons {
        display: flex;
        justify-content: center;
        gap: 2rem;
        margin-bottom: 2rem;
    }
    
    .correct-button, .incorrect-button {
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        border: none;
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
            height: 250px;
        }
        
        .card-text {
            font-size: 1.5rem;
        }
        
        .study-stats {
            flex-direction: column;
            gap: 1rem;
        }
    }
</style>
