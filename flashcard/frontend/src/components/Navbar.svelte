<script>
    import { Link } from 'svelte-routing';
</script>

<nav>
    <div class="nav-container">
        <div class="brand">
            <Link to="/">Language Learning Flashcards</Link>
        </div>
        <div class="nav-links">
            <Link to="/">Home</Link>
            <Link to="/decks">My Decks</Link>
            <Link to="/create">Create Flashcards</Link>
        </div>
    </div>
</nav>

<style>
    nav {
        padding: 1rem;
        background: #f0f0f0;
        margin-bottom: 1rem;
    }
    .nav-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        max-width: 1200px;
        margin: 0 auto;
    }
    .brand {
        font-weight: bold;
        font-size: 1.2rem;
    }
    .nav-links :global(a) {
        margin-left: 1rem;
        text-decoration: none;
        color: #333;
    }
    .nav-links :global(a:hover) {
        text-decoration: underline;
    }
</style>
