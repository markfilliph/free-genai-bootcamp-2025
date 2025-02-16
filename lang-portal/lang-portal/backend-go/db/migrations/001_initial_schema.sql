
-- Drop existing tables if they exist
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS word_review_items;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS words_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words;
SET FOREIGN_KEY_CHECKS=1;

-- Create words table
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    japanese TEXT NOT NULL,
    romaji TEXT NOT NULL,
    english TEXT NOT NULL,
    parts TEXT NOT NULL -- JSON field
);

-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- Create words_groups join table
CREATE TABLE IF NOT EXISTS words_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    UNIQUE(word_id, group_id)
);

-- Create study_sessions table
CREATE TABLE IF NOT EXISTS study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    study_activity_id INTEGER,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE SET NULL
);

-- Create study_activities table
CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER,
    group_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- Add foreign key for study_activities.study_session_id after study_sessions table exists
ALTER TABLE study_activities ADD FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE;

-- Create word_review_items table
CREATE TABLE IF NOT EXISTS word_review_items (
    word_id INTEGER NOT NULL,
    study_session_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (word_id, study_session_id),
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX IF NOT EXISTS idx_words_groups_group_id ON words_groups(group_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_activity_id ON study_sessions(study_activity_id);
CREATE INDEX IF NOT EXISTS idx_study_activities_session_id ON study_activities(study_session_id);
CREATE INDEX IF NOT EXISTS idx_word_review_items_session_id ON word_review_items(study_session_id);