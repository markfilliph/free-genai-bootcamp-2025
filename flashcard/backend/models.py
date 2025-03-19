from sqlalchemy import Column, Integer, String
from .database import Base

class Deck(Base):
    __tablename__ = "decks"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
