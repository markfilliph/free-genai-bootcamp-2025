import { writable } from 'svelte/store';

// Initial mock data for decks
const initialDecks = [
    {
        id: '1',
        name: 'Spanish Verbs',
        description: 'Common Spanish verbs with conjugations',
        cardCount: 5,
        lastReviewed: '2025-03-20'
    },
    {
        id: '2',
        name: 'Food Vocabulary',
        description: 'Words related to food and dining',
        cardCount: 3,
        lastReviewed: '2025-03-22'
    },
    {
        id: '3',
        name: 'Travel Phrases',
        description: 'Useful phrases for traveling',
        cardCount: 2,
        lastReviewed: '2025-03-18'
    }
];

// Initial mock data for flashcards
const initialFlashcards = {
    '1': [
        { id: '101', frontText: 'hablar', backText: 'to speak', examples: 'Yo hablo español. (I speak Spanish.)', notes: 'Regular -ar verb' },
        { id: '102', frontText: 'comer', backText: 'to eat', examples: 'Ellos comen pizza. (They eat pizza.)', notes: 'Regular -er verb' },
        { id: '103', frontText: 'vivir', backText: 'to live', examples: 'Nosotros vivimos en Madrid. (We live in Madrid.)', notes: 'Regular -ir verb' },
        { id: '104', frontText: 'ser', backText: 'to be (permanent)', examples: 'Ella es doctora. (She is a doctor.)', notes: 'Irregular verb' },
        { id: '105', frontText: 'estar', backText: 'to be (temporary)', examples: 'Estoy cansado. (I am tired.)', notes: 'Irregular verb' }
    ],
    '2': [
        { id: '201', frontText: 'la manzana', backText: 'apple', examples: 'Me gusta comer manzanas. (I like to eat apples.)', notes: 'Feminine noun' },
        { id: '202', frontText: 'el pan', backText: 'bread', examples: 'Quiero comprar pan. (I want to buy bread.)', notes: 'Masculine noun' },
        { id: '203', frontText: 'la leche', backText: 'milk', examples: 'Bebo leche cada mañana. (I drink milk every morning.)', notes: 'Feminine noun' }
    ],
    '3': [
        { id: '301', frontText: '¿Dónde está...?', backText: 'Where is...?', examples: '¿Dónde está el baño? (Where is the bathroom?)', notes: 'Question phrase' },
        { id: '302', frontText: '¿Cuánto cuesta?', backText: 'How much does it cost?', examples: '¿Cuánto cuesta este libro? (How much does this book cost?)', notes: 'Question phrase' }
    ]
};

// Create a singleton store pattern to ensure all components share the same data
// This is a key improvement to fix the data persistence issue between components

// Storage keys
const DECKS_STORAGE_KEY = 'flashcardDecks';
const FLASHCARDS_STORAGE_KEY = 'flashcardCards';

// Helper function to safely access localStorage
const storage = {
    get: (key, defaultValue) => {
        try {
            if (typeof window === 'undefined' || !window.localStorage) {
                return defaultValue;
            }
            
            const value = localStorage.getItem(key);
            if (!value) return defaultValue;
            
            return JSON.parse(value);
        } catch (e) {
            console.error(`Error reading ${key} from localStorage:`, e);
            return defaultValue;
        }
    },
    
    set: (key, value) => {
        try {
            if (typeof window === 'undefined' || !window.localStorage) {
                return false;
            }
            
            localStorage.setItem(key, JSON.stringify(value));
            return true;
        } catch (e) {
            console.error(`Error writing ${key} to localStorage:`, e);
            return false;
        }
    }
};

// Create a custom store factory with persistence
function createPersistentStore(key, initialValue) {
    // Get initial data from storage or use provided initial value
    const storedValue = storage.get(key, initialValue);
    
    // Create the writable store with the initial value
    const store = writable(storedValue);
    
    // Subscribe to changes and update localStorage
    const { subscribe, set, update } = store;
    
    return {
        subscribe,
        
        set: (value) => {
            // Save to localStorage first
            storage.set(key, value);
            // Then update the store
            set(value);
            console.log(`Store '${key}' updated:`, value);
        },
        
        update: (updater) => {
            // Use update to get the current value, modify it, and set it back
            update(currentValue => {
                const newValue = updater(currentValue);
                // Save to localStorage
                storage.set(key, newValue);
                console.log(`Store '${key}' updated via updater:`, newValue);
                return newValue;
            });
        },
        
        // Force a refresh from localStorage (useful when components mount)
        refresh: () => {
            const refreshedValue = storage.get(key, initialValue);
            set(refreshedValue);
            console.log(`Store '${key}' refreshed from localStorage:`, refreshedValue);
            return refreshedValue;
        }
    };
}

// User store
export const currentUser = writable(null);

// Active deck store
export const activeDeck = writable(null);

// Create persistent stores
const decksStore = createPersistentStore(DECKS_STORAGE_KEY, initialDecks);
const flashcardsStore = createPersistentStore(FLASHCARDS_STORAGE_KEY, initialFlashcards);

// Decks store with methods for adding, updating, and deleting decks
export const decks = {
    // Basic store subscription
    subscribe: decksStore.subscribe,
    
    // Force refresh from localStorage (call this when components mount)
    refresh: () => decksStore.refresh(),
    
    // Add a new deck
    addDeck: (deck) => {
        decksStore.update(currentDecks => {
            const newDeck = {
                ...deck,
                id: Date.now().toString(),
                cardCount: 0,
                lastReviewed: 'Never'
            };
            
            console.log('Adding new deck:', newDeck);
            return [...currentDecks, newDeck];
        });
        
        // Return the current state after update
        return decksStore.refresh();
    },
    
    // Update an existing deck
    updateDeck: (id, data) => {
        console.log('Updating deck:', id, data);
        
        decksStore.update(currentDecks => {
            return currentDecks.map(deck => 
                deck.id === id ? { ...deck, ...data } : deck
            );
        });
    },
    
    // Delete a deck and its flashcards
    deleteDeck: (id) => {
        console.log('Deleting deck:', id);
        
        // First update the decks store
        decksStore.update(currentDecks => {
            return currentDecks.filter(deck => deck.id !== id);
        });
        
        // Then remove associated flashcards
        flashcardsStore.update(currentFlashcards => {
            const updatedFlashcards = { ...currentFlashcards };
            if (updatedFlashcards[id]) {
                delete updatedFlashcards[id];
            }
            return updatedFlashcards;
        });
    },
    
    // Get all decks
    getAllDecks: () => {
        // Force a refresh from localStorage first
        return decksStore.refresh();
    },
    
    // Reset to initial state
    reset: () => {
        decksStore.set(initialDecks);
        flashcardsStore.set(initialFlashcards);
    }
};

// Flashcards store with methods for adding, updating, and deleting flashcards
export const flashcards = {
    // Basic store subscription
    subscribe: flashcardsStore.subscribe,
    
    // Force refresh from localStorage (call this when components mount)
    refresh: () => flashcardsStore.refresh(),
    
    // Add a new flashcard to a deck
    addFlashcard: (deckId, card) => {
        const newCard = {
            ...card,
            id: Date.now().toString(),
            createdAt: new Date().toISOString()
        };
        
        console.log('Adding flashcard to deck:', deckId, newCard);
        
        // Add the flashcard to the specified deck
        flashcardsStore.update(currentFlashcards => {
            const updatedFlashcards = { ...currentFlashcards };
            
            // Make sure the deck exists in our flashcards object
            if (!updatedFlashcards[deckId]) {
                updatedFlashcards[deckId] = [];
            }
            
            updatedFlashcards[deckId] = [...updatedFlashcards[deckId], newCard];
            return updatedFlashcards;
        });
        
        // Update the card count in the deck
        const currentFlashcards = flashcardsStore.refresh();
        const deckCards = currentFlashcards[deckId] || [];
        decks.updateDeck(deckId, { cardCount: deckCards.length });
        
        return newCard;
    },
    
    // Get all flashcards for a specific deck
    getFlashcardsByDeck: (deckId) => {
        const allFlashcards = flashcardsStore.refresh();
        return allFlashcards[deckId] || [];
    },
    
    // Delete a flashcard from a deck
    deleteFlashcard: (deckId, cardId) => {
        console.log('Deleting flashcard:', deckId, cardId);
        
        // Remove the flashcard
        flashcardsStore.update(currentFlashcards => {
            const updatedFlashcards = { ...currentFlashcards };
            
            // Make sure the deck exists
            if (!updatedFlashcards[deckId]) {
                return updatedFlashcards;
            }
            
            updatedFlashcards[deckId] = updatedFlashcards[deckId].filter(
                card => card.id !== cardId
            );
            return updatedFlashcards;
        });
        
        // Update the card count in the deck
        const currentFlashcards = flashcardsStore.refresh();
        const deckCards = currentFlashcards[deckId] || [];
        decks.updateDeck(deckId, { cardCount: deckCards.length });
    },
    
    // Reset to initial state
    reset: () => {
        flashcardsStore.set(initialFlashcards);
    }
};
