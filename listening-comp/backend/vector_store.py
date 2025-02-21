import chromadb
from chromadb.utils import embedding_functions
import json
import os
from typing import Dict, List, Optional
from sentence_transformers import SentenceTransformer

class SentenceTransformerEmbedding(embedding_functions.EmbeddingFunction):
    def __init__(self, model_name="all-MiniLM-L6-v2"):
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
        self.section2_collection = self.client.create_collection(
            name="section2_questions",
            embedding_function=self.embedding_fn
        )

        self.section3_collection = self.client.create_collection(
            name="section3_questions",
            embedding_function=self.embedding_fn
        )

        # Load and add example questions
        self._load_example_questions()

        # Create new collections
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

    def get_question_by_id(self, section_num: int, question_id: str) -> Optional[Dict]:
        """Retrieve a specific question by its ID"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")
            
        collection = self.collections[f"section{section_num}"]
        
        result = collection.get(
            ids=[question_id],
            include=['metadatas']
        )
        
        if result['metadatas']:
            return json.loads(result['metadatas'][0]['full_structure'])
        return None

    def parse_questions_from_file(self, filename: str) -> List[Dict]:
        """Parse questions from a structured text file"""
        questions = []
        current_question = {}
        
        try:
            with open(filename, 'r', encoding='utf-8') as f:
                lines = f.readlines()
                
            i = 0
            while i < len(lines):
                line = lines[i].strip()
                
                if line.startswith('<question>'):
                    current_question = {}
                elif line.startswith('Introduction:'):
                    i += 1
                    if i < len(lines):
                        current_question['Introduction'] = lines[i].strip()
                elif line.startswith('Conversation:'):
                    i += 1
                    if i < len(lines):
                        current_question['Conversation'] = lines[i].strip()
                elif line.startswith('Situation:'):
                    i += 1
                    if i < len(lines):
                        current_question['Situation'] = lines[i].strip()
                elif line.startswith('Question:'):
                    i += 1
                    if i < len(lines):
                        current_question['Question'] = lines[i].strip()
                elif line.startswith('Options:'):
                    options = []
                    for _ in range(4):
                        i += 1
                        if i < len(lines):
                            option = lines[i].strip()
                            if option.startswith('1.') or option.startswith('2.') or option.startswith('3.') or option.startswith('4.'):
                                options.append(option[2:].strip())
                    current_question['Options'] = options
                elif line.startswith('</question>'):
                    if current_question:
                        questions.append(current_question)
                        current_question = {}
                i += 1
            return questions
        except Exception as e:
            print(f"Error parsing questions from {filename}: {str(e)}")
            return []

    def index_questions_file(self, filename: str, section_num: int):
        """Index all questions from a file into the vector store"""
        # Extract video ID from filename
        video_id = os.path.basename(filename).split('_section')[0]
        
        # Parse questions from file
        questions = self.parse_questions_from_file(filename)
        
        # Add to vector store
        if questions:
            self.add_questions(section_num, questions, video_id)
            print(f"Indexed {len(questions)} questions from {filename}")

if __name__ == "__main__":
    # Example usage
    store = QuestionVectorStore()
    
    # Index questions from files
    question_files = [
        ("backend/data/questions/sY7L5cfCWno_section2.txt", 2),
        ("backend/data/questions/sY7L5cfCWno_section3.txt", 3)
    ]
    
    for filename, section_num in question_files:
        if os.path.exists(filename):
            store.index_questions_file(filename, section_num)
    
    # Search for similar questions
    similar = store.search_similar_questions(2, "誕生日について質問", n_results=1)
