export const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8000';

export async function apiFetch(path, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${path}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include',
            ...options
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(`API Error (${response.status}): ${error}`);
        }

        return await response.json();
    } catch (error) {
        if (error.message.startsWith('API Error')) {
            throw error;
        }
        throw new Error(error.message);
    }
}

// User authentication functions
export async function login(username, password) {
    return apiFetch('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
    });
}

export async function register(username, email, password) {
    return apiFetch('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify({ username, email, password })
    });
}

export async function logout() {
    return apiFetch('/api/auth/logout', { method: 'POST' });
}

// Deck management functions
export async function getDecks() {
    return apiFetch('/api/decks');
}

export async function createDeck(name, description) {
    return apiFetch('/api/decks', {
        method: 'POST',
        body: JSON.stringify({ name, description })
    });
}

export async function updateDeck(id, data) {
    return apiFetch(`/api/decks/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

export async function deleteDeck(id) {
    return apiFetch(`/api/decks/${id}`, { method: 'DELETE' });
}

// Flashcard management functions
export async function getFlashcards(deckId) {
    return apiFetch(`/api/decks/${deckId}/flashcards`);
}

export async function createFlashcard(deckId, flashcardData) {
    return apiFetch(`/api/decks/${deckId}/flashcards`, {
        method: 'POST',
        body: JSON.stringify(flashcardData)
    });
}

export async function updateFlashcard(deckId, cardId, data) {
    return apiFetch(`/api/decks/${deckId}/flashcards/${cardId}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

export async function deleteFlashcard(deckId, cardId) {
    return apiFetch(`/api/decks/${deckId}/flashcards/${cardId}`, { method: 'DELETE' });
}

// AI Generation functions
export async function generateContent(word, isVerb = false) {
    console.log('Calling generateContent API with:', { word, is_verb: isVerb });
    try {
        const result = await apiFetch('/api/generate', {
            method: 'POST',
            body: JSON.stringify({ word, is_verb: isVerb })
        });
        console.log('Generate content API response:', result);
        return result;
    } catch (error) {
        console.error('Generate content API error:', error);
        throw error;
    }
}

export async function generateExampleSentences(word) {
    const result = await generateContent(word, false);
    return result.example_sentences;
}

export async function generateVerbConjugations(verb) {
    const result = await generateContent(verb, true);
    return result.conjugations;
}

export async function generateCulturalNote(word) {
    const result = await generateContent(word, false);
    return result.cultural_note;
}
