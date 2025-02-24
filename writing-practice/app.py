import streamlit as st
import requests
import json
from PIL import Image
from typing import Dict, List
import random  # Add this at the top with other imports

# At the top of your file, add mock data
MOCK_VOCABULARY = [
    {"japanese": "犬", "english": "dog"},
    {"japanese": "猫", "english": "cat"},
    {"japanese": "魚", "english": "fish"},
    {"japanese": "本", "english": "book"},
    {"japanese": "車", "english": "car"},
    {"japanese": "水", "english": "water"},
    {"japanese": "友達", "english": "friend"},
    {"japanese": "学校", "english": "school"}
]

# Initialize session state if not already done
if 'current_state' not in st.session_state:
    st.session_state.current_state = 'setup'
if 'current_sentence' not in st.session_state:
    st.session_state.current_sentence = ''
if 'vocabulary' not in st.session_state:
    st.session_state.vocabulary = []

class AppStates:
    SETUP = 'setup'
    PRACTICE = 'practice'
    REVIEW = 'review'

def fetch_vocabulary(group_id: str) -> List[Dict]:
    """Fetch vocabulary from the API"""
    # For testing, return mock data instead of making API call
    return MOCK_VOCABULARY

def generate_sentence(word: str) -> str:
    """Generate a sentence using the LLM prompt"""
    prompt = f"""Generate a simple sentence using the following word: {word}
    The grammar should be scoped to JLPTN5 grammar.
    You can use the following vocabulary to construct a simple sentence:
    - simple objects eg. book, car, ramen, sushi
    - simple verbs, to drink, to eat, to meet
    - simple times eg. tomorrow, today, yesterday"""
    
    return f"I will eat sushi with {word} tomorrow."

def grade_submission(image) -> Dict:
    """Grade the submitted image using the Grading System"""
    return {
        "transcription": "私は明日寿司を食べます",
        "translation": "I will eat sushi tomorrow",
        "grade": {
            "score": "A",
            "feedback": "Good attempt! The sentence structure is correct."
        }
    }

def setup_state():
    st.title("Japanese Learning App")
    
    if st.button("Generate Sentence"):
        # Randomly select a word from vocabulary
        if st.session_state.vocabulary:
            word = random.choice(st.session_state.vocabulary)  # Random selection instead of [0]
            st.session_state.current_sentence = generate_sentence(word['english'])
            st.session_state.current_state = AppStates.PRACTICE
        else:
            st.error("No vocabulary available. Please check if the API server is running.")

def practice_state():
    st.title("Practice Writing")
    st.write(st.session_state.current_sentence)
    
    uploaded_file = st.file_uploader("Upload your written answer", type=['png', 'jpg', 'jpeg'])
    
    if uploaded_file is not None and st.button("Submit for Review"):
        image = Image.open(uploaded_file)
        st.session_state.current_grade = grade_submission(image)
        st.session_state.current_state = AppStates.REVIEW

def review_state():
    st.title("Review")
    st.write(f"Original sentence: {st.session_state.current_sentence}")
    
    grade_data = st.session_state.current_grade
    
    st.subheader("Your Submission")
    st.write(f"Transcription: {grade_data['transcription']}")
    st.write(f"Translation: {grade_data['translation']}")
    
    st.subheader("Grade")
    st.write(f"Score: {grade_data['grade']['score']}")
    st.write(f"Feedback: {grade_data['grade']['feedback']}")
    
    if st.button("Next Question"):
        st.session_state.current_state = AppStates.SETUP

def main():
    # Fetch vocabulary if not already done
    if not st.session_state.vocabulary:
        st.session_state.vocabulary = fetch_vocabulary("default-group-id")
    
    # Route to appropriate state
    if st.session_state.current_state == AppStates.SETUP:
        setup_state()
    elif st.session_state.current_state == AppStates.PRACTICE:
        practice_state()
    elif st.session_state.current_state == AppStates.REVIEW:
        review_state()

if __name__ == "__main__":
    main()