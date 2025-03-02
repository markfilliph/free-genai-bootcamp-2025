from fastapi import FastAPI, HTTPException
from sentence_transformers import SentenceTransformer
from pydantic import BaseModel
import uvicorn
import torch

app = FastAPI()
model = SentenceTransformer('all-MiniLM-L6-v2')

class TextRequest(BaseModel):
    text: str

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.post("/embed")
async def embed_text(request: TextRequest):
    try:
        embeddings = model.encode([request.text])
        return {"embeddings": embeddings[0].tolist()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)
