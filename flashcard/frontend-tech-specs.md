# Frontend Technical Specifications

## Overview
The frontend of the Language Learning Flashcard Generator is built using **Svelte** and **Svelte Material UI**. It provides a user-friendly interface for creating, organizing, and reviewing flashcards. The frontend integrates with the backend APIs and uses **ResponsiveVoice.js** for text-to-speech functionality.

---

## Technology Stack
- **Framework**: Svelte
- **UI Components**: Svelte Material UI
- **Localization**: svelte-i18n
- **Text-to-Speech**: ResponsiveVoice.js
- **State Management**: Svelte stores
- **Routing**: svelte-spa-router

---

## Key Features
1. **Flashcard Creation**:
   - Users can input a word, phrase, or verb in Spanish.
   - The app displays generated example sentences, conjugations, and cultural notes.

2. **Verb Conjugation Support**:
   - Display conjugations for all tenses and moods.

3. **Cultural Context**:
   - Show cultural notes and idiomatic expressions.

4. **Text-to-Speech (TTS)**:
   - Users can listen to pronunciations using ResponsiveVoice.js.

5. **Flashcard Organization**:
   - Organize flashcards into decks and tag them.

6. **Review Mode**:
   - Flashcards are reviewed using a spaced repetition system (SuperMemo2).

7. **Export Flashcards**:
   - Export flashcards as PDF or CSV.

---

## Pages and Components

### Pages
1. **Home Page**:
   - Welcome message and quick links to create flashcards or review.

2. **Login/Register Page**:
   - Forms for user authentication.

3. **Flashcard Creation Page**:
   - Input fields for word/phrase, generated content, and customization options.

4. **Deck Management Page**:
   - List of decks, options to create/edit/delete decks.

5. **Review Page**:
   - Flashcards are displayed for review, with options to mark as "easy," "medium," or "hard."

6. **Export Page**:
   - Options to export flashcards as PDF or CSV.

---

### Components
1. **Navbar**:
   - Navigation links to all pages.

2. **Flashcard Form**:
   - Form for creating/editing flashcards.

3. **Deck List**:
   - Displays all decks with options to manage them.

4. **Flashcard Viewer**:
   - Displays flashcards with example sentences, conjugations, and cultural notes.

5. **Review Card**:
   - Displays a flashcard for review with TTS and review options.

6. **Export Options**:
   - Buttons to export flashcards as PDF or CSV.

---

## State Management
- **Svelte Stores**:
  - `userStore`: Manages user authentication state.
  - `deckStore`: Manages deck and flashcard data.
  - `reviewStore`: Manages review session state.

---

## Text-to-Speech (TTS)
- **ResponsiveVoice.js**:
  - Integrated into the `Flashcard Viewer` and `Review Card` components.
  - Users can click a button to hear the pronunciation of words and sentences.

---

## Localization
- **svelte-i18n**:
  - Supports multiple languages for the UI.
  - Language files are stored in `src/locales`.

---

## Routing
- **svelte-spa-router**:
  - Handles client-side routing for navigation between pages.

---

## Environment Variables
- `API_BASE_URL`: Base URL for backend API.
- `TTS_ENABLED`: Enable/disable text-to-speech functionality.

---

## Setup Instructions
1. Install dependencies:
   ```bash
   npm install