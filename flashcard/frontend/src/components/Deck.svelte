<script>
    import { Link } from 'svelte-routing';
    import { formatDate } from '../lib/utils.js';
    
    export let deck;
    export let showActions = true;
    
    // Format the created_at date
    $: formattedDate = formatDate(deck.created_at);
</script>

<div class="deck-card" data-testid="deck-card">
    <div class="deck-info">
        <h3 class="deck-name">{deck.name}</h3>
        <p class="deck-date">Created: {formattedDate}</p>
    </div>
    
    {#if showActions}
        <div class="deck-actions">
            <Link to={`/decks/${deck.id}`} class="view-btn">View</Link>
            <Link to={`/decks/${deck.id}/study`} class="study-btn">Study</Link>
            <Link to={`/decks/${deck.id}/edit`} class="edit-btn">Edit</Link>
        </div>
    {/if}
</div>

<style>
    .deck-card {
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 15px;
        margin-bottom: 15px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        background-color: #fff;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    
    .deck-name {
        margin: 0 0 5px 0;
        font-size: 18px;
    }
    
    .deck-date {
        margin: 0;
        color: #666;
        font-size: 14px;
    }
    
    .deck-actions {
        display: flex;
        gap: 10px;
    }
    
    .view-btn, .study-btn, .edit-btn {
        padding: 5px 10px;
        border-radius: 4px;
        text-decoration: none;
        font-size: 14px;
    }
    
    .view-btn {
        background-color: #2196f3;
        color: white;
    }
    
    .study-btn {
        background-color: #4caf50;
        color: white;
    }
    
    .edit-btn {
        background-color: #ff9800;
        color: white;
    }
</style>
