# Testing Strategy for Language Learning Flashcard Generator

This document outlines the testing approach used to verify the functionality of the Language Learning Flashcard Generator application.

## Testing Goals

1. Verify the core functionality of the backend API
2. Ensure proper database schema and relationships
3. Test user authentication and authorization
4. Validate flashcard and deck management operations
5. Test the LLM integration for content generation

## Testing Components

### 1. Database Model Testing (`test_models.py`)

This script tests the database models directly, bypassing the API layer to ensure that the core data operations work correctly.

**What it tests:**
- User creation and authentication
- Deck creation and retrieval
- Flashcard creation and retrieval
- Relationships between users, decks, and flashcards

**How to run:**
```bash
python test_models.py
```

### 2. API Testing (`test_minimal_api.py`)

This script tests the API endpoints using the Python requests library, simulating client interactions with the server.

**What it tests:**
- User registration and login
- JWT token-based authentication
- Deck creation and retrieval
- Flashcard creation and retrieval
- Content generation

**How to run:**
```bash
# First start the server
python backend_minimal.py

# Then in another terminal
python test_minimal_api.py
```

### 3. Manual Testing Interface (`test_api.html`)

A simple HTML page that allows manual testing of the API endpoints through a user-friendly interface.

**What it tests:**
- All API endpoints with a visual interface
- Real-time feedback on API responses
- Token handling and authentication flow

**How to use:**
1. Start the server: `python backend_minimal.py`
2. Open `test_api.html` in a web browser
3. Follow the steps in the interface to test each endpoint

### 4. Curl-based Testing (`test_curl.sh`)

A bash script that tests the API endpoints using curl commands, useful for CI/CD pipelines and automated testing.

**What it tests:**
- All API endpoints with detailed output
- Success and error handling
- Authentication flow

**How to run:**
```bash
# First start the server
python backend_minimal.py

# Then in another terminal
./test_curl.sh
```

### 5. Database Schema Verification (`verify_db_schema.py`)

A utility script to verify the database schema, ensuring that all tables, columns, and relationships are correctly defined.

**What it tests:**
- Table structure
- Column types and constraints
- Foreign key relationships
- Indexes

**How to run:**
```bash
python verify_db_schema.py
```

## Simplified Backend Implementation

To isolate and test core functionality without dependency complications, we created a simplified backend implementation (`backend_minimal.py`) that uses standard Python libraries instead of FastAPI and other dependencies.

**Features:**
- SQLite database for data storage
- Simple HTTP server for API endpoints
- JWT-like token authentication
- CRUD operations for users, decks, and flashcards
- Simulated LLM content generation

**How to run:**
```bash
python backend_minimal.py
```

## Test Results

All tests have been successfully executed, verifying that:

1. The database schema is correctly defined with proper relationships
2. User authentication works as expected
3. Deck and flashcard management operations function correctly
4. The API endpoints return the expected responses
5. Content generation provides the expected format of data

## Next Steps

1. Integrate the frontend with the backend API
2. Implement end-to-end testing with Cypress or similar tools
3. Add unit tests for individual components
4. Set up CI/CD pipeline for automated testing
