# Language Learning Flashcard Generator (Spanish)

## Project Overview
The Language Learning Flashcard Generator is an app designed to help learners of Spanish create personalized flashcards with example sentences, verb conjugations, and cultural context. The app leverages Ollama, a local LLM, to generate content and provides tools for effective learning and review.

## Features
- **Flashcard Creation**: Generate example sentences, verb conjugations, and translations.
- **Verb Conjugation Support**: Conjugations for all tenses and moods.
- **Cultural Context**: Cultural notes and idiomatic expressions.
- **Text-to-Speech (TTS)**: Listen to pronunciations using ResponsiveVoice.js. (not deployed)
- **Flashcard Organization**: Organize flashcards into decks and tag them.
- **Review Mode**: Spaced repetition system (SRS) using SuperMemo2.
- **Export Flashcards**: Export as PDF or CSV.
- **User Accounts**: Save and sync flashcards across devices.

## Getting Started

### Prerequisites
- Python 3.8 or higher
- Node.js 14 or higher
- npm or yarn
- Ollama (for LLM integration)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd flashcard
   ```

2. Set up the backend:
   ```bash
   # Install backend dependencies
   pip install -r requirements.txt
   
   # Initialize the database (if needed)
   python -m backend.setup_db
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

### Starting the Servers

1. Start the backend server using the unified server launcher:
   ```bash
   # From the project root directory
   python run_server.py --backend main --port 8000
   ```
   The backend API will be available at http://localhost:8000
   
   You can also choose alternative backend implementations:
   ```bash
   # Simple FastAPI backend
   python run_server.py --backend simple
   
   # Minimal backend (no external dependencies)
   python run_server.py --backend minimal
   ```
   
   For more information about the different backend implementations, see [BACKEND_IMPLEMENTATIONS.md](./BACKEND_IMPLEMENTATIONS.md)

2. Start the frontend development server:
   ```bash
   # From the project root directory
   cd frontend
   npm run dev
   ```
   The frontend application will be available at http://localhost:8080

3. Access the application:
   Open your browser and navigate to http://localhost:8080

### API Documentation
Once the backend server is running, you can access the API documentation at:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Technology Stack
- **Frontend**: Svelte, Svelte Material UI, ResponsiveVoice.js.
- **Backend**: FastAPI, SQLite3.
- **LLM Integration**: Ollama.
- **Spaced Repetition**: SuperMemo2.
- **Exporting Flashcards**: ReportLab (PDF), Pandas (CSV).
- **State Management**: Custom persistent store implementation with localStorage.

## Recent Improvements

### Data Persistence Between Components
We've implemented a robust state management solution to address the previous issue where data created in one component (like new decks) wasn't accessible in other components. The improvements include:

1. **Singleton Store Pattern**: Ensures all components share the same data source
2. **Persistent Storage**: Reliable localStorage integration with proper error handling
3. **Reactive Updates**: Components now automatically reflect changes made elsewhere in the application
4. **Refresh Mechanism**: Components can force-refresh from localStorage when mounted

This fixes the core issue where newly created decks in DeckManagement.svelte weren't appearing in CreateFlashcards.svelte.

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
- `backend/alternatives/minimal_backend.py`: A lightweight implementation of the backend using standard Python libraries.
- Helps isolate and troubleshoot issues with the core functionality without dependency complications.

### 6. Database Schema Verification
- `verify_db_schema.py`: Verifies that the database schema matches the expected structure.
- Checks tables, columns, foreign keys, and indexes.

### 7. Performance Testing
- `test_performance.py`: Measures API response times for various endpoints.
- Tests both individual and concurrent requests to identify bottlenecks.

### 8. Frontend Testing
- Jest and Testing Library for Svelte components testing.
- Custom testing approach for Svelte components.
- Unit tests for individual components (FlashcardForm, Navbar, DeckList, etc.).
- Integration tests for API interactions and state management.
- Mock implementations for external dependencies.

For more details on our frontend testing approach, see:
- [Frontend Testing Guide](./frontend/TESTING_GUIDE.md) - General testing setup and organization
- [Svelte Testing Approach](./frontend/SVELTE_TESTING_APPROACH.md) - Detailed explanation of our custom approach to testing Svelte components

To run the tests, use the following commands:

```bash
# Test database models
python test_models.py

# Test API endpoints
python test_minimal_api.py

# Verify database schema
python verify_db_schema.py

# Run performance tests
python test_performance.py

# Run frontend tests
cd frontend
npm test

# Run frontend tests with coverage
cd frontend
npm run test:coverage

# Start the minimal backend server
python backend_minimal.py

# Run curl-based tests
./test_curl.sh
```

## Screenshots
See the screenshots in the Screenshot folder.

## Troubleshooting

### Backend Issues
- Ensure all dependencies are installed: `pip install -r requirements.txt`
- Check if the database exists and is properly initialized
- Verify that port 8000 is not in use by another application
- Check the backend logs for specific error messages
- If using the minimal backend, ensure it's running on the expected port

### Frontend Issues
- Ensure all dependencies are installed: `npm install` in the frontend directory
- Verify that port 8080 is not in use by another application
- Check the browser console for JavaScript errors
- Make sure the backend server is running and accessible
- If data isn't persisting between components, try refreshing the page to force localStorage synchronization
- Clear browser cache and localStorage if you encounter persistent data issues

### CORS Issues
- The backend is configured to allow requests from the following origins:
  - http://localhost:5173
  - http://localhost:8080
  - http://localhost:8083
- If you're running the frontend on a different port, update the CORS configuration in `backend/main.py`