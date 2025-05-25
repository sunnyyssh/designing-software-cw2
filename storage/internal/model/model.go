package model

type MD5 string

type FileMeta struct {
	ID      int64  `json:"id"`
	MD5Hash MD5    `json:"md5_hash"`
	S3Key   string `json:"s3_key"`
	URL     string `json:"url"`
}
