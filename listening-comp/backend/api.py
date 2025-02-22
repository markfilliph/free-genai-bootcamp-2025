from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Optional
from backend.question_generator import QuestionGenerator

app = FastAPI()
question_generator = QuestionGenerator()

class QuestionRequest(BaseModel):
    section_num: int
    topic: str

@app.post("/api/generate_question")
async def generate_question(request: QuestionRequest) -> Dict:
    """Generate a new question based on section number and topic"""
    print(f"Received request for section {request.section_num}, topic: {request.topic}")
    try:
        # Initialize question generator if needed
        if not question_generator._initialized:
            print("Initializing question generator...")
            question_generator._generator = question_generator.generator
            question_generator._vector_store = question_generator.vector_store
            question_generator._initialized = True
            print("Question generator initialized successfully")

        print("Generating question...")
        question = question_generator.generate_similar_question(
            request.section_num,
            request.topic
        )
        if not question:
            print("Failed to generate question - returned None")
            raise HTTPException(
                status_code=500,
                detail="Failed to generate question - no question generated"
            )
        print("Question generated successfully")
        return {"question": question}
    except Exception as e:
        print(f"Error in generate_question endpoint: {str(e)}")
        raise HTTPException(
            status_code=500,
            detail=f"Error generating question: {str(e)}"
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
