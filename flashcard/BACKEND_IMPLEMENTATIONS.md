# Backend Implementations Guide

This document explains the different backend implementations available in the Language Learning Flashcard Generator project and when to use each one.

## Overview

The project includes three different backend implementations:

1. **Main FastAPI Backend** (`backend/main.py`)
2. **Minimal Backend** (`backend/alternatives/minimal_backend.py`)
3. **Simple FastAPI Backend** (`backend/alternatives/simple_backend.py`)

Each implementation serves a specific purpose in development, testing, and production scenarios.

## Unified Server Launcher

A unified server launcher script (`run_server.py`) is provided to make it easy to switch between backend implementations:

```bash
# Run the main FastAPI backend (default)
python run_server.py

# Run the simple FastAPI backend
python run_server.py --backend simple

# Run the minimal backend
python run_server.py --backend minimal

# Specify a custom port
python run_server.py --backend main --port 8080
```

## 1. Main FastAPI Backend

**Location**: `backend/main.py`

**Purpose**: This is the primary backend implementation intended for production use.

**Features**:
- Full-featured FastAPI implementation
- Modular architecture with separate routes, models, and services
- Complete authentication system
- Ollama integration for language generation
- Database models with SQLAlchemy ORM

**When to use**:
- For production deployment
- When developing new features that require the full backend stack
- When testing the complete application flow

**How to start**:
```bash
python run_server.py --backend main
```
or
```bash
uvicorn backend.main:app --reload --port 8000
```

## 2. Minimal Backend

**Location**: `backend/alternatives/minimal_backend.py`

**Purpose**: A lightweight implementation using standard Python libraries, designed for testing without external dependencies.

**Features**:
- Uses only standard Python libraries (no FastAPI, SQLAlchemy, etc.)
- Direct SQLite3 connection
- Simplified API that mimics the main backend
- No external dependencies required

**When to use**:
- For testing core functionality without dependency complications
- When you need to isolate backend issues
- In environments where installing all dependencies is problematic
- For CI/CD pipelines with minimal requirements

**How to start**:
```bash
python run_server.py --backend minimal
```

## 3. Simple FastAPI Backend

**Location**: `backend/alternatives/simple_backend.py`

**Purpose**: A simplified FastAPI implementation that serves as a middle ground between the full and minimal backends.

**Features**:
- Uses FastAPI but with a simplified structure
- Direct SQLite3 connection (no ORM)
- All code in a single file for easier debugging
- Fewer dependencies than the main backend

**When to use**:
- For development when you need FastAPI features but want a simpler codebase
- When debugging specific FastAPI-related issues
- As a reference implementation for the main backend

**How to start**:
```bash
python run_server.py --backend simple
```

## Choosing the Right Backend

- **For normal development**: Use the main FastAPI backend
- **For testing**: Use the minimal backend
- **For debugging or learning**: Use the simple FastAPI backend

## Port Configuration

All backends are configured to run on port 8000 by default and include CORS settings to allow requests from the frontend running on ports 5173, 8080, and 8083.

You can specify a custom port when using the unified server launcher:
```bash
python run_server.py --backend main --port 8080
```

If you're using a different port than the ones configured in the CORS settings, you'll need to update:
1. The CORS configuration in the backend files
2. The API URL in the frontend configuration
