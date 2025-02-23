import asyncio
from get_transcript import YouTubeTranscriptDownloader
from structured_data import JLPTTranscriptParser
import json
from pathlib import Path


async def test_transcript_parsing(video_url: str):
    """Test parsing JLPT questions from a YouTube video transcript"""
    
    print(f"\nTesting with video URL: {video_url}")
    print("-" * 50)
    
    # Initialize components
    downloader = YouTubeTranscriptDownloader()
    parser = JLPTTranscriptParser()
    
    try:
        # Step 1: Get transcript
        print("\n1. Fetching transcript...")
        transcript_data = await downloader.get_transcript(video_url)
        if not transcript_data:
            print("❌ Error: No transcript found")
            return
        
        video_id = transcript_data['video_id']
        print(f"✓ Found transcript for video {video_id}")
        print(f"✓ Total entries: {len(transcript_data['transcript'])}")
        
        # Debug: Show transcript content
        print("\nTranscript content:")
        print("-" * 50)
        for entry in transcript_data['transcript'][:10]:  # Show first 10 entries
            print(f"[{entry.get('start', 0):.1f}s] {entry.get('text', '')}")
        
        # Step 2: Parse questions
        print("\n2. Parsing questions...")
        questions = parser.parse_transcript(transcript_data['transcript'])
        print(f"✓ Found {len(questions)} questions")
        
        # Step 3: Save questions
        print("\n3. Saving questions...")
        await parser.save_questions(video_id, questions)
        print(f"✓ Saved questions to {parser._get_storage_path(video_id)}")
        
        # Step 4: Load questions
        print("\n4. Loading saved questions...")
        loaded_questions = await parser.load_questions(video_id)
        if loaded_questions:
            print(f"✓ Successfully loaded {len(loaded_questions)} questions")
        
        # Step 5: Display parsed questions
        print("\n5. Question details:")
        print("-" * 50)
        for i, question in enumerate(questions, 1):
            print(f"\nQuestion {i}:")
            print("Introduction:")
            print(question.introduction)
            print("\nConversation:")
            print(question.conversation)
            print("\nQuestion:")
            print(question.question)
            print("-" * 50)
        
        # Step 6: Save pretty output
        output_dir = Path("data/parsed_output")
        output_dir.mkdir(parents=True, exist_ok=True)
        output_path = output_dir / f"{video_id}_parsed.txt"
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(f"JLPT Listening Practice Questions\n")
            f.write(f"Video ID: {video_id}\n")
            f.write("-" * 50 + "\n\n")
            
            for i, question in enumerate(questions, 1):
                f.write(f"Question {i}:\n")
                f.write("Introduction:\n")
                f.write(question.introduction + "\n\n")
                f.write("Conversation:\n")
                f.write(question.conversation + "\n\n")
                f.write("Question:\n")
                f.write(question.question + "\n")
                f.write("-" * 50 + "\n\n")
        
        print(f"\n✓ Saved formatted output to {output_path}")
        
    except Exception as e:
        print(f"❌ Error: {str(e)}")


async def test_with_sample_data():
    """Test parsing with sample JLPT transcript data"""
    parser = JLPTTranscriptParser()
    
    # Load sample data
    with open('sample_transcript.json', 'r', encoding='utf-8') as f:
        transcript_data = json.load(f)
    
    print("\nTesting with sample JLPT transcript data")
    print("-" * 50)
    
    # Parse questions
    print("\n1. Parsing questions...")
    questions = parser.parse_transcript(transcript_data['transcript'])
    print(f"✓ Found {len(questions)} questions")
    
    # Save questions
    print("\n2. Saving questions...")
    txt_path = await parser.save_questions(
        video_id=transcript_data['video_id'],
        questions=questions,
        title="Sample JLPT N4 Listening Practice"
    )
    print(f"✓ Saved questions to JSON: {parser._get_storage_path(transcript_data['video_id'])}")
    print(f"✓ Saved formatted transcript to: {txt_path}")
    
    # Load and verify questions
    print("\n3. Loading saved questions...")
    loaded_questions = await parser.load_questions(transcript_data['video_id'])
    if loaded_questions:
        print(f"✓ Successfully loaded {len(loaded_questions)} questions")
        
        # Verify content matches
        all_match = all(
            q1.to_dict() == q2.to_dict()
            for q1, q2 in zip(questions, loaded_questions)
        )
        print(f"✓ Loaded content {'matches' if all_match else 'differs from'} original")

if __name__ == "__main__":
    # Test with sample JLPT transcript data
    asyncio.run(test_with_sample_data())
