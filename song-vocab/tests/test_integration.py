import asyncio
import logging
import pytest
from pathlib import Path
from tools.search_web_serp import search_web_serp, SearchError, NetworkError, QueryError, ConfigError, SearchResult
from tools.save_results import save_results
from database.db import Database
import shutil
import os
from unittest.mock import patch

# Configure logging
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

@pytest.fixture
def test_db():
    """Create a test database and clean it up after tests."""
    test_db_path = Path("test_data/test.db")
    test_db_path.parent.mkdir(parents=True, exist_ok=True)
    
    # Remove existing test database if it exists
    if test_db_path.exists():
        os.remove(test_db_path)
    
    db = Database(str(test_db_path))
    yield db
    
    # Cleanup
    shutil.rmtree("test_data", ignore_errors=True)

@pytest.fixture
def test_paths():
    """Create test paths for lyrics and vocabulary files."""
    test_dir = Path("test_data")
    lyrics_path = test_dir / "lyrics"
    vocab_path = test_dir / "vocabulary"
    
    # Create directories
    lyrics_path.mkdir(parents=True, exist_ok=True)
    vocab_path.mkdir(parents=True, exist_ok=True)
    
    yield lyrics_path, vocab_path
    
    # Cleanup handled by test_db fixture

@pytest.mark.asyncio
@pytest.mark.skipif(not os.getenv('GOOGLE_API_KEY') or not os.getenv('SEARCH_ENGINE_ID'),
                    reason="API credentials required for integration test")
async def test_search_and_save_flow(test_db, test_paths):
    """Test the entire flow from search to save."""
    try:
        lyrics_path, vocab_path = test_paths
        
        # Test the complete search and save flow
        results = await search_web_serp("上を向いて歩こう", max_results=1)
        assert len(results) > 0, "Search should return at least one result"
        
        # Use the first result for saving
        result = results[0]
        assert isinstance(result, SearchResult), "Search should return SearchResult instances"
        
        # Test save with real data
        song_id = "test_song_001"
        title = result.title
        lyrics = "上を向いて歩こう\n涙がこぼれないように"
        vocabulary = [
            {
                "kanji": "上",
                "romaji": "ue",
                "english": "up, above",
                "parts": {"pos": "noun"}
            }
        ]
        
        # Test save_results
        try:
            success, errors = save_results(
                song_id=song_id,
                title=title,
                lyrics=lyrics,
                vocabulary=vocabulary,
                lyrics_path=lyrics_path,
                vocabulary_path=vocab_path,
                db=test_db
            )
            assert success, f"Save operation failed with errors: {errors}"
            assert len(errors) == 0, f"Unexpected errors during save: {errors}"
        except Exception as e:
            logger.error(f"Save operation failed: {e}")
            raise
        
        # Verify database entries
        try:
            song = test_db.get_song(song_id)
            assert song is not None, "Song should be in database"
            assert song['title'] == title, "Song title should match"
            assert song['lyrics'] == lyrics, "Song lyrics should match"
            
            vocab_items = test_db.get_vocabulary(song_id)
            assert len(vocab_items) == 1, "Should have one vocabulary item"
            assert vocab_items[0]['kanji'] == "上", "Vocabulary kanji should match"
        except Exception as e:
            logger.error(f"Database verification failed: {e}")
            raise
        
        # Verify files
        try:
            lyrics_file = lyrics_path / f"{song_id}.txt"
            vocab_file = vocab_path / f"{song_id}.json"
            assert lyrics_file.exists(), "Lyrics file should exist"
            assert vocab_file.exists(), "Vocabulary file should exist"
        except Exception as e:
            logger.error(f"File verification failed: {e}")
            raise
        
        # Test file cleanup
        try:
            lyrics_file = lyrics_path / f"{song_id}.txt"
            vocab_file = vocab_path / f"{song_id}.json"
            lyrics_file.unlink(missing_ok=True)
            vocab_file.unlink(missing_ok=True)
        except Exception as e:
            logger.error(f"Cleanup failed: {e}")
            raise
    
    except Exception as e:
        logger.error(f"Test failed: {e}")
        raise

if __name__ == "__main__":
    pytest.main([__file__, "-v"])
