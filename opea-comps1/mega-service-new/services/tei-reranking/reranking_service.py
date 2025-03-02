from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
import uvicorn
from transformers import AutoModelForSequenceClassification, AutoTokenizer
import torch

app = FastAPI()

# Load model and tokenizer
model_name = "cross-encoder/ms-marco-MiniLM-L-6-v2"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = AutoModelForSequenceClassification.from_pretrained(model_name)

class RerankRequest(BaseModel):
    query: str
    contexts: List[str]

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.post("/rerank")
async def rerank_context(request: RerankRequest):
    try:
        pairs = [[request.query, doc] for doc in request.contexts]
        features = tokenizer.batch_encode_plus(
            pairs,
            max_length=512,
            padding=True,
            truncation=True,
            return_tensors="pt"
        )
        
        with torch.no_grad():
            scores = model(**features).logits.squeeze()
        
        # Sort contexts by score
        ranked_pairs = list(zip(request.contexts, scores.tolist()))
        ranked_pairs.sort(key=lambda x: x[1], reverse=True)
        
        return {
            "reranked_contexts": [context for context, _ in ranked_pairs]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8082)
