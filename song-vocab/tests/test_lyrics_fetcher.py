import pytest
import os
from tools.lyrics_fetcher import fetch_lyrics, LyricsResult, LyricsError, APIError, LyricsNotFoundError

@pytest.mark.asyncio
@pytest.mark.skipif(not os.getenv('MUSIXMATCH_API_KEY'),
                    reason="MUSIXMATCH_API_KEY required for lyrics test")
async def test_fetch_lyrics():
    """Test fetching lyrics for a known Japanese song."""
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
    original_key = os.environ.get('MUSIXMATCH_API_KEY')
    try:
        os.environ.pop('MUSIXMATCH_API_KEY', None)
        with pytest.raises(APIError, match="MUSIXMATCH_API_KEY must be set"):
            await fetch_lyrics("test song")
    finally:
        if original_key:
            os.environ['MUSIXMATCH_API_KEY'] = original_key

@pytest.mark.asyncio
@pytest.mark.skipif(not os.getenv('MUSIXMATCH_API_KEY'),
                    reason="MUSIXMATCH_API_KEY required for lyrics test")
async def test_fetch_lyrics_not_found():
    """Test behavior when lyrics are not found."""
    with pytest.raises(LyricsNotFoundError):
        await fetch_lyrics("この曲は存在しない曲です12345")
