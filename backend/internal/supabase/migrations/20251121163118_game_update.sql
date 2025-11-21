CREATE TYPE exercise_type AS ENUM ('game', 'pdf');

ALTER TABLE game_content
ADD COLUMN exercise_type exercise_type NOT NULL DEFAULT 'game';

ALTER TABLE game_content
ADD COLUMN applicable_game_types question_type[] DEFAULT '{}';
