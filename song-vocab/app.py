from fastapi import FastAPI, HTTPException, Response
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import List, Optional, Union, Dict, Any
from agent import SongLyricsAgent

app = FastAPI(
    title="Song Vocabulary Extractor",
    description="Extract lyrics and vocabulary from Japanese songs"
)

class LyricsRequest(BaseModel):
    message_request: str = Field(
        description="A string that describes the song and/or artist to get lyrics for"
    )

class VocabularyPart(BaseModel):
    kanji: str
    romaji: List[str]

class VocabularyItem(BaseModel):
    kanji: str
    romaji: str
    parts: List[VocabularyPart]

    def to_simple_format(self) -> str:
        """Convert to simple string format."""
        return f"{self.kanji} ({self.romaji})"

class LyricsResponseDetailed(BaseModel):
    """Detailed response with structured vocabulary"""
    lyrics: str
    vocabulary: List[VocabularyItem]
    metadata: Optional[Dict[str, Any]] = None

class LyricsResponseSimple(BaseModel):
    """Simple response format as per tech specs"""
    lyrics: str
    vocabulary: List[str]

@app.post(
    "/api/agent",
    response_model=Union[LyricsResponseSimple, LyricsResponseDetailed],
    responses={
        200: {
            "description": "Successfully retrieved lyrics and vocabulary",
            "content": {
                "application/json": {
                    "examples": {
                        "simple": {
                            "value": {
                                "lyrics": "Japanese lyrics...",
                                "vocabulary": ["言葉 (kotoba)", "歌 (uta)"]
                            }
                        }
                    }
                }
            }
        },
        500: {"description": "Internal server error"}
    }
)
async def get_lyrics(
    request: LyricsRequest,
    response: Response,
    detailed: bool = False
) -> Union[LyricsResponseSimple, LyricsResponseDetailed]:
    """Get lyrics and vocabulary for a song.
    
    Args:
        request: The lyrics request containing the song/artist description
        response: FastAPI response object
        detailed: If True, return detailed vocabulary structure. If False, return simple format.
    """
    try:
        agent = SongLyricsAgent()
        result = await agent.process_request(request.message_request)
        
        # Format vocabulary items
        formatted_vocabulary = []
        for vocab in result.get("vocabulary", []):
            if isinstance(vocab, dict):
                formatted_vocabulary.append(VocabularyItem(**vocab))
            elif isinstance(vocab, VocabularyItem):
                formatted_vocabulary.append(vocab)
        
        # Return appropriate response format
        if detailed:
            return LyricsResponseDetailed(
                lyrics=result.get("lyrics", ""),
                vocabulary=formatted_vocabulary,
                metadata=result.get("metadata", {})
            )
        else:
            # Convert to simple format
            simple_vocabulary = [item.to_simple_format() for item in formatted_vocabulary]
            return LyricsResponseSimple(
                lyrics=result.get("lyrics", ""),
                vocabulary=simple_vocabulary
            )
            
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
