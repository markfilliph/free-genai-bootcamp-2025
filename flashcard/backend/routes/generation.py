from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List

from backend.database.db import get_db
from backend import schemas, models
from backend.auth.token import get_current_user
from backend.services.ollama_service import ollama_service

router = APIRouter(tags=["Generation"])

@router.post("/generate", response_model=schemas.GenerationResponse)
async def generate_content(
    request: schemas.GenerationRequest,
    db: Session = Depends(get_db)
    # Temporarily removed for testing
    # current_user: models.User = Depends(get_current_user)
):
    """
    Generate example sentences, conjugations, and cultural notes for a word.
    If the word is a verb, also generate conjugations.
    """
    try:
        # Generate example sentences
        example_sentences_data = await ollama_service.generate_example_sentences(request.word)
        
        # Format example sentences
        example_sentences = []
        for item in example_sentences_data:
            if isinstance(item, dict) and 'spanish' in item and 'english' in item:
                example_sentences.append(f"{item['spanish']} - {item['english']}")
            else:
                # Fallback if the structure is unexpected
                example_sentences.append(str(item))
        
        # Generate conjugations if it's a verb
        conjugations = None
        if request.is_verb:
            conjugations = await ollama_service.generate_verb_conjugations(request.word)
        
        # Generate cultural note
        cultural_note = await ollama_service.generate_cultural_note(request.word)
        
        return {
            "example_sentences": example_sentences,
            "conjugations": conjugations,
            "cultural_note": cultural_note
        }
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Error generating content: {str(e)}"
        )
