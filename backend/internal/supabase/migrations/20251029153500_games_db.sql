-- GameType ENUM
CREATE TYPE game_category AS ENUM (
    'receptive_language',
    'expressive_language',
    'social_pragmatic_language',
    'speech'
);

-- QuestionType Enum
CREATE TYPE question_type AS ENUM (
    'sequencing',
    'following_directions',
    'wh_questions',
    'true_false',
    'concepts_sorting',
    'fill_in_the_blank',
    'categorical_language',
    'emotions',
    'teamwork_talk',
    'express_excitement_interest',
    'fluency',
    'articulation_s',
    'articulation_l'
);

-- Create GameContent Table
CREATE TABLE game_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_id UUID NOT NULL,
    week INT NOT NULL CHECK ( week >= 0 AND week <= 6 ),
    category game_category,
    question_type question_type NOT NULL,
    difficulty_level INT NOT NULL CHECK ( difficulty_level >= 1 ),
    question TEXT NOT NULL,
    options TEXT[] NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    FOREIGN KEY (theme_id) REFERENCES theme(id) ON DELETE RESTRICT
);

-- Create GameResult Table
CREATE TABLE game_result (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_student_id INT NOT NULL,
    content_id UUID NOT NULL,
    time_taken_sec INTEGER NOT NULL CHECK ( time_taken_sec >= 0 ),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    count_of_incorrect_attempts INTEGER NOT NULL DEFAULT 0 CHECK ( count_of_incorrect_attempts >= 0 ),
    incorrect_attempts TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    FOREIGN KEY (session_student_id) REFERENCES session_student(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES game_content(id) ON DELETE RESTRICT
);

CREATE TRIGGER update_game_content_updated_at BEFORE UPDATE ON game_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_result_updated_at BEFORE UPDATE ON game_result
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();