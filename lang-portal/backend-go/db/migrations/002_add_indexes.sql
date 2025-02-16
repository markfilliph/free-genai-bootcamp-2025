-- Add indexes for better query performance

-- Words table indexes
CREATE INDEX idx_words_created_at ON words(created_at);
CREATE INDEX idx_words_updated_at ON words(updated_at);
CREATE INDEX idx_words_original ON words(original);
CREATE INDEX idx_words_translation ON words(translation);

-- Groups table indexes
CREATE INDEX idx_groups_created_at ON groups(created_at);
CREATE INDEX idx_groups_updated_at ON groups(updated_at);
CREATE INDEX idx_groups_name ON groups(name);

-- Study sessions table indexes
CREATE INDEX idx_study_sessions_created_at ON study_sessions(created_at);
CREATE INDEX idx_study_sessions_ended_at ON study_sessions(ended_at);
CREATE INDEX idx_study_sessions_group_id ON study_sessions(group_id);

-- Word reviews table indexes
CREATE INDEX idx_word_reviews_created_at ON word_reviews(created_at);
CREATE INDEX idx_word_reviews_word_id ON word_reviews(word_id);
CREATE INDEX idx_word_reviews_study_session_id ON word_reviews(study_session_id);
CREATE INDEX idx_word_reviews_correct ON word_reviews(correct);

-- Group words table indexes
CREATE INDEX idx_group_words_group_id_word_id ON group_words(group_id, word_id);

-- Down migration
-- DROP INDEX idx_words_created_at ON words;
-- DROP INDEX idx_words_updated_at ON words;
-- DROP INDEX idx_words_original ON words;
-- DROP INDEX idx_words_translation ON words;
-- DROP INDEX idx_groups_created_at ON groups;
-- DROP INDEX idx_groups_updated_at ON groups;
-- DROP INDEX idx_groups_name ON groups;
-- DROP INDEX idx_study_sessions_created_at ON study_sessions;
-- DROP INDEX idx_study_sessions_ended_at ON study_sessions;
-- DROP INDEX idx_study_sessions_group_id ON study_sessions;
-- DROP INDEX idx_word_reviews_created_at ON word_reviews;
-- DROP INDEX idx_word_reviews_word_id ON word_reviews;
-- DROP INDEX idx_word_reviews_study_session_id ON word_reviews;
-- DROP INDEX idx_word_reviews_correct ON word_reviews;
-- DROP INDEX idx_group_words_group_id_word_id ON group_words;
