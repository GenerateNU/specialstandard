CREATE TYPE exercise_type AS ENUM ('game', 'pdf');
CREATE TYPE game_type AS ENUM ('drag and drop', 'spinner', 'word/image matching', 'flashcards');

ALTER TABLE game_content
ADD COLUMN exercise_type exercise_type NOT NULL DEFAULT 'game';

ALTER TABLE game_content
ADD COLUMN applicable_game_types game_type[] DEFAULT '{}';
