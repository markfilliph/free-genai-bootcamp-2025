import sqlite3
from contextlib import contextmanager
from pathlib import Path
import logging
from typing import List, Dict, Optional
import json
from datetime import datetime
import shutil
from queue import Queue
from threading import Lock

# Configure logging
logger = logging.getLogger(__name__)

class ConnectionPool:
    def __init__(self, db_path: str, max_connections: int = 5):
        self.db_path = db_path
        self.max_connections = max_connections
        self.connections = Queue(maxsize=max_connections)
        self.lock = Lock()
        
        # Initialize connection pool
        for _ in range(max_connections):
            conn = sqlite3.connect(db_path, check_same_thread=False)
            conn.row_factory = sqlite3.Row
            self.connections.put(conn)
    
    @contextmanager
    def get_connection(self):
        conn = self.connections.get()
        try:
            yield conn
        finally:
            self.connections.put(conn)
    
    def close_all(self):
        while not self.connections.empty():
            conn = self.connections.get()
            conn.close()

class Database:
    def __init__(self, db_path: str = None, max_connections: int = 5):
        if db_path is None:
            # Use the database directory in the project root
            db_dir = Path(__file__).parent.parent / 'data'
            db_dir.mkdir(parents=True, exist_ok=True)
            self.db_path = str(db_dir / 'song_vocab.db')
        else:
            self.db_path = db_path
        logger.info(f"Using database at: {self.db_path}")
        
        # Initialize database
        self._init_db()
        
        # Run migrations
        from .migrations import run_migrations
        run_migrations(self.db_path)
        
        # Initialize connection pool
        self.pool = ConnectionPool(self.db_path, max_connections)
    
    def backup_database(self, backup_dir: Optional[Path] = None) -> str:
        """Create a backup of the database."""
        if backup_dir is None:
            backup_dir = Path(self.db_path).parent / 'backups'
        backup_dir.mkdir(parents=True, exist_ok=True)
        
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        backup_path = backup_dir / f'song_vocab_{timestamp}.db'
        
        try:
            # Create a new connection for backup to avoid conflicts
            with sqlite3.connect(self.db_path) as source:
                with sqlite3.connect(str(backup_path)) as dest:
                    source.backup(dest)
            logger.info(f"Database backed up to {backup_path}")
            return str(backup_path)
        except Exception as e:
            logger.error(f"Backup failed: {e}")
            raise
    
    def _init_db(self):
        """Initialize the database with required tables."""
        logger.info(f"Initializing database at {self.db_path}")
        
        with self.get_connection() as conn:
            cursor = conn.cursor()
            
            # Create songs table
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS songs (
                    id TEXT PRIMARY KEY,
                    title TEXT NOT NULL,
                    artist TEXT,
                    lyrics TEXT NOT NULL,
                    romaji_lyrics TEXT,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            
            # Create vocabulary table
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS vocabulary (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    song_id TEXT NOT NULL,
                    kanji TEXT NOT NULL,
                    romaji TEXT NOT NULL,
                    english TEXT NOT NULL,
                    parts JSON NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (song_id) REFERENCES songs (id),
                    UNIQUE (song_id, kanji)
                )
            """)
            
            conn.commit()
            logger.info("Database initialized successfully")
    
    @contextmanager
    def get_connection(self):
        """Context manager for database connections."""
        conn = None
        try:
            conn = sqlite3.connect(self.db_path)
            conn.row_factory = sqlite3.Row
            yield conn
        finally:
            if conn:
                conn.close()
    
    def save_song(self, song_id: str, title: str, lyrics: str, artist: Optional[str] = None, romaji_lyrics: Optional[str] = None) -> bool:
        """Save a song to the database."""
        logger.info(f"Saving song: {title} ({song_id})")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("""
                    INSERT OR REPLACE INTO songs (id, title, artist, lyrics, romaji_lyrics)
                    VALUES (?, ?, ?, ?, ?)
                """, (song_id, title, artist, lyrics, romaji_lyrics))
                conn.commit()
                logger.info(f"Song saved successfully: {song_id}")
                return True
        except Exception as e:
            logger.error(f"Error saving song: {e}")
            return False
    
    def save_vocabulary(self, song_id: str, vocabulary_items: List[Dict]) -> bool:
        """Save vocabulary items for a song."""
        logger.info(f"Saving {len(vocabulary_items)} vocabulary items for song {song_id}")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                for item in vocabulary_items:
                    cursor.execute("""
                        INSERT OR REPLACE INTO vocabulary 
                        (song_id, kanji, romaji, english, parts)
                        VALUES (?, ?, ?, ?, ?)
                    """, (
                        song_id,
                        item['kanji'],
                        item['romaji'],
                        item['english'],
                        json.dumps(item['parts'])
                    ))
                conn.commit()
                logger.info(f"Vocabulary items saved successfully for song {song_id}")
                return True
        except Exception as e:
            logger.error(f"Error saving vocabulary: {e}")
            return False
    
    def get_song(self, song_id: str) -> Optional[Dict]:
        """Retrieve a song by its ID."""
        logger.info(f"Retrieving song: {song_id}")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("SELECT * FROM songs WHERE id = ?", (song_id,))
                row = cursor.fetchone()
                return dict(row) if row else None
        except Exception as e:
            logger.error(f"Error retrieving song: {e}")
            return None
    
    def get_vocabulary(self, song_id: str) -> List[Dict]:
        """Retrieve vocabulary items for a song."""
        logger.info(f"Retrieving vocabulary for song: {song_id}")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("SELECT * FROM vocabulary WHERE song_id = ?", (song_id,))
                rows = cursor.fetchall()
                vocabulary = []
                for row in rows:
                    item = dict(row)
                    item['parts'] = json.loads(item['parts'])
                    vocabulary.append(item)
                return vocabulary
        except Exception as e:
            logger.error(f"Error retrieving vocabulary: {e}")
            return []
    
    def search_songs(self, query: str) -> List[Dict]:
        """Search songs by title or artist."""
        logger.info(f"Searching songs with query: {query}")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("""
                    SELECT * FROM songs 
                    WHERE title LIKE ? OR artist LIKE ?
                """, (f"%{query}%", f"%{query}%"))
                return [dict(row) for row in cursor.fetchall()]
        except Exception as e:
            logger.error(f"Error searching songs: {e}")
            return []
    
    def search_vocabulary(self, query: str) -> List[Dict]:
        """Search vocabulary items by kanji, romaji, or English."""
        logger.info(f"Searching vocabulary with query: {query}")
        try:
            with self.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("""
                    SELECT v.*, s.title, s.artist 
                    FROM vocabulary v
                    JOIN songs s ON v.song_id = s.id
                    WHERE v.kanji LIKE ? 
                    OR v.romaji LIKE ? 
                    OR v.english LIKE ?
                """, (f"%{query}%", f"%{query}%", f"%{query}%"))
                vocabulary = []
                for row in cursor.fetchall():
                    item = dict(row)
                    item['parts'] = json.loads(item['parts'])
                    vocabulary.append(item)
                return vocabulary
        except Exception as e:
            logger.error(f"Error searching vocabulary: {e}")
            return []
