<script>
    import { onMount, createEventDispatcher } from 'svelte';
    import FlashcardReview from './FlashcardReview.svelte';
    
    export let flashcards = [];
    export let deckName = 'Flashcards';
    
    const dispatch = createEventDispatcher();
    
    let currentIndex = 0;
    let sessionStats = {
        total: flashcards.length,
        completed: 0,
        ratings: {
            difficult: 0,
            good: 0,
            easy: 0
        }
    };
    
    $: currentFlashcard = flashcards[currentIndex];
    $: isSessionComplete = currentIndex >= flashcards.length;
    $: progress = Math.round((sessionStats.completed / sessionStats.total) * 100) || 0;
    
    function handleRate(event) {
        const { rating } = event.detail;
        
        // Update stats based on rating
        if (rating === 1) sessionStats.ratings.difficult++;
        else if (rating === 2) sessionStats.ratings.good++;
        else if (rating === 3) sessionStats.ratings.easy++;
        
        sessionStats.completed++;
        
        // Move to next card
        currentIndex++;
        
        // If session is complete, dispatch event
        if (isSessionComplete) {
            dispatch('complete', sessionStats);
        }
    }
    
    function restartSession() {
        currentIndex = 0;
        sessionStats = {
            total: flashcards.length,
            completed: 0,
            ratings: {
                difficult: 0,
                good: 0,
                easy: 0
            }
        };
    }
</script>

<div class="study-session">
    <div class="session-header">
        <h2>{deckName}</h2>
        <div class="progress-bar">
            <div class="progress-fill" style="width: {progress}%"></div>
        </div>
        <div class="progress-text">
            {sessionStats.completed} / {sessionStats.total} cards
        </div>
    </div>
    
    <div class="session-content">
        {#if isSessionComplete}
            <div class="session-complete">
                <h3>Session Complete!</h3>
                <div class="stats">
                    <div class="stat">
                        <span class="label">Difficult:</span>
                        <span class="value">{sessionStats.ratings.difficult}</span>
                    </div>
                    <div class="stat">
                        <span class="label">Good:</span>
                        <span class="value">{sessionStats.ratings.good}</span>
                    </div>
                    <div class="stat">
                        <span class="label">Easy:</span>
                        <span class="value">{sessionStats.ratings.easy}</span>
                    </div>
                </div>
                <button on:click={restartSession} class="restart-btn">Restart Session</button>
            </div>
        {:else if currentFlashcard}
            <FlashcardReview 
                flashcard={currentFlashcard} 
                on:rate={handleRate} 
            />
        {:else}
            <p>No flashcards available for this deck.</p>
        {/if}
    </div>
</div>

<style>
    .study-session {
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
    }
    
    .session-header {
        margin-bottom: 20px;
    }
    
    .progress-bar {
        height: 10px;
        background-color: #eee;
        border-radius: 5px;
        margin: 10px 0;
        overflow: hidden;
    }
    
    .progress-fill {
        height: 100%;
        background-color: #4caf50;
        transition: width 0.3s ease;
    }
    
    .progress-text {
        text-align: right;
        font-size: 14px;
        color: #666;
    }
    
    .session-complete {
        text-align: center;
        padding: 30px;
        background-color: #f5f5f5;
        border-radius: 8px;
    }
    
    .stats {
        display: flex;
        justify-content: space-around;
        margin: 20px 0;
    }
    
    .stat {
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    
    .label {
        font-size: 14px;
        color: #666;
    }
    
    .value {
        font-size: 24px;
        font-weight: bold;
    }
    
    .restart-btn {
        background-color: #2196f3;
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 16px;
    }
</style>
