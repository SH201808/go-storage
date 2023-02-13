package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetOffsetFromHeader(header http.Header) int64 {
	byteRange := header.Get("range")
	if len(byteRange) < 7 {
		log.Println(len(byteRange))
		return 0
	}
	if byteRange[:6] != "bytes=" {
		log.Println(byteRange[:6])
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}

func GetSizeFromHeader(header http.Header) int64 {
	fileSize := header.Get("content-length")
	size, err := strconv.Atoi(fileSize)
	if err != nil {
		log.Println("size illegal")
		return 0
	}
	return int64(size)
}

func GetHashFromHeader(header http.Header) string {
	// todo 上传携带空格的数据存在转义问题
	hash := header.Get("Digests")
	return hash
}

func GetObjectFromHeader(header http.Header) string {
	object := header.Get("object")
	return object
}

func GetEncodingFromHeader(header http.Header) []string {
	encoding := header["Accept-Encoding"]
	return encoding
}
