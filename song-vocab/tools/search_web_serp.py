import aiohttp
from typing import List, Dict, Optional, Union
import logging
from urllib.parse import quote
import asyncio
import json
from dataclasses import dataclass
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Google Custom Search API configuration
API_KEY = os.getenv('GOOGLE_API_KEY')
SEARCH_ENGINE_ID = os.getenv('SEARCH_ENGINE_ID')
BASE_URL = 'https://www.googleapis.com/customsearch/v1'

@dataclass
class SearchResult:
    title: str
    url: str
    snippet: str

class SearchError(Exception):
    """Base exception for search-related errors"""
    pass

class NetworkError(SearchError):
    """Raised when network-related issues occur"""
    pass

class QueryError(SearchError):
    """Raised when there are issues with the search query"""
    pass

class ConfigError(SearchError):
    """Raised when there are configuration issues"""
    pass

class RateLimitError(SearchError):
    """Raised when API rate limit is exceeded"""
    pass

async def search_web_serp(query: str, max_results: int = 5) -> Union[List[SearchResult], List[Dict[str, str]]]:
    """
    Search for Japanese song lyrics using Google Custom Search.
    
    Args:
        query (str): Search query for the song lyrics
        max_results (int): Maximum number of search results to return
        
    Returns:
        List[Dict[str, str]]: List of search results with title and url
    """
    if not query or not query.strip():
        raise QueryError("Search query cannot be empty")

    if max_results < 1:
        raise ValueError("max_results must be greater than 0")

    if not API_KEY or not SEARCH_ENGINE_ID:
        raise ConfigError("API_KEY and SEARCH_ENGINE_ID must be set in environment variables")

    try:
        logger.info(f"Starting search for: {query}")
        
        # Add Japanese-specific keywords to improve results and avoid translations
        japanese_keywords = ["歌詞", "原曲", "日本語", "-英訳", "-英語"]
        enhanced_query = f"{query} {' '.join(japanese_keywords)}"
        logger.info(f"Enhanced query: {enhanced_query}")
        
        # Prepare API request parameters
        params = {
            'key': API_KEY,
            'cx': SEARCH_ENGINE_ID,
            'q': enhanced_query,
            'num': max_results,
            'lr': 'lang_ja',  # Restrict to Japanese results
            'safe': 'off'
        }
        
        try:
            async with aiohttp.ClientSession() as session:
                try:
                    async with session.get(BASE_URL, params=params, timeout=10) as response:
                        if response.status == 429:
                            logger.error("API rate limit exceeded")
                            raise RateLimitError("Google Custom Search API rate limit exceeded")
                        
                        if response.status != 200:
                            error_text = await response.text()
                            logger.error(f"API request failed with status {response.status}: {error_text}")
                            raise NetworkError(f"API request failed with status {response.status}")
                        
                        data = await response.json()
                        
                        if 'error' in data:
                            error_msg = data['error'].get('message', 'Unknown API error')
                            logger.error(f"API error: {error_msg}")
                            raise SearchError(f"API error: {error_msg}")
                        
                        if 'items' not in data:
                            logger.warning(f"No results found for query: {query}")
                            return []
                        
                        results = []
                        for item in data['items']:
                            result = SearchResult(
                                title=item.get('title', ''),
                                url=item.get('link', ''),
                                snippet=item.get('snippet', '')
                            )
                            results.append(result)
                        
                        logger.info(f"Found {len(results)} results")
                        return results[:max_results]
                        
                except asyncio.TimeoutError:
                    logger.error("API request timed out")
                    raise NetworkError("API request timed out")
                    
                except aiohttp.ClientError as e:
                    logger.error(f"Network error during API request: {str(e)}")
                    raise NetworkError(f"Network error during API request: {str(e)}")
                
        except Exception as e:
            logger.error(f"Error processing search results: {str(e)}")
            raise SearchError(f"Failed to process search results: {str(e)}")
            
    except Exception as e:
        logger.error(f"Unexpected error during search: {str(e)}")
        raise SearchError(f"Search operation failed: {str(e)}")
