import os
import sqlite3
from datetime import datetime
import json

# Database setup
DB_PATH = "test_models.db"

def get_db_connection():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

def init_db():
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Create users table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TEXT NOT NULL
    )
    ''')
    
    # Create decks table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS decks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    )
    ''')
    
    # Create flashcards table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS flashcards (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word TEXT NOT NULL,
        example_sentence TEXT NOT NULL,
        translation TEXT NOT NULL,
        conjugation TEXT,
        cultural_note TEXT,
        deck_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        last_reviewed TEXT,
        ease_factor INTEGER DEFAULT 250,
        interval INTEGER DEFAULT 1,
        FOREIGN KEY (deck_id) REFERENCES decks (id)
    )
    ''')
    
    conn.commit()
    conn.close()

def create_user(username, email, password):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Hash password (simplified for testing)
    password_hash = password + "_hashed"
    
    cursor.execute(
        "INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
        (username, email, password_hash, datetime.now().isoformat())
    )
    
    user_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    return user_id

def create_deck(name, user_id):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        "INSERT INTO decks (name, user_id, created_at) VALUES (?, ?, ?)",
        (name, user_id, datetime.now().isoformat())
    )
    
    deck_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    return deck_id

def create_flashcard(word, example_sentence, translation, conjugation, cultural_note, deck_id):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        """
        INSERT INTO flashcards 
        (word, example_sentence, translation, conjugation, cultural_note, deck_id, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """,
        (
            word, 
            example_sentence, 
            translation, 
            conjugation, 
            cultural_note, 
            deck_id, 
            datetime.now().isoformat()
        )
    )
    
    flashcard_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    return flashcard_id

def get_user_decks(user_id):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        "SELECT id, name, user_id, created_at FROM decks WHERE user_id = ?",
        (user_id,)
    )
    
    decks = [dict(deck) for deck in cursor.fetchall()]
    conn.close()
    
    return decks

def get_deck_flashcards(deck_id):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        """
        SELECT id, word, example_sentence, translation, conjugation, cultural_note, 
        deck_id, created_at, last_reviewed, ease_factor, interval 
        FROM flashcards WHERE deck_id = ?
        """,
        (deck_id,)
    )
    
    flashcards = [dict(fc) for fc in cursor.fetchall()]
    conn.close()
    
    return flashcards

def run_tests():
    # Remove test database if it exists
    if os.path.exists(DB_PATH):
        os.remove(DB_PATH)
    
    print("Creating database tables...")
    init_db()
    
    print("\nTesting user creation...")
    user_id = create_user("testuser", "test@example.com", "password123")
    print(f"Created user with ID: {user_id}")
    
    print("\nTesting deck creation...")
    deck_id = create_deck("Spanish Basics", user_id)
    print(f"Created deck with ID: {deck_id}")
    
    print("\nTesting flashcard creation...")
    flashcard_id = create_flashcard(
        "hola", 
        "Hola, ¿cómo estás?", 
        "Hello, how are you?", 
        None, 
        "Common greeting in Spanish-speaking countries", 
        deck_id
    )
    print(f"Created flashcard with ID: {flashcard_id}")
    
    print("\nTesting data retrieval...")
    decks = get_user_decks(user_id)
    print(f"Found {len(decks)} decks for user {user_id}")
    print(json.dumps(decks, indent=2))
    
    for deck in decks:
        flashcards = get_deck_flashcards(deck["id"])
        print(f"\nFound {len(flashcards)} flashcards for deck {deck['id']}")
        print(json.dumps(flashcards, indent=2))
    
    print("\nTests completed successfully!")

if __name__ == "__main__":
    run_tests()
