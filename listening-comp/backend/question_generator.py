import json
import os
import time
from typing import Dict, List, Optional
import outlines
from transformers import AutoTokenizer
from backend.vector_store import QuestionVectorStore

class QuestionGenerator:
    def __init__(self):
        """Initialize question generator and vector store"""
        self._model = None  # Will be initialized when needed
        self._generator = None  # Will be initialized when needed
        self._vector_store = None
        self._initialized = False
        
        # Define our JSON schema for JLPT questions
        self._question_schema = {
            "type": "object",
            "properties": {
                "Introduction": {"type": "string"},
                "Conversation": {"type": "string"},
                "Question": {"type": "string"},
                "Options": {
                    "type": "array",
                    "items": {"type": "string"},
                    "minItems": 4,
                    "maxItems": 4
                },
                "CorrectAnswer": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4
                },
                "Explanation": {"type": "string"}
            },
            "required": ["Introduction", "Conversation", "Question", "Options", "CorrectAnswer", "Explanation"]
        }
        
    @property
    def model(self):
        """Lazy load the model"""
        if self._model is None:
            try:
                print("Initializing text generation model...")
                from text_generation import Client
                
                # Get API token from environment variable
                api_token = os.getenv('HUGGINGFACE_API_TOKEN')
                if not api_token:
                    raise ValueError('HUGGINGFACE_API_TOKEN environment variable not set')
                
                # Initialize TGI client with API token
                headers = {"Authorization": f"Bearer {api_token}"}
                endpoint = "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.2"
                client = Client(
                    endpoint,
                    headers=headers,
                    timeout=120  # Overall timeout
                )
                
                # Create text generation model with direct client
                self._model = client
                print("Model initialized successfully")
            except Exception as e:
                print(f"Error initializing model: {str(e)}")
                raise RuntimeError(f"Failed to initialize model: {str(e)}")
        return self._model
        
    @property
    def generator(self):
        """Lazy load the text generation model"""
        if self._generator is None:
            try:
                print("Initializing text generator...")
                # Create a structured generator with our schema
                schema_str = json.dumps(self._question_schema)
                self._generator = outlines.generate.json(self.model, schema_str)
                print("Generator initialized successfully")
            except Exception as e:
                print(f"Error initializing generator: {str(e)}")
                raise RuntimeError(f"Failed to initialize generator: {str(e)}")
        return self._generator
    
    @property
    def vector_store(self):
        """Lazy load the vector store and initialize with examples if needed"""
        if self._vector_store is None:
            self._vector_store = QuestionVectorStore()
            
            # Add example questions if store is empty
            if not self._vector_store.search_similar_questions(2, "Daily Conversation", n_results=1):
                example_questions = [
                    {
                        "Introduction": "At a train station",
                        "Conversation": "A: すみません、新宿駅はどこですか？\nB: あ、この先の角を右に曲がってください。",
                        "Question": "Where is Shinjuku station?",
                        "Options": ["Turn right at the corner ahead", "Turn left at the traffic light", "Go straight for 5 minutes", "Take the escalator"],
                        "CorrectAnswer": 1,
                        "Explanation": "Person B tells A to turn right at the corner ahead to get to Shinjuku station."
                    },
                    {
                        "Introduction": "At a restaurant",
                        "Conversation": "A: いらっしゃいませ。\nB: すみません、メニューをお願いします。",
                        "Question": "What does the customer want?",
                        "Options": ["A menu", "The bill", "Water", "Chopsticks"],
                        "CorrectAnswer": 1,
                        "Explanation": "The customer says 'メニューをお願いします' which means 'Please give me a menu.'"
                    }
                ]
                
                for q in example_questions:
                    self._vector_store.add_question(2, q, "Daily Conversation")
        return self._vector_store

    def _try_extract_json(self, text: str) -> Optional[str]:
        """Try various methods to extract JSON from text"""
        print(f"Attempting to extract JSON from text:\n{text}\n")

        # Method 1: Try direct parsing
        try:
            parsed = json.loads(text)
            return json.dumps(parsed)
        except json.JSONDecodeError:
            print("Direct parsing failed, trying other methods...")

        # Method 2: Find JSON-like structure with balanced braces
        stack = []
        start = -1
        for i, char in enumerate(text):
            if char == '{':
                if not stack:  # First opening brace
                    start = i
                stack.append(char)
            elif char == '}':
                if stack:
                    stack.pop()
                    if not stack and start != -1:  # Found a complete JSON object
                        try:
                            json_str = text[start:i+1]
                            parsed = json.loads(json_str)
                            return json.dumps(parsed)
                        except json.JSONDecodeError:
                            continue

        # Method 3: Try to fix common formatting issues
        try:
            # Remove any text before first { and after last }
            start_idx = text.find('{')
            end_idx = text.rfind('}')
            if start_idx != -1 and end_idx != -1:
                text = text[start_idx:end_idx + 1]

            # Fix common issues
            text = text.replace('\\n', '\n')
            text = text.replace('\t', ' ')
            text = text.replace('\"', '"')
            text = text.replace('\'', '"')
            text = ' '.join(text.split())  # Normalize whitespace

            parsed = json.loads(text)
            return json.dumps(parsed)
        except json.JSONDecodeError:
            print("All JSON extraction methods failed")
            return None

    def _clean_json_text(self, text: str) -> str:
        """Clean text to extract valid JSON"""
        print(f"Raw text to clean: {text}")
        
        # Try to extract JSON using various methods
        result = self._try_extract_json(text)
        if result is None:
            raise ValueError("No valid JSON found in text")
            
        # Validate the extracted JSON
        try:
            parsed = json.loads(result)
            if not isinstance(parsed, dict):
                raise ValueError("Extracted JSON is not an object")
            return json.dumps(parsed)
        except json.JSONDecodeError as e:
            raise ValueError(f"Invalid JSON structure: {e}")

    def _generate_question(self, prompt: str) -> Optional[str]:
        """Generate a question using the text generation model.

        Args:
            prompt: Prompt to generate question from

        Returns:
            Generated question text or None if generation failed
        """
        try:
            # Format prompt for Mistral
            template = '''<s>[INST] You are a JLPT question generator. Generate a JLPT question about {topic}. Use natural Japanese with proper keigo (polite language). Include a brief introduction, a realistic conversation, a clear question, and 4 answer options where only one is correct. Return ONLY valid JSON in this exact format:
            {{
                "Introduction": "Brief context like At a restaurant",
                "Conversation": "A: Japanese dialogue here\nB: Response in Japanese here",
                "Question": "Clear question in English about the dialogue",
                "Options": ["Correct answer", "Wrong but plausible answer", "Another wrong answer", "Another wrong answer"],
                "CorrectAnswer": 1,
                "Explanation": "Clear explanation of why the correct answer is right"
            }}
            [/INST]\n'''
            full_prompt = template.format(topic=prompt)
            
            print("Generating text with model...")
            print(f"Using prompt:\n{full_prompt}\n")

            # Generate with TGI client with retries
            max_retries = 3
            retry_delay = 2  # seconds
            
            for attempt in range(max_retries):
                try:
                    print(f"Attempt {attempt + 1} of {max_retries}...")
                    # Try different generation parameters on each retry
                    if attempt == 0:
                        params = {
                            "max_new_tokens": 1024,
                            "temperature": 0.7,
                            "top_p": 0.95
                        }
                    elif attempt == 1:
                        params = {
                            "max_new_tokens": 1024,
                            "temperature": 0.8,
                            "top_k": 50
                        }
                    else:
                        params = {
                            "max_new_tokens": 1024,
                            "temperature": 0.9,
                            "top_p": 0.99
                        }
                    
                    response = self.model.generate(full_prompt, **params)
                    if not response or not response.generated_text:
                        raise ValueError("Empty response from model")
                    
                    result = response.generated_text
                    break
                except Exception as e:
                    print(f"Error on attempt {attempt + 1}: {str(e)}")
                    if attempt < max_retries - 1:
                        print(f"Retrying in {retry_delay} seconds...")
                        time.sleep(retry_delay)
                        retry_delay *= 2  # Exponential backoff
                    else:
                        raise RuntimeError(f"Failed after {max_retries} attempts: {str(e)}")
            
            # Clean up the response
            result = result.strip()
            if result.startswith('```json'):
                result = result[7:]
            if result.endswith('```'):
                result = result[:-3]
            
            # Try to parse and re-serialize to ensure valid JSON
            try:
                parsed = json.loads(result)
                return json.dumps(parsed, ensure_ascii=False)
            except json.JSONDecodeError:
                # If parsing fails, return the cleaned string for further processing
                return result

        except Exception as e:
            print(f"Error in _generate_question: {str(e)}")
            return None
                    
    def _validate_question(self, question: Dict) -> Dict:
        """Validate the generated question"""
        try:
            # Validate the question format
            if 'Introduction' in question:
                required_fields = ["Introduction", "Conversation", "Question", "Options", "CorrectAnswer", "Explanation"]
            else:
                required_fields = ["Situation", "Content", "Question", "Options", "CorrectAnswer", "Explanation"]
            
            missing_fields = [field for field in required_fields if field not in question]
            if missing_fields:
                raise ValueError(f"Generated question missing required fields: {missing_fields}")
            
            if not isinstance(question["Options"], list) or len(question["Options"]) != 4:
                raise ValueError("Options must be a list of exactly 4 items")
            
            if not isinstance(question["CorrectAnswer"], int) or not 1 <= question["CorrectAnswer"] <= 4:
                raise ValueError("CorrectAnswer must be an integer between 1 and 4")
            
            # If we got here, we have valid JSON
            print("Successfully generated valid JSON")
            return question
        
        except Exception as e:
            print(f"Error validating question: {str(e)}")
            raise Exception(f"Failed to validate question: {str(e)}")

    def generate_similar_question(self, section_num: int, topic: str) -> Optional[Dict]:
        """Generate a new question similar to existing ones on a given topic"""
        try:
            print(f"Generating question for section {section_num}, topic: {topic}")
            
            # Initialize vector store if needed
            if self._vector_store is None:
                print("Initializing vector store...")
                self._vector_store = self.vector_store
                print("Vector store initialized successfully")
            
            # Get similar questions for context
            print("Searching for similar questions...")
            similar_questions = self.vector_store.search_similar_questions(section_num, topic)
            
            if not similar_questions:
                print("No similar questions found, using default example")
                # Return a default example question if no similar questions found
                if section_num == 2:
                    return {
                        "Introduction": "At a restaurant",
                        "Conversation": "A: すみません、メニューをお願いします。\nB: はい、少々お待ちください。",
                        "Question": "What is the customer asking for?",
                        "Options": ["A menu", "The bill", "Water", "A table"],
                        "CorrectAnswer": 1,
                        "Explanation": "The customer says 'メニューをお願いします' which means 'Please give me a menu.'"
                    }
                else:
                    return {
                        "Situation": "At a train station",
                        "Content": "電車は10分遅れています。",
                        "Question": "What is the announcement about?",
                        "Options": ["The train is 10 minutes late", "The train is arriving", "The train is cancelled", "The train is on time"],
                        "CorrectAnswer": 1,
                        "Explanation": "The announcement '電車は10分遅れています' means 'The train is 10 minutes late.'"
                    }
            
            print("Found similar questions, preparing context")
            # Use similar questions as context for generating new question
            context = json.dumps(similar_questions, indent=2, ensure_ascii=False)
            
            # Create prompt for generating new question
            prompt = f"""You are a JLPT listening comprehension question generator. Create a new question about {topic} following this exact format:

            For section 2, the question must be a valid JSON object with these exact fields:
            {{
                "Introduction": "Brief context like 'At a restaurant' or 'At a train station'",
                "Conversation": "A: Japanese dialogue here\nB: Response in Japanese here",
                "Question": "Clear question in English about the dialogue",
                "Options": ["Correct answer", "Wrong but plausible answer", "Another wrong answer", "Another wrong answer"],
                "CorrectAnswer": 1,
                "Explanation": "Clear explanation of why the correct answer is right"
            }}

            For other sections, use this format:
            {{
                "Situation": "Brief context of the situation",
                "Content": "Japanese content or announcement",
                "Question": "Clear question in English",
                "Options": ["Correct answer", "Wrong but plausible answer", "Another wrong answer", "Another wrong answer"],
                "CorrectAnswer": 1,
                "Explanation": "Clear explanation of why the correct answer is right"
            }}

            Here are some example questions for reference:
            {context}

            Generate a new question that:
            1. Is different from the examples
            2. Uses natural Japanese appropriate for the JLPT level
            3. Has a clear correct answer
            4. Has plausible but clearly incorrect alternative options
            5. Follows the exact JSON format shown above
            6. Contains ONLY the JSON object, no other text
            """

            print("Generating new question...")
            # Generate new question using the model
            try:
                response = self._generate_question(prompt)
                if not response:
                    print("No response from question generator")
                    raise Exception("Failed to generate question - no response")

                print("Parsing generated response...")
                # Parse the generated response into a question dict
                question = json.loads(response)
                
                # Validate the question format
                validated_question = self._validate_question(question)
                
                print("Adding question to vector store...")
                # Add the question to the vector store for future use
                self.vector_store.add_question(section_num, validated_question, topic)
                
                print("Question generated and stored successfully")
                return validated_question
                
            except json.JSONDecodeError as e:
                print(f"Error parsing generated question JSON: {str(e)}\nResponse was: {response}")
                raise Exception(f"Failed to parse generated question: {str(e)}")
            
        except Exception as e:
            print(f"Error in generate_similar_question: {str(e)}")
            raise Exception(f"Failed to generate question: {str(e)}")

    def get_feedback(self, question: Dict, selected_answer: int) -> Dict:
        """Generate feedback for the selected answer"""
        if not question or 'Options' not in question:
            return None

        # Create prompt for generating feedback
        prompt = f"""Given this JLPT listening question and the selected answer, provide feedback explaining if it's correct 
        and why. Keep the explanation clear and concise.
        
        """
        if 'Introduction' in question:
            prompt += f"Introduction: {question['Introduction']}\n"
            prompt += f"Conversation: {question['Conversation']}\n"
        else:
            prompt += f"Situation: {question['Situation']}\n"
        
        prompt += f"Question: {question['Question']}\n"
        prompt += "Options:\n"
        for i, opt in enumerate(question['Options'], 1):
            prompt += f"{i}. {opt}\n"
        
        prompt += f"\nSelected Answer: {selected_answer}\n"
        prompt += "\nProvide feedback in JSON format with these fields:\n"
        prompt += "- correct: true/false\n"
        prompt += "- explanation: brief explanation of why the answer is correct/incorrect\n"
        prompt += "- correct_answer: the number of the correct option (1-4)\n"

        # Compare with correct answer
        correct_answer = question.get('CorrectAnswer', 1)
        is_correct = selected_answer == correct_answer
        
        # Get explanation from question if available
        explanation = question.get('Explanation', '')
        if not explanation:
            if is_correct:
                explanation = "That's correct! Good job!"
            else:
                explanation = f"The correct answer is option {correct_answer}."
        
        return {
            "correct": is_correct,
            "explanation": explanation,
            "correct_answer": correct_answer
        }
