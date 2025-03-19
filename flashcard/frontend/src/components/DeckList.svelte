<script>
    export let decks = [];
</script>

<ul>
    {#each decks as deck}
        <li>
            <a href="/decks/{deck.id}">{deck.title}</a>
        </li>
    {/each}
</ul>
