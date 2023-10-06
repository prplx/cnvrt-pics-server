CREATE TABLE IF NOT EXISTS files (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    job_id bigserial NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    name varchar(255) NOT NULL
);
