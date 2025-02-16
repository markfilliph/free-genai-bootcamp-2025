# Language Learning Portal Backend

## Development Setup

### Prerequisites

1. **Go 1.21 or later**
   ```bash
   go version
   ```

2. **MySQL 9.2 or later**
   - Make sure MySQL server is running
   - Create a database named `lang_portal`:
     ```sql
     CREATE DATABASE IF NOT EXISTS lang_portal;
     ```

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

2. **Configure environment**
   ```bash
   # Copy the environment template
   cp .env.template .env
   
   # Edit .env with your MySQL credentials
   # Default values are:
   # DB_HOST=localhost
   # DB_PORT=3306
   # DB_USER=root
   # DB_PASSWORD=
   # DB_NAME=lang_portal
   ```

3. **Install dependencies**
   ```bash
   mage install
   ```

4. **Initialize the database**
   ```bash
   mage initdb
   ```

5. **Start the server**
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
   - The database is created as `lang_portal` in the MySQL server
   - Backups are stored in `db/backups` with timestamps

2. **Using WSL**
   - After WSL setup, the project will be accessible at `/mnt/d/GenAI/free-genai-bootcamp-2025/lang-portal`
   - Run all commands from within WSL for better MySQL support

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

### Language Portal API Testing Guide

## Running the Tests

1. Start the server:
   ```bash
   go run cmd/main.go
   ```

2. Run the test script:
   ```powershell
   ./test_endpoints.ps1
   ```

## Available Endpoints

### Dashboard
- GET `/api/dashboard/last-study-session`
- GET `/api/dashboard/study-progress`
- GET `/api/dashboard/quick-stats`

### Study Activities
- GET `/api/study-activities/:id`
- GET `/api/study-activities/:id/sessions`
- POST `/api/study-activities`

### Study Sessions
- GET `/api/study-sessions`
- GET `/api/study-sessions/:id`
- GET `/api/study-sessions/:id/words`
- POST `/api/study-sessions/:id/words/:word_id/review`

### Words
- GET `/api/words`
- GET `/api/words/:id`

### Groups
- GET `/api/groups`
- GET `/api/groups/:id`
- GET `/api/groups/:id/words`
- GET `/api/groups/:id/study-sessions`

## Manual Testing

You can also test individual endpoints using curl:

```bash
# Get all words with pagination
curl "http://localhost:8080/api/words?page=1&page_size=10"

# Create a study activity
curl -X POST http://localhost:8080/api/study-activities \
  -H "Content-Type: application/json" \
  -d '{"group_id": 1}'

# Review a word
curl -X POST http://localhost:8080/api/study-sessions/1/words/1/review \
  -H "Content-Type: application/json" \
  -d '{"correct": true}'
```

## Common Response Formats

All endpoints follow these response formats:

1. List endpoints:
   ```json
   {
     "items": [...],
     "pagination": {
       "current_page": 1,
       "page_size": 10,
       "total_items": 100,
       "total_pages": 10
     }
   }
   ```

2. Single item endpoints:
   ```json
   {
     "id": 1,
     "created_at": "2025-02-15T21:26:48-03:00",
     ...
   }
   ```

3. Error responses:
   ```json
   {
     "error": "Error message here"
   }
   ```

## Troubleshooting

1. If you get connection refused:
   - Make sure the server is running on port 8080
   - Check if there are any firewall issues

2. If you get database errors:
   - Verify MySQL is running
   - Check database connection settings in config.go
   - Try running `api/reset/full` to reset the database

3. For pagination issues:
   - Ensure page and page_size parameters are positive integers
   - Default page_size is 10 if not specified
