from fastapi import APIRouter, Depends, HTTPException
from backend.database.db import get_db
from backend import schemas, models
from sqlalchemy.orm import Session
from sqlalchemy.exc import SQLAlchemyError

router = APIRouter()

@router.post("/decks/", response_model=schemas.Deck)
def create_deck(deck: schemas.DeckCreate, db: Session = Depends(get_db)):
    db_deck = models.Deck(name=deck.name)
    db.add(db_deck)
    db.commit()
    db.refresh(db_deck)
    return db_deck

@router.get("/decks/", response_model=list[schemas.Deck])
def read_decks(skip: int = 0, limit: int = 100, db: Session = Depends(get_db)):
    decks = db.query(models.Deck).offset(skip).limit(limit).all()
    return decks

@router.get("/decks/{deck_id}", response_model=schemas.Deck)
def read_deck(deck_id: int, db: Session = Depends(get_db)):
    db_deck = db.query(models.Deck).filter(models.Deck.id == deck_id).first()
    if db_deck is None:
        raise HTTPException(status_code=404, detail="Deck not found")
    return db_deck