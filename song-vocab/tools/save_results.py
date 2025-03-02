from typing import List, Dict, Any, Optional, Tuple
import json
import logging
import traceback
from pathlib import Path
from database.db import Database

# Configure logging
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger('song_vocab')

class SaveError(Exception):
    """Base exception for save-related errors"""
    pass

class DatabaseError(SaveError):
    """Raised when database operations fail"""
    pass

class FileSystemError(SaveError):
    """Raised when file operations fail"""
    pass

def save_results(song_id: str, 
                title: str,
                lyrics: str, 
                vocabulary: List[Dict[str, Any]], 
                lyrics_path: Optional[Path] = None, 
                vocabulary_path: Optional[Path] = None,
                artist: Optional[str] = None,
                romaji_lyrics: Optional[str] = None,
                db: Optional[Database] = None) -> Tuple[bool, List[str]]:
    """
    Save lyrics and vocabulary to database and optionally to files.
    
    Args:
        song_id (str): ID of the song
        title (str): Title of the song
        lyrics (str): Japanese lyrics text
        vocabulary (List[Dict[str, Any]]): List of vocabulary items
        lyrics_path (Path, optional): Directory to save lyrics files
        vocabulary_path (Path, optional): Directory to save vocabulary files
        artist (str, optional): Name of the artist
        romaji_lyrics (str, optional): Romaji version of lyrics
    
    Returns:
        bool: True if all save operations succeeded, False otherwise
    """
    logger.info(f"Saving results for song: {title} ({song_id})")
    success = True
    errors = []

    # Input validation
    if not song_id or not title or not lyrics:
        raise ValueError("song_id, title, and lyrics are required parameters")
    
    if not isinstance(vocabulary, list):
        raise ValueError("vocabulary must be a list")

    try:
        # Initialize database if not provided
        if db is None:
            try:
                db = Database()
            except Exception as e:
                error_msg = f"Failed to initialize database: {str(e)}"
                logger.error(error_msg)
                raise DatabaseError(error_msg)

        # Save to database with detailed error handling
        try:
            if not db.save_song(song_id, title, lyrics, artist, romaji_lyrics):
                error_msg = "Failed to save song to database"
                logger.error(error_msg)
                errors.append(error_msg)
                success = False

            if not db.save_vocabulary(song_id, vocabulary):
                error_msg = "Failed to save vocabulary to database"
                logger.error(error_msg)
                errors.append(error_msg)
                success = False

        except Exception as e:
            error_msg = f"Database operation failed: {str(e)}"
            logger.error(error_msg)
            errors.append(error_msg)
            success = False

        # File system operations with path validation and error handling
        if lyrics_path:
            try:
                lyrics_file = lyrics_path / f"{song_id}.txt"
                lyrics_path.mkdir(parents=True, exist_ok=True)
                lyrics_file.write_text(lyrics, encoding='utf-8')
                logger.info(f"Lyrics saved to {lyrics_file}")
            except Exception as e:
                error_msg = f"Error saving lyrics file: {str(e)}"
                logger.error(f"{error_msg}\n{traceback.format_exc()}")
                errors.append(error_msg)
                success = False

        if vocabulary_path:
            try:
                vocab_file = vocabulary_path / f"{song_id}.json"
                vocabulary_path.mkdir(parents=True, exist_ok=True)
                with open(vocab_file, 'w', encoding='utf-8') as f:
                    json.dump(vocabulary, f, ensure_ascii=False, indent=2)
                logger.info(f"Vocabulary saved to {vocab_file}")
            except Exception as e:
                error_msg = f"Error saving vocabulary file: {str(e)}"
                logger.error(f"{error_msg}\n{traceback.format_exc()}")
                errors.append(error_msg)
                success = False

        return success, errors

    except (DatabaseError, FileSystemError) as e:
        logger.error(f"Critical error in save_results: {str(e)}\n{traceback.format_exc()}")
        errors.append(str(e))
        return False, errors
    except Exception as e:
        error_msg = f"Unexpected error in save_results: {str(e)}"
        logger.error(f"{error_msg}\n{traceback.format_exc()}")
        errors.append(error_msg)
        return False, errors
