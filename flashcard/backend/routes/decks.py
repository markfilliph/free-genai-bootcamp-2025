from fastapi import APIRouter, Depends, HTTPException, status
from backend.database.db import get_db
from backend import schemas, models
from backend.auth.token import get_current_user
from sqlalchemy.orm import Session
from sqlalchemy.exc import SQLAlchemyError
from typing import List

router = APIRouter(tags=["Decks"])

@router.post("/decks/", response_model=schemas.Deck)
async def create_deck(
    deck: schemas.DeckCreate, 
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    db_deck = models.Deck(name=deck.name, user_id=current_user.id)
    db.add(db_deck)
    db.commit()
    db.refresh(db_deck)
    return db_deck

@router.get("/decks/", response_model=List[schemas.Deck])
async def read_decks(
    skip: int = 0, 
    limit: int = 100, 
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Only return decks that belong to the current user
    decks = db.query(models.Deck).filter(
        models.Deck.user_id == current_user.id
    ).offset(skip).limit(limit).all()
    return decks

@router.get("/decks/{deck_id}", response_model=schemas.Deck)
async def read_deck(
    deck_id: int, 
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Only return the deck if it belongs to the current user
    db_deck = db.query(models.Deck).filter(
        models.Deck.id == deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if db_deck is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, 
            detail="Deck not found or you don't have access to this deck"
        )
    return db_deck

@router.put("/decks/{deck_id}", response_model=schemas.Deck)
async def update_deck(
    deck_id: int,
    deck_update: schemas.DeckCreate,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Find the deck and verify ownership
    db_deck = db.query(models.Deck).filter(
        models.Deck.id == deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if db_deck is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, 
            detail="Deck not found or you don't have access to this deck"
        )
    
    # Update the deck name
    db_deck.name = deck_update.name
    db.commit()
    db.refresh(db_deck)
    return db_deck

@router.delete("/decks/{deck_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_deck(
    deck_id: int,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(get_current_user)
):
    # Find the deck and verify ownership
    db_deck = db.query(models.Deck).filter(
        models.Deck.id == deck_id,
        models.Deck.user_id == current_user.id
    ).first()
    
    if db_deck is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, 
            detail="Deck not found or you don't have access to this deck"
        )
    
    # Delete the deck (and all associated flashcards due to cascade)
    db.delete(db_deck)
    db.commit()
    return None