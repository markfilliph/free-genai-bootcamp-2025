#!/bin/bash

# Start the FastAPI backend server
echo "Starting backend server..."
uvicorn backend.api:app --host 0.0.0.0 --port 8000 &

# Wait a bit for the backend to start
sleep 2

# Start the Streamlit frontend
echo "Starting frontend..."
streamlit run frontend/main.py
