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

// Check if localStorage is available (browser environment)
const isLocalStorageAvailable = typeof window !== 'undefined' && window.localStorage;

// Helper function to check if localStorage is actually working
function testLocalStorage() {
    if (!isLocalStorageAvailable) return false;
    
    try {
        localStorage.setItem('test', 'test');
        const testValue = localStorage.getItem('test');
        localStorage.removeItem('test');
        return testValue === 'test';
    } catch (e) {
        console.error('localStorage test failed:', e);
        return false;
    }
}

// Flag to indicate if localStorage is working
const localStorageWorks = testLocalStorage();

// Load data from localStorage or use initial data
let savedDecks;
let savedFlashcards;

if (localStorageWorks) {
    try {
        const storedDecks = localStorage.getItem('flashcardDecks');
        const storedFlashcards = localStorage.getItem('flashcardCards');
        
        if (storedDecks && storedFlashcards) {
            try {
                savedDecks = JSON.parse(storedDecks);
                savedFlashcards = JSON.parse(storedFlashcards);
                console.log('Loaded from localStorage:', { decks: savedDecks, flashcards: savedFlashcards });
            } catch (parseError) {
                console.error('Error parsing localStorage data:', parseError);
                savedDecks = initialDecks;
                savedFlashcards = initialFlashcards;
            }
        } else {
            console.log('No data found in localStorage, using initial data');
            savedDecks = initialDecks;
            savedFlashcards = initialFlashcards;
        }
    } catch (error) {
        console.error('Error accessing localStorage:', error);
        savedDecks = initialDecks;
        savedFlashcards = initialFlashcards;
    }
} else {
    console.warn('localStorage is not available or not working, using initial data');
    savedDecks = initialDecks;
    savedFlashcards = initialFlashcards;
}

// Create global shared data that will be accessible across components
let globalDecks = [...savedDecks];
let globalFlashcards = {...savedFlashcards};

// Function to save data to localStorage
function saveToLocalStorage() {
    if (localStorageWorks) {
        try {
            // Create string versions first to catch any serialization errors
            const decksString = JSON.stringify(globalDecks);
            const flashcardsString = JSON.stringify(globalFlashcards);
            
            // Only proceed if serialization was successful
            if (decksString && flashcardsString) {
                localStorage.setItem('flashcardDecks', decksString);
                localStorage.setItem('flashcardCards', flashcardsString);
                console.log('Saved to localStorage successfully');
            } else {
                console.error('Failed to serialize data for localStorage');
            }
        } catch (error) {
            console.error('Error saving to localStorage:', error);
        }
    } else {
        console.warn('localStorage is not available or not working, data will not persist');
    }
}

// Create a custom store with global data and localStorage persistence
function createStore(initialData, globalRef) {
    // Initialize the global reference if it's empty
    if (globalRef === globalDecks && globalDecks.length === 0) {
        globalDecks = [...initialData];
    } else if (globalRef === globalFlashcards && Object.keys(globalFlashcards).length === 0) {
        globalFlashcards = {...initialData};
    }
    
    const store = writable(globalRef);
    
    // Subscribe to changes and update the global reference and localStorage
    store.subscribe(value => {
        if (globalRef === globalDecks) {
            globalDecks = value;
            saveToLocalStorage();
        } else if (globalRef === globalFlashcards) {
            globalFlashcards = value;
            saveToLocalStorage();
        }
    });
    
    return store;
}

// User store
export const currentUser = writable(null);

// Active deck store
export const activeDeck = writable(null);

// Create stores with global data
const decksStore = createStore(initialDecks, globalDecks);
const flashcardsStore = createStore(initialFlashcards, globalFlashcards);

// Decks store with methods for adding, updating, and deleting decks
export const decks = {
    subscribe: decksStore.subscribe,
    
    addDeck: (deck) => {
        const newDeck = {
            ...deck,
            id: Date.now().toString(),
            cardCount: 0,
            lastReviewed: 'Never'
        };
        
        console.log('Adding new deck:', newDeck);
        
        // Add to global data
        globalDecks = [...globalDecks, newDeck];
        
        // Update the store and save to localStorage
        decksStore.set(globalDecks);
        saveToLocalStorage();
        
        return newDeck;
    },
    
    updateDeck: (id, data) => {
        console.log('Updating deck:', id, data);
        
        // Update in global data
        globalDecks = globalDecks.map(deck => 
            deck.id === id ? { ...deck, ...data } : deck
        );
        
        // Update the store and save to localStorage
        decksStore.set(globalDecks);
        saveToLocalStorage();
    },
    
    deleteDeck: (id) => {
        console.log('Deleting deck:', id);
        
        // Delete from global data
        globalDecks = globalDecks.filter(deck => deck.id !== id);
        
        // Delete flashcards for this deck
        if (globalFlashcards[id]) {
            delete globalFlashcards[id];
        }
        
        // Update the stores and save to localStorage
        decksStore.set(globalDecks);
        flashcardsStore.set(globalFlashcards);
        saveToLocalStorage();
    },
    
    getAllDecks: () => globalDecks,
    
    reset: () => {
        globalDecks = [...initialDecks];
        decksStore.set(globalDecks);
        saveToLocalStorage();
    }
};

// Flashcards store with methods for adding, updating, and deleting flashcards
export const flashcards = {
    subscribe: flashcardsStore.subscribe,
    
    addFlashcard: (deckId, card) => {
        const newCard = {
            ...card,
            id: Date.now().toString(),
            createdAt: new Date().toISOString()
        };
        
        console.log('Adding flashcard to deck:', deckId, newCard);
        
        // Make sure the deck exists in the global data
        if (!globalFlashcards[deckId]) {
            globalFlashcards[deckId] = [];
        }
        
        // Add to global data
        globalFlashcards[deckId] = [...globalFlashcards[deckId], newCard];
        
        // Update the store and save to localStorage
        flashcardsStore.set(globalFlashcards);
        saveToLocalStorage();
        
        // Update the deck's card count
        const deckCards = globalFlashcards[deckId] || [];
        decks.updateDeck(deckId, { cardCount: deckCards.length });
        
        return newCard;
    },
    
    getFlashcardsByDeck: (deckId) => {
        return globalFlashcards[deckId] || [];
    },
    
    deleteFlashcard: (deckId, cardId) => {
        console.log('Deleting flashcard:', deckId, cardId);
        
        // Make sure the deck exists
        if (globalFlashcards[deckId]) {
            // Delete from global data
            globalFlashcards[deckId] = globalFlashcards[deckId].filter(
                card => card.id !== cardId
            );
            
            // Update the store and save to localStorage
            flashcardsStore.set(globalFlashcards);
            saveToLocalStorage();
            
            // Update the deck's card count
            const deckCards = globalFlashcards[deckId] || [];
            decks.updateDeck(deckId, { cardCount: deckCards.length });
        }
    },
    
    reset: () => {
        globalFlashcards = {...initialFlashcards};
        flashcardsStore.set(globalFlashcards);
        saveToLocalStorage();
    }
};
