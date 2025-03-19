<script>
    export let deckId;
    let question = '';
    let answer = '';
</script>

<form on:submit|preventDefault>
    <input bind:value={question} placeholder="Question" class="input-field" />
    <textarea bind:value={answer} placeholder="Answer" class="input-field" />
    <button type="submit" class="submit-btn">Save Card</button>
</form>
