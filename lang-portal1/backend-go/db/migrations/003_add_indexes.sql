-- Add indexes for words table
CREATE INDEX IF NOT EXISTS idx_words_japanese ON words(japanese);
CREATE INDEX IF NOT EXISTS idx_words_english ON words(english);

-- Add indexes for words_groups table
CREATE INDEX IF NOT EXISTS idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX IF NOT EXISTS idx_words_groups_group_id ON words_groups(group_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_words_groups_unique ON words_groups(word_id, group_id);

-- Add indexes for study_sessions
CREATE INDEX IF NOT EXISTS idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_created_at ON study_sessions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_study_sessions_activity ON study_sessions(study_activity_id, created_at DESC);

-- Add indexes for study_activities
CREATE INDEX IF NOT EXISTS idx_study_activities_group ON study_activities(group_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_study_activities_session ON study_activities(study_session_id);

-- Add indexes for word_review_items
CREATE INDEX IF NOT EXISTS idx_word_reviews_session ON word_review_items(study_session_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_word_reviews_word ON word_review_items(word_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_word_reviews_stats ON word_review_items(word_id, correct);
