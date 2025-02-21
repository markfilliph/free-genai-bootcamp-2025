from youtube_transcript_api import YouTubeTranscriptApi
from typing import Optional, List, Dict, Any, Union
import asyncio
import aiofiles
import json
import os
from datetime import datetime
from pathlib import Path
from functools import lru_cache


class YouTubeTranscriptDownloader:
    def __init__(self, languages: List[str] = ["ja", "en"]):
        self.languages = languages
        self.cache_dir = Path("data/transcripts")
        self.cache_dir.mkdir(parents=True, exist_ok=True)

    def extract_video_id(self, url: str) -> Optional[str]:
        """Extract video ID from YouTube URL
        
        Args:
            url (str): YouTube URL
            
        Returns:
            Optional[str]: Video ID if found, None otherwise
        """
        if "v=" in url:
            return url.split("v=")[1][:11]
        elif "youtu.be/" in url:
            return url.split("youtu.be/")[1][:11]
        return None

    def _get_cache_path(self, video_id: str) -> Path:
        """Get the cache file path for a video ID"""
        return self.cache_dir / f"{video_id}.json"

    async def _read_cache(self, video_id: str) -> Optional[Dict[str, Any]]:
        """Read transcript from cache if it exists"""
        cache_path = self._get_cache_path(video_id)
        if not cache_path.exists():
            return None

        try:
            async with aiofiles.open(cache_path, 'r', encoding='utf-8') as f:
                data = json.loads(await f.read())
                # Check if cache is still valid (7 days)
                cached_time = datetime.fromisoformat(data['cached_at'])
                if (datetime.now() - cached_time).days > 7:
                    return None
                return data
        except Exception as e:
            print(f"Cache read error: {str(e)}")
            return None

    async def _write_cache(self, video_id: str, transcript: List[Dict[str, Union[str, float]]]) -> None:
        """Write transcript to cache"""
        cache_path = self._get_cache_path(video_id)
        cache_data = {
            'video_id': video_id,
            'transcript': transcript,
            'cached_at': datetime.now().isoformat(),
            'language': self.languages[0]
        }

        try:
            async with aiofiles.open(cache_path, 'w', encoding='utf-8') as f:
                await f.write(json.dumps(cache_data, ensure_ascii=False, indent=2))
        except Exception as e:
            print(f"Cache write error: {str(e)}")

    def _fetch_transcript(self, video_id: str) -> Optional[List[Dict[str, Union[str, float]]]]:
        """Fetch transcript from YouTube"""
        try:
            return YouTubeTranscriptApi.get_transcript(video_id, languages=self.languages)
        except Exception as e:
            print(f"YouTube API error: {str(e)}")
            return None

    async def get_transcript(self, video_id: str) -> Optional[Dict[str, Any]]:
        """Get transcript with caching
        
        Args:
            video_id (str): YouTube video ID or URL
            
        Returns:
            Optional[Dict[str, Any]]: Transcript data if successful, None otherwise
        """
        # Extract video ID if full URL is provided
        if "youtube.com" in video_id or "youtu.be" in video_id:
            video_id = self.extract_video_id(video_id)
            
        if not video_id:
            print("Invalid video ID or URL")
            return None

        # Try cache first
        cached = await self._read_cache(video_id)
        if cached:
            return cached

        # Fetch from YouTube
        transcript = await asyncio.to_thread(self._fetch_transcript, video_id)
        
        if not transcript:
            return None

        # Write to cache
        await self._write_cache(video_id, transcript)

        return {
            'video_id': video_id,
            'transcript': transcript,
            'cached_at': datetime.now().isoformat(),
            'language': self.languages[0]
        }

    async def save_transcript(self, transcript_data: Dict[str, Any], format: str = 'json') -> bool:
        """Save transcript to file
        
        Args:
            transcript_data (Dict[str, Any]): Transcript data including metadata
            format (str): Output format ('json' or 'txt')
            
        Returns:
            bool: True if successful, False otherwise
        """
        video_id = transcript_data['video_id']
        transcript = transcript_data['transcript']
        
        try:
            if format == 'json':
                output_path = self.cache_dir / f"{video_id}.json"
                async with aiofiles.open(output_path, 'w', encoding='utf-8') as f:
                    await f.write(json.dumps(transcript_data, ensure_ascii=False, indent=2))
            else:
                output_path = self.cache_dir / f"{video_id}.txt"
                async with aiofiles.open(output_path, 'w', encoding='utf-8') as f:
                    for entry in transcript:
                        await f.write(f"{entry['text']}\n")
            return True
        except Exception as e:
            print(f"Error saving transcript: {str(e)}")
            return False


async def main(video_url: str, print_transcript: bool = False) -> None:
    """Main function to demonstrate usage"""
    try:
        # Initialize downloader
        downloader = YouTubeTranscriptDownloader()
        
        # Get transcript
        transcript_data = await downloader.get_transcript(video_url)
        if not transcript_data:
            print("Failed to get transcript")
            return

        # Save transcript
        if await downloader.save_transcript(transcript_data):
            print(f"Transcript saved successfully for video {transcript_data['video_id']}")
            
            if print_transcript:
                for entry in transcript_data['transcript']:
                    print(f"{entry['text']}")
        else:
            print("Failed to save transcript")

    except Exception as e:
        print(f"An error occurred: {str(e)}")


if __name__ == "__main__":
    video_url = "https://www.youtube.com/watch?v=sY7L5cfCWno"
    asyncio.run(main(video_url, print_transcript=True))
