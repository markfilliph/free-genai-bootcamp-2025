import asyncio
import logging
import pytest
from unittest.mock import patch, MagicMock
from tools.search_web_serp import search_web_serp, SearchError, NetworkError, QueryError, ConfigError, SearchResult
from dotenv import load_dotenv
import os

# Configure logging
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

load_dotenv()

@pytest.mark.asyncio
async def test_missing_credentials():
    """Test behavior when API credentials are missing"""
    with patch.dict('os.environ', {}, clear=True):
        with patch('tools.search_web_serp.API_KEY', None):
            with patch('tools.search_web_serp.SEARCH_ENGINE_ID', None):
                with pytest.raises(ConfigError, match="API_KEY and SEARCH_ENGINE_ID must be set"):
                    await search_web_serp("test query")

@pytest.mark.asyncio
async def test_empty_query():
    """Test behavior with empty query"""
    with pytest.raises(QueryError, match="Search query cannot be empty"):
        await search_web_serp("")
    with pytest.raises(QueryError, match="Search query cannot be empty"):
        await search_web_serp("   ")

@pytest.mark.asyncio
async def test_invalid_max_results():
    """Test behavior with invalid max_results"""
    with pytest.raises(ValueError, match="max_results must be greater than 0"):
        await search_web_serp("test", max_results=0)
    with pytest.raises(ValueError, match="max_results must be greater than 0"):
        await search_web_serp("test", max_results=-1)

@pytest.mark.asyncio
async def test_successful_search():
    """Test successful search with mocked API response"""
    mock_response = {
        'items': [{
            'title': '上を向いて歩こう',
            'link': 'https://example.com/song/1',
            'snippet': 'Test snippet'
        }]
    }
    
    class MockResponse:
        def __init__(self):
            self.status = 200
        
        async def json(self):
            return mock_response
        
        async def __aenter__(self):
            return self
            
        async def __aexit__(self, exc_type, exc_val, exc_tb):
            pass
    
    with patch('tools.search_web_serp.API_KEY', 'mock_key'):
        with patch('tools.search_web_serp.SEARCH_ENGINE_ID', 'mock_id'):
            with patch('aiohttp.ClientSession.get', return_value=MockResponse()):
                results = await search_web_serp("上を向いて歩こう", max_results=1)
                assert len(results) == 1
                assert isinstance(results[0], SearchResult)
                assert results[0].title == '上を向いて歩こう'