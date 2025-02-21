import streamlit as st
import sys
import os
import json
from datetime import datetime
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.question_generator import QuestionGenerator
from backend.chat import Chat

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
    if selected_stage == "Chat with Claude":
        st.subheader("Chat with Claude")
        st.write("Ask questions about Japanese language and get instant answers.")
        
        # Initialize chat history
        if 'chat_history' not in st.session_state:
            st.session_state.chat_history = []
        
        # Display chat history
        for msg in st.session_state.chat_history:
            if msg['type'] == 'user':
                st.write(f"🧑 **You:** {msg['content']}")
            else:
                st.write(f"🤖 **Assistant:** {msg['content']}")
        
        # Text input and send button in the same row
        col1, col2 = st.columns([4, 1])
        with col1:
            # Handle example question selection
            if 'example_question' in st.session_state:
                user_input = st.text_input("Ask a question:", value=st.session_state.example_question)
                # Clear the example question after using it
                del st.session_state.example_question
            else:
                user_input = st.text_input("Ask a question:")
        
        with col2:
            send_pressed = st.button("Send", type="primary")
        
        # Handle user input
        if user_input and send_pressed:
            # Add user message to chat
            st.session_state.chat_history.append({"type": "user", "content": user_input})
            
            # Generate response based on question type
            if "train station" in user_input.lower():
                response = "To ask 'Where is the train station?' in Japanese, you can say: 駅はどこですか？ (eki wa doko desu ka?)\n\nBreaking it down:\n- 駅 (eki) = train station\n- は (wa) = topic marker\n- どこ (doko) = where\n- ですか (desu ka) = polite question marker"
            elif "は and が" in user_input:
                response = "は (wa) and が (ga) are both particles but serve different purposes:\n\n1. は marks the topic of the sentence\n2. が marks the subject and shows emphasis\n\nExample:\n- 私は学生です (watashi wa gakusei desu) = As for me, I'm a student\n- 私が学生です (watashi ga gakusei desu) = I (specifically) am the student"
            elif "食べる" in user_input:
                response = "The polite form of 食べる (taberu) is 食べます (tabemasu).\n\nOther forms:\n- Plain present: 食べる (taberu)\n- Polite present: 食べます (tabemasu)\n- Super polite: 召し上がります (meshiagarimasu)"
            elif "count" in user_input.lower():
                response = "In Japanese, different types of objects use different counter words. Here are some common ones:\n\n1. Small objects (鉛筆1本): ～本 (-hon)\n2. Flat objects (紙1枚): ～枚 (-mai)\n3. Small animals (猫2匹): ～匹 (-hiki)\n4. People (学生3人): ～人 (-nin)"
            else:
                response = "I can help you with Japanese language questions! Try asking about grammar, vocabulary, or cultural aspects of Japanese."
            
            # Add assistant response to chat
            st.session_state.chat_history.append({"type": "assistant", "content": response})
            
            # Rerun to update chat display
            st.rerun()
            
        # Show example questions at the bottom
        st.markdown("---")
        st.subheader("Try These Examples")
        example_questions = [
            "How do I say 'Where is the train station?' in Japanese?",
            "Explain the difference between は and が",
            "What's the polite form of 食べる?",
            "How do I count objects in Japanese?"
        ]
        
        # Display examples in two columns
        col1, col2 = st.columns(2)
        for i, q in enumerate(example_questions):
            with col1 if i < len(example_questions)//2 else col2:
                if st.button(q, key=f"example_{q}"):
                    st.session_state.example_question = q
                    st.rerun()
            
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
                    import requests
                    response = requests.post(
                        "http://localhost:8000/api/generate_question",
                        json={"section_num": section_num, "topic": topic}
                    )
                    
                    if response.status_code != 200:
                        st.error(f"Failed to generate question: {response.text}")
                        return
                        
                    new_question = response.json()["question"]
                    if not new_question:
                        st.error("Failed to generate a new question. Please try again.")
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
        
        # Display current question if available
        if 'current_question' in st.session_state:
            question = st.session_state.current_question
            
            try:
                # Display in tabs
                question_tab, scenario_tab, answer_tab = st.tabs(["Question", "Scenario", "Answer"])
                
                with question_tab:
                    st.write(question.get('Question', 'No question available'))
                    options = question.get('Options', [])
                    for i, option in enumerate(options, 1):
                        st.write(f"{i}. {option}")
                
                with scenario_tab:
                    if section_num == 2:
                        st.write(f"**Introduction:** {question.get('Introduction', '')}")
                        st.write(f"**Conversation:**\n{question.get('Conversation', '')}")
                    else:
                        st.write(f"**Situation:** {question.get('Situation', '')}")
                        if 'Content' in question:
                            st.write(f"**Content:**\n{question['Content']}")
                
                with answer_tab:
                    options = question.get('Options', [])
                    correct_answer_idx = question.get('CorrectAnswer', 1) - 1
                    if 0 <= correct_answer_idx < len(options):
                        st.write(f"**Correct Answer:** {options[correct_answer_idx]}")
                    if 'Explanation' in question:
                        st.write(f"**Explanation:** {question['Explanation']}")
            except Exception as e:
                st.error(f"Error displaying question: {str(e)}")
    
    # Display debug information expandable section
    with st.expander("Debug Information"):
        st.json({
            "Current Stage": st.session_state.current_stage,
            "Practice Type": practice_type,
            "Topic": topic,
            "Has Question": st.session_state.current_question is not None
        })
    
    if st.session_state.current_question:
        st.markdown("---")
        st.subheader("Practice Question")
        
        # Create tabs for different parts of the question
        scenario_tab, question_tab, answer_tab = st.tabs(["Scenario", "Question", "Answer"])
        
        with scenario_tab:
            if practice_type == "Dialogue Practice":
                st.write("**Introduction:**")
                st.info(st.session_state.current_question['Introduction'])
                st.write("**Conversation:**")
                st.info(st.session_state.current_question['Conversation'])
            else:
                st.write("**Situation:**")
                st.info(st.session_state.current_question['Situation'])
        
        with question_tab:
            st.write("**Question:**")
            st.info(st.session_state.current_question['Question'])
            
            # Display options
            options = st.session_state.current_question['Options']
            
            # If we have feedback, show which answers were correct/incorrect
            if st.session_state.feedback:
                correct = st.session_state.feedback.get('correct', False)
                correct_answer = st.session_state.feedback.get('correct_answer', 1) - 1
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
            else:
                # Display options as radio buttons when no feedback yet
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
                    st.session_state.feedback = st.session_state.question_generator.get_feedback(
                        st.session_state.current_question,
                        selected_index
                    )
                    st.rerun()
        
        with answer_tab:
            if st.session_state.feedback:
                # Show explanation
                st.write("**Explanation:**")
                explanation = st.session_state.feedback.get('explanation', 'No feedback available')
                if st.session_state.feedback.get('correct', False):
                    st.success(explanation)
                else:
                    st.error(explanation)
                
                # Add button to try new question
                if st.button("Try Another Question", type="primary"):
                    st.session_state.feedback = None
                    st.rerun()
            else:
                st.info("Submit your answer in the Question tab to see the explanation.")
    else:
        st.info("Click 'Generate New Question' to start practicing!")

def main():
    st.title("JLPT Listening Practice")
    render_interactive_stage()

if __name__ == "__main__":
    main()
