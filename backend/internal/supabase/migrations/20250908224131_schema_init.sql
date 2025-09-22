-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS sessions;

-- Create Therapist table
CREATE TABLE therapist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create Theme table (needs to be created before Resource due to FK)
CREATE TABLE theme (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_name VARCHAR(255) NOT NULL,
    month INTEGER CHECK (month >= 1 AND month <= 12),
    year INTEGER CHECK (year >= 2000 AND year <= 2500),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create Student table
CREATE TABLE student (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    therapist_id UUID NOT NULL,
    grade VARCHAR(20),
    dob DATE,
    iep TEXT DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (therapist_id) REFERENCES therapist(id) ON DELETE RESTRICT
);

-- Create Session table
CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_datetime TIMESTAMPTZ NOT NULL,
    end_datetime TIMESTAMPTZ NOT NULL,
    therapist_id UUID NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (therapist_id) REFERENCES therapist(id) ON DELETE RESTRICT,
    CHECK (end_datetime > start_datetime)
);

-- Create SessionStudent junction table
CREATE TABLE session_student (
    session_id UUID,
    student_id UUID,
    present BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (session_id, student_id),
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE
);

-- Create Resource table
CREATE TABLE resource (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_id UUID NOT NULL,
    grade_level VARCHAR(255),
    date DATE,
    type VARCHAR(50),
    title VARCHAR(100),
    category VARCHAR(100),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (theme_id) REFERENCES theme(id) ON DELETE RESTRICT
);

-- Create indexes for better performance
CREATE INDEX idx_student_therapist ON student(therapist_id);
CREATE INDEX idx_session_therapist ON session(therapist_id);
CREATE INDEX idx_session_datetime ON session(start_datetime, end_datetime);
CREATE INDEX idx_resource_theme ON resource(theme_id);
CREATE INDEX idx_resource_grade ON resource(grade_level);
CREATE INDEX idx_theme_month_year ON theme(month, year);

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic updated_at timestamp updates
CREATE TRIGGER update_therapist_updated_at BEFORE UPDATE ON therapist
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_theme_updated_at BEFORE UPDATE ON theme
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_updated_at BEFORE UPDATE ON student
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_session_updated_at BEFORE UPDATE ON session
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_session_student_updated_at BEFORE UPDATE ON session_student
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resource_updated_at BEFORE UPDATE ON resource
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();