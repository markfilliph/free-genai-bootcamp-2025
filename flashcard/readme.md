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