<script>
    import { onMount } from 'svelte';
    import { apiFetch } from '../lib/api.js';
    import DeckList from '../components/DeckList.svelte';
    
    let decks = [];
    let newDeckName = '';
    let loading = true;
    let error = null;
    
    onMount(async () => {
        try {
            loading = true;
            decks = await apiFetch('/decks');
        } catch (err) {
            error = err.message;
        } finally {
            loading = false;
        }
    });
    
    async function createDeck() {
        if (!newDeckName.trim()) return;
        
        try {
            const newDeck = await apiFetch('/decks', {
                method: 'POST',
                body: JSON.stringify({ name: newDeckName })
            });
            
            decks = [...decks, newDeck];
            newDeckName = '';
        } catch (err) {
            error = err.message;
        }
    }
</script>

<h1>Manage Your Decks</h1>

{#if loading}
    <p>Loading decks...</p>
{:else if error}
    <p class="error">{error}</p>
{:else}
    <DeckList {decks} />
{/if}

<form on:submit|preventDefault={createDeck}>
    <input bind:value={newDeckName} placeholder="New Deck Name" />
    <button type="submit">Create Deck</button>
</form>

<style>
    .error {
        color: red;
    }
</style>
