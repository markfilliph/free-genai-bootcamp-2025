from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Optional
from .question_generator import QuestionGenerator

app = FastAPI()
question_generator = QuestionGenerator()

class QuestionRequest(BaseModel):
    section_num: int
    topic: str

@app.post("/api/generate_question")
async def generate_question(request: QuestionRequest) -> Dict:
    """Generate a new question based on section number and topic"""
    try:
        question = question_generator.generate_similar_question(
            request.section_num,
            request.topic
        )
        if not question:
            raise HTTPException(
                status_code=500,
                detail="Failed to generate question"
            )
        return {"question": question}
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error generating question: {str(e)}"
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
