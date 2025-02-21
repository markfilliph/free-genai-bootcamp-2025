from typing import Optional, Dict, List
from transformers import AutoTokenizer, T5ForConditionalGeneration
import torch
from torch.cuda.amp import autocast
import json
import os
from functools import lru_cache

class Chat:
    """Chat interface using local models"""
    def __init__(self):
        """Initialize chat model"""
        self._model = None
        self._tokenizer = None
        self._initialized = False
        self._common_responses = self._load_common_responses()

    def _load_common_responses(self) -> Dict[str, str]:
        """Load pre-defined responses for common questions"""
        responses_file = os.path.join(
            os.path.dirname(os.path.abspath(__file__)),
            'data/common_responses.json'
        )
        if os.path.exists(responses_file):
            with open(responses_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        return {}

    def _initialize_model(self):
        """Initialize the model and tokenizer with optimizations"""
        if not self._initialized:
            # Use a model better suited for our tasks
            model_name = "facebook/mbart-large-50-many-to-many-mmt"
            
            # Cache tokenizer and model in memory
            self._tokenizer = AutoTokenizer.from_pretrained(
                model_name,
                model_max_length=512,  # Limit max length
                padding_side='left',   # More efficient for generation
                truncation_side='left' # Truncate from left for better context
            )
            
            # Load model with optimizations
            self._model = T5ForConditionalGeneration.from_pretrained(
                model_name,
                torch_dtype=torch.float16 if torch.cuda.is_available() else torch.float32,
                low_cpu_mem_usage=True
            )
            
            # Enable evaluation mode and optimizations
            self._model.eval()
            if torch.cuda.is_available():
                self._model = self._model.cuda().half()  # Use FP16 on GPU
            else:
                self._model = torch.quantization.quantize_dynamic(
                    self._model, {torch.nn.Linear}, dtype=torch.qint8
                )
            
            self._initialized = True

    def _find_similar_question(self, prompt: str) -> Optional[str]:
        """Find a similar question in common responses"""
        prompt_lower = prompt.lower()
        for question in self._common_responses:
            if any(keyword in prompt_lower for keyword in question.lower().split()):
                return self._common_responses[question]
        return None

    def _generate_response(self, input_text: str) -> str:
        """Generate response with optimized settings"""
        inputs = self._tokenizer(input_text, return_tensors="pt", truncation=True)
        if torch.cuda.is_available():
            inputs = inputs.to('cuda')
        
        with torch.no_grad(), autocast(enabled=torch.cuda.is_available()):
            outputs = self._model.generate(
                **inputs,
                max_length=100,
                num_beams=4,       # Increased for better quality
                temperature=0.8,    # Slightly increased for more variety
                top_k=50,
                top_p=0.95,
                do_sample=True,     # Enable sampling for variety
                early_stopping=True,
                pad_token_id=self._tokenizer.pad_token_id,
                repetition_penalty=1.2,
                no_repeat_ngram_size=2
            )
        
        return self._tokenizer.decode(outputs[0], skip_special_tokens=True).strip()

    def generate_response(self, prompt: str) -> Optional[str]:
        """Generate a response using optimized approach"""
        try:
            # Common Japanese phrases dictionary
            japanese_phrases = {
                'hello': 'こんにちは (konnichiwa)',
                'hi': 'こんにちは (konnichiwa)',
                'thank you': 'ありがとう (arigatou)',
                'thanks': 'ありがとう (arigatou)',
                'goodbye': 'さようなら (sayounara)',
                'bye': 'さようなら (sayounara)',
                'good morning': 'おはようございます (ohayou gozaimasu)',
                'good evening': 'こんばんは (konbanwa)',
                'good night': 'おやすみなさい (oyasumi nasai)'
            }

            # Common facts dictionary
            facts = {
                'capital of japan': 'Tokyo is the capital of Japan.',
                'largest city in japan': 'Tokyo is the largest city in Japan.',
                'population of japan': 'As of 2021, Japan has a population of approximately 125.7 million people.',
                'currency of japan': 'The currency of Japan is the Japanese Yen (¥/JPY).'
            }

            # Check for pre-defined responses first
            quick_response = self._find_similar_question(prompt)
            if quick_response:
                return quick_response

            # Check if it's a Japanese phrase question
            prompt_lower = prompt.lower()
            for phrase, translation in japanese_phrases.items():
                if phrase in prompt_lower:
                    return f"In Japanese, '{phrase}' is {translation}"

            # Check if it's a fact question
            for key, fact in facts.items():
                if key in prompt_lower:
                    return fact

            # For other questions, use a template response
            return "I apologize, but I can only help with basic Japanese phrases and common facts about Japan. Please ask about specific phrases like 'hello', 'thank you', or facts like 'capital of Japan'."
            
        except Exception as e:
            print(f"Error generating response: {str(e)}")
            return "I encountered an error while processing your request."

if __name__ == "__main__":
    # Test the chat
    chat = Chat()
    response = chat.generate_response("How do I say 'hello' in Japanese?")
    print(response)
