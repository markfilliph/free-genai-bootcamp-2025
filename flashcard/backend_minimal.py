import sqlite3
import hashlib
import secrets
from datetime import datetime, timedelta
import json
from http.server import HTTPServer, BaseHTTPRequestHandler
import cgi
import urllib.parse
import os

# Database setup
DB_PATH = "flashcards.db"

def get_db_connection():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

def init_db():
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # Create users table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TEXT NOT NULL
    )
    ''')
    
    # Create decks table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS decks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    )
    ''')
    
    # Create flashcards table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS flashcards (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word TEXT NOT NULL,
        example_sentence TEXT NOT NULL,
        translation TEXT NOT NULL,
        conjugation TEXT,
        cultural_note TEXT,
        deck_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        last_reviewed TEXT,
        ease_factor INTEGER DEFAULT 250,
        interval INTEGER DEFAULT 1,
        FOREIGN KEY (deck_id) REFERENCES decks (id)
    )
    ''')
    
    # Create tokens table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS tokens (
        token TEXT PRIMARY KEY,
        user_id INTEGER NOT NULL,
        expires_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    )
    ''')
    
    conn.commit()
    conn.close()

# Initialize database
init_db()

# Helper functions
def hash_password(password):
    return hashlib.sha256(password.encode()).hexdigest()

def verify_password(plain_password, hashed_password):
    return hash_password(plain_password) == hashed_password

def create_access_token(user_id):
    token = secrets.token_hex(32)
    expires_at = (datetime.now() + timedelta(minutes=30)).isoformat()
    
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        "INSERT INTO tokens (token, user_id, expires_at) VALUES (?, ?, ?)",
        (token, user_id, expires_at)
    )
    
    conn.commit()
    conn.close()
    
    return token

def verify_token(token):
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute(
        "SELECT user_id, expires_at FROM tokens WHERE token = ?",
        (token,)
    )
    
    token_data = cursor.fetchone()
    
    if not token_data:
        conn.close()
        return None
    
    expires_at = datetime.fromisoformat(token_data["expires_at"])
    if expires_at < datetime.now():
        conn.close()
        return None
    
    cursor.execute(
        "SELECT id, username, email, created_at FROM users WHERE id = ?",
        (token_data["user_id"],)
    )
    
    user = cursor.fetchone()
    conn.close()
    
    if not user:
        return None
    
    return dict(user)

# HTTP Request Handler
class FlashcardAPIHandler(BaseHTTPRequestHandler):
    def _set_headers(self, status_code=200, content_type="application/json"):
        self.send_response(status_code)
        self.send_header("Content-type", content_type)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        self.end_headers()
    
    def do_OPTIONS(self):
        self._set_headers()
    
    def _get_json_data(self):
        content_length = int(self.headers.get("Content-Length", 0))
        post_data = self.rfile.read(content_length).decode("utf-8")
        return json.loads(post_data)
    
    def _get_form_data(self):
        form = cgi.FieldStorage(
            fp=self.rfile,
            headers=self.headers,
            environ={"REQUEST_METHOD": "POST"}
        )
        return form
    
    def _get_current_user(self):
        auth_header = self.headers.get("Authorization", "")
        if not auth_header.startswith("Bearer "):
            return None
        
        token = auth_header[7:]  # Remove "Bearer " prefix
        return verify_token(token)
    
    def do_GET(self):
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path
        
        # Root endpoint
        if path == "/":
            self._set_headers()
            response = {
                "message": "Welcome to Language Learning Flashcard Generator API",
                "version": "1.0.0"
            }
            self.wfile.write(json.dumps(response).encode())
            return
        
        # Get decks endpoint
        if path == "/api/decks":
            user = self._get_current_user()
            if not user:
                self._set_headers(401)
                self.wfile.write(json.dumps({"detail": "Unauthorized"}).encode())
                return
            
            conn = get_db_connection()
            cursor = conn.cursor()
            
            cursor.execute(
                "SELECT id, name, user_id, created_at FROM decks WHERE user_id = ?",
                (user["id"],)
            )
            
            decks = [dict(deck) for deck in cursor.fetchall()]
            conn.close()
            
            self._set_headers()
            self.wfile.write(json.dumps(decks).encode())
            return
        
        # Get flashcards by deck endpoint
        if path.startswith("/api/decks/") and "/flashcards" in path:
            user = self._get_current_user()
            if not user:
                self._set_headers(401)
                self.wfile.write(json.dumps({"detail": "Unauthorized"}).encode())
                return
            
            # Extract deck_id from path
            deck_id = int(path.split("/")[3])
            
            conn = get_db_connection()
            cursor = conn.cursor()
            
            # Check if deck exists and belongs to user
            cursor.execute(
                "SELECT id FROM decks WHERE id = ? AND user_id = ?",
                (deck_id, user["id"])
            )
            
            if not cursor.fetchone():
                conn.close()
                self._set_headers(404)
                self.wfile.write(json.dumps({"detail": "Deck not found or you don't have permission"}).encode())
                return
            
            # Get flashcards
            cursor.execute(
                """
                SELECT id, word, example_sentence, translation, conjugation, cultural_note, 
                deck_id, created_at, last_reviewed, ease_factor, interval 
                FROM flashcards WHERE deck_id = ?
                """,
                (deck_id,)
            )
            
            flashcards = [dict(fc) for fc in cursor.fetchall()]
            conn.close()
            
            self._set_headers()
            self.wfile.write(json.dumps(flashcards).encode())
            return
        
        # Handle 404 for unknown endpoints
        self._set_headers(404)
        self.wfile.write(json.dumps({"detail": "Not found"}).encode())
    
    def do_POST(self):
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path
        
        # Register endpoint
        if path == "/api/auth/register":
            try:
                data = self._get_json_data()
                username = data.get("username")
                email = data.get("email")
                password = data.get("password")
                
                if not all([username, email, password]):
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Missing required fields"}).encode())
                    return
                
                conn = get_db_connection()
                cursor = conn.cursor()
                
                # Check if username or email already exists
                cursor.execute(
                    "SELECT id FROM users WHERE username = ? OR email = ?",
                    (username, email)
                )
                
                if cursor.fetchone():
                    conn.close()
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Username or email already registered"}).encode())
                    return
                
                # Hash password
                hashed_password = hash_password(password)
                
                # Insert new user
                cursor.execute(
                    "INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
                    (username, email, hashed_password, datetime.now().isoformat())
                )
                
                user_id = cursor.lastrowid
                conn.commit()
                conn.close()
                
                self._set_headers()
                self.wfile.write(json.dumps({"message": "User registered successfully", "user_id": user_id}).encode())
                return
            except Exception as e:
                self._set_headers(500)
                self.wfile.write(json.dumps({"detail": str(e)}).encode())
                return
        
        # Login endpoint
        if path == "/api/auth/login":
            try:
                form = self._get_form_data()
                username = form.getvalue("username")
                password = form.getvalue("password")
                
                if not all([username, password]):
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Missing username or password"}).encode())
                    return
                
                conn = get_db_connection()
                cursor = conn.cursor()
                
                # Find user by username
                cursor.execute(
                    "SELECT id, password_hash FROM users WHERE username = ?",
                    (username,)
                )
                
                user = cursor.fetchone()
                conn.close()
                
                if not user or not verify_password(password, user["password_hash"]):
                    self._set_headers(401)
                    self.wfile.write(json.dumps({"detail": "Incorrect username or password"}).encode())
                    return
                
                # Create access token
                access_token = create_access_token(user["id"])
                
                self._set_headers()
                self.wfile.write(json.dumps({"access_token": access_token, "token_type": "bearer"}).encode())
                return
            except Exception as e:
                self._set_headers(500)
                self.wfile.write(json.dumps({"detail": str(e)}).encode())
                return
        
        # Create deck endpoint
        if path == "/api/decks":
            user = self._get_current_user()
            if not user:
                self._set_headers(401)
                self.wfile.write(json.dumps({"detail": "Unauthorized"}).encode())
                return
            
            try:
                data = self._get_json_data()
                name = data.get("name")
                
                if not name:
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Missing deck name"}).encode())
                    return
                
                conn = get_db_connection()
                cursor = conn.cursor()
                
                # Insert new deck
                cursor.execute(
                    "INSERT INTO decks (name, user_id, created_at) VALUES (?, ?, ?)",
                    (name, user["id"], datetime.now().isoformat())
                )
                
                deck_id = cursor.lastrowid
                conn.commit()
                
                # Get the created deck
                cursor.execute(
                    "SELECT id, name, user_id, created_at FROM decks WHERE id = ?",
                    (deck_id,)
                )
                
                new_deck = cursor.fetchone()
                conn.close()
                
                self._set_headers()
                self.wfile.write(json.dumps(dict(new_deck)).encode())
                return
            except Exception as e:
                self._set_headers(500)
                self.wfile.write(json.dumps({"detail": str(e)}).encode())
                return
        
        # Create flashcard endpoint
        if path == "/api/flashcards":
            user = self._get_current_user()
            if not user:
                self._set_headers(401)
                self.wfile.write(json.dumps({"detail": "Unauthorized"}).encode())
                return
            
            try:
                data = self._get_json_data()
                word = data.get("word")
                example_sentence = data.get("example_sentence")
                translation = data.get("translation")
                conjugation = data.get("conjugation")
                cultural_note = data.get("cultural_note")
                deck_id = data.get("deck_id")
                
                if not all([word, example_sentence, translation, deck_id]):
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Missing required fields"}).encode())
                    return
                
                conn = get_db_connection()
                cursor = conn.cursor()
                
                # Check if deck exists and belongs to user
                cursor.execute(
                    "SELECT id FROM decks WHERE id = ? AND user_id = ?",
                    (deck_id, user["id"])
                )
                
                if not cursor.fetchone():
                    conn.close()
                    self._set_headers(404)
                    self.wfile.write(json.dumps({"detail": "Deck not found or you don't have permission"}).encode())
                    return
                
                # Insert new flashcard
                cursor.execute(
                    """
                    INSERT INTO flashcards 
                    (word, example_sentence, translation, conjugation, cultural_note, deck_id, created_at) 
                    VALUES (?, ?, ?, ?, ?, ?, ?)
                    """,
                    (
                        word, 
                        example_sentence, 
                        translation, 
                        conjugation, 
                        cultural_note, 
                        deck_id, 
                        datetime.now().isoformat()
                    )
                )
                
                flashcard_id = cursor.lastrowid
                conn.commit()
                
                # Get the created flashcard
                cursor.execute(
                    """
                    SELECT id, word, example_sentence, translation, conjugation, cultural_note, 
                    deck_id, created_at, last_reviewed, ease_factor, interval 
                    FROM flashcards WHERE id = ?
                    """,
                    (flashcard_id,)
                )
                
                new_flashcard = cursor.fetchone()
                conn.close()
                
                self._set_headers()
                self.wfile.write(json.dumps(dict(new_flashcard)).encode())
                return
            except Exception as e:
                self._set_headers(500)
                self.wfile.write(json.dumps({"detail": str(e)}).encode())
                return
        
        # Generate content endpoint
        if path == "/api/generate":
            user = self._get_current_user()
            if not user:
                self._set_headers(401)
                self.wfile.write(json.dumps({"detail": "Unauthorized"}).encode())
                return
            
            try:
                data = self._get_json_data()
                word = data.get("word")
                is_verb = data.get("is_verb", False)
                
                if not word:
                    self._set_headers(400)
                    self.wfile.write(json.dumps({"detail": "Missing word"}).encode())
                    return
                
                # Simulate LLM generation
                example_sentences = [
                    f"Ejemplo con '{word}': Esta es una oración de ejemplo.",
                    f"Otro ejemplo con '{word}': Segunda oración de ejemplo."
                ]
                
                conjugations = None
                if is_verb:
                    conjugations = f"Conjugaciones para '{word}':\nPresente: yo {word}o, tú {word}es..."
                
                cultural_note = f"Nota cultural sobre '{word}': Este término es comúnmente usado en España y Latinoamérica."
                
                response = {
                    "example_sentences": example_sentences,
                    "conjugations": conjugations,
                    "cultural_note": cultural_note
                }
                
                self._set_headers()
                self.wfile.write(json.dumps(response).encode())
                return
            except Exception as e:
                self._set_headers(500)
                self.wfile.write(json.dumps({"detail": str(e)}).encode())
                return
        
        # Handle 404 for unknown endpoints
        self._set_headers(404)
        self.wfile.write(json.dumps({"detail": "Not found"}).encode())

def run_server(port=8000):
    server_address = ("", port)
    httpd = HTTPServer(server_address, FlashcardAPIHandler)
    print(f"Starting server on port {port}...")
    httpd.serve_forever()

if __name__ == "__main__":
    run_server()
