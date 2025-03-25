<script>
    import { Link } from 'svelte-routing';
</script>

<div class="home-container">
    <section class="hero">
        <h1>Language Learning Flashcard Generator</h1>
        <p class="subtitle">Master Spanish vocabulary, verbs, and phrases with our intelligent flashcard system</p>
        
        <div class="cta-buttons">
            <Link to="/decks" class="cta-button primary">My Flashcard Decks</Link>
            <Link to="/create" class="cta-button secondary">Create New Flashcards</Link>
        </div>
    </section>
    
    <section class="features">
        <h2>Key Features</h2>
        <div class="feature-grid">
            <div class="feature-card">
                <h3>Flashcard Creation</h3>
                <p>Create flashcards with words, phrases, or verbs in Spanish with automatic example generation</p>
            </div>
            <div class="feature-card">
                <h3>Verb Conjugation</h3>
                <p>View conjugations for all tenses and moods to master Spanish verbs</p>
            </div>
            <div class="feature-card">
                <h3>Cultural Context</h3>
                <p>Learn cultural notes and idiomatic expressions for better understanding</p>
            </div>
            <div class="feature-card">
                <h3>Text-to-Speech</h3>
                <p>Listen to pronunciations to improve your speaking skills</p>
            </div>
        </div>
    </section>
</div>

<style>
    .home-container {
        max-width: 1200px;
        margin: 0 auto;
    }
    
    .hero {
        text-align: center;
        padding: 3rem 1rem;
        margin-bottom: 2rem;
    }
    
    h1 {
        font-size: 2.5rem;
        margin-bottom: 1rem;
        color: #333;
    }
    
    .subtitle {
        font-size: 1.2rem;
        color: #666;
        margin-bottom: 2rem;
        max-width: 800px;
        margin-left: auto;
        margin-right: auto;
    }
    
    .cta-buttons {
        display: flex;
        justify-content: center;
        gap: 1rem;
        margin-bottom: 2rem;
    }
    
    .cta-buttons :global(.cta-button) {
        display: inline-block;
        padding: 0.8rem 1.5rem;
        border-radius: 4px;
        text-decoration: none;
        font-weight: bold;
        transition: all 0.3s ease;
    }
    
    .cta-buttons :global(.primary) {
        background: #007bff;
        color: white;
    }
    
    .cta-buttons :global(.primary:hover) {
        background: #0069d9;
    }
    
    .cta-buttons :global(.secondary) {
        background: #f8f9fa;
        color: #333;
        border: 1px solid #ddd;
    }
    
    .cta-buttons :global(.secondary:hover) {
        background: #e2e6ea;
    }
    
    .features {
        padding: 2rem 1rem;
    }
    
    h2 {
        text-align: center;
        margin-bottom: 2rem;
        color: #333;
    }
    
    .feature-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 2rem;
    }
    
    .feature-card {
        background: #f8f9fa;
        border-radius: 8px;
        padding: 1.5rem;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    
    h3 {
        color: #007bff;
        margin-bottom: 0.5rem;
    }
    
    @media (max-width: 768px) {
        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }
        
        .feature-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
