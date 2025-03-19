from sqlalchemy import Column, Integer, String, ForeignKey, DateTime, Boolean, Text
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from .database import Base
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class User(Base):
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True)
    email = Column(String, unique=True, index=True)
    password_hash = Column(String)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    decks = relationship("Deck", back_populates="owner", cascade="all, delete-orphan")
    
    def verify_password(self, password):
        return pwd_context.verify(password, self.password_hash)
    
    @staticmethod
    def get_password_hash(password):
        return pwd_context.hash(password)

class Deck(Base):
    __tablename__ = "decks"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    user_id = Column(Integer, ForeignKey("users.id"))
    
    owner = relationship("User", back_populates="decks")
    flashcards = relationship("Flashcard", back_populates="deck", cascade="all, delete-orphan")

class Flashcard(Base):
    __tablename__ = "flashcards"
    
    id = Column(Integer, primary_key=True, index=True)
    word = Column(String, index=True)
    example_sentence = Column(Text)
    translation = Column(Text)
    conjugation = Column(Text, nullable=True)
    cultural_note = Column(Text, nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    last_reviewed = Column(DateTime(timezone=True), nullable=True)
    ease_factor = Column(Integer, default=250)  # SuperMemo2 ease factor (multiplied by 100)
    interval = Column(Integer, default=1)       # SuperMemo2 interval in days
    deck_id = Column(Integer, ForeignKey("decks.id"))
    
    deck = relationship("Deck", back_populates="flashcards")
    
class Tag(Base):
    __tablename__ = "tags"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True, index=True)
    
    flashcard_tags = relationship("FlashcardTag", back_populates="tag")

class FlashcardTag(Base):
    __tablename__ = "flashcard_tags"
    
    flashcard_id = Column(Integer, ForeignKey("flashcards.id"), primary_key=True)
    tag_id = Column(Integer, ForeignKey("tags.id"), primary_key=True)
    
    tag = relationship("Tag", back_populates="flashcard_tags")
