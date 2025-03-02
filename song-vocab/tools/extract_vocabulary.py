from typing import List, Dict
import ollama
import logging
import json
from pathlib import Path

# Configure logging
logger = logging.getLogger(__name__)

async def extract_vocabulary(text: str) -> List[Dict]:
    """
    Extract ALL vocabulary from Japanese text using LLM.
    
    Args:
        text (str): The text to extract vocabulary from
        
    Returns:
        List[Dict]: Complete list of vocabulary items in Japanese format with kanji, romaji, and parts
    """
    logger.info("Starting vocabulary extraction")
    logger.debug(f"Input text length: {len(text)} characters")
    
    try:
        # Initialize Ollama client
        logger.debug("Initializing Ollama client")
        client = ollama.Client()
        
        # Load the prompt from the prompts directory
        prompt_path = Path(__file__).parent.parent / "prompts" / "Extract-Vocabulary.md"
        logger.debug(f"Loading prompt from {prompt_path}")
        with open(prompt_path, 'r', encoding='utf-8') as f:
            prompt_template = f.read()
        
        # Construct the full prompt with the text to analyze
        prompt = f"{prompt_template}\n\nText to analyze:\n{text}"
        logger.debug(f"Constructed prompt of length {len(prompt)}")
        
        # We'll use multiple calls to ensure we get all vocabulary
        all_vocabulary = set()
        max_attempts = 3
        
        for attempt in range(max_attempts):
            logger.info(f"Making LLM call attempt {attempt + 1}/{max_attempts}")
            try:
                response = client.chat(
                    model="mistral",
                    messages=[
                        {"role": "system", "content": "You are a Japanese language expert. Extract vocabulary from the given text and format it as a JSON array with kanji, romaji, english meaning, and parts breakdown. Each vocabulary item should have this structure: {\"kanji\": string, \"romaji\": string, \"english\": string, \"parts\": [{\"kanji\": string, \"romaji\": [string]}]}"},
                        {"role": "user", "content": prompt}
                    ]
                )
                
                # Parse the JSON response
                content = response['message']['content']
                # Find the JSON array in the response
                start = content.find('[{')
                end = content.rfind('}]') + 2
                if start >= 0 and end > start:
                    json_str = content[start:end]
                    try:
                        items = json.loads(json_str)
                        # Add new vocabulary items to our set
                        for item in items:
                            item_tuple = tuple(sorted(item.items()))
                            all_vocabulary.add(item_tuple)
                        logger.info(f"Attempt {attempt + 1} added {len(items)} items")
                    except json.JSONDecodeError as e:
                        logger.error(f"Failed to parse JSON in attempt {attempt + 1}: {str(e)}")
                else:
                    logger.warning(f"No JSON array found in attempt {attempt + 1}")
                
            except Exception as e:
                logger.error(f"Error in attempt {attempt + 1}: {str(e)}")
                if attempt == max_attempts - 1:
                    raise  # Re-raise on last attempt
        
        # Convert back to list of dicts
        result = [dict(item) for item in all_vocabulary]
        logger.info(f"Extracted {len(result)} unique vocabulary items")
        return result
        
    except Exception as e:
        logger.error(f"Failed to extract vocabulary: {str(e)}", exc_info=True)
        raise