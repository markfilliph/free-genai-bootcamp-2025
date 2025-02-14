-- Insert sample groups
INSERT INTO groups (name) VALUES 
    ('Basic Greetings'),
    ('Numbers'),
    ('Colors'),
    ('Family Members');

-- Insert sample words
INSERT INTO words (japanese, romaji, english, parts) VALUES 
    ('こんにちは', 'konnichiwa', 'hello', '{"type": "greeting", "formality": "neutral"}'),
    ('さようなら', 'sayounara', 'goodbye', '{"type": "greeting", "formality": "formal"}'),
    ('おはよう', 'ohayou', 'good morning', '{"type": "greeting", "formality": "informal"}'),
    ('一', 'ichi', 'one', '{"type": "number", "category": "cardinal"}'),
    ('二', 'ni', 'two', '{"type": "number", "category": "cardinal"}'),
    ('三', 'san', 'three', '{"type": "number", "category": "cardinal"}'),
    ('赤', 'aka', 'red', '{"type": "color", "category": "basic"}'),
    ('青', 'ao', 'blue', '{"type": "color", "category": "basic"}'),
    ('緑', 'midori', 'green', '{"type": "color", "category": "basic"}');

-- Link words to groups
INSERT INTO words_groups (word_id, group_id) VALUES 
    (1, 1), -- konnichiwa -> Basic Greetings
    (2, 1), -- sayounara -> Basic Greetings
    (3, 1), -- ohayou -> Basic Greetings
    (4, 2), -- ichi -> Numbers
    (5, 2), -- ni -> Numbers
    (6, 2), -- san -> Numbers
    (7, 3), -- aka -> Colors
    (8, 3), -- ao -> Colors
    (9, 3); -- midori -> Colors