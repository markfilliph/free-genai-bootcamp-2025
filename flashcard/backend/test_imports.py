"""
Test script to check imports and identify issues with FastAPI and Pydantic.
"""
import sys
print(f"Python version: {sys.version}")

try:
    import fastapi
    print(f"FastAPI version: {fastapi.__version__}")
except ImportError as e:
    print(f"Error importing FastAPI: {e}")

try:
    import pydantic
    print(f"Pydantic version: {pydantic.__version__}")
except ImportError as e:
    print(f"Error importing Pydantic: {e}")

try:
    from backend import schemas
    print("Successfully imported schemas")
except Exception as e:
    print(f"Error importing schemas: {e}")

try:
    from backend import models
    print("Successfully imported models")
except Exception as e:
    print(f"Error importing models: {e}")

try:
    from backend.routes import auth, decks, flashcards, generation
    print("Successfully imported routes")
except Exception as e:
    print(f"Error importing routes: {e}")

try:
    from backend.services.ollama_service import ollama_service
    print("Successfully imported ollama_service")
except Exception as e:
    print(f"Error importing ollama_service: {e}")

print("Import test complete")
