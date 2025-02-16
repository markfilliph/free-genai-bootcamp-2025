
-- Drop existing tables if they exist
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS word_review_items;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS words_groups;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS word_groups;
SET FOREIGN_KEY_CHECKS=1;

-- Create words table
CREATE TABLE IF NOT EXISTS words (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    japanese TEXT NOT NULL,
    romaji TEXT NOT NULL,
    english TEXT NOT NULL,
    parts JSON
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create word_groups table
CREATE TABLE IF NOT EXISTS word_groups (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create words_groups join table
CREATE TABLE IF NOT EXISTS words_groups (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    word_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create study_sessions table
CREATE TABLE IF NOT EXISTS study_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create study_activities table
CREATE TABLE IF NOT EXISTS study_activities (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    study_session_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create word_review_items table
CREATE TABLE IF NOT EXISTS word_review_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    word_id BIGINT NOT NULL,
    study_activity_id BIGINT NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create indexes for better query performance
CREATE INDEX idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX idx_words_groups_group_id ON words_groups(group_id);
CREATE INDEX idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX idx_word_review_items_word_id ON word_review_items(word_id);
CREATE INDEX idx_word_review_items_activity_id ON word_review_items(study_activity_id);