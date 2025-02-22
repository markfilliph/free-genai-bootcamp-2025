import streamlit as st
import sys
import os
import json
from datetime import datetime
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.question_generator import QuestionGenerator
from backend.chat import Chat
import requests

# Page config
st.set_page_config(
    page_title="Japanese Learning Assistant",
    page_icon="🎌",
    layout="wide",
    initial_sidebar_state="expanded"
)

def load_stored_questions():
    """Load previously stored questions from JSON file"""
    questions_file = os.path.join(
        os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
        "backend/data/stored_questions.json"
    )
    if os.path.exists(questions_file):
        with open(questions_file, 'r', encoding='utf-8') as f:
            return json.load(f)
    return {}

def save_question(question, practice_type, topic):
    """Save a generated question to JSON file"""
    questions_file = os.path.join(
        os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
        "backend/data/stored_questions.json"
    )
    
    # Load existing questions
    stored_questions = load_stored_questions()
    
    # Create a unique ID for the question using timestamp
    question_id = datetime.now().strftime("%Y%m%d_%H%M%S")
    
    # Add metadata
    question_data = {
        "question": question,
        "practice_type": practice_type,
        "topic": topic,
        "created_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    }
    
    # Add to stored questions
    stored_questions[question_id] = question_data
    
    # Save back to file
    os.makedirs(os.path.dirname(questions_file), exist_ok=True)
    with open(questions_file, 'w', encoding='utf-8') as f:
        json.dump(stored_questions, f, ensure_ascii=False, indent=2)
    
    return question_id

def render_interactive_stage():
    """Render the interactive learning stage"""
    # Initialize session state
    if 'question_generator' not in st.session_state:
        st.session_state.question_generator = QuestionGenerator()
    if 'current_question' not in st.session_state:
        st.session_state.current_question = None
    if 'feedback' not in st.session_state:
        st.session_state.feedback = None
    if 'current_practice_type' not in st.session_state:
        st.session_state.current_practice_type = None
    if 'current_topic' not in st.session_state:
        st.session_state.current_topic = None
    if 'current_stage' not in st.session_state:
        st.session_state.current_stage = "Chat"
        
    # Create sidebar with development stages
    with st.sidebar:
        st.subheader("Development Stages")
        # Reset current_stage if it's an old value
        if st.session_state.current_stage not in ["Chat", "Question Generation", "Interactive Learning"]:
            st.session_state.current_stage = "Chat"
        
        selected_stage = st.radio(
            "Select Stage:",
            ["Chat", "Question Generation", "Interactive Learning"],
            key="current_stage"
        )
        
        st.subheader("Current Focus:")
        st.markdown("- Basic Japanese learning")
        st.markdown("- Understanding LLM capabilities")
        st.markdown("- Identifying limitations")
        

    
    # Main content area
    st.title("Japanese Learning Assistant")
    st.caption("Transform YouTube transcripts into interactive Japanese learning experiences.")
    
    st.markdown("This tool demonstrates:")
    st.markdown("- Local LLM Integration (google/flan-t5-small)")
    st.markdown("- Vector Search with ChromaDB")
    st.markdown("- Sentence Transformers (all-MiniLM-L6-v2)")
    st.markdown("- Interactive Question Generation")
    
    # Handle different stages
    if selected_stage == "Chat":
        st.subheader("Chat with Japanese Learning Assistant")
        st.write("Ask questions about Japanese language and get helpful responses.")
        
        # Initialize chat if not already done
        if 'chat' not in st.session_state:
            st.session_state.chat = Chat()
        
        # Example questions
        st.subheader("Try These Examples")
        example_questions = [
            "How do I say 'Where is the train station?' in Japanese?",
            "Explain the difference between は and が",
            "What's the polite form of 食べる?",
            "How do I count objects in Japanese?"
        ]
        
        # Create two columns for a cleaner layout
        col1, col2 = st.columns(2)
        for i, q in enumerate(example_questions):
            # Alternate between columns
            with col1 if i % 2 == 0 else col2:
                if st.button(q, key=f"example_{q}"):
                    st.session_state.user_input = q
                    st.rerun()
        
        # Chat input
        user_input = st.text_input(
            "Ask a question:",
            key="user_input",
            value=st.session_state.get('user_input', ''),
            placeholder="Type your question here..."
        )
        
        # Initialize chat history in session state
        if 'chat_history' not in st.session_state:
            st.session_state.chat_history = []

        # Display chat history
        for msg in st.session_state.chat_history:
            role = msg['role']
            content = msg['content']
            with st.chat_message(role):
                st.write(content)

        if user_input:
            # Add user message to chat history
            st.session_state.chat_history.append({"role": "user", "content": user_input})
            
            # Display user message
            with st.chat_message("user"):
                st.write(user_input)

            # Generate and display assistant response
            with st.chat_message("assistant"):
                with st.spinner("Thinking..."):
                    response = st.session_state.chat.generate_response(user_input)
                    if response:
                        st.write(response)
                        # Add assistant response to chat history
                        st.session_state.chat_history.append({"role": "assistant", "content": response})
                    else:
                        st.error("Failed to generate response. Please try again.")
            
            # Clear input after sending
            st.session_state.user_input = ''
    
    elif selected_stage == "Question Generation" or selected_stage == "Interactive Learning":
        st.subheader("Practice Area")
        
        # Initialize question generator if not already done
        if 'question_generator' not in st.session_state:
            st.session_state.question_generator = QuestionGenerator()
        
        # Practice type selection in a cleaner format
        col1, col2 = st.columns(2)
        with col1:
            practice_type = st.selectbox(
                "Select Practice Type",
                ["Dialogue Practice", "Phrase Matching"],
                key="practice_type"
            )
        
        # Topic selection
        topics = {
            "Dialogue Practice": ["Daily Conversation", "Shopping", "Restaurant", "Travel", "School/Work"],
            "Phrase Matching": ["Announcements", "Instructions", "Weather Reports", "News Updates"]
        }
        
        with col2:
            topic = st.selectbox(
                "Select Topic",
                topics[practice_type],
                key="topic"
            )
        
        # Generate new question button
        if st.button("Generate New Question", type="primary"):
            section_num = 2 if practice_type == "Dialogue Practice" else 3
            try:
                with st.spinner("Generating question..."):
                    # Call backend API to generate question
                    try:
                        response = requests.post(
                            "http://localhost:8000/api/generate_question",
                            json={"section_num": section_num, "topic": topic}
                        )
                        
                        if response.status_code != 200:
                            st.error(f"Failed to generate question: {response.text}")
                            return
                            
                        response_data = response.json()
                    except requests.RequestException as e:
                        st.error(f"Error connecting to the backend service: {str(e)}")
                        return
                    if not response_data or 'question' not in response_data:
                        st.error("Invalid response from server")
                        return
                            
                    new_question = response_data["question"]
                    if not isinstance(new_question, dict):
                        st.error("Invalid question format received")
                        return
                            
                    # Validate required fields
                    required_fields = ['Question', 'Options', 'CorrectAnswer']
                    if not all(field in new_question for field in required_fields):
                        st.error("Question is missing required fields")
                        return
                        
                    st.session_state.current_question = new_question
                    st.session_state.current_practice_type = practice_type
                    st.session_state.current_topic = topic
                    st.session_state.feedback = None
                        
                    # Save the generated question
                    save_question(new_question, practice_type, topic)
            except Exception as e:
                st.error(f"Error generating question: {str(e)}. Please try again.")
                return
        

    
    # Display debug information expandable section
    with st.expander("Debug Information"):
        st.json({
            "Current Stage": st.session_state.current_stage,
            "Practice Type": st.session_state.get('current_practice_type', 'Not Selected'),
            "Topic": st.session_state.get('current_topic', 'Not Selected'),
            "Has Question": st.session_state.current_question is not None
        })
    
    if hasattr(st.session_state, 'current_question') and st.session_state.current_question:
        st.markdown("---")
        st.subheader("Practice Question")
        
        # Create tabs for different parts of the question
        scenario_tab, question_tab, answer_tab = st.tabs(["Scenario", "Question", "Answer"])
        
        # Get the current question safely
        question = st.session_state.current_question
        
        with scenario_tab:
            if practice_type == "Dialogue Practice":
                st.write("**Introduction:**")
                st.info(question.get('Introduction', 'No introduction available'))
                st.write("**Conversation:**")
                st.info(question.get('Conversation', 'No conversation available'))
            else:
                st.write("**Situation:**")
                st.info(question.get('Situation', 'No situation available'))
                if 'Content' in question:
                    st.write("**Content:**")
                    st.info(question['Content'])
        
        with question_tab:
            st.write("**Question:**")
            st.info(question.get('Question', 'No question available'))
            
            # Display options
            options = question.get('Options', [])
            if not options:
                st.warning("No options available for this question")
                return
            
            # If we have feedback, show which answers were correct/incorrect
            if hasattr(st.session_state, 'feedback') and st.session_state.feedback:
                try:
                    correct = st.session_state.feedback.get('correct', False)
                    correct_answer = question.get('CorrectAnswer', 1) - 1
                    selected_index = st.session_state.selected_answer - 1 if hasattr(st.session_state, 'selected_answer') else -1
                    
                    st.write("\n**Your Answer:**")
                    for i, option in enumerate(options):
                        if i == correct_answer and i == selected_index:
                            st.success(f"{i+1}. {option} ✓ (Correct!)")
                        elif i == correct_answer:
                            st.success(f"{i+1}. {option} ✓ (This was the correct answer)")
                        elif i == selected_index:
                            st.error(f"{i+1}. {option} ✗ (Your answer)")
                        else:
                            st.write(f"{i+1}. {option}")
                except Exception as e:
                    st.error(f"Error displaying feedback: {str(e)}")
                    st.session_state.feedback = None
            else:
                # Display options as radio buttons when no feedback yet
                try:
                    selected = st.radio(
                        "Choose your answer:",
                        options,
                        index=None,
                        format_func=lambda x: f"{options.index(x) + 1}. {x}"
                    )
                    
                    # Submit answer button
                    if selected and st.button("Submit Answer", type="primary"):
                        selected_index = options.index(selected) + 1
                        st.session_state.selected_answer = selected_index
                        
                        # Create feedback dictionary
                        correct_answer = question.get('CorrectAnswer', 1)
                        is_correct = selected_index == correct_answer
                        
                        st.session_state.feedback = {
                            'correct': is_correct,
                            'explanation': question.get('Explanation', 'No explanation available')
                        }
                        st.rerun()
                except Exception as e:
                    st.error(f"Error handling answer submission: {str(e)}")
        
        with answer_tab:
            if hasattr(st.session_state, 'feedback') and st.session_state.feedback:
                try:
                    # Show explanation
                    st.write("**Explanation:**")
                    explanation = st.session_state.feedback.get('explanation', 'No explanation available')
                    if st.session_state.feedback.get('correct', False):
                        st.success(explanation)
                    else:
                        st.error(explanation)
                    
                    # Add button to try new question
                    if st.button("Try Another Question", type="primary"):
                        st.session_state.feedback = None
                        st.rerun()
                except Exception as e:
                    st.error(f"Error displaying explanation: {str(e)}")
                    st.session_state.feedback = None
            else:
                st.info("Submit your answer in the Question tab to see the explanation.")
    else:
        st.info("Click 'Generate New Question' to start practicing!")

def main():
    st.title("JLPT Listening Practice")
    render_interactive_stage()

if __name__ == "__main__":
    main()
