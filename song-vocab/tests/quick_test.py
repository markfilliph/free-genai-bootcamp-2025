import asyncio
from tools.lyrics_fetcher import fetch_lyrics

async def test():
    try:
        result = await fetch_lyrics("上を向いて歩こう", "坂本九")
        print(f"Success! Found lyrics for {result.title} by {result.artist}")
        print("\nFirst few lines of lyrics:")
        print("\n".join(result.lyrics.split("\n")[:5]))
    except Exception as e:
        print(f"Error: {str(e)}")

if __name__ == "__main__":
    asyncio.run(test())
