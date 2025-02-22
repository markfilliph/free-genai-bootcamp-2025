import json
import os
from typing import Dict, List, Optional
from transformers import pipeline
from vector_store import QuestionVectorStore

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

    def _generate_question(self, prompt: str) -> Optional[str]:
        """Generate a question using local T5 model"""
        try:
            # Add specific instruction for JLPT question generation
            full_prompt = f"Generate a JLPT listening comprehension question: {prompt}"
            
            # Generate text
            response = self.generator(
                full_prompt,
                max_length=512,
                num_return_sequences=1,
                temperature=0.7
            )
            
            return response[0]['generated_text']
        except Exception as e:
            print(f"Error generating question: {str(e)}")
            return None

    def generate_similar_question(self, section_num: int, topic: str) -> Optional[Dict]:
        """Generate a new question similar to existing ones on a given topic"""
        try:
            # Get similar questions for context
            similar_questions = self.vector_store.search_similar_questions(section_num, topic)
            
            if not similar_questions:
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
            
            # Use similar questions as context for generating new question
            context = json.dumps(similar_questions, indent=2, ensure_ascii=False)
            
            # Create prompt for generating new question
            prompt = f"""Based on the following example JLPT listening questions, create a new question about {topic}.
            The question should follow the same format but be different from the examples.
            Make sure the question tests listening comprehension and has a clear correct answer.
            
            Example Questions:
            {context}
            
            Generate a new question following the exact same format as above. Include all components (Introduction/Situation, 
            Conversation/Question, and Options). Make sure the question is challenging but fair, and the options are plausible 
            but with only one clearly correct answer.
            """

            # Generate new question using the model
            response = self._generate_question(prompt)
            if not response:
                return None

            try:
                # Parse the generated response into a question dict
                question = json.loads(response)
                
                # Validate the question format
                if section_num == 2:
                    required_fields = ["Introduction", "Conversation", "Question", "Options", "CorrectAnswer", "Explanation"]
                else:
                    required_fields = ["Situation", "Content", "Question", "Options", "CorrectAnswer", "Explanation"]
                
                if not all(field in question for field in required_fields):
                    print(f"Generated question missing required fields: {required_fields}")
                    return None
                
                # Add the question to the vector store for future use
                self.vector_store.add_question(section_num, question, topic)
                
                return question
            except json.JSONDecodeError as e:
                print(f"Error parsing generated question: {str(e)}")
                return None
            
        except Exception as e:
            print(f"Error generating question: {str(e)}")
            return None
        except Exception as e:
            print(f"Error parsing generated question: {str(e)}")
            return None

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
