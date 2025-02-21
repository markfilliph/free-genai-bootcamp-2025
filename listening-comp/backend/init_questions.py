import os
import sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.vector_store import QuestionVectorStore
from backend.config import *  # Import OpenAI API key

# Sample questions for initialization
sample_questions = {
    "section2": [
        {
            "Introduction": "You are at a restaurant with a friend.",
            "Conversation": "Man: What would you like to order?\nWoman: I'm thinking about the pasta, but I'm not sure.\nMan: The pasta here is really good. I had it last time.\nWoman: Really? What kind did you try?\nMan: I had the seafood pasta. It was amazing!",
            "Question": "What did the man recommend to the woman?",
            "Options": [
                "The seafood pasta",
                "The vegetarian pasta",
                "The meat pasta",
                "The chicken pasta"
            ],
            "CorrectAnswer": 1,
            "Topic": "Restaurant"
        },
        {
            "Introduction": "You are at a train station.",
            "Conversation": "Woman: Excuse me, which platform is the train to Tokyo?\nMan: Let me check... It's platform 3.\nWoman: Thank you. Do you know when it leaves?\nMan: It leaves in 10 minutes at 2:30.",
            "Question": "When does the train leave?",
            "Options": [
                "At 2:00",
                "At 2:15",
                "At 2:30",
                "At 2:45"
            ],
            "CorrectAnswer": 3,
            "Topic": "Travel"
        }
    ],
    "section3": [
        {
            "Situation": "This is an announcement at a department store.",
            "Content": "Attention shoppers. We are having a special sale in our clothing department on the second floor. All winter items are 50% off. This sale will end at 7:00 PM today. Don't miss this great opportunity!",
            "Question": "Until what time is the sale?",
            "Options": [
                "5:00 PM",
                "6:00 PM",
                "7:00 PM",
                "8:00 PM"
            ],
            "CorrectAnswer": 3,
            "Topic": "Announcements"
        },
        {
            "Situation": "You are listening to the weather forecast.",
            "Content": "Good morning. Today will be mostly sunny with a high of 25 degrees. However, there is a 60% chance of rain in the evening. Tomorrow will be cloudy with temperatures around 20 degrees.",
            "Question": "What is the weather forecast for tomorrow?",
            "Options": [
                "Sunny and warm",
                "Cloudy and mild",
                "Rainy and cold",
                "Windy and dry"
            ],
            "CorrectAnswer": 2,
            "Topic": "Weather Reports"
        }
    ]
}

def init_vector_store():
    """Initialize the vector store with sample questions"""
    store = QuestionVectorStore()
    
    # Add section 2 questions
    for question in sample_questions["section2"]:
        store.add_question(2, question)
    
    # Add section 3 questions
    for question in sample_questions["section3"]:
        store.add_question(3, question)

if __name__ == "__main__":
    init_vector_store()
