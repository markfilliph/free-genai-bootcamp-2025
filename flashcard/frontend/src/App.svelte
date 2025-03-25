<script>
    import { Router, Route, Link } from 'svelte-routing'
    import Navbar from './components/Navbar.svelte'
    import Home from './routes/Home.svelte'
    import DeckManagement from './routes/DeckManagement.svelte'
    import CreateFlashcards from './routes/CreateFlashcards.svelte'
    import StudyDeck from './routes/StudyDeck.svelte'
    
    export let url = '';
</script>

<Router {url}>
    <div class="app-container">
        <Navbar />
        <main>
            <Route path="/" component={Home} />
            <Route path="/decks" component={DeckManagement} />
            <Route path="/create" component={CreateFlashcards} />
            <Route path="/study" component={StudyDeck} />
        </main>
    </div>
</Router>

<style>
    .app-container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem;
    }
    main {
        padding: 1rem 0;
    }
</style>
