#!/bin/bash

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${YELLOW}[DEV]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a port is in use
is_port_in_use() {
    netstat -tuln | grep ":$1 " > /dev/null
    return $?
}

# Function to kill processes
kill_processes() {
    print_status "Cleaning up previous processes..."
    
    # Kill any running streamlit processes
    if pgrep -f "streamlit run" > /dev/null; then
        pkill -9 -f "streamlit run"
        print_success "Killed Streamlit processes"
    fi
    
    # Kill any running backend processes
    if pgrep -f "python.*main\.py" > /dev/null; then
        pkill -9 -f "python.*main\.py"
        print_success "Killed backend processes"
    fi
    
    # Small delay to ensure processes are killed
    sleep 2
    
    # Check common ports
    local ports=(8000 8501 8502 8503 8504 8505 8506 8507 8508)
    local busy_ports=()
    
    for port in "${ports[@]}"; do
        if is_port_in_use "$port"; then
            busy_ports+=($port)
        fi
    done
    
    if [ ${#busy_ports[@]} -ne 0 ]; then
        print_error "Some ports are still in use: ${busy_ports[*]}"
        print_status "Attempting to force kill processes on these ports..."
        for port in "${busy_ports[@]}"; do
            fuser -k "$port/tcp" 2>/dev/null
        done
        sleep 2
    fi
}

# Function to start servers
start_servers() {
    print_status "Starting backend server..."
    cd backend
    python main.py &
    backend_pid=$!
    cd ..
    
    # Wait for backend to start
    print_status "Waiting for backend to initialize..."
    for i in {1..30}; do
        if netstat -tuln | grep ':8000 ' > /dev/null; then
            # Give it a moment to fully initialize
            sleep 2
            print_success "Backend started successfully"
            break
        fi
        if [ $i -eq 30 ]; then
            print_error "Backend failed to start"
            kill_processes
            exit 1
        fi
        sleep 1
    done
    
    print_status "Starting frontend..."
    cd frontend
    streamlit run main.py &
    frontend_pid=$!
    cd ..
    
    print_success "Development environment started!"
    print_status "Backend running on http://localhost:8000"
    print_status "Frontend will be available on http://localhost:8501"
    print_status "Press Ctrl+C to stop all servers"
    
    # Wait for Ctrl+C
    trap 'kill_processes; exit 0' SIGINT
    wait
}

# Main execution
case "$1" in
    "stop")
        kill_processes
        print_success "All processes stopped"
        ;;
    "start")
        kill_processes
        start_servers
        ;;
    "restart")
        kill_processes
        start_servers
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        echo "  start   - Start the development environment"
        echo "  stop    - Stop all related processes"
        echo "  restart - Restart the development environment"
        exit 1
        ;;
esac
