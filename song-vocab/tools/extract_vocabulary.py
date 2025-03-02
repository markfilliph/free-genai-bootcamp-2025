from typing import List, Dict
import ollama
import logging
import json

# Configure logging
logger = logging.getLogger(__name__)

PROMPT = """
Analyze these Japanese lyrics and extract unique vocabulary items. For each word or phrase:
1. Include the original Japanese (kanji/kana)
2. Add the romaji pronunciation
3. Provide a simple English translation

Format each item as a JSON object with these exact fields:
{{
    "kanji": "Japanese text",
    "romaji": "romaji pronunciation",
    "english": "English meaning"
}}

Combine all items into a JSON array.

Text to analyze:
{text}
"""

async def extract_vocabulary(text: str) -> List[Dict]:
    """Extract vocabulary from Japanese text using LLM.
    
    Args:
        text: Japanese text to extract vocabulary from
        
    Returns:
        List of vocabulary items with kanji, romaji, and english translations
    """
    logger.info("Starting vocabulary extraction")
    
    try:
        # Initialize Ollama client
        client = ollama.Client()
        
        # Call Ollama API
        response = client.chat(
            model="mistral",
            messages=[
                {"role": "system", "content": "You are a Japanese language expert. Extract vocabulary from the given text and format it as a JSON array."},
                {"role": "user", "content": PROMPT.format(text=text)}
            ]
        )
        
        # Parse the JSON response
        content = response['message']['content']
        logger.info(f"Raw LLM response:\n{content}")
        
        # Find the JSON array
        start = content.find('[{')
        end = content.rfind('}]') + 2
        
        if start >= 0 and end > start:
            # Extract JSON and clean it up
            json_str = content[start:end]
            json_lines = []
            for line in json_str.split('\n'):
                line = line.strip()
                # Skip comments and empty lines
                if not line or line.startswith('//'):
                    continue
                # Remove trailing commas
                if line.rstrip().endswith(',') and line.strip() != '{':
                    line = line.rstrip().rstrip(',')
                json_lines.append(line)
            json_str = '\n'.join(json_lines)
            
            try:
                items = json.loads(json_str)
                
                # Validate and normalize items
                vocabulary = []
                for item in items:
                    if all(k in item for k in ['kanji', 'romaji', 'english']):
                        # Ensure romaji is a string
                        if isinstance(item['romaji'], list):
                            item['romaji'] = ' '.join(item['romaji'])
                        vocabulary.append(item)
                
                logger.info(f"Extracted {len(vocabulary)} vocabulary items")
                return vocabulary
                
            except json.JSONDecodeError as e:
                logger.error(f"Failed to parse JSON: {str(e)}")
        else:
            logger.warning("No JSON array found in response")
            
        # Return empty list if we couldn't extract vocabulary
        return []
            
    except Exception as e:
        logger.error(f"Failed to extract vocabulary: {str(e)}")
        return []