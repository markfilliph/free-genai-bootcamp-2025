import os
import logging
import aiohttp
import json
from typing import Dict, Optional, List
from dataclasses import dataclass
from urllib.parse import quote
from time import sleep

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Constants
BASE_URL = "https://api.musixmatch.com/ws/1.1"

def get_api_key() -> str:
    """Get the Musixmatch API key from environment variables.
    
    Returns:
        str: The API key if set
        
    Raises:
        APIError: If API key is not set
    """
    api_key = os.getenv('MUSIXMATCH_API_KEY')
    if not api_key:
        raise APIError("MUSIXMATCH_API_KEY must be set in environment variables")
    return api_key

@dataclass
class LyricsResult:
    """Represents the result of a lyrics fetch operation.
    
    Attributes:
        title: The title of the song
        artist: The artist name
        lyrics: The song lyrics text
        language: The language code of the lyrics (e.g. 'ja' for Japanese)
    """
    title: str
    artist: str
    lyrics: str
    language: Optional[str] = None

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
    """Search for a track using Musixmatch API.
    
    Args:
        title: The song title to search for
        artist: Optional artist name to refine the search
        
    Raises:
        APIError: If API request fails or credentials are missing
        LyricsNotFoundError: If no matching tracks are found
        
    Returns:
        Dict: Track information from Musixmatch API
    """
    if not title or not title.strip():
        raise ValueError("title cannot be empty")

    try:
        api_key = get_api_key()
    except APIError:
        raise

    # Build search query
    params = {
        'q_track': title,
        'f_has_lyrics': 1,
        'f_lyrics_language': 'ja',
        'apikey': api_key
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
                try:
                    data = json.loads(text)
                except json.JSONDecodeError as e:
                    raise APIError(f"Invalid JSON response: {str(e)}")
                
                # Add delay for rate limiting
                sleep(0.1)
                
                if data['message']['header']['status_code'] != 200:
                    raise APIError(f"API error: {data['message']['header']}")
                
                track_list = data['message']['body']['track_list']
                if not track_list:
                    raise LyricsNotFoundError(f"No tracks found for {title}")
                
                return track_list[0]['track']
                
    except aiohttp.ClientError as e:
        raise APIError(f"Network error: {str(e)}")
    except LyricsNotFoundError:
        raise
    except Exception as e:
        raise LyricsError(f"Unexpected error: {str(e)}")

async def get_lyrics(track_id: int) -> str:
    """Get lyrics for a track using Musixmatch API.
    
    Args:
        track_id: The track ID from Musixmatch API
        
    Raises:
        ValueError: If track_id is invalid
        APIError: If API request fails
        LyricsError: For unexpected errors
        
    Returns:
        str: The lyrics text
    """
    if not isinstance(track_id, int) or track_id <= 0:
        raise ValueError("track_id must be a positive integer")
    """Get lyrics for a track using Musixmatch API.
    
    Args:
        track_id: The track ID from Musixmatch API
        
    Raises:
        APIError: If track_id is invalid or API request fails
        LyricsError: For unexpected errors
        
    Returns:
        str: The lyrics text
    """
    if not isinstance(track_id, int) or track_id <= 0:
        raise ValueError("track_id must be a positive integer")
    """Get lyrics for a track using Musixmatch API."""
    try:
        api_key = get_api_key()
    except APIError:
        raise

    params = {
        'track_id': track_id,
        'apikey': api_key
    }

    try:
        headers = {'Accept': 'application/json'}
        async with aiohttp.ClientSession() as session:
            async with session.get(f"{BASE_URL}/track.lyrics.get", params=params, headers=headers) as response:
                if response.status != 200:
                    raise APIError(f"API request failed with status {response.status}")
                
                # Read response text and parse as JSON regardless of content-type
                text = await response.text()
                try:
                    data = json.loads(text)
                except json.JSONDecodeError as e:
                    raise APIError(f"Invalid JSON response: {str(e)}")
                
                # Add delay for rate limiting
                sleep(0.1)
                
                if data['message']['header']['status_code'] != 200:
                    raise APIError(f"API error: {data['message']['header']}")
                
                return data['message']['body']['lyrics']['lyrics_body']
                
    except aiohttp.ClientError as e:
        raise APIError(f"Network error: {str(e)}")
    except LyricsNotFoundError:
        raise
    except Exception as e:
        raise LyricsError(f"Unexpected error: {str(e)}")

async def fetch_lyrics(title: str, artist: Optional[str] = None) -> LyricsResult:
    """Main function to fetch lyrics for a song.
    
    Args:
        title: The song title to search for
        artist: Optional artist name to refine the search
        
    Raises:
        ValueError: If title is empty
        APIError: If API request fails
        LyricsNotFoundError: If lyrics cannot be found
        LyricsError: For unexpected errors
        
    Returns:
        LyricsResult: Object containing song title, artist, and lyrics
    """
    if not title or not title.strip():
        raise ValueError("title cannot be empty")
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
