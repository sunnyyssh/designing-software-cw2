package model

type AnalysisResult struct {
	FileID       int64   `json:"file_id"`
	PlagiatedIDs []int64 `json:"plagiated_ids"`
	WordCount    int     `json:"word_count"`
	SymbolCount  int     `json:"symbol_count"`
	WordCloudURL string  `json:"word_cloud_url"`
}

type MD5 string

type FileMeta struct {
	ID      int64  `json:"id"`
	MD5Hash MD5    `json:"md5_hash"`
	S3Key   string `json:"s3_key"`
	URL     string `json:"url"`
}
