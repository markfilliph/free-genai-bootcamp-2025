from typing import List, Dict, Optional
from dataclasses import dataclass
from datetime import datetime
import json
import aiofiles
from pathlib import Path


@dataclass
class JLPTQuestion:
    """Represents a single JLPT listening question"""
    introduction: str
    conversation: str
    question: str
    
    def to_dict(self) -> Dict:
        """Convert to dictionary format"""
        return {
            "introduction": self.introduction,
            "conversation": self.conversation,
            "question": self.question
        }
    
    @classmethod
    def from_dict(cls, data: Dict) -> 'JLPTQuestion':
        """Create from dictionary format"""
        return cls(
            introduction=data["introduction"],
            conversation=data["conversation"],
            question=data["question"]
        )


class JLPTTranscriptParser:
    """Parser for extracting JLPT questions from transcripts"""
    
    def __init__(self):
        # Directory for JSON storage
        self.questions_dir = Path("data/questions")
        self.questions_dir.mkdir(parents=True, exist_ok=True)
        
        # Directory for text output
        self.output_dir = Path("data/parsed_output")
        self.output_dir.mkdir(parents=True, exist_ok=True)
    
    def _get_storage_path(self, video_id: str) -> Path:
        """Get the storage file path for parsed questions"""
        return self.questions_dir / f"{video_id}_questions.json"
    
    def parse_transcript(self, transcript_data: List[Dict]) -> List[JLPTQuestion]:
        """Parse transcript entries into JLPT questions
        
        Args:
            transcript_data: List of transcript entries with 'text' and 'start' fields
            
        Returns:
            List of JLPTQuestion objects
        """
        questions = []
        current_question = []
        in_question = False
        
        for entry in transcript_data:
            text = entry['text'].strip()
            
            # Skip empty entries
            if not text:
                continue
            
            # Start of a new question
            if any(marker in text.lower() for marker in ["問題", "question", "もんだい"]):
                if current_question:  # Save previous question if exists
                    questions.append(self._create_question(current_question))
                current_question = [text]
                in_question = True
                continue
            
            # Add text to current question if we're in one
            if in_question:
                current_question.append(text)
        
        # Add final question
        if current_question:
            questions.append(self._create_question(current_question))
        
        return questions
    
    def _create_question(self, section: List[str]) -> JLPTQuestion:
        """Create a JLPTQuestion object from a section of text"""
        # Initialize parts
        intro_texts = []
        conv_texts = []
        question_texts = []
        current_part = intro_texts
        found_intro = False
        
        for text in section:
            # Skip empty text
            if not text.strip():
                continue
                
            # Skip question number but mark start of introduction
            if any(marker in text for marker in ["問題", "question", "もんだい"]):
                found_intro = True
                continue
            
            # Check for conversation marker
            if any(marker in text.lower() for marker in ["会話", "conversation", "かいわ"]):
                current_part = conv_texts
                # Add this instruction to introduction
                if found_intro:
                    intro_texts.append(text)
                continue
            
            # Check for question marker
            if any(marker in text.lower() for marker in ["質問", "question", "what", "why", "how", "where", "しつもん"]):
                current_part = question_texts
                # Include this text as it contains the question
                current_part.append(text)
                continue
            
            # Add text to current part
            current_part.append(text)
        
        # If no clear markers were found, use position-based splitting
        if not conv_texts and not question_texts:
            third = len(section) // 3
            intro_texts = section[:third]
            conv_texts = section[third:2*third]
            question_texts = section[2*third:]
        
        # Clean up the text
        def clean_text(text_list):
            return " ".join(t for t in text_list if t.strip())
        
        # Create sections
        introduction = clean_text(intro_texts)
        conversation = clean_text(conv_texts)
        question = clean_text(question_texts)
        
        return JLPTQuestion(
            introduction=introduction,
            conversation=conversation,
            question=question
        )
    
    async def save_questions(self, video_id: str, questions: List[JLPTQuestion], title: str = ""):
        """Save parsed questions to JSON and text files
        
        Args:
            video_id: YouTube video ID or other unique identifier
            questions: List of parsed questions
            title: Optional title for the transcript (e.g. video title)
        """
        # Save as JSON
        storage_path = self._get_storage_path(video_id)
        data = {
            "video_id": video_id,
            "title": title,
            "parsed_at": datetime.now().isoformat(),
            "questions": [q.to_dict() for q in questions]
        }
        
        async with aiofiles.open(storage_path, 'w', encoding='utf-8') as f:
            await f.write(json.dumps(data, ensure_ascii=False, indent=2))
            
        # Save as formatted text
        txt_path = self.output_dir / f"{video_id}_parsed.txt"
        
        async with aiofiles.open(txt_path, 'w', encoding='utf-8') as f:
            # Write header
            await f.write("JLPT Listening Practice Questions\n")
            if title:
                await f.write(f"Title: {title}\n")
            await f.write(f"Video ID: {video_id}\n")
            await f.write(f"Parsed at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            await f.write("-" * 50 + "\n\n")
            
            # Write each question
            for i, question in enumerate(questions, 1):
                await f.write(f"Question {i}:\n")
                
                # Write introduction if present
                if question.introduction.strip():
                    await f.write("Introduction:\n")
                    await f.write(question.introduction + "\n\n")
                
                # Write conversation
                await f.write("Conversation:\n")
                # Split conversation into individual speaker turns
                conversation = question.conversation
                speaker_turns = []
                current_turn = []
                
                # Split by spaces but preserve Japanese text
                words = conversation.replace("。", "。 ").split(" ")
                
                for word in words:
                    if ": " in word:  # New speaker
                        if current_turn:  # Save previous turn
                            speaker_turns.append(" ".join(current_turn))
                        current_turn = [word]
                    else:
                        current_turn.append(word)
                
                # Add final turn
                if current_turn:
                    speaker_turns.append(" ".join(current_turn))
                
                # Write each speaker's turn on a new line with proper indentation
                await f.write(speaker_turns[0] + "\n")  # First speaker
                for turn in speaker_turns[1:]:
                    await f.write("  " + turn + "\n")  # Indent subsequent speakers
                await f.write("\n")
                
                # Write question
                await f.write("Question:\n")
                await f.write(question.question + "\n")
                await f.write("-" * 50 + "\n\n")
        
        return txt_path
    
    async def load_questions(self, video_id: str) -> Optional[List[JLPTQuestion]]:
        """Load previously parsed questions from storage"""
        storage_path = self._get_storage_path(video_id)
        if not storage_path.exists():
            return None
            
        try:
            async with aiofiles.open(storage_path, 'r', encoding='utf-8') as f:
                data = json.loads(await f.read())
                return [JLPTQuestion.from_dict(q) for q in data["questions"]]
        except Exception as e:
            print(f"Error loading questions: {e}")
            return None
