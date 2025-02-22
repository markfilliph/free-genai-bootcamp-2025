import json
import os
from typing import Dict, List, Optional
from transformers import pipeline
from backend.vector_store import QuestionVectorStore

class QuestionGenerator:
    def __init__(self):
        """Initialize question generator and vector store"""
        self._generator = None
        self._vector_store = None
        self._initialized = False
        
    @property
    def generator(self):
        """Lazy load the generator model"""
        if self._generator is None:
            try:
                print("Initializing text generation model...")
                self._generator = pipeline(
                    'text2text-generation',
                    model='google/flan-t5-small',
                    device=-1  # Use CPU
                )
                print("Model initialized successfully")
            except Exception as e:
                print(f"Error initializing model: {str(e)}")
                raise RuntimeError(f"Failed to initialize text generation model: {str(e)}")
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

    def _clean_json_text(self, text: str) -> str:
        """Clean text to extract valid JSON"""
        # Find the first { and last } to extract JSON object
        start_idx = text.find('{')
        end_idx = text.rfind('}')
        
        if start_idx == -1 or end_idx == -1:
            raise ValueError("No JSON object found in text")
            
        # Extract the potential JSON object
        json_str = text[start_idx:end_idx + 1]
        
        # Fix common JSON formatting issues
        json_str = json_str.replace('\\n', '\n')  # Fix escaped newlines
        json_str = json_str.replace('\t', ' ')    # Replace tabs with spaces
        
        # Remove any extra whitespace between tokens
        json_str = ' '.join(json_str.split())
        
        # Attempt to parse and re-serialize to ensure valid JSON
        parsed = json.loads(json_str)
        return json.dumps(parsed)

    def _generate_question(self, prompt: str) -> Optional[str]:
        """Generate a question using local T5 model"""
        try:
            print("Initializing generator if needed...")
            if self._generator is None:
                self._generator = self.generator
                print("Generator initialized successfully")

            # Create a very explicit prompt for JSON generation
            full_prompt = f"""You are a JLPT question generator API that must output ONLY valid JSON.
            Generate a question about: {prompt}
            
            Your response must be a valid JSON object with exactly these fields:
            {{
                "Introduction": "(brief situation context)",
                "Conversation": "A: (Japanese text)\nB: (Japanese response)",
                "Question": "(clear English question)",
                "Options": [
                    "(correct answer in English)",
                    "(wrong but plausible answer)",
                    "(another wrong answer)",
                    "(another wrong answer)"
                ],
                "CorrectAnswer": 1,
                "Explanation": "(why the correct answer is right)"
            }}
            
            Important:
            1. Output ONLY the JSON object
            2. Use real Japanese text appropriate for JLPT
            3. Make sure the conversation is natural
            4. Ensure all Japanese has appropriate keigo/politeness
            5. Make the question challenging but fair
            """
            
            print("Generating text with model...")
            # Generate multiple attempts and take the best one
            for attempt in range(3):
                try:
                    response = self.generator(
                        full_prompt,
                        max_length=512,
                        num_return_sequences=1,
                        temperature=0.7
                    )
                    
                    generated_text = response[0]['generated_text'].strip()
                    print(f"Attempt {attempt + 1} generated text: {generated_text}")
                    
                    # Try to clean and validate the JSON
                    cleaned_json = self._clean_json_text(generated_text)
                    
                    # If we got here, we have valid JSON
                    print(f"Successfully generated valid JSON on attempt {attempt + 1}")
                    return cleaned_json
                    
                except Exception as e:
                    print(f"Attempt {attempt + 1} failed: {str(e)}")
                    if attempt == 2:  # Last attempt
                        raise
            
        except Exception as e:
            print(f"Error in _generate_question: {str(e)}")
            raise Exception(f"Failed to generate question text: {str(e)}")

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
                if section_num == 2:
                    required_fields = ["Introduction", "Conversation", "Question", "Options", "CorrectAnswer", "Explanation"]
                else:
                    required_fields = ["Situation", "Content", "Question", "Options", "CorrectAnswer", "Explanation"]
                
                missing_fields = [field for field in required_fields if field not in question]
                if missing_fields:
                    print(f"Generated question missing required fields: {missing_fields}")
                    raise Exception(f"Generated question missing required fields: {missing_fields}")
                
                print("Adding question to vector store...")
                # Add the question to the vector store for future use
                self.vector_store.add_question(section_num, question, topic)
                
                print("Question generated and stored successfully")
                return question
                
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
