CREATE TABLE file_analysis (
    file_id BIGINT NOT NULL PRIMARY KEY,
    plagiated_ids BIGINT[],
    word_count INTEGER NOT NULL,
    symbol_count INTEGER NOT NULL,
    word_cloud_url VARCHAR(256) NOT NULL
);
