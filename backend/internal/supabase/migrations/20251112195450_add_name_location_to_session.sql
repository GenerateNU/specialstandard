-- Adding Name & Location to Session
ALTER TABLE session
ADD COLUMN session_name VARCHAR(255) NOT NULL,
ADD COLUMN location VARCHAR(255);