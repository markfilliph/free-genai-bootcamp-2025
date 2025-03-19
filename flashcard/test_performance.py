import requests
import time
import statistics
import concurrent.futures
import json
import argparse

# Base URL for the API
BASE_URL = "http://localhost:8000/api"

# Global variables for authentication
access_token = None
user_id = None
deck_id = None

def register_and_login():
    """Register a new user and login to get an access token."""
    global access_token, user_id
    
    # Generate a unique username to avoid conflicts
    timestamp = int(time.time())
    username = f"testuser_{timestamp}"
    
    # Register user
    register_data = {
        "username": username,
        "email": f"{username}@example.com",
        "password": "password123"
    }
    
    response = requests.post(f"{BASE_URL}/auth/register", json=register_data)
    if response.status_code == 200:
        user_id = response.json().get("user_id")
        print(f"User registered with ID: {user_id}")
    else:
        print(f"Failed to register user: {response.text}")
        return False
    
    # Login
    login_data = {
        "username": username,
        "password": "password123"
    }
    
    response = requests.post(f"{BASE_URL}/auth/login", data=login_data)
    if response.status_code == 200:
        access_token = response.json().get("access_token")
        print(f"Login successful, received token")
        return True
    else:
        print(f"Failed to login: {response.text}")
        return False

def create_deck():
    """Create a new deck."""
    global deck_id
    
    headers = {"Authorization": f"Bearer {access_token}"}
    deck_data = {"name": f"Test Deck {int(time.time())}"}
    
    response = requests.post(f"{BASE_URL}/decks", json=deck_data, headers=headers)
    if response.status_code == 200:
        deck_id = response.json().get("id")
        print(f"Deck created with ID: {deck_id}")
        return True
    else:
        print(f"Failed to create deck: {response.text}")
        return False

def create_flashcard():
    """Create a new flashcard."""
    headers = {"Authorization": f"Bearer {access_token}"}
    
    # Generate a unique word to avoid conflicts
    timestamp = int(time.time())
    word = f"test_{timestamp}"
    
    flashcard_data = {
        "word": word,
        "example_sentence": f"This is an example sentence with {word}.",
        "translation": f"This is a translation for {word}.",
        "cultural_note": f"Cultural note for {word}.",
        "deck_id": deck_id
    }
    
    response = requests.post(f"{BASE_URL}/flashcards", json=flashcard_data, headers=headers)
    if response.status_code == 200:
        return True
    else:
        print(f"Failed to create flashcard: {response.text}")
        return False

def get_decks():
    """Get all decks."""
    headers = {"Authorization": f"Bearer {access_token}"}
    
    response = requests.get(f"{BASE_URL}/decks", headers=headers)
    if response.status_code == 200:
        return True
    else:
        print(f"Failed to get decks: {response.text}")
        return False

def get_flashcards():
    """Get all flashcards for a deck."""
    headers = {"Authorization": f"Bearer {access_token}"}
    
    response = requests.get(f"{BASE_URL}/decks/{deck_id}/flashcards", headers=headers)
    if response.status_code == 200:
        return True
    else:
        print(f"Failed to get flashcards: {response.text}")
        return False

def generate_content():
    """Generate content for a word."""
    headers = {"Authorization": f"Bearer {access_token}"}
    
    generation_data = {
        "word": "test",
        "is_verb": True
    }
    
    response = requests.post(f"{BASE_URL}/generate", json=generation_data, headers=headers)
    if response.status_code == 200:
        return True
    else:
        print(f"Failed to generate content: {response.text}")
        return False

def measure_endpoint_performance(endpoint_func, num_requests=10):
    """Measure the performance of an endpoint."""
    times = []
    
    for _ in range(num_requests):
        start_time = time.time()
        success = endpoint_func()
        end_time = time.time()
        
        if success:
            times.append(end_time - start_time)
    
    if not times:
        return {
            "min": None,
            "max": None,
            "avg": None,
            "median": None,
            "p95": None,
            "success_rate": 0
        }
    
    times.sort()
    p95_index = int(len(times) * 0.95)
    
    return {
        "min": min(times),
        "max": max(times),
        "avg": sum(times) / len(times),
        "median": statistics.median(times),
        "p95": times[p95_index] if p95_index < len(times) else max(times),
        "success_rate": len(times) / num_requests
    }

def concurrent_test(endpoint_func, num_concurrent=5, num_requests_per_thread=5):
    """Test an endpoint with concurrent requests."""
    start_time = time.time()
    success_count = 0
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=num_concurrent) as executor:
        futures = []
        for _ in range(num_concurrent):
            future = executor.submit(lambda: sum(1 for _ in range(num_requests_per_thread) if endpoint_func()))
            futures.append(future)
        
        for future in concurrent.futures.as_completed(futures):
            success_count += future.result()
    
    end_time = time.time()
    total_requests = num_concurrent * num_requests_per_thread
    
    return {
        "total_time": end_time - start_time,
        "total_requests": total_requests,
        "successful_requests": success_count,
        "success_rate": success_count / total_requests,
        "requests_per_second": total_requests / (end_time - start_time)
    }

def run_performance_tests(num_requests=10, num_concurrent=5, num_requests_per_thread=5):
    """Run performance tests on all endpoints."""
    print("\n=== Performance Testing ===\n")
    
    # Setup: Register and login
    if not register_and_login():
        print("Failed to setup test user, aborting tests.")
        return
    
    # Setup: Create a deck
    if not create_deck():
        print("Failed to create test deck, aborting tests.")
        return
    
    # Test each endpoint
    endpoints = {
        "Create Flashcard": create_flashcard,
        "Get Decks": get_decks,
        "Get Flashcards": get_flashcards,
        "Generate Content": generate_content
    }
    
    results = {}
    
    print("\n--- Individual Endpoint Performance ---\n")
    for name, func in endpoints.items():
        print(f"Testing {name}...")
        perf = measure_endpoint_performance(func, num_requests)
        
        print(f"  Min: {perf['min']:.4f}s")
        print(f"  Max: {perf['max']:.4f}s")
        print(f"  Avg: {perf['avg']:.4f}s")
        print(f"  Median: {perf['median']:.4f}s")
        print(f"  P95: {perf['p95']:.4f}s")
        print(f"  Success Rate: {perf['success_rate'] * 100:.2f}%")
        print()
        
        results[name] = perf
    
    print("\n--- Concurrent Performance ---\n")
    for name, func in endpoints.items():
        print(f"Testing {name} with {num_concurrent} concurrent clients...")
        perf = concurrent_test(func, num_concurrent, num_requests_per_thread)
        
        print(f"  Total Time: {perf['total_time']:.4f}s")
        print(f"  Total Requests: {perf['total_requests']}")
        print(f"  Successful Requests: {perf['successful_requests']}")
        print(f"  Success Rate: {perf['success_rate'] * 100:.2f}%")
        print(f"  Requests Per Second: {perf['requests_per_second']:.2f}")
        print()
        
        results[f"{name} (Concurrent)"] = perf
    
    # Save results to file
    with open("performance_results.json", "w") as f:
        json.dump(results, f, indent=2)
    
    print("Performance test results saved to performance_results.json")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Performance testing for the Flashcard API")
    parser.add_argument("--requests", type=int, default=10, help="Number of requests per endpoint")
    parser.add_argument("--concurrent", type=int, default=5, help="Number of concurrent clients")
    parser.add_argument("--requests-per-thread", type=int, default=5, help="Number of requests per concurrent client")
    
    args = parser.parse_args()
    
    run_performance_tests(args.requests, args.concurrent, args.requests_per_thread)
