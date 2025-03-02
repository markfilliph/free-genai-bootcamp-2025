import pytest
import os
import importlib
import tools.lyrics_fetcher
from tools.lyrics_fetcher import LyricsResult, LyricsError, APIError, LyricsNotFoundError

@pytest.mark.asyncio
@pytest.mark.skipif(not os.getenv('MUSIXMATCH_API_KEY'),
                    reason="MUSIXMATCH_API_KEY required for lyrics test")
async def test_fetch_lyrics():
    """Test fetching lyrics for a known Japanese song."""
    from tools.lyrics_fetcher import fetch_lyrics
    result = await fetch_lyrics("上を向いて歩こう", "坂本九")
    
    # Verify result structure
    assert isinstance(result, LyricsResult)
    assert result.title == "上を向いて歩こう"
    assert "坂本九" in result.artist
    assert "上を向いて歩こう" in result.lyrics
    assert "涙がこぼれないように" in result.lyrics
    
@pytest.mark.asyncio
async def test_fetch_lyrics_no_credentials():
    """Test behavior when API credentials are missing."""
    # Store original key and function
    original_key = os.environ.get('MUSIXMATCH_API_KEY')
    
    # Remove key from environment
    if 'MUSIXMATCH_API_KEY' in os.environ:
        del os.environ['MUSIXMATCH_API_KEY']
    
    # Force reload of lyrics_fetcher to clear cached environment variables
    tools.lyrics_fetcher = importlib.reload(tools.lyrics_fetcher)
    from tools.lyrics_fetcher import fetch_lyrics
    
    try:
        with pytest.raises(APIError, match="MUSIXMATCH_API_KEY must be set in environment variables"):
            await fetch_lyrics("test song")
    finally:
        # Restore original key and reload module
        if original_key:
            os.environ['MUSIXMATCH_API_KEY'] = original_key
            tools.lyrics_fetcher = importlib.reload(tools.lyrics_fetcher)

@pytest.mark.asyncio
@pytest.mark.skipif(not os.getenv('MUSIXMATCH_API_KEY'),
                    reason="MUSIXMATCH_API_KEY required for lyrics test")
async def test_fetch_lyrics_not_found():
    """Test behavior when lyrics are not found."""
    from tools.lyrics_fetcher import fetch_lyrics
    with pytest.raises(LyricsNotFoundError):
        await fetch_lyrics("この曲は存在しない曲です12345")
