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
    
    <section class="features-section">
        <h2>Our Powerful Learning Tools</h2>
        <div class="features-container">
            <div class="feature-item">
                <div class="feature-icon">📝</div>
                <div class="feature-content">
                    <h3>Smart Flashcard Creation</h3>
                    <p>Create custom flashcards with Spanish words, phrases, and verbs with examples</p>
                </div>
            </div>
            
            <div class="feature-item">
                <div class="feature-icon">🔄</div>
                <div class="feature-content">
                    <h3>Complete Verb Conjugations</h3>
                    <p>Master all tenses and moods with comprehensive verb tables</p>
                </div>
            </div>
            
            <div class="feature-item">
                <div class="feature-icon">🌍</div>
                <div class="feature-content">
                    <h3>Rich Cultural Context</h3>
                    <p>Learn cultural notes and idiomatic expressions for authentic understanding</p>
                </div>
            </div>
            
            <div class="feature-item">
                <div class="feature-icon">⏱️</div>
                <div class="feature-content">
                    <h3>Spaced Repetition System</h3>
                    <p>Study cards at scientifically optimized intervals for maximum retention</p>
                </div>
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
    
    .features-section {
        padding: 3rem 1rem;
        background-color: #f9f9f9;
        border-radius: 12px;
        margin: 2rem 1rem;
    }
    
    h2 {
        text-align: center;
        margin-bottom: 3rem;
        color: #333;
        font-size: 2rem;
    }
    
    .features-container {
        display: flex;
        flex-direction: column;
        gap: 2rem;
        max-width: 900px;
        margin: 0 auto;
    }
    
    .feature-item {
        display: flex;
        background: white;
        border-radius: 12px;
        padding: 1.5rem;
        box-shadow: 0 4px 8px rgba(0,0,0,0.08);
        transition: transform 0.3s ease, box-shadow 0.3s ease;
    }
    
    .feature-item:hover {
        transform: translateY(-5px);
        box-shadow: 0 6px 12px rgba(0,0,0,0.12);
    }
    
    .feature-icon {
        font-size: 2.5rem;
        margin-right: 1.5rem;
        display: flex;
        align-items: center;
    }
    
    .feature-content {
        flex: 1;
    }
    
    h3 {
        color: #007bff;
        margin-bottom: 0.75rem;
        font-size: 1.3rem;
    }
    
    @media (max-width: 768px) {
        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }
        
        .feature-item {
            flex-direction: column;
            text-align: center;
        }
        
        .feature-icon {
            margin-right: 0;
            margin-bottom: 1rem;
            justify-content: center;
        }
        
        .features-section {
            padding: 2rem 1rem;
            margin: 1rem 0.5rem;
        }
    }
</style>
