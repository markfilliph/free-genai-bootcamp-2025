from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Dict, List, Optional
import uvicorn
import asyncio

import sys
import os

# Add parent directory to Python path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.chat import Chat
from backend.get_transcript import YouTubeTranscriptDownloader
from backend.vector_store import TranscriptVectorStore
from backend.question_generator import QuestionGenerator

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize components
chat = Chat()
transcript_downloader = YouTubeTranscriptDownloader()
vector_store = TranscriptVectorStore()
question_generator = QuestionGenerator()

class TranscriptRequest(BaseModel):
    video_url: str

class ChatRequest(BaseModel):
    message: str
    history: Optional[List[Dict[str, str]]] = None

class SearchRequest(BaseModel):
    query: str
    n_results: Optional[int] = 5

class QuestionRequest(BaseModel):
    section_num: int
    topic: str

@app.post("/api/transcript")
async def get_transcript(request: TranscriptRequest):
    """Get transcript for a YouTube video and add it to vector store"""
    try:
        transcript_data = await transcript_downloader.get_transcript(request.video_url)
        if not transcript_data:
            raise HTTPException(status_code=404, detail="No transcript found")
            
        # Add to vector store
        video_id = transcript_data['video_id']
        await vector_store.add_transcript(video_id, transcript_data)
        
        return {"status": "success", "video_id": video_id}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/search")
async def search_transcript(request: SearchRequest):
    """Search for similar transcript segments"""
    try:
        results = await vector_store.find_similar(request.query, request.n_results)
        return {"results": results}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/chat")
async def chat_endpoint(request: ChatRequest):
    """Chat with the language learning assistant"""
    try:
        response = chat.get_response(request.message, request.history)
        return {"response": response}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/generate_question")
async def generate_question(request: QuestionRequest):
    """Generate a new question"""
    try:
        question = question_generator.generate_similar_question(request.section_num, request.topic)
        if not question:
            raise HTTPException(status_code=404, detail="Failed to generate question")
        return {"question": question}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
