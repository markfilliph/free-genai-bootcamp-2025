from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from backend.routes import decks, flashcards, auth, generation
from backend.database import engine
from backend import models
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Create FastAPI app
app = FastAPI(
    title="Language Learning Flashcard Generator API",
    description="API for generating and managing language learning flashcards",
    version="1.0.0"
)

# CORS Configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:8083", "http://localhost:8080"],  # Frontend URLs
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Create database tables
models.Base.metadata.create_all(bind=engine)

# Include routers
app.include_router(auth.router, prefix="/api/auth", tags=["Authentication"])
app.include_router(decks.router, prefix="/api")
app.include_router(flashcards.router, prefix="/api")
app.include_router(generation.router, prefix="/api")

@app.get("/")
async def read_root():
    return {
        "message": "Welcome to Language Learning Flashcard Generator API",
        "docs": "/docs",
        "version": "1.0.0"
    }