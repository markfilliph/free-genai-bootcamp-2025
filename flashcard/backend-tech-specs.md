# Backend Technical Specifications

## Overview
The backend of the Language Learning Flashcard Generator is built using **FastAPI** and **SQLite3**. It provides APIs for flashcard generation, user authentication, and data storage. The backend integrates with **Ollama** (a local LLM) for generating example sentences, verb conjugations, and cultural context.

---

## Technology Stack
- **Framework**: FastAPI (Python)
- **Database**: SQLite3
- **LLM Integration**: Ollama
- **Spaced Repetition**: SuperMemo2 (Python implementation)
- **Exporting Flashcards**:
  - PDF: ReportLab
  - CSV: Pandas
- **Authentication**: JWT (JSON Web Tokens)

---

## API Endpoints

### User Authentication
1. **POST /register**
   - Register a new user.
   - Request Body:
     ```json
     {
       "username": "string",
       "email": "string",
       "password": "string"
     }
     ```
   - Response:
     ```json
     {
       "user_id": "int",
       "message": "User registered successfully"
     }
     ```

2. **POST /login**
   - Authenticate a user and return a JWT token.
   - Request Body:
     ```json
     {
       "email": "string",
       "password": "string"
     }
     ```
   - Response:
     ```json
     {
       "access_token": "string",
       "token_type": "bearer"
     }
     ```

### Flashcard Management
1. **POST /flashcards**
   - Create a new flashcard.
   - Request Body:
     ```json
     {
       "deck_id": "int",
       "word": "string",
       "example_sentence": "string",
       "translation": "string",
       "conjugation": "string",
       "cultural_note": "string"
     }
     ```
   - Response:
     ```json
     {
       "flashcard_id": "int",
       "message": "Flashcard created successfully"
     }
     ```

2. **GET /flashcards/{deck_id}**
   - Retrieve all flashcards in a deck.
   - Response:
     ```json
     [
       {
         "flashcard_id": "int",
         "word": "string",
         "example_sentence": "string",
         "translation": "string",
         "conjugation": "string",
         "cultural_note": "string"
       }
     ]
     ```

3. **DELETE /flashcards/{flashcard_id}**
   - Delete a flashcard.
   - Response:
     ```json
     {
       "message": "Flashcard deleted successfully"
     }
     ```

### Deck Management
1. **POST /decks**
   - Create a new deck.
   - Request Body:
     ```json
     {
       "user_id": "int",
       "deck_name": "string"
     }
     ```
   - Response:
     ```json
     {
       "deck_id": "int",
       "message": "Deck created successfully"
     }
     ```

2. **GET /decks/{user_id}**
   - Retrieve all decks for a user.
   - Response:
     ```json
     [
       {
         "deck_id": "int",
         "deck_name": "string",
         "created_at": "string"
       }
     ]
     ```

### LLM Integration (Ollama)
1. **POST /generate**
   - Generate example sentences, conjugations, and cultural notes using Ollama.
   - Request Body:
     ```json
     {
       "word": "string",
       "is_verb": "boolean"
     }
     ```
   - Response:
     ```json
     {
       "example_sentences": ["string"],
       "conjugations": "string",
       "cultural_note": "string"
     }
     ```

---

## Database Schema
### Tables
1. **Users**:
   - `user_id` (Primary Key)
   - `username`
   - `email`
   - `password_hash`

2. **Decks**:
   - `deck_id` (Primary Key)
   - `user_id` (Foreign Key)
   - `deck_name`
   - `created_at`

3. **Flashcards**:
   - `flashcard_id` (Primary Key)
   - `deck_id` (Foreign Key)
   - `word`
   - `example_sentence`
   - `translation`
   - `conjugation` (for verbs)
   - `cultural_note`
   - `created_at`

4. **Tags**:
   - `tag_id` (Primary Key)
   - `tag_name`

5. **Flashcard_Tags**:
   - `flashcard_id` (Foreign Key)
   - `tag_id` (Foreign Key)

---

## Authentication
- **JWT (JSON Web Tokens)**:
  - Tokens are issued upon successful login.
  - Tokens are validated for protected routes (e.g., creating flashcards, decks).

---

## Exporting Flashcards
1. **PDF**:
   - Use ReportLab to generate PDFs of flashcards.
   - Endpoint: `POST /export/pdf`
   - Request Body:
     ```json
     {
       "deck_id": "int"
     }
     ```
   - Response: PDF file.

2. **CSV**:
   - Use Pandas to export flashcards as CSV.
   - Endpoint: `POST /export/csv`
   - Request Body:
     ```json
     {
       "deck_id": "int"
     }
     ```
   - Response: CSV file.

---

## Spaced Repetition
- **SuperMemo2**:
  - Flashcards are reviewed based on user performance (easy, medium, hard).
  - Review intervals are calculated using the SuperMemo2 algorithm.

---

## Environment Variables
- `DATABASE_URL`: SQLite database file path.
- `JWT_SECRET_KEY`: Secret key for JWT token generation.
- `OLLAMA_API_URL`: URL for Ollama API.

---

## Setup Instructions
1. Install dependencies:
   ```bash
   pip install -r requirements.txt