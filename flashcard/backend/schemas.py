from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime

# Use str instead of EmailStr for compatibility
class EmailStr(str):
    pass

# User schemas
class UserBase(BaseModel):
    username: str
    email: EmailStr

class UserCreate(UserBase):
    password: str

class User(UserBase):
    id: int
    created_at: datetime
    
    class Config:
        orm_mode = True

# Deck schemas
class DeckBase(BaseModel):
    name: str

class DeckCreate(DeckBase):
    pass

class Deck(DeckBase):
    id: int
    created_at: datetime
    user_id: int
    
    class Config:
        orm_mode = True

# Flashcard schemas
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
    created_at: datetime
    last_reviewed: Optional[datetime] = None
    ease_factor: int = 250
    interval: int = 1
    deck_id: int
    
    class Config:
        orm_mode = True

# Tag schemas
class TagBase(BaseModel):
    name: str

class TagCreate(TagBase):
    pass

class Tag(TagBase):
    id: int
    
    class Config:
        orm_mode = True

# LLM Generation schemas
class GenerationRequest(BaseModel):
    word: str
    is_verb: bool = False

class GenerationResponse(BaseModel):
    example_sentences: List[str]
    conjugations: Optional[str] = None
    cultural_note: Optional[str] = None

# Authentication schemas
class Token(BaseModel):
    access_token: str
    token_type: str

class TokenData(BaseModel):
    username: Optional[str] = None


