import os
import logging
import aiohttp
from typing import Dict, Optional, List
from dataclasses import dataclass
from urllib.parse import quote
import json

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Constants
MUSIXMATCH_API_KEY = os.getenv('MUSIXMATCH_API_KEY')
BASE_URL = "https://api.musixmatch.com/ws/1.1"

@dataclass
class LyricsResult:
    title: str
    artist: str
    lyrics: str
    language: str = "ja"

class LyricsError(Exception):
    """Base exception for lyrics-related errors"""
    pass

class APIError(LyricsError):
    """Raised when there are API-related issues"""
    pass

class LyricsNotFoundError(LyricsError):
    """Raised when lyrics cannot be found"""
    pass

async def search_track(title: str, artist: Optional[str] = None) -> Dict:
    """Search for a track using Musixmatch API."""
    if not MUSIXMATCH_API_KEY:
        raise APIError("MUSIXMATCH_API_KEY must be set in environment variables")

    # Build search query
    params = {
        'q_track': title,
        'f_has_lyrics': 1,
        'f_lyrics_language': 'ja',
        'apikey': MUSIXMATCH_API_KEY
    }
    
    if artist:
        params['q_artist'] = artist

    try:
        headers = {'Accept': 'application/json'}
        async with aiohttp.ClientSession() as session:
            async with session.get(f"{BASE_URL}/track.search", params=params, headers=headers) as response:
                if response.status != 200:
                    raise APIError(f"API request failed with status {response.status}")
                
                # Read response text and parse as JSON regardless of content-type
                text = await response.text()
                data = loads(text)
                
                if data['message']['header']['status_code'] != 200:
                    raise APIError(f"API error: {data['message']['header']}")
                
                track_list = data['message']['body']['track_list']
                if not track_list:
                    raise LyricsNotFoundError(f"No tracks found for {title}")
                
                return track_list[0]['track']
                
    except aiohttp.ClientError as e:
        raise APIError(f"Network error: {str(e)}")
    except Exception as e:
        raise LyricsError(f"Unexpected error: {str(e)}")

async def get_lyrics(track_id: int) -> str:
    """Get lyrics for a track using Musixmatch API."""
    if not MUSIXMATCH_API_KEY:
        raise APIError("MUSIXMATCH_API_KEY must be set in environment variables")

    params = {
        'track_id': track_id,
        'apikey': MUSIXMATCH_API_KEY
    }

    try:
        headers = {'Accept': 'application/json'}
        async with aiohttp.ClientSession() as session:
            async with session.get(f"{BASE_URL}/track.lyrics.get", params=params, headers=headers) as response:
                if response.status != 200:
                    raise APIError(f"API request failed with status {response.status}")
                
                # Read response text and parse as JSON regardless of content-type
                text = await response.text()
                data = loads(text)
                
                if data['message']['header']['status_code'] != 200:
                    raise APIError(f"API error: {data['message']['header']}")
                
                return data['message']['body']['lyrics']['lyrics_body']
                
    except aiohttp.ClientError as e:
        raise APIError(f"Network error: {str(e)}")
    except Exception as e:
        raise LyricsError(f"Unexpected error: {str(e)}")

async def fetch_lyrics(title: str, artist: Optional[str] = None) -> LyricsResult:
    """Main function to fetch lyrics for a song."""
    try:
        # Search for the track
        track = await search_track(title, artist)
        
        # Get the lyrics
        lyrics = await get_lyrics(track['track_id'])
        
        return LyricsResult(
            title=track['track_name'],
            artist=track['artist_name'],
            lyrics=lyrics
        )
        
    except (APIError, LyricsNotFoundError) as e:
        logger.error(f"Failed to fetch lyrics: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Failed to fetch lyrics: {str(e)}")
        raise LyricsError(f"Unexpected error: {str(e)}")
