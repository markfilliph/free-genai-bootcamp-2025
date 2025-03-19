from fastapi import FastAPI, Depends, HTTPException, status, Body
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
import uvicorn
import sqlite3
import json
import os
import time
from datetime import datetime, timedelta
import secrets
import hashlib

# Create a FastAPI app
app = FastAPI(title="Language Learning Flashcard Generator API")

# CORS Configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],  # Frontend URL
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Database setup
DB_PATH = "flashcards.db"

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

# Initialize database
init_db()

# Models
class UserBase(BaseModel):
    username: str
    email: str

class UserCreate(UserBase):
    password: str

class User(UserBase):
    id: int
    created_at: str

class Token(BaseModel):
    access_token: str
    token_type: str

class DeckBase(BaseModel):
    name: str

class DeckCreate(DeckBase):
    pass

class Deck(DeckBase):
    id: int
    user_id: int
    created_at: str

class FlashcardBase(BaseModel):
    word: str
    example_sentence: str
    translation: str
    conjugation: Optional[str] = None
    cultural_note: Optional[str] = None

class FlashcardCreate(FlashcardBase):
    deck_id: int

class Flashcard(FlashcardBase):
    id: int
    deck_id: int
    created_at: str
    last_reviewed: Optional[str] = None
    ease_factor: int = 250
    interval: int = 1

class GenerationRequest(BaseModel):
    word: str
    is_verb: bool = False

class GenerationResponse(BaseModel):
    example_sentences: List[str]
    conjugations: Optional[str] = None
    cultural_note: Optional[str] = None

# Authentication
SECRET_KEY = "your-secret-key-here"  # In production, use a proper secret key
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="api/auth/login")

def hash_password(password: str) -> str:
    return hashlib.sha256(password.encode()).hexdigest()

def verify_password(plain_password: str, hashed_password: str) -> bool:
    return hash_password(plain_password) == hashed_password

def create_access_token(data: dict) -> str:
    # Simple token generation - in production use JWT
    token = secrets.token_hex(32)
    
    # Store token in a simple in-memory cache (in production use Redis)
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Store token with user_id
    cursor.execute(
        "INSERT INTO tokens (token, user_id, expires_at) VALUES (?, ?, ?)",
        (token, data["user_id"], (datetime.now() + timedelta(minutes=30)).isoformat())
    )
    
    conn.commit()
    conn.close()
    
    return token

def get_current_user(token: str = Depends(oauth2_scheme)):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Check if token exists and is valid
    cursor.execute(
        "SELECT user_id, expires_at FROM tokens WHERE token = ?",
        (token,)
    )
    
    token_data = cursor.fetchone()
    
    if not token_data:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Check if token has expired
    expires_at = datetime.fromisoformat(token_data["expires_at"])
    if expires_at < datetime.now():
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token has expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Get user data
    cursor.execute(
        "SELECT id, username, email, created_at FROM users WHERE id = ?",
        (token_data["user_id"],)
    )
    
    user = cursor.fetchone()
    conn.close()
    
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="User not found",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    return dict(user)

# Create tokens table
conn = get_db_connection()
cursor = conn.cursor()
cursor.execute('''
CREATE TABLE IF NOT EXISTS tokens (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
)
''')
conn.commit()
conn.close()

# Routes
@app.post("/api/auth/register", response_model=Dict[str, Any])
async def register(user: UserCreate):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Check if username or email already exists
    cursor.execute(
        "SELECT id FROM users WHERE username = ? OR email = ?",
        (user.username, user.email)
    )
    
    if cursor.fetchone():
        conn.close()
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Username or email already registered"
        )
    
    # Hash password
    hashed_password = hash_password(user.password)
    
    # Insert new user
    cursor.execute(
        "INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
        (user.username, user.email, hashed_password, datetime.now().isoformat())
    )
    
    user_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    return {"message": "User registered successfully", "user_id": user_id}

@app.post("/api/auth/login", response_model=Token)
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Find user by username
    cursor.execute(
        "SELECT id, password_hash FROM users WHERE username = ?",
        (form_data.username,)
    )
    
    user = cursor.fetchone()
    conn.close()
    
    if not user or not verify_password(form_data.password, user["password_hash"]):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Create access token
    access_token = create_access_token({"user_id": user["id"]})
    
    return {"access_token": access_token, "token_type": "bearer"}

@app.post("/api/decks", response_model=Deck)
async def create_deck(deck: DeckCreate, current_user: dict = Depends(get_current_user)):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Insert new deck
    cursor.execute(
        "INSERT INTO decks (name, user_id, created_at) VALUES (?, ?, ?)",
        (deck.name, current_user["id"], datetime.now().isoformat())
    )
    
    deck_id = cursor.lastrowid
    conn.commit()
    
    # Get the created deck
    cursor.execute(
        "SELECT id, name, user_id, created_at FROM decks WHERE id = ?",
        (deck_id,)
    )
    
    new_deck = cursor.fetchone()
    conn.close()
    
    return dict(new_deck)

@app.get("/api/decks", response_model=List[Deck])
async def get_decks(current_user: dict = Depends(get_current_user)):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Get all decks for the current user
    cursor.execute(
        "SELECT id, name, user_id, created_at FROM decks WHERE user_id = ?",
        (current_user["id"],)
    )
    
    decks = [dict(deck) for deck in cursor.fetchall()]
    conn.close()
    
    return decks

@app.post("/api/flashcards", response_model=Flashcard)
async def create_flashcard(flashcard: FlashcardCreate, current_user: dict = Depends(get_current_user)):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Check if deck exists and belongs to user
    cursor.execute(
        "SELECT id FROM decks WHERE id = ? AND user_id = ?",
        (flashcard.deck_id, current_user["id"])
    )
    
    if not cursor.fetchone():
        conn.close()
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Deck not found or you don't have permission to access it"
        )
    
    # Insert new flashcard
    cursor.execute(
        """
        INSERT INTO flashcards 
        (word, example_sentence, translation, conjugation, cultural_note, deck_id, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """,
        (
            flashcard.word, 
            flashcard.example_sentence, 
            flashcard.translation, 
            flashcard.conjugation, 
            flashcard.cultural_note, 
            flashcard.deck_id, 
            datetime.now().isoformat()
        )
    )
    
    flashcard_id = cursor.lastrowid
    conn.commit()
    
    # Get the created flashcard
    cursor.execute(
        """
        SELECT id, word, example_sentence, translation, conjugation, cultural_note, 
        deck_id, created_at, last_reviewed, ease_factor, interval 
        FROM flashcards WHERE id = ?
        """,
        (flashcard_id,)
    )
    
    new_flashcard = cursor.fetchone()
    conn.close()
    
    return dict(new_flashcard)

@app.get("/api/decks/{deck_id}/flashcards", response_model=List[Flashcard])
async def get_flashcards_by_deck(deck_id: int, current_user: dict = Depends(get_current_user)):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Check if deck exists and belongs to user
    cursor.execute(
        "SELECT id FROM decks WHERE id = ? AND user_id = ?",
        (deck_id, current_user["id"])
    )
    
    if not cursor.fetchone():
        conn.close()
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Deck not found or you don't have permission to access it"
        )
    
    # Get all flashcards for the deck
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

@app.post("/api/generate", response_model=GenerationResponse)
async def generate_content(request: GenerationRequest, current_user: dict = Depends(get_current_user)):
    # Simulate LLM generation
    example_sentences = [
        f"Ejemplo con '{request.word}': Esta es una oración de ejemplo.",
        f"Otro ejemplo con '{request.word}': Segunda oración de ejemplo."
    ]
    
    conjugations = None
    if request.is_verb:
        conjugations = f"Conjugaciones para '{request.word}':\nPresente: yo {request.word}o, tú {request.word}es..."
    
    cultural_note = f"Nota cultural sobre '{request.word}': Este término es comúnmente usado en España y Latinoamérica."
    
    return {
        "example_sentences": example_sentences,
        "conjugations": conjugations,
        "cultural_note": cultural_note
    }

@app.get("/")
async def read_root():
    return {
        "message": "Welcome to Language Learning Flashcard Generator API",
        "docs": "/docs",
        "version": "1.0.0"
    }

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
