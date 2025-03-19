import requests
import json
import time

# Base URL for the API
BASE_URL = "http://localhost:8000/api"

def print_response(response):
    print(f"Status Code: {response.status_code}")
    try:
        print(f"Response: {json.dumps(response.json(), indent=2)}")
    except:
        print(f"Response: {response.text}")
    print("-" * 50)

def test_api():
    # Test 1: Register a new user
    print("\nTest 1: Register a new user")
    register_data = {
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123"
    }
    response = requests.post(f"{BASE_URL}/auth/register", json=register_data)
    print_response(response)
    
    # Test 2: Login with the new user
    print("\nTest 2: Login with the new user")
    login_data = {
        "username": "testuser",
        "password": "password123"
    }
    response = requests.post(f"{BASE_URL}/auth/login", data=login_data)
    print_response(response)
    
    # Extract the access token
    if response.status_code == 200:
        access_token = response.json().get("access_token")
        headers = {"Authorization": f"Bearer {access_token}"}
    else:
        print("Login failed, cannot continue tests")
        return
    
    # Test 3: Create a new deck
    print("\nTest 3: Create a new deck")
    deck_data = {"name": "Spanish Basics"}
    response = requests.post(f"{BASE_URL}/decks", json=deck_data, headers=headers)
    print_response(response)
    
    # Extract the deck ID
    if response.status_code == 200:
        deck_id = response.json().get("id")
    else:
        print("Deck creation failed, cannot continue tests")
        return
    
    # Test 4: Get all decks
    print("\nTest 4: Get all decks")
    response = requests.get(f"{BASE_URL}/decks", headers=headers)
    print_response(response)
    
    # Test 5: Create a flashcard
    print("\nTest 5: Create a flashcard")
    flashcard_data = {
        "word": "hola",
        "example_sentence": "Hola, ¿cómo estás?",
        "translation": "Hello, how are you?",
        "cultural_note": "Common greeting in Spanish-speaking countries",
        "deck_id": deck_id
    }
    response = requests.post(f"{BASE_URL}/flashcards", json=flashcard_data, headers=headers)
    print_response(response)
    
    # Extract the flashcard ID
    if response.status_code == 200:
        flashcard_id = response.json().get("id")
    else:
        print("Flashcard creation failed, cannot continue tests")
        return
    
    # Test 6: Get flashcards by deck
    print("\nTest 6: Get flashcards by deck")
    response = requests.get(f"{BASE_URL}/decks/{deck_id}/flashcards", headers=headers)
    print_response(response)
    
    # Test 7: Test LLM generation
    print("\nTest 7: Test LLM generation")
    generation_data = {
        "word": "hablar",
        "is_verb": True
    }
    response = requests.post(f"{BASE_URL}/generate", json=generation_data, headers=headers)
    print_response(response)
    
    print("\nAll tests completed!")

if __name__ == "__main__":
    # Wait for the server to start
    print("Waiting for server to start...")
    time.sleep(2)
    test_api()
