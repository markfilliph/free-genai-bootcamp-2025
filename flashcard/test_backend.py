from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel
from typing import List, Optional
import uvicorn

app = FastAPI()

# Simple in-memory database
users_db = {}
decks_db = {}
flashcards_db = {}
current_user_id = 0
current_deck_id = 0
current_flashcard_id = 0

# Models
class User(BaseModel):
    id: int
    username: str
    email: str
    password: str

class UserCreate(BaseModel):
    username: str
    email: str
    password: str

class Token(BaseModel):
    access_token: str
    token_type: str

class Deck(BaseModel):
    id: int
    name: str
    user_id: int

class DeckCreate(BaseModel):
    name: str

class Flashcard(BaseModel):
    id: int
    word: str
    example_sentence: str
    translation: str
    conjugation: Optional[str] = None
    cultural_note: Optional[str] = None
    deck_id: int

class FlashcardCreate(BaseModel):
    word: str
    example_sentence: str
    translation: str
    conjugation: Optional[str] = None
    cultural_note: Optional[str] = None
    deck_id: int

class GenerationRequest(BaseModel):
    word: str
    is_verb: bool = False

class GenerationResponse(BaseModel):
    example_sentences: List[str]
    conjugations: Optional[str] = None
    cultural_note: Optional[str] = None

# Auth
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

def get_current_user(token: str = Depends(oauth2_scheme)):
    user = users_db.get(token)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
    return user

# Auth routes
@app.post("/token", response_model=Token)
async def login_for_access_token(form_data: OAuth2PasswordRequestForm = Depends()):
    for user_id, user in users_db.items():
        if user.username == form_data.username and user.password == form_data.password:
            return {"access_token": str(user.id), "token_type": "bearer"}
    raise HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Incorrect username or password",
        headers={"WWW-Authenticate": "Bearer"},
    )

@app.post("/register")
async def register_user(user: UserCreate):
    global current_user_id
    current_user_id += 1
    new_user = User(
        id=current_user_id,
        username=user.username,
        email=user.email,
        password=user.password
    )
    users_db[str(current_user_id)] = new_user
    return {"message": "User created successfully", "user_id": current_user_id}

# Deck routes
@app.post("/decks", response_model=Deck)
async def create_deck(deck: DeckCreate, current_user: User = Depends(get_current_user)):
    global current_deck_id
    current_deck_id += 1
    new_deck = Deck(
        id=current_deck_id,
        name=deck.name,
        user_id=current_user.id
    )
    decks_db[current_deck_id] = new_deck
    return new_deck

@app.get("/decks", response_model=List[Deck])
async def get_decks(current_user: User = Depends(get_current_user)):
    return [deck for deck in decks_db.values() if deck.user_id == current_user.id]

@app.get("/decks/{deck_id}", response_model=Deck)
async def get_deck(deck_id: int, current_user: User = Depends(get_current_user)):
    deck = decks_db.get(deck_id)
    if not deck:
        raise HTTPException(status_code=404, detail="Deck not found")
    if deck.user_id != current_user.id:
        raise HTTPException(status_code=403, detail="Not authorized to access this deck")
    return deck

# Flashcard routes
@app.post("/flashcards", response_model=Flashcard)
async def create_flashcard(flashcard: FlashcardCreate, current_user: User = Depends(get_current_user)):
    # Check if deck exists and belongs to user
    deck = decks_db.get(flashcard.deck_id)
    if not deck:
        raise HTTPException(status_code=404, detail="Deck not found")
    if deck.user_id != current_user.id:
        raise HTTPException(status_code=403, detail="Not authorized to access this deck")
    
    global current_flashcard_id
    current_flashcard_id += 1
    new_flashcard = Flashcard(
        id=current_flashcard_id,
        word=flashcard.word,
        example_sentence=flashcard.example_sentence,
        translation=flashcard.translation,
        conjugation=flashcard.conjugation,
        cultural_note=flashcard.cultural_note,
        deck_id=flashcard.deck_id
    )
    flashcards_db[current_flashcard_id] = new_flashcard
    return new_flashcard

@app.get("/decks/{deck_id}/flashcards", response_model=List[Flashcard])
async def get_flashcards_by_deck(deck_id: int, current_user: User = Depends(get_current_user)):
    # Check if deck exists and belongs to user
    deck = decks_db.get(deck_id)
    if not deck:
        raise HTTPException(status_code=404, detail="Deck not found")
    if deck.user_id != current_user.id:
        raise HTTPException(status_code=403, detail="Not authorized to access this deck")
    
    return [fc for fc in flashcards_db.values() if fc.deck_id == deck_id]

# LLM Generation route
@app.post("/generate", response_model=GenerationResponse)
async def generate_content(request: GenerationRequest, current_user: User = Depends(get_current_user)):
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

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
