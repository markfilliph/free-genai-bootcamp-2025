from sqlalchemy import Column, Integer, String, ForeignKey
from database.db import Base

class Flashcard(Base):
    __tablename__ = "flashcards"
    id = Column(Integer, primary_key=True, index=True)
    question = Column(String)
    answer = Column(String)
    deck_id = Column(Integer, ForeignKey("decks.id"))
