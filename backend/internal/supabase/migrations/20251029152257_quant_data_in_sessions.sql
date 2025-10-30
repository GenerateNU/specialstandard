CREATE TYPE category AS ENUM ('visual_cue', 'verbal_cue', 'gestural_cue', 'engagement');
CREATE TYPE response_level AS ENUM ('minimal', 'moderate', 'maximal', 'low', 'high');

ALTER TABLE session_student
DROP CONSTRAINT "session_student_pkey",
ADD COLUMN id SERIAL PRIMARY KEY;

CREATE TABLE session_rating (
  id SERIAL PRIMARY KEY,
  session_student_id INT REFERENCES session_student(id),
  category category,
  level response_level,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE session_rating
ADD CONSTRAINT unique_session_student_category 
UNIQUE (session_student_id, category);

CREATE INDEX idx_session_rating_session_student_id ON session_rating (session_student_id);