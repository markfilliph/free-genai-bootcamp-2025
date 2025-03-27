import requests
import json
from typing import Dict, List, Optional
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class OllamaService:
    def __init__(self, base_url: str = "http://localhost:11434"):
        self.base_url = base_url
        self.api_url = f"{base_url}/api/generate"
        self.model = "mistral"  # Default model based on available models
    
    def set_model(self, model_name: str):
        """Set the model to use for generation."""
        self.model = model_name
    
    async def generate_example_sentences(self, word: str, count: int = 3) -> List[str]:
        """Generate example sentences for a Spanish word."""
        prompt = f"""Generate {count} example sentences in Spanish using the word '{word}'. 
        Each sentence should be natural, conversational, and demonstrate proper usage of the word.
        For each sentence, also provide an English translation.
        Format the output as a JSON array of objects with 'spanish' and 'english' keys.
        """
        
        response = self._call_ollama(prompt)
        try:
            # Try to parse as JSON
            sentences = json.loads(response)
            return sentences
        except json.JSONDecodeError:
            # Fallback to text parsing if JSON parsing fails
            logger.warning("Failed to parse JSON from Ollama response, falling back to text parsing")
            sentences = []
            lines = response.strip().split('\n')
            for line in lines:
                if line.strip() and ':' in line:
                    parts = line.split(':', 1)
                    if len(parts) == 2:
                        spanish = parts[0].strip()
                        english = parts[1].strip()
                        sentences.append({"spanish": spanish, "english": english})
            
            return sentences[:count]
    
    async def generate_verb_conjugations(self, verb: str) -> str:
        """Generate verb conjugations for a Spanish verb."""
        prompt = f"""Generate the conjugation table for the Spanish verb '{verb}' in present, preterite, 
        imperfect, future, and conditional tenses. Format it in a clear, readable way with tense names as headers.
        """
        
        return self._call_ollama(prompt)
    
    async def generate_cultural_note(self, word: str) -> str:
        """Generate cultural context or note for a Spanish word or phrase."""
        prompt = f"""Provide a brief cultural note or context about the Spanish word '{word}'. 
        Include any regional variations, idiomatic usage, or cultural significance.
        Keep it concise but informative, around 2-3 sentences.
        """
        
        return self._call_ollama(prompt)
    
    def _call_ollama(self, prompt: str) -> str:
        """Make a call to the Ollama API."""
        try:
            payload = {
                "model": self.model,
                "prompt": prompt,
                "stream": False
            }
            
            response = requests.post(self.api_url, json=payload)
            response.raise_for_status()
            
            result = response.json()
            return result.get("response", "")
        except requests.RequestException as e:
            logger.error(f"Error calling Ollama API: {str(e)}")
            return f"Error: {str(e)}"
        except Exception as e:
            logger.error(f"Unexpected error: {str(e)}")
            return f"Unexpected error: {str(e)}"

# Create a singleton instance
ollama_service = OllamaService()
