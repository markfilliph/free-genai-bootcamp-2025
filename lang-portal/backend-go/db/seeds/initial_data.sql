-- Insert initial study activities
INSERT INTO study_activities (name, thumbnail_url, description) VALUES
    ('Vocabulary Quiz', 'https://example.com/vocab-quiz.jpg', 'Practice your vocabulary with flashcards'),
    ('Writing Practice', 'https://example.com/writing.jpg', 'Practice writing Japanese characters'),
    ('Listening Exercise', 'https://example.com/listening.jpg', 'Improve your listening comprehension');

-- Insert sample groups
INSERT INTO groups (name) VALUES
    ('Basic Greetings'),
    ('Numbers 1-10'),
    ('Colors'),
    ('Days of the Week');

-- Insert sample words
INSERT INTO words (japanese, romaji, english, parts) VALUES
    ('こんにちは', 'konnichiwa', 'hello', '{"type": "greeting", "formality": "neutral"}'),
    ('さようなら', 'sayounara', 'goodbye', '{"type": "greeting", "formality": "formal"}'),
    ('おはよう', 'ohayou', 'good morning', '{"type": "greeting", "formality": "informal"}'),
    ('一', 'ichi', 'one', '{"type": "number"}'),
    ('二', 'ni', 'two', '{"type": "number"}'),
    ('三', 'san', 'three', '{"type": "number"}'),
    ('赤', 'aka', 'red', '{"type": "color"}'),
    ('青', 'ao', 'blue', '{"type": "color"}'),
    ('緑', 'midori', 'green', '{"type": "color"}');

-- Link words to groups
INSERT INTO words_groups (word_id, group_id)
SELECT w.id, g.id
FROM words w, groups g
WHERE 
    (w.english IN ('hello', 'goodbye', 'good morning') AND g.name = 'Basic Greetings')
    OR (w.english IN ('one', 'two', 'three') AND g.name = 'Numbers 1-10')
    OR (w.english IN ('red', 'blue', 'green') AND g.name = 'Colors');
