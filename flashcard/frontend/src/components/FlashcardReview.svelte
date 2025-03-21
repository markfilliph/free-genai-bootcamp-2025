<script>
    import { createEventDispatcher } from 'svelte';
    
    export let flashcard;
    export let showAnswer = false;
    
    const dispatch = createEventDispatcher();
    
    function toggleAnswer() {
        showAnswer = !showAnswer;
    }
    
    function rateCard(rating) {
        dispatch('rate', { 
            flashcardId: flashcard.id, 
            rating 
        });
        showAnswer = false;
    }
</script>

<div class="flashcard">
    <div class="card-content">
        <div class="word">{flashcard.word}</div>
        
        {#if showAnswer}
            <div class="answer">
                <div class="example-sentence">{flashcard.example_sentence}</div>
                <div class="translation">{flashcard.translation}</div>
                
                {#if flashcard.conjugation}
                    <div class="conjugation">{flashcard.conjugation}</div>
                {/if}
                
                {#if flashcard.cultural_note}
                    <div class="cultural-note">{flashcard.cultural_note}</div>
                {/if}
                
                <div class="rating-buttons">
                    <button on:click={() => rateCard(1)} class="rating-btn difficult">Difficult</button>
                    <button on:click={() => rateCard(2)} class="rating-btn good">Good</button>
                    <button on:click={() => rateCard(3)} class="rating-btn easy">Easy</button>
                </div>
            </div>
        {:else}
            <button on:click={toggleAnswer} class="show-answer-btn">Show Answer</button>
        {/if}
    </div>
</div>

<style>
    .flashcard {
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 20px;
        margin-bottom: 20px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    
    .word {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 15px;
    }
    
    .example-sentence {
        font-style: italic;
        margin-bottom: 10px;
    }
    
    .translation {
        margin-bottom: 10px;
    }
    
    .conjugation, .cultural-note {
        background-color: #f5f5f5;
        padding: 10px;
        border-radius: 4px;
        margin-bottom: 10px;
    }
    
    .rating-buttons {
        display: flex;
        justify-content: space-between;
        margin-top: 20px;
    }
    
    .rating-btn {
        padding: 8px 16px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
    
    .difficult {
        background-color: #f44336;
        color: white;
    }
    
    .good {
        background-color: #ffeb3b;
    }
    
    .easy {
        background-color: #4caf50;
        color: white;
    }
    
    .show-answer-btn {
        background-color: #2196f3;
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 16px;
    }
</style>
