# Language Learning Flashcard Generator (Spanish)

## Project Overview
The Language Learning Flashcard Generator is an app designed to help learners of Spanish create personalized flashcards with example sentences, verb conjugations, and cultural context. The app leverages Ollama, a local LLM, to generate content and provides tools for effective learning and review.

## Features
- **Flashcard Creation**: Generate example sentences, verb conjugations, and translations.
- **Verb Conjugation Support**: Conjugations for all tenses and moods.
- **Cultural Context**: Cultural notes and idiomatic expressions.
- **Text-to-Speech (TTS)**: Listen to pronunciations using ResponsiveVoice.js.
- **Flashcard Organization**: Organize flashcards into decks and tag them.
- **Review Mode**: Spaced repetition system (SRS) using SuperMemo2.
- **Export Flashcards**: Export as PDF or CSV.
- **User Accounts**: Save and sync flashcards across devices.

## Technology Stack
- **Frontend**: Svelte, Svelte Material UI, ResponsiveVoice.js.
- **Backend**: FastAPI, SQLite3.
- **LLM Integration**: Ollama.
- **Spaced Repetition**: SuperMemo2.
- **Exporting Flashcards**: ReportLab (PDF), Pandas (CSV).

## Testing Approach
We've implemented several testing strategies to ensure the reliability of our application:

### 1. Database Model Testing
- `test_models.py`: Directly tests the database models and CRUD operations without relying on the API layer.
- Verifies user creation, deck management, and flashcard operations.

### 2. API Testing
- `test_minimal_api.py`: Tests the API endpoints using the requests library.
- Covers user registration, authentication, deck and flashcard management, and content generation.
- Provides detailed output of each API call for debugging.

### 3. Manual Testing Interface
- `test_api.html`: A simple HTML interface for manually testing the API endpoints.
- Allows interactive testing of all API features with a user-friendly UI.

### 4. Curl-based Testing
- `test_curl.sh`: A bash script that tests the API endpoints using curl commands.
- Useful for CI/CD pipelines and automated testing.

### 5. Simplified Backend
- `backend_minimal.py`: A lightweight implementation of the backend using standard Python libraries.
- Helps isolate and troubleshoot issues with the core functionality without dependency complications.

To run the tests, use the following commands:

```bash
# Test database models
python test_models.py

# Test API endpoints
python test_minimal_api.py

# Start the minimal backend server
python backend_minimal.py

# Run curl-based tests
./test_curl.sh
```