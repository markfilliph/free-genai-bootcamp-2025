from typing import Optional, Dict, List
from transformers import AutoTokenizer, T5ForConditionalGeneration
import torch
import json
import os

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
        """Initialize the model and tokenizer"""
        if not self._initialized:
            # Use smaller model for faster responses
            model_name = "google/flan-t5-small"
            self._tokenizer = AutoTokenizer.from_pretrained(model_name)
            self._model = T5ForConditionalGeneration.from_pretrained(model_name)
            
            # Enable model optimizations
            self._model.eval()
            if torch.cuda.is_available():
                self._model = self._model.cuda()
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

    def generate_response(self, prompt: str) -> Optional[str]:
        """Generate a response using local model"""
        try:
            # Check for pre-defined responses first
            quick_response = self._find_similar_question(prompt)
            if quick_response:
                return quick_response

            # Initialize model if needed
            self._initialize_model()

            # Prepare input
            context = "Answer this Japanese language question concisely: "
            input_text = context + prompt
            
            # Tokenize and generate
            inputs = self._tokenizer(input_text, return_tensors="pt", max_length=100, truncation=True)
            if torch.cuda.is_available():
                inputs = inputs.to('cuda')
            
            with torch.no_grad():
                outputs = self._model.generate(
                    **inputs,
                    max_length=100,
                    num_beams=2,
                    temperature=0.7,
                    top_k=50,
                    top_p=0.9,
                    do_sample=True,
                    early_stopping=True
                )
            
            response = self._tokenizer.decode(outputs[0], skip_special_tokens=True)
            return response.strip()
            
        except Exception as e:
            print(f"Error generating response: {str(e)}")
            return None

if __name__ == "__main__":
    # Test the chat
    chat = Chat()
    response = chat.generate_response("How do I say 'hello' in Japanese?")
    print(response)
