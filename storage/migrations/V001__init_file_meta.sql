CREATE TABLE IF NOT EXISTS file_meta (
    id SERIAL PRIMARY KEY,
    md5_hash CHAR(24),
    url VARCHAR(256),
);
