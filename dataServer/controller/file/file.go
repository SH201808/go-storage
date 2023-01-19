package file

import (
	"crypto/sha1"
	"encoding/base64"
	"file-server/dataServer/locate"
	response "file-server/models/Response"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	log.Println(fileName)

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("form file err: ", err)
		c.JSON(http.StatusInternalServerError, response.Err("dataServer form file err"))
		return
	}
	dst := "../../dataServer/fileStore/" + c.Request.Header.Get("Digests")
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, response.Success("save file success"))
}

func Download(c *gin.Context) {
	fileHash := c.Request.Header.Get("Hash")
	log.Println("download hash:" + fileHash)
	filePath := getFile(fileHash)
	if filePath == "" {
		c.Status(http.StatusNotFound)
		return
	}
	sendFile(c.Writer, filePath)
}

func getFile(name string) string {
	filePath := locate.FileLoc + name + ".*"
	files, _ := filepath.Glob(filePath)
	if len(files) != 1 {
		log.Println("未找到")
		return ""
	}
	file := files[0]
	h := sha1.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, name+".")[1]
	if d != hash {
		log.Printf("object hash mismatch:%s\n, getHash:%s\n", hash, d)
		locate.Delete(file)
		return ""
	}
	return file
}

func sendFile(w io.Writer, filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()
	io.Copy(w, file)
}
