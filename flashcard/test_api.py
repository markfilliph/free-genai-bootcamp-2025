import requests
import json
import time
import os

# Base URL for the API
BASE_URL = "http://localhost:8000/api"

# Test data
test_user = {
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}

test_deck = {
    "name": "Spanish Basics"
}

test_flashcard = {
    "word": "hola",
    "example_sentence": "Hola, ¿cómo estás?",
    "translation": "Hello, how are you?",
    "conjugation": None,
    "cultural_note": "Common greeting in Spanish-speaking countries"
}

# Helper function to print responses
def print_response(response, message):
    print(f"\n{message}")
    print(f"Status Code: {response.status_code}")
    try:
        print(f"Response: {json.dumps(response.json(), indent=2)}")
    except:
        print(f"Response: {response.text}")

def test_api():
    # Step 1: Register a new user
    print("\n--- Testing User Registration ---")
    register_url = f"{BASE_URL}/auth/register"
    response = requests.post(register_url, json=test_user)
    print_response(response, "Register User Response:")
    
    # Step 2: Login with the new user
    print("\n--- Testing User Login ---")
    login_url = f"{BASE_URL}/auth/login"
    login_data = {
        "username": test_user["username"],
        "password": test_user["password"]
    }
    response = requests.post(login_url, data=login_data)
    print_response(response, "Login Response:")
    
    if response.status_code != 200:
        print("Login failed. Cannot continue testing.")
        return
    
    # Get the access token
    token_data = response.json()
    access_token = token_data["access_token"]
    headers = {"Authorization": f"Bearer {access_token}"}
    
    # Step 3: Create a new deck
    print("\n--- Testing Deck Creation ---")
    deck_url = f"{BASE_URL}/decks"
    response = requests.post(deck_url, json=test_deck, headers=headers)
    print_response(response, "Create Deck Response:")
    
    if response.status_code != 200:
        print("Deck creation failed. Cannot continue testing.")
        return
    
    # Get the deck ID
    deck_id = response.json()["id"]
    
    # Step 4: Get all decks for the user
    print("\n--- Testing Get All Decks ---")
    response = requests.get(deck_url, headers=headers)
    print_response(response, "Get All Decks Response:")
    
    # Step 5: Create a flashcard in the deck
    print("\n--- Testing Flashcard Creation ---")
    flashcard_url = f"{BASE_URL}/flashcards"
    flashcard_data = {**test_flashcard, "deck_id": deck_id}
    response = requests.post(flashcard_url, json=flashcard_data, headers=headers)
    print_response(response, "Create Flashcard Response:")
    
    if response.status_code != 200:
        print("Flashcard creation failed. Cannot continue testing.")
        return
    
    # Get the flashcard ID
    flashcard_id = response.json()["id"]
    
    # Step 6: Get all flashcards in the deck
    print("\n--- Testing Get Flashcards by Deck ---")
    deck_flashcards_url = f"{BASE_URL}/decks/{deck_id}/flashcards"
    response = requests.get(deck_flashcards_url, headers=headers)
    print_response(response, "Get Flashcards by Deck Response:")
    
    # Step 7: Test LLM generation
    print("\n--- Testing LLM Generation ---")
    generation_url = f"{BASE_URL}/generate"
    generation_data = {"word": "hablar", "is_verb": True}
    response = requests.post(generation_url, json=generation_data, headers=headers)
    print_response(response, "LLM Generation Response:")
    
    print("\n--- All tests completed ---")

if __name__ == "__main__":
    # Check if the server is running
    try:
        response = requests.get("http://localhost:8000/")
        print("Server is running. Starting tests...")
        test_api()
    except requests.exceptions.ConnectionError:
        print("Error: Cannot connect to the server. Make sure it's running on http://localhost:8000/")
