CREATE TABLE IF NOT EXISTS jobs (
    id uuid PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
