<script>
    import { Router, Link } from 'svelte-routing'
    import Navbar from './components/Navbar.svelte'
    import Home from './routes/Home.svelte'
    import Login from './routes/Login.svelte'
    import DeckManagement from './routes/DeckManagement.svelte'
</script>

<Navbar />
<Router>
    <Home path="/" />
    <Login path="/login" />
    <DeckManagement path="/decks" />
</Router>
