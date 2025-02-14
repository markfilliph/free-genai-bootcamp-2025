# Database Backups Directory

This directory contains SQLite database backups created by the Mage backup task.

## Backup Format
- Backups are named in the format: `words_YYYYMMDD_HHMMSS.db`
- Example: `words_20250214_113000.db`

## Managing Backups

Use the following Mage commands to manage backups:

```bash
# Create a backup
mage backup

# Restore from the most recent backup
mage restore

# Reset database (creates a backup before resetting)
mage reset

# Check database status
mage status
```

Note: Backups are automatically created before potentially destructive operations like database reset.
