from pydantic import BaseModel
from typing import List
from .tag import Tag

class FlashcardBase(BaseModel):
    question: str
    answer: str

class FlashcardCreate(FlashcardBase):
    tag_names: List[str] = []

class Flashcard(FlashcardBase):
    id: int
    deck_id: int
    tags: List[Tag]
    
    class Config:
        orm_mode = True
