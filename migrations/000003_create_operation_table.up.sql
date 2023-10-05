CREATE TABLE IF NOT EXISTS operations (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    job_id uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    file_id bigserial NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    format varchar(6) NOT NULL,
    quality smallint NOT NULL,
    fileName varchar(50) NOT NULL,
    width int,
    height int
);
