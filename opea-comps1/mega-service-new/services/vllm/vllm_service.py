# Current Status:
# ❌ Service is failing due to device type inference errors
# ❌ Original VLLM implementation required GPU
# ⚠️ Attempted fix: Switched to Transformers library with CPU support
# ⚠️ Using TinyLlama model instead of Llama-2 for lighter resource usage
# TODO: Test new implementation with CPU-only setup

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import uvicorn
from transformers import AutoModelForCausalLM, AutoTokenizer
import torch
import os

app = FastAPI()

# Initialize model and tokenizer
model_id = os.getenv("LLM_MODEL_ID", "TinyLlama/TinyLlama-1.1B-Chat-v1.0")
tokenizer = AutoTokenizer.from_pretrained(model_id)
model = AutoModelForCausalLM.from_pretrained(
    model_id,
    torch_dtype=torch.float32,
    device_map='auto',
    load_in_8bit=True
)

class Message(BaseModel):
    role: str
    content: str

class CompletionRequest(BaseModel):
    messages: List[Message]
    temperature: Optional[float] = 0.7
    max_tokens: Optional[int] = 150

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.post("/v1/completions")
async def complete(request: CompletionRequest):
    try:
        # Convert messages to prompt
        prompt = "\n".join([f"{msg.role}: {msg.content}" for msg in request.messages])
        
        # Tokenize input
        inputs = tokenizer(prompt, return_tensors="pt", truncation=True, max_length=512)
        inputs = inputs.to(model.device)
        
        # Generate completion
        with torch.no_grad():
            outputs = model.generate(
                **inputs,
                max_new_tokens=request.max_tokens,
                temperature=request.temperature,
                do_sample=True,
                pad_token_id=tokenizer.pad_token_id
            )
        
        # Decode output
        generated_text = tokenizer.decode(outputs[0][inputs["input_ids"].shape[1]:], skip_special_tokens=True)
        
        return {
            "text": generated_text,
            "finish_reason": "length" if len(generated_text.split()) >= request.max_tokens else "stop",
            "usage": {
                "prompt_tokens": inputs["input_ids"].shape[1],
                "completion_tokens": len(generated_text.split()),
                "total_tokens": inputs["input_ids"].shape[1] + len(generated_text.split())
            }
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=80)
