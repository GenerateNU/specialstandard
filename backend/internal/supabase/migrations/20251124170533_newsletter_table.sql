CREATE TABLE IF NOT EXISTS newsletter (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    s3_url TEXT NOT NULL
);
