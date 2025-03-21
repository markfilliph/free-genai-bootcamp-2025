<script>
    import { createEventDispatcher } from 'svelte';
    
    export let deckId;
    let question = '';
    let answer = '';
    
    const dispatch = createEventDispatcher();
    
    function handleSubmit() {
        if (question && answer) {
            dispatch('submit', {
                deckId,
                question,
                answer
            });
            
            // Reset form
            question = '';
            answer = '';
        }
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <input bind:value={question} placeholder="Question" class="input-field" />
    <textarea bind:value={answer} placeholder="Answer" class="input-field" />
    <button type="submit" class="submit-btn">Save Card</button>
</form>
