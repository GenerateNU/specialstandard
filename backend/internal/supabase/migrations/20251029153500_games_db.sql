-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create GameContent Table
CREATE TABLE game_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(255) NOT NULL CHECK (category IN ( 'sequencing', 'following_directions',
                                                       'wh_questions', 'true_false',
                                                       'concepts_sorting' )),
    level INT NOT NULL CHECK ( level >= 0 AND level <= 12 ),
    options TEXT[] NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (category, level)
);

-- Create GameResult Table
CREATE TABLE game_result (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL,
    student_id UUID NOT NULL,
    content_id UUID NOT NULL,
    time_taken INTEGER NOT NULL CHECK ( time_taken >= 0 ),
    completed BOOLEAN DEFAULT FALSE,
    incorrect_tries INTEGER DEFAULT 0 CHECK ( incorrect_tries >= 0 ),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    FOREIGN KEY (session_id, student_id) REFERENCES session_student(session_id, student_id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES game_content(id) ON DELETE RESTRICT
);

CREATE TRIGGER update_game_content_updated_at BEFORE UPDATE ON game_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_result_updated_at BEFORE UPDATE ON game_result
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
