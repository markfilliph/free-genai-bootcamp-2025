import time
from chat import Chat

def test_response_time(chat, prompt):
    start_time = time.time()
    response = chat.generate_response(prompt)
    end_time = time.time()
    return response, end_time - start_time

def main():
    chat = Chat()
    test_prompts = [
        "How do I say 'hello' in Japanese?",
        "What is 'thank you' in Japanese?",
        "How do I say 'good morning' in Japanese?",
        "What is the capital of Japan?",  # Non-translation question
    ]
    
    print("Testing chat response times...")
    print("-" * 50)
    
    for i, prompt in enumerate(test_prompts, 1):
        response, duration = test_response_time(chat, prompt)
        print(f"\nQuery {i}: '{prompt}'")
        print(f"Response: {response}")
        print(f"Time taken: {duration:.2f} seconds")
        
        if i == 1:
            print("\nNow that model is loaded, testing subsequent queries...")
            print("-" * 50)

if __name__ == "__main__":
    main()
