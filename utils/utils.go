package utils

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"os"
)

func GenSha1(data []byte) string {
	_sha1 := sha1.New()
	return hex.EncodeToString(_sha1.Sum(data))
}

func CalculateSha1(reader io.Reader) string {
	_sha1 := sha1.New()
	n, err := io.Copy(_sha1, reader)
	log.Println("CalculateSha1 n:", n)
	if err != nil {
		log.Println("CalculateSha1 err:", err)
	}
	return url.PathEscape(base64.StdEncoding.EncodeToString(_sha1.Sum(nil)))
}

func GenByte(fileDst string) {
	file, err := os.Open(fileDst)
	if err != nil {
		log.Println("打开文件错误")
		return
	}
	defer file.Close()

	bfRd := bufio.NewReader(file)
	buf := make([]byte, 1024)

	for {
		n, err := bfRd.Read(buf)
		if err != nil {
			log.Print("read buf err:", err)
			return
		}
		if n <= 0 {
			break
		}

		fmt.Println(string(buf))
		for _, b := range buf {
			fmt.Println(b)
			fmt.Println(string(b))
		}

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		// http.Post(targetUrl+"&index=1", "multipart/form-data", bytes.NewReader(bufCopied))
	}
}

func DeleteQuery(req *http.Request, delKey ...string) {
	temp := req.URL.Query()
	for _, key := range delKey {
		temp.Del(key)
	}
	req.URL.RawQuery = temp.Encode()
}

func AddQuery(req *http.Request, values ...string) {
	length := len(values)
	if length%2 == 1 {
		return
	}
	temp := req.URL.Query()
	for i := 0; i < length/2; i++ {
		temp.Add(values[i], values[i+1])
	}
	req.URL.RawQuery = temp.Encode()
}
