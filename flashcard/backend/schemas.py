from pydantic import BaseModel

class DeckBase(BaseModel):
    name: str

class DeckCreate(DeckBase):
    pass

class Deck(DeckBase):
    id: int
    
    class Config:
        orm_mode = True
