-- Insert sample word_groups
INSERT INTO word_groups (name) VALUES 
    ('Basic Greetings'),
    ('Numbers'),
    ('Colors'),
    ('Family Members'),
    ('Animals'),
    ('Actions'),
    ('Colors'),
    ('Basic Adjectives');

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
    ('緑', 'midori', 'green', '{"type": "color", "category": "basic"}'),
    ('犬', 'inu', 'dog', '{"type": "noun", "category": "animals"}'),
    ('猫', 'neko', 'cat', '{"type": "noun", "category": "animals"}'),
    ('鳥', 'tori', 'bird', '{"type": "noun", "category": "animals"}'),
    ('食べる', 'taberu', 'to eat', '{"type": "verb", "category": "actions", "group": "ichidan"}'),
    ('飲む', 'nomu', 'to drink', '{"type": "verb", "category": "actions", "group": "godan"}'),
    ('走る', 'hashiru', 'to run', '{"type": "verb", "category": "actions", "group": "godan"}'),
    ('赤い', 'akai', 'red', '{"type": "adjective", "category": "colors"}'),
    ('青い', 'aoi', 'blue', '{"type": "adjective", "category": "colors"}'),
    ('黄色い', 'kiiroi', 'yellow', '{"type": "adjective", "category": "colors"}'),
    ('大きい', 'ookii', 'big', '{"type": "adjective", "category": "size"}');

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
    (9, 3), -- midori -> Colors
    (10, 4), -- dog -> Animals
    (11, 4), -- cat -> Animals
    (12, 4), -- bird -> Animals
    (13, 5), -- eat -> Actions
    (14, 5), -- drink -> Actions
    (15, 5), -- run -> Actions
    (16, 6), -- red -> Colors
    (17, 6), -- blue -> Colors
    (18, 6), -- yellow -> Colors
    (19, 7); -- big -> Basic Adjectives

-- Create a sample study activity
INSERT INTO study_activities (group_id) VALUES (4); -- Study activity for Animals group

-- Create a study session for the activity
INSERT INTO study_sessions (group_id, study_activity_id) 
SELECT 4, id FROM study_activities ORDER BY id DESC LIMIT 1;

-- Record some word reviews
INSERT INTO word_review_items (word_id, study_session_id, correct) 
SELECT 10, s.id, 1 FROM study_sessions s ORDER BY s.id DESC LIMIT 1
UNION ALL
SELECT 11, s.id, 0 FROM study_sessions s ORDER BY s.id DESC LIMIT 1
UNION ALL
SELECT 12, s.id, 1 FROM study_sessions s ORDER BY s.id DESC LIMIT 1;