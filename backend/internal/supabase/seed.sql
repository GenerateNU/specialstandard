INSERT INTO sessions (therapist_id, session_date, start_time, end_time, notes)
VALUES
    (uuid_generate_v4(), '2023-09-01', '10:00:00', '11:00:00', 'Initial consultation'),
    (uuid_generate_v4(), '2023-09-02', '14:00:00', '15:00:00', 'Follow-up session');
