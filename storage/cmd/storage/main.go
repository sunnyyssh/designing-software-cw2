package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

func main() {
	hash := md5.Sum([]byte{2})
	enc := base64.StdEncoding.EncodeToString(hash[:])
	fmt.Println(enc, len(enc))
}
