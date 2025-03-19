#!/bin/bash

# Base URL for the API
BASE_URL="http://localhost:8000/api"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Testing Language Learning Flashcard Generator API${NC}"
echo "====================================================="

# Test 1: Register a new user
echo -e "\n${BLUE}Test 1: Register a new user${NC}"
echo "Sending request to $BASE_URL/auth/register"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo -e "Response: $REGISTER_RESPONSE"

if [[ $REGISTER_RESPONSE == *"User registered successfully"* ]]; then
  echo -e "${GREEN}✓ User registration successful${NC}"
else
  echo -e "${RED}✗ User registration failed${NC}"
fi

# Test 2: Login with the new user
echo -e "\n${BLUE}Test 2: Login with the new user${NC}"
echo "Sending request to $BASE_URL/auth/login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=testuser&password=password123")

echo -e "Response: $LOGIN_RESPONSE"

# Extract the access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')

if [[ -n "$ACCESS_TOKEN" ]]; then
  echo -e "${GREEN}✓ Login successful, received token${NC}"
else
  echo -e "${RED}✗ Login failed, no token received${NC}"
  exit 1
fi

# Test 3: Create a new deck
echo -e "\n${BLUE}Test 3: Create a new deck${NC}"
echo "Sending request to $BASE_URL/decks"
DECK_RESPONSE=$(curl -s -X POST "$BASE_URL/decks" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "Spanish Basics"
  }')

echo -e "Response: $DECK_RESPONSE"

# Extract the deck ID
DECK_ID=$(echo $DECK_RESPONSE | grep -o '"id":[0-9]*' | sed 's/"id"://')

if [[ -n "$DECK_ID" ]]; then
  echo -e "${GREEN}✓ Deck creation successful, received deck ID: $DECK_ID${NC}"
else
  echo -e "${RED}✗ Deck creation failed, no deck ID received${NC}"
  exit 1
fi

# Test 4: Get all decks
echo -e "\n${BLUE}Test 4: Get all decks${NC}"
echo "Sending request to $BASE_URL/decks"
DECKS_RESPONSE=$(curl -s -X GET "$BASE_URL/decks" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo -e "Response: $DECKS_RESPONSE"

if [[ $DECKS_RESPONSE == *"$DECK_ID"* ]]; then
  echo -e "${GREEN}✓ Successfully retrieved decks${NC}"
else
  echo -e "${RED}✗ Failed to retrieve decks${NC}"
fi

# Test 5: Create a flashcard
echo -e "\n${BLUE}Test 5: Create a flashcard${NC}"
echo "Sending request to $BASE_URL/flashcards"
FLASHCARD_RESPONSE=$(curl -s -X POST "$BASE_URL/flashcards" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "word": "hola",
    "example_sentence": "Hola, ¿cómo estás?",
    "translation": "Hello, how are you?",
    "cultural_note": "Common greeting in Spanish-speaking countries",
    "deck_id": '$DECK_ID'
  }')

echo -e "Response: $FLASHCARD_RESPONSE"

# Extract the flashcard ID
FLASHCARD_ID=$(echo $FLASHCARD_RESPONSE | grep -o '"id":[0-9]*' | sed 's/"id"://')

if [[ -n "$FLASHCARD_ID" ]]; then
  echo -e "${GREEN}✓ Flashcard creation successful, received flashcard ID: $FLASHCARD_ID${NC}"
else
  echo -e "${RED}✗ Flashcard creation failed, no flashcard ID received${NC}"
  exit 1
fi

# Test 6: Get flashcards by deck
echo -e "\n${BLUE}Test 6: Get flashcards by deck${NC}"
echo "Sending request to $BASE_URL/decks/$DECK_ID/flashcards"
FLASHCARDS_RESPONSE=$(curl -s -X GET "$BASE_URL/decks/$DECK_ID/flashcards" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo -e "Response: $FLASHCARDS_RESPONSE"

if [[ $FLASHCARDS_RESPONSE == *"$FLASHCARD_ID"* ]]; then
  echo -e "${GREEN}✓ Successfully retrieved flashcards${NC}"
else
  echo -e "${RED}✗ Failed to retrieve flashcards${NC}"
fi

# Test 7: Test LLM generation
echo -e "\n${BLUE}Test 7: Test LLM generation${NC}"
echo "Sending request to $BASE_URL/generate"
GENERATION_RESPONSE=$(curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "word": "hablar",
    "is_verb": true
  }')

echo -e "Response: $GENERATION_RESPONSE"

if [[ $GENERATION_RESPONSE == *"example_sentences"* ]]; then
  echo -e "${GREEN}✓ Successfully generated content${NC}"
else
  echo -e "${RED}✗ Failed to generate content${NC}"
fi

echo -e "\n${GREEN}All tests completed!${NC}"
