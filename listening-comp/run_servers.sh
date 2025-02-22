#!/bin/bash

# Get the directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Export the parent directory to PYTHONPATH
export PYTHONPATH=$DIR:$PYTHONPATH

# Start the backend server
cd $DIR/backend
python main.py &

# Wait a bit for the backend to start
sleep 2

# Start the frontend
cd $DIR/frontend
streamlit run main.py
