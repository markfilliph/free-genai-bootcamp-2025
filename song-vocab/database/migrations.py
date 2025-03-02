import sqlite3
import logging
from pathlib import Path
from datetime import datetime

logger = logging.getLogger(__name__)

class Migration:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self._init_migrations_table()
    
    def _init_migrations_table(self):
        """Initialize the migrations table if it doesn't exist."""
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS migrations (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    version TEXT NOT NULL,
                    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            conn.commit()
    
    def get_applied_migrations(self) -> set:
        """Get list of already applied migrations."""
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT version FROM migrations")
            return {row[0] for row in cursor.fetchall()}
    
    def apply_migration(self, version: str, sql: str):
        """Apply a single migration if not already applied."""
        try:
            with sqlite3.connect(self.db_path) as conn:
                cursor = conn.cursor()
                # Check if migration was already applied
                cursor.execute("SELECT 1 FROM migrations WHERE version = ?", (version,))
                if cursor.fetchone():
                    logger.info(f"Migration {version} already applied")
                    return
                
                # Apply migration
                logger.info(f"Applying migration {version}")
                cursor.executescript(sql)
                
                # Record migration
                cursor.execute(
                    "INSERT INTO migrations (version) VALUES (?)",
                    (version,)
                )
                conn.commit()
                logger.info(f"Successfully applied migration {version}")
        except Exception as e:
            logger.error(f"Failed to apply migration {version}: {e}")
            raise

def get_migrations() -> list:
    """Get list of available migrations."""
    migrations_dir = Path(__file__).parent / 'migrations'
    migrations_dir.mkdir(exist_ok=True)
    
    migrations = []
    for file in sorted(migrations_dir.glob('*.sql')):
        version = file.stem
        with open(file, 'r') as f:
            sql = f.read()
        migrations.append((version, sql))
    
    return migrations

def run_migrations(db_path: str):
    """Run all pending migrations."""
    logger.info("Starting database migrations")
    migration = Migration(db_path)
    
    for version, sql in get_migrations():
        try:
            migration.apply_migration(version, sql)
        except Exception as e:
            logger.error(f"Migration failed: {e}")
            raise
    
    logger.info("Completed database migrations")
