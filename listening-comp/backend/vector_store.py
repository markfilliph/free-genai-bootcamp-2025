import chromadb
from chromadb.utils import embedding_functions
import json
import os
from typing import Dict, List, Optional
from sentence_transformers import SentenceTransformer
from langdetect import detect
from googletrans import Translator

class SentenceTransformerEmbedding(embedding_functions.EmbeddingFunction):
    def __init__(self, model_name="paraphrase-multilingual-mpnet-base-v2"):
        """Initialize sentence-transformers embedding function"""
        self.model_name = model_name
        self._model = None
        
    @property
    def model(self):
        """Lazy load the model"""
        if self._model is None:
            self._model = SentenceTransformer(self.model_name)
        return self._model

    def __call__(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for a list of texts using sentence-transformers"""
        try:
            embeddings = self.model.encode(texts, convert_to_tensor=False)
            return embeddings.tolist()
        except Exception as e:
            print(f"Error generating embeddings: {str(e)}")
            # Return zero vectors as fallback
            return [[0.0] * 384] * len(texts)  # MiniLM embeddings are 384-dimensional

class QuestionVectorStore:
    def __init__(self, persist_directory: str = "backend/data/vectorstore"):
        """Initialize the vector store for JLPT listening questions"""
        self.persist_directory = persist_directory
        
        # Initialize ChromaDB client
        self.client = chromadb.PersistentClient(path=persist_directory)
        
        # Use sentence-transformers embedding model
        self.embedding_fn = SentenceTransformerEmbedding()
        
        # Reset collections to ensure correct embedding dimensions
        self._reset_collections()
        
    def _reset_collections(self):
        """Reset collections to ensure correct embedding dimensions"""
        # Delete existing collections if they exist
        try:
            self.client.delete_collection("section2_questions")
            self.client.delete_collection("section3_questions")
        except Exception:
            pass

        # Create collections with embedding function
        self.collections = {
            "section2": self.client.create_collection(
                name="section2_questions",
                embedding_function=self.embedding_fn,
                metadata={"description": "JLPT listening comprehension questions - Section 2"}
            ),
            "section3": self.client.create_collection(
                name="section3_questions",
                embedding_function=self.embedding_fn,
                metadata={"description": "JLPT listening comprehension questions - Section 3"}
            )
        }

        # Load and add example questions
        self._load_example_questions()

    def add_question(self, section_num: int, question: Dict, topic: str) -> None:
        """Add a question to the vector store
        
        Args:
            section_num: Section number (2 or 3)
            question: Question dictionary containing all fields
            topic: Topic category for the question
        """
        collection = self.collections[f"section{section_num}"]
        
        # Create document text based on section type
        if section_num == 2:
            doc_text = f"Introduction: {question['Introduction']}\n"
            doc_text += f"Conversation: {question['Conversation']}\n"
            doc_text += f"Question: {question['Question']}\n"
            doc_text += f"Options: {', '.join(question['Options'])}"
        else:
            doc_text = f"Situation: {question['Situation']}\n"
            doc_text += f"Content: {question.get('Content', '')}\n"
            doc_text += f"Question: {question['Question']}\n"
            doc_text += f"Options: {', '.join(question['Options'])}"
        
        # Add to collection
        collection.add(
            documents=[doc_text],
            metadatas=[{
                "topic": topic,
                "section": section_num,
                **{k: str(v) for k, v in question.items()}
            }],
            ids=[f"q_{section_num}_{len(collection.get()['ids']) + 1}"]
        )

    def add_questions(self, section_num: int, questions: List[Dict], video_id: str):
        """Add questions to the vector store"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")
            
        collection = self.collections[f"section{section_num}"]
        
        ids = []
        documents = []
        metadatas = []
        
        for idx, question in enumerate(questions):
            # Create a unique ID for each question
            question_id = f"{video_id}_{section_num}_{idx}"
            ids.append(question_id)
            
            # Store the full question structure as metadata
            metadatas.append({
                "video_id": video_id,
                "section": section_num,
                "question_index": idx,
                "full_structure": json.dumps(question)
            })
            
            # Create a searchable document from the question content
            if section_num == 2:
                document = f"""
                Situation: {question['Introduction']}
                Dialogue: {question['Conversation']}
                Question: {question['Question']}
                """
            else:  # section 3
                document = f"""
                Situation: {question['Situation']}
                Question: {question['Question']}
                """
            documents.append(document)
        
        # Add to collection
        collection.add(
            ids=ids,
            documents=documents,
            metadatas=metadatas
        )

    def search_similar_questions(
        self, 
        section_num: int, 
        query: str, 
        n_results: int = 5
    ) -> List[Dict]:
        """Search for similar questions in the vector store"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")
            
        collection = self.collections[f"section{section_num}"]
        
        results = collection.query(
            query_texts=[query],
            n_results=n_results
        )
        
        # Convert results to more usable format
        questions = []
        for idx, metadata in enumerate(results['metadatas'][0]):
            # Reconstruct question from metadata fields
            if section_num == 2:
                question_data = {
                    "Introduction": metadata.get("Introduction", ""),
                    "Conversation": metadata.get("Conversation", ""),
                    "Question": metadata.get("Question", ""),
                    "Options": metadata.get("Options", "").split(", ") if metadata.get("Options") else [],
                    "CorrectAnswer": int(metadata.get("CorrectAnswer", "1")),
                    "Explanation": metadata.get("Explanation", "")
                }
            else:
                question_data = {
                    "Situation": metadata.get("Situation", ""),
                    "Content": metadata.get("Content", ""),
                    "Question": metadata.get("Question", ""),
                    "Options": metadata.get("Options", "").split(", ") if metadata.get("Options") else [],
                    "CorrectAnswer": int(metadata.get("CorrectAnswer", "1")),
                    "Explanation": metadata.get("Explanation", "")
                }
            question_data['similarity_score'] = results['distances'][0][idx]
            questions.append(question_data)
            
        return questions

    def _load_example_questions(self):
        """Load example questions from JSON files"""
        example_dir = os.path.join(self.persist_directory, "examples")
        if not os.path.exists(example_dir):
            return
            
        for filename in os.listdir(example_dir):
            if not filename.endswith(".json"):
                continue
                
            section_num = int(filename[7])  # e.g. section2_examples.json
            with open(os.path.join(example_dir, filename)) as f:
                questions = json.load(f)
                
            video_id = "example"
            self.add_questions(section_num, questions, video_id)
            print(f"Indexed {len(questions)} questions from {filename}")

class TranscriptVectorStore:
    def __init__(self, persist_directory: str = "backend/data/vectorstore"):
        """Initialize the vector store for YouTube transcripts"""
        self.persist_directory = persist_directory
        
        # Initialize ChromaDB client
        self.client = chromadb.PersistentClient(path=persist_directory)
        
        # Use multilingual sentence-transformers model
        self.embedding_fn = SentenceTransformerEmbedding()
        
        # Initialize translator
        self.translator = Translator()
        
        # Get or create collection
        try:
            self.collection = self.client.get_collection(
                name="youtube_transcripts",
                embedding_function=self.embedding_fn
            )
        except Exception:
            self.collection = self.client.create_collection(
                name="youtube_transcripts",
                embedding_function=self.embedding_fn,
                metadata={"description": "YouTube video transcripts for listening practice"}
            )
    
    async def _detect_and_translate(self, text: str, target_lang='en') -> tuple[str, str, str]:
        """Detect language and translate text if needed
        
        Args:
            text: Text to process
            target_lang: Target language code (default: 'en')
            
        Returns:
            Tuple of (original_text, translated_text, source_language)
        """
        try:
            source_lang = detect(text)
            if source_lang != target_lang:
                translated = await self.translator.translate(text, dest=target_lang)
                return text, translated.text, source_lang
            return text, text, source_lang
        except Exception as e:
            print(f"Translation error: {e}")
            return text, text, 'unknown'

    async def _chunk_transcript(self, transcript: List[Dict], chunk_size: int = 5, overlap: int = 2) -> List[Dict]:
        """Split transcript into overlapping chunks for better context
        
        Args:
            transcript: List of transcript entries with 'text' and 'start' fields
            chunk_size: Number of entries per chunk
            overlap: Number of overlapping entries between chunks
            
        Returns:
            List of chunks with combined text and metadata
        """
        chunks = []
        for i in range(0, len(transcript), chunk_size - overlap):
            chunk_entries = transcript[i:i + chunk_size]
            if len(chunk_entries) < 2:  # Skip very small chunks
                continue
                
            # Combine text and collect metadata
            text = " ".join(entry['text'] for entry in chunk_entries)
            start_time = chunk_entries[0]['start']
            end_time = chunk_entries[-1]['start'] + chunk_entries[-1].get('duration', 0)
            
            # Detect language and translate
            original_text, translated_text, source_lang = await self._detect_and_translate(text)
            
            chunks.append({
                'original_text': original_text,
                'translated_text': translated_text,
                'source_language': source_lang,
                'start_time': start_time,
                'end_time': end_time,
                'entries': chunk_entries
            })
        
        return chunks

    async def add_transcript(self, video_id: str, transcript_data: Dict) -> None:
        """Add a transcript to the vector store
        
        Args:
            video_id: YouTube video ID
            transcript_data: Transcript data from YouTubeTranscriptDownloader
        """
        # Get transcript entries
        entries = transcript_data['transcript']
        
        # Split into chunks with translation
        chunks = await self._chunk_transcript(entries)
        
        # Prepare data for ChromaDB
        ids = [f"{video_id}_{i}" for i in range(len(chunks))]
        texts = [chunk['translated_text'] for chunk in chunks]  # Use translated text for matching
        metadatas = [{
            'video_id': video_id,
            'start_time': chunk['start_time'],
            'end_time': chunk['end_time'],
            'original_text': chunk['original_text'],
            'source_language': chunk['source_language']
        } for chunk in chunks]
        
        # Add to collection
        self.collection.add(
            ids=ids,
            documents=texts,
            metadatas=metadatas
        )

    async def find_similar(self, query: str, video_id: str = None, n_results: int = 5) -> List[Dict]:
        """Find transcript chunks similar to the query
        
        Args:
            query: Search query
            video_id: Optional video ID to filter results
            n_results: Number of results to return
            
        Returns:
            List of similar chunks with metadata
        """
        # Translate query if needed
        _, translated_query, query_lang = await self._detect_and_translate(query)
        
        # Prepare where clause if video_id is provided
        where = {"video_id": video_id} if video_id else None
        
        # Search using translated query
        results = self.collection.query(
            query_texts=[translated_query],
            where=where,
            n_results=n_results,
            include=["documents", "metadatas", "distances"]
        )
        
        # Format results
        formatted_results = []
        if len(results['ids'][0]) > 0:
            # Find min and max distances for normalization
            distances = results['distances'][0]
            min_dist = min(distances)
            max_dist = max(distances)
            dist_range = max_dist - min_dist if max_dist > min_dist else 1
            
            for i in range(len(results['ids'][0])):
                metadata = results['metadatas'][0][i]
                # Normalize distance to [0, 1] range and invert (closer = higher score)
                similarity = 1 - ((distances[i] - min_dist) / dist_range)
                
                formatted_results.append({
                    'text': metadata['original_text'],  # Return original text
                    'translated_text': results['documents'][0][i],  # Also include translation
                    'metadata': {
                        'video_id': metadata['video_id'],
                        'start_time': metadata['start_time'],
                        'end_time': metadata['end_time'],
                        'source_language': metadata['source_language']
                    },
                    'similarity': similarity  # Normalized similarity score
                })
        
        return formatted_results
