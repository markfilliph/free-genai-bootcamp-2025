import asyncio
from get_transcript import YouTubeTranscriptDownloader
from vector_store import TranscriptVectorStore

async def main():
    # Initialize components
    transcript_downloader = YouTubeTranscriptDownloader()
    vector_store = TranscriptVectorStore()
    
    # Test video URLs (JLPT N5 listening samples)
    test_videos = [
        "https://www.youtube.com/watch?v=sY7L5cfCWno",  # JLPT N5 Sample 1
        "https://www.youtube.com/watch?v=HD9IMmwHEGw",  # JLPT N5 Sample 2
    ]
    
    # Download and index transcripts
    for video_url in test_videos:
        print(f"\nProcessing video: {video_url}")
        
        # Get transcript
        transcript_data = await transcript_downloader.get_transcript(video_url)
        if not transcript_data:
            print("Failed to get transcript")
            continue
            
        # Add to vector store
        video_id = transcript_data['video_id']
        await vector_store.add_transcript(video_id, transcript_data)
        print(f"Added transcript to vector store: {video_id}")
        
        # Test similarity search
        test_queries = [
            "listening test instructions",
            "how to answer questions",
            "conversation at a restaurant",
            "asking for directions",
            "train station conversation"
        ]
        
        print("\nTesting similarity search:")
        for query in test_queries:
            print(f"\nQuery: {query}")
            results = await vector_store.find_similar(query, n_results=2)
            
            for i, result in enumerate(results, 1):
                print(f"\nResult {i}:")
                print(f"Original Text: {result['text']}")
                print(f"Translated Text: {result['translated_text']}")
                print(f"Language: {result['metadata']['source_language']}")
                print(f"Video: {result['metadata']['video_id']}")
                print(f"Time: {result['metadata']['start_time']:.1f}s - {result['metadata']['end_time']:.1f}s")
                print(f"Similarity: {result['similarity']:.2%}")

if __name__ == "__main__":
    asyncio.run(main())
