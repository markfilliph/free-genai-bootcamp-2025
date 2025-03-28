#!/usr/bin/env python3
"""
Unified server launcher for Language Learning Flashcard Generator

This script provides a single entry point to launch any of the backend implementations:
- Main FastAPI backend (default)
- Simple FastAPI backend
- Minimal backend

Usage:
    python run_server.py [--backend <backend_type>] [--port <port>]

Options:
    --backend    Type of backend to run: 'main' (default), 'simple', or 'minimal'
    --port       Port to run the server on (default: 8000)
"""

import argparse
import os
import sys
import importlib.util
import subprocess

def run_main_backend(port):
    """Run the main FastAPI backend"""
    print(f"Starting main FastAPI backend on port {port}...")
    subprocess.run([
        sys.executable, "-m", "uvicorn", 
        "backend.main:app", 
        "--reload", 
        f"--port={port}"
    ])

def run_simple_backend(port):
    """Run the simple FastAPI backend"""
    print(f"Starting simple FastAPI backend on port {port}...")
    # Import the simple backend module
    sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
    from backend.alternatives.simple_backend import app
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=port, reload=True)

def run_minimal_backend(port):
    """Run the minimal backend (no FastAPI)"""
    print(f"Starting minimal backend on port {port}...")
    # Import the minimal backend module
    sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
    
    # Modify the port in the minimal backend before importing
    minimal_backend_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)),
        "backend", "alternatives", "minimal_backend.py"
    )
    
    with open(minimal_backend_path, 'r') as f:
        code = f.read()
    
    # Replace the port in the run_server function
    if "run_server(port=" in code:
        code = code.replace("run_server(port=8000)", f"run_server(port={port})")
    
    # Execute the modified code
    exec(compile(code, minimal_backend_path, 'exec'))

def main():
    parser = argparse.ArgumentParser(description="Run the Language Learning Flashcard Generator backend server")
    parser.add_argument(
        "--backend", 
        choices=["main", "simple", "minimal"], 
        default="main",
        help="Type of backend to run (default: main)"
    )
    parser.add_argument(
        "--port", 
        type=int, 
        default=8000,
        help="Port to run the server on (default: 8000)"
    )
    
    args = parser.parse_args()
    
    if args.backend == "main":
        run_main_backend(args.port)
    elif args.backend == "simple":
        run_simple_backend(args.port)
    elif args.backend == "minimal":
        run_minimal_backend(args.port)
    else:
        print(f"Unknown backend type: {args.backend}")
        sys.exit(1)

if __name__ == "__main__":
    main()
