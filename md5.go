package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// EncodeStringMD5 は、MD5エンコードした文字列を返します。
func encode(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	encodeStr := hex.EncodeToString(h.Sum(nil))

	return encodeStr
}
