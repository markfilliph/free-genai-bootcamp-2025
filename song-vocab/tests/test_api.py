import pytest
import os
import logging
from tools.search_web_serp import search_web_serp, SearchResult

# Configure logging
logging.basicConfig(level=logging.INFO,
                   format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

@pytest.mark.skipif(not os.getenv('GOOGLE_API_KEY') or not os.getenv('SEARCH_ENGINE_ID'),
                    reason="API credentials required for API test")
@pytest.mark.asyncio
async def test_api_search():
    """Test actual API search with real credentials."""
    try:
        # Test with a well-known Japanese song
        results = await search_web_serp("上を向いて歩こう 歌詞", max_results=3)
        
        # Verify we got results
        assert len(results) > 0, "Should return at least one result"
        
        # Verify result structure
        for result in results:
            assert isinstance(result, SearchResult), "Each result should be a SearchResult instance"
            assert result.title, "Result should have a title"
            assert result.url, "Result should have a URL"
            assert result.snippet, "Result should have a snippet"
            
            # Log the results for inspection
            logger.info(f"\nTitle: {result.title}")
            logger.info(f"URL: {result.url}")
            logger.info(f"Snippet: {result.snippet}")
            
    except Exception as e:
        logger.error(f"API test failed: {str(e)}")
        raise
