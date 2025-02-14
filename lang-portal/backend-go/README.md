# Language Learning Portal Backend

## Development Setup

### Prerequisites

1. **Go 1.21 or later**
   ```bash
   go version
   ```

2. **SQLite Dependencies**
   
   For Windows, you have two options:

   a. **Using WSL (Recommended)**:
   ```bash
   # Install WSL if not already installed
   wsl --install

   # After WSL is installed and you've rebooted, open WSL and run:
   sudo apt update
   sudo apt install gcc sqlite3
   ```

   b. **Using MinGW-w64 on Windows**:
   - Download and install MinGW-w64 from [https://www.mingw-w64.org/](https://www.mingw-w64.org/)
   - Add the MinGW-w64 bin directory to your PATH

3. **Mage (Build Tool)**
   ```bash
   go install github.com/magefile/mage@latest
   ```

### Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/markfilliph/free-genai-bootcamp-2025.git
   cd free-genai-bootcamp-2025/lang-portal/backend-go
   ```

2. **Install dependencies**
   ```bash
   mage install
   ```

3. **Initialize the database**
   ```bash
   mage initdb
   ```

4. **Start the server**
   ```bash
   mage run
   ```

### Available Mage Commands

- `mage install` - Install project dependencies
- `mage initdb` - Initialize database with schema and seed data
- `mage backup` - Create a database backup
- `mage restore` - Restore from the most recent backup
- `mage reset` - Reset database to initial state (creates backup first)
- `mage status` - Check database status and statistics
- `mage clean` - Remove generated files
- `mage run` - Start the server

### Development Notes

1. **Database Location**
   - The SQLite database is created as `words.db` in the project root
   - Backups are stored in `db/backups` with timestamps

2. **Using WSL**
   - After WSL setup, the project will be accessible at `/mnt/d/GenAI/free-genai-bootcamp-2025/lang-portal`
   - Run all commands from within WSL for better SQLite support

3. **Code Organization**
   - `cmd/server` - Main application entry point
   - `internal/models` - Data structures and database operations
   - `internal/handlers` - HTTP handlers
   - `internal/service` - Business logic
   - `db/migrations` - Database schema
   - `db/seeds` - Initial data

### Contributing

1. Create a feature branch
2. Make your changes
3. Run tests
4. Submit a pull request
