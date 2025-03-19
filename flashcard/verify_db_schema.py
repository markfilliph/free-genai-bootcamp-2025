import sqlite3
import os
import json

def get_table_schema(db_path, table_name):
    """Get the schema of a specific table in the database."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    cursor.execute(f"PRAGMA table_info({table_name})")
    columns = cursor.fetchall()
    
    schema = []
    for col in columns:
        schema.append({
            "cid": col[0],
            "name": col[1],
            "type": col[2],
            "notnull": col[3],
            "default_value": col[4],
            "pk": col[5]
        })
    
    conn.close()
    return schema

def get_all_tables(db_path):
    """Get a list of all tables in the database."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = [row[0] for row in cursor.fetchall()]
    
    conn.close()
    return tables

def get_foreign_keys(db_path, table_name):
    """Get all foreign keys for a specific table."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    cursor.execute(f"PRAGMA foreign_key_list({table_name})")
    foreign_keys = cursor.fetchall()
    
    fk_info = []
    for fk in foreign_keys:
        fk_info.append({
            "id": fk[0],
            "seq": fk[1],
            "table": fk[2],
            "from": fk[3],
            "to": fk[4],
            "on_update": fk[5],
            "on_delete": fk[6],
            "match": fk[7]
        })
    
    conn.close()
    return fk_info

def get_indexes(db_path, table_name):
    """Get all indexes for a specific table."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    cursor.execute(f"PRAGMA index_list({table_name})")
    indexes = cursor.fetchall()
    
    idx_info = []
    for idx in indexes:
        idx_info.append({
            "seq": idx[0],
            "name": idx[1],
            "unique": idx[2]
        })
    
    conn.close()
    return idx_info

def verify_database_schema(db_path):
    """Verify the schema of the database."""
    if not os.path.exists(db_path):
        print(f"Database file {db_path} does not exist.")
        return
    
    tables = get_all_tables(db_path)
    
    print(f"Database: {db_path}")
    print(f"Tables found: {len(tables)}")
    print("=" * 50)
    
    for table in tables:
        print(f"\nTable: {table}")
        print("-" * 30)
        
        # Get table schema
        schema = get_table_schema(db_path, table)
        print("Columns:")
        for col in schema:
            pk_str = "PRIMARY KEY" if col["pk"] else ""
            null_str = "NOT NULL" if col["notnull"] else "NULL"
            default = f"DEFAULT {col['default_value']}" if col["default_value"] is not None else ""
            print(f"  - {col['name']} ({col['type']}) {null_str} {default} {pk_str}")
        
        # Get foreign keys
        foreign_keys = get_foreign_keys(db_path, table)
        if foreign_keys:
            print("\nForeign Keys:")
            for fk in foreign_keys:
                print(f"  - {fk['from']} -> {fk['table']}({fk['to']})")
        
        # Get indexes
        indexes = get_indexes(db_path, table)
        if indexes:
            print("\nIndexes:")
            for idx in indexes:
                unique_str = "UNIQUE" if idx["unique"] else ""
                print(f"  - {idx['name']} {unique_str}")
    
    print("\nDatabase schema verification completed.")

if __name__ == "__main__":
    # Verify the test models database
    print("\n=== VERIFYING TEST MODELS DATABASE ===\n")
    verify_database_schema("test_models.db")
    
    # Verify the main database
    print("\n=== VERIFYING MAIN DATABASE ===\n")
    verify_database_schema("flashcards.db")
