from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List, Optional
from datetime import datetime

from backend.database.db import get_db
from backend import schemas, models
from backend.auth.token import get_current_user

router = APIRouter(tags=["Flashcards"])

@router.post("/flashcards/", response_model=schemas.Flashcard)
async def create_flashcard(
    flashcard: schemas.FlashcardCreate,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Verify the deck exists and belongs to the current user
    deck = db.query(models.Deck).filter(
        models.Deck.id == flashcard.deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if not deck:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Deck not found or you don't have access to this deck"
        )
    
    # Create new flashcard
    db_flashcard = models.Flashcard(
        word=flashcard.word,
        example_sentence=flashcard.example_sentence,
        translation=flashcard.translation,
        conjugation=flashcard.conjugation,
        cultural_note=flashcard.cultural_note,
        deck_id=flashcard.deck_id
    )
    
    db.add(db_flashcard)
    db.commit()
    db.refresh(db_flashcard)
    return db_flashcard

@router.get("/flashcards/deck/{deck_id}", response_model=List[schemas.Flashcard])
async def get_flashcards_by_deck(
    deck_id: int,
    skip: int = 0,
    limit: int = 100,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Verify the deck exists and belongs to the current user
    deck = db.query(models.Deck).filter(
        models.Deck.id == deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if not deck:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Deck not found or you don't have access to this deck"
        )
    
    # Get flashcards for the deck
    flashcards = db.query(models.Flashcard).filter(
        models.Flashcard.deck_id == deck_id
    ).offset(skip).limit(limit).all()
    
    return flashcards

@router.get("/flashcards/{flashcard_id}", response_model=schemas.Flashcard)
async def get_flashcard(
    flashcard_id: int,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Get the flashcard
    flashcard = db.query(models.Flashcard).filter(models.Flashcard.id == flashcard_id).first()
    
    if not flashcard:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Flashcard not found"
        )
    
    # Verify the deck belongs to the current user
    deck = db.query(models.Deck).filter(
        models.Deck.id == flashcard.deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if not deck:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You don't have access to this flashcard"
        )
    
    return flashcard

@router.put("/flashcards/{flashcard_id}", response_model=schemas.Flashcard)
async def update_flashcard(
    flashcard_id: int,
    flashcard_update: schemas.FlashcardBase,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Get the flashcard
    db_flashcard = db.query(models.Flashcard).filter(models.Flashcard.id == flashcard_id).first()
    
    if not db_flashcard:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Flashcard not found"
        )
    
    # Verify the deck belongs to the current user
    deck = db.query(models.Deck).filter(
        models.Deck.id == db_flashcard.deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if not deck:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You don't have access to this flashcard"
        )
    
    # Update flashcard fields
    for key, value in flashcard_update.dict().items():
        setattr(db_flashcard, key, value)
    
    db.commit()
    db.refresh(db_flashcard)
    return db_flashcard

@router.delete("/flashcards/{flashcard_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_flashcard(
    flashcard_id: int,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Get the flashcard
    db_flashcard = db.query(models.Flashcard).filter(models.Flashcard.id == flashcard_id).first()
    
    if not db_flashcard:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Flashcard not found"
        )
    
    # Verify the deck belongs to the current user
    deck = db.query(models.Deck).filter(
        models.Deck.id == db_flashcard.deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if not deck:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You don't have access to this flashcard"
        )
    
    # Delete the flashcard
    db.delete(db_flashcard)
    db.commit()
    return None
