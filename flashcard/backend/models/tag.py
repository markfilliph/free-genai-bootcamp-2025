from sqlalchemy import Column, Integer, String, Table, ForeignKey
from sqlalchemy.orm import relationship
from database.db import Base

flashcard_tag = Table(
    "flashcard_tag", Base.metadata,
    Column("flashcard_id", Integer, ForeignKey("flashcards.id")),
    Column("tag_id", Integer, ForeignKey("tags.id"))
)

class Tag(Base):
    __tablename__ = "tags"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True)
    flashcards = relationship("Flashcard", secondary=flashcard_tag, back_populates="tags")
