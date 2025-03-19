from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from backend.routes import decks
from backend.database import engine
from backend import models

app = FastAPI()

# CORS Configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

models.Base.metadata.create_all(bind=engine)

app.include_router(decks.router, prefix="/api")

@app.get("/")
async def read_root():
    return {"message": "Welcome to Flashcard API"}